/**
 * @file main.c
 * @author Taka Wang
 * @brief main flow with zmq related functions
*/

#include "main.h"
#include "mb.h"
#include <czmq.h>

/* ==================================================
 *  global variable
================================================== */

extern long tcp_conn_timeout_usec;                  // from mbtcp.c
int enable_syslog = 0;                              // syslog flag
static char *ipc_sub = "ipc:///tmp/to.modbus";
static char *ipc_pub = "ipc:///tmp/from.modbus";

static const char *env_conf_path = "CONF_MODBUSD";  // environment variable name
static cJSON *config_json;                          // config in cJSON object format
static char *config_fname = NULL;                   // config filename
void *zmq_pub = NULL;                               // zmq publish socket

/* ==================================================
 *  static functions
================================================== */

/**
 * @brief Load configuration file
 *
 * @param fname Configuration file name.
 * @param ptr_config Pointer to config cJSON object.
 * @return Void.
 */
static void load_config (const char *fname, cJSON **ptr_config)
{
    BEGIN (enable_syslog);

    if (file_to_json (fname, ptr_config) < 0)
    {
        ERR (enable_syslog, "Failed to parse setting json: %s!", config_fname);
    }
    else
    {
        enable_syslog = json_get_int (config_json, KEY_LOGGER);
        cJSON * ipc = cJSON_GetObjectItem (config_json, KEY_IPC_TYPE);
        ipc_sub = json_get_char (ipc, KEY_IPC_SUB);
        ipc_pub = json_get_char (ipc, KEY_IPC_PUB);
        cJSON * mbtcp = cJSON_GetObjectItem (config_json, KEY_MB_TYPE);
        tcp_conn_timeout_usec = json_get_long (mbtcp, KEY_MBTCP_CONN_TIMEOUT);
    }

    END (enable_syslog);
}

/**
 * @brief Save configuration to file
 *
 * @param fname Configuration file name.
 * @param config Config cJSON object.
 * @return Void.
 */
static void save_config (const char *fname, cJSON *config)
{
    BEGIN (enable_syslog);

    cJSON * mbtcp = cJSON_GetObjectItem (config, KEY_MB_TYPE);
    json_set_double (mbtcp, KEY_MBTCP_CONN_TIMEOUT, (double) tcp_conn_timeout_usec);
    json_to_file (fname, config);

    END (enable_syslog);
}

/**
 * @brief Generic zmq response sender for modbus
 *
 * @param cmd Request command number.
 * @param json_resp Response string in JSON format.
 * @return Void.
 */
static void send_modbus_resp (cmd_s cmd, char *json_resp)
{
    BEGIN (enable_syslog);

    if (zmq_pub != NULL)
    {
        zmsg_t * zmq_resp = zmsg_new ();
        zmsg_addstrf (zmq_resp, "%d", cmd);  // frame 1: cmd
        zmsg_addstr (zmq_resp, json_resp);   // frame 2: resp
        // send zmq msg
        zmsg_send (&zmq_resp, zmq_pub);
        // cleanup zmsg
        zmsg_destroy (&zmq_resp);
    }
    else
    {
        ERR (enable_syslog, "NULL publisher");
    }

    END (enable_syslog);
}


// entry
int main (int argc, char *argv[])
{
    LOG (enable_syslog, "modbusd version: %s", VERSION);

    // @get environemnt variable; 12-Factor
    char *env = getenv (env_conf_path);
    if (env != NULL)
    {
        config_fname = env;
        LOG (enable_syslog, "Get config path from environment variable: %s", env);
    }
    else
    {
        config_fname = argc > 1 ? argv[1] : "./modbusd.json";
        LOG (enable_syslog, "Get config path from flag: %s", config_fname);
    }

    // @load config
    load_config (config_fname, &config_json);

    // @setup zmq
    zctx_t *zmq_context = zctx_new ();                  // init zmq context
    void *zmq_sub = zsocket_new (zmq_context, ZMQ_SUB); // init zmq subscriber: zmq_sub
    zsocket_bind (zmq_sub, ipc_sub);                    // bind zmq subscriber
    zsocket_set_subscribe (zmq_sub, "");                // set zmq subscriber filter
    
    zmq_pub = zsocket_new (zmq_context, ZMQ_PUB);       // init zmq publisher: zmq_pub
    zsocket_bind (zmq_pub, ipc_pub);                    // bind zmq publisher

    LOG (enable_syslog, "start request listener");

    while (!zctx_interrupted) // handle ctrl+c
    {
        zmsg_t *msg = zmsg_recv (zmq_sub); // recv zmsg
        if (msg != NULL)
        {
            // get request mode (ex. tcp, rtu, others)
            zframe_t *frame_mode = zmsg_pop (msg);
            char *mode = zframe_strdup (frame_mode);

            // get request json string
            zframe_t *frame_json = zmsg_pop (msg);
            char *req_json_string = zframe_strdup (frame_json);

            // cleanup zmsg releated resources
            zmsg_destroy (&msg);
            zframe_destroy (&frame_mode);
            zframe_destroy (&frame_json);
            
            LOG (enable_syslog, "recv msg: %s, %s\n", mode, req_json_string);

            // parse json string
            cJSON *req_json_obj = cJSON_Parse (req_json_string);
            char * tid = json_get_char (req_json_obj, "tid");
            
            if (req_json_obj != NULL)
            {
                cmd_s cmd = json_get_int (req_json_obj, "cmd");
                
                // @handle modbus tcp requests
                if (strcmp (mode, "tcp") == 0)
                {
                    LOG (enable_syslog, "@@@req: %d", cmd);

                    switch (cmd)
                    {
                        case FC1:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_read_bit_req));
                            break;
                        case FC2:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_read_bit_req));
                            break;
                        case FC3:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_read_reg_req));
                            break;
                        case FC4:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_read_reg_req));
                            break;
                        case FC5:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_single_write_req));
                            break;
                        case FC6:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_single_write_req));
                            break;
                        case FC15:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_multi_write_req));
                            break;
                        case FC16:
                            send_modbus_resp (cmd, mbtcp_cmd_hanlder (cmd, req_json_obj, mbtcp_multi_write_req));
                            break;
                        case SET_TCP_TIMEOUT:
                            send_modbus_resp (cmd, mbtcp_set_response_timeout (tid, json_get_long (req_json_obj, "timeout")));
                            break;
                        case GET_TCP_TIMEOUT:
                            send_modbus_resp (cmd, mbtcp_get_response_timeout (tid));
                            break;
                        default: 
                            send_modbus_resp (cmd, set_modbus_fail_resp_str (tid, "unsupport request"));
                            break;
                    }
                }
                // @handle modbus rtu requests
                else if (strcmp (mode, "rtu") == 0)
                {
                    LOG (enable_syslog, "rtu:%d", cmd);
                    send_modbus_resp (-1, set_modbus_fail_resp_str (tid, "TODO"));
                }
                // @unkonw mode
                else
                {
                    send_modbus_resp (-1, set_modbus_fail_resp_str (tid, "unsupport mode"));
                }
            }
            else
            {
                send_modbus_resp (-2, set_modbus_fail_resp_str (tid, "Fail to parse command string"));
            }
            
            // @cleanup cJson object (auto mode)
            cJSON_Delete (req_json_obj);
        }
        else
        {
            //ERR(enable_syslog, "Recv null message");
        }
    }
    
    // @resource clean up
    LOG (enable_syslog, "housekeeping");
    zctx_destroy (&zmq_context);
    save_config (config_fname, config_json); 
}
