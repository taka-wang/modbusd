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

int enable_syslog  = 1;                 // syslog flag
static cJSON * config_json;             // config in cJSON object format
static char *config_fname = NULL;       // config filename
static char *ipc_sub = "ipc:///tmp/to.modbus";
static char *ipc_pub = "ipc:///tmp/from.modbus";
extern long tcp_conn_timeout_usec;  // from mb.c


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
static void load_config(const char *fname, cJSON ** ptr_config)
{
    BEGIN(enable_syslog);
    if (file_to_json(fname, ptr_config) < 0)
    {
        ERR(enable_syslog, "Failed to parse setting json: %s! Bye!", config_fname);
        exit(EXIT_FAILURE);
    }
    else
    {
        enable_syslog = json_get_int(config_json, "syslog");
        cJSON * zmq = cJSON_GetObjectItem(config_json, "zmq");
        ipc_sub = json_get_char(zmq, "sub");
        ipc_pub = json_get_char(zmq, "pub");
        cJSON * mbtcp = cJSON_GetObjectItem(config_json, "mbtcp");
        tcp_conn_timeout_usec = json_get_long(mbtcp, "connect_timeout");
    }
    END(enable_syslog);
}

/**
 * @brief Save configuration to file
 *
 * @param fname Configuration file name.
 * @param config Config cJSON object.
 * @return Void.
 */
static void save_config(const char *fname, cJSON * config)
{
    BEGIN(enable_syslog);
    cJSON * mbtcp = cJSON_GetObjectItem(config, "mbtcp");
    json_set_double(mbtcp, "connect_timeout", (double)tcp_conn_timeout_usec);
    json_to_file(fname, config);
    END(enable_syslog);
}

/**
 * @brief Generic zmq response sender for modbus
 *
 * @param pub Zmq publisher.
 * @param cmd Request command number.
 * @param json_resp Response string in JSON format.
 * @return Void.
 */
static void send_modbus_zmq_resp(void * pub, cmd_t cmd, char *json_resp)
{
    BEGIN(enable_syslog);
    if (pub != NULL)
    {
        zmsg_t * zmq_resp = zmsg_new();
        zmsg_addstrf(zmq_resp, "%d", cmd);// frame 1: cmd
        zmsg_addstr(zmq_resp, json_resp); // frame 2: resp
        zmsg_send(&zmq_resp, pub);        // send zmq msg
        zmsg_destroy(&zmq_resp);          // cleanup zmsg
    }
    else
    {
        ERR(enable_syslog, "NULL publisher");
    }
    END(enable_syslog);
}

// entry
int main(int argc, char *argv[])
{
    LOG(enable_syslog, "modbusd version: %s", VERSION);

    // @load config
    config_fname = argc > 1 ? argv[1] : "./modbusd.json";
    load_config(config_fname, &config_json);

    // @setup zmq
    zctx_t *zmq_context = zctx_new ();                  // init zmq context
    void *zmq_sub = zsocket_new (zmq_context, ZMQ_SUB); // init zmq subscriber: zmq_sub
    zsocket_bind (zmq_sub, ipc_sub);                    // bind zmq subscriber
    zsocket_set_subscribe (zmq_sub, "");                // set zmq subscriber filter
    void *zmq_pub = zsocket_new (zmq_context, ZMQ_PUB); // init zmq publisher: zmq_pub
    zsocket_bind (zmq_pub, ipc_pub);                    // bind zmq publisher

    LOG(enable_syslog, "start request listener");
    while (!zctx_interrupted) // handle ctrl+c
    {
        zmsg_t *msg = zmsg_recv(zmq_sub); // recv zmsg
        if (msg != NULL)
        {
            // get request mode (ex. tcp, rtu, others)
            zframe_t *frame_mode = zmsg_pop(msg);
            char *mode = zframe_strdup(frame_mode);

            // get request json string
            zframe_t *frame_json = zmsg_pop(msg);
            char *req_json_string = zframe_strdup(frame_json);

            // cleanup zmsg releated resources
            zmsg_destroy(&msg);
            zframe_destroy(&frame_mode);
            zframe_destroy(&frame_json);
            
            LOG(enable_syslog, "recv msg: %s, %s\n", mode, req_json_string);

            // parse json string
            cJSON *req_json_obj = cJSON_Parse(req_json_string);
            char * tid = json_get_char(req_json_obj, "tid");
            
            if (req_json_obj != NULL)
            {
                cmd_t cmd = json_get_int(req_json_obj, "cmd");
                
                // @handle modbus tcp requests
                if (strcmp(mode, "tcp") == 0)
                {
                    LOG(enable_syslog, "@@@req: %d", cmd);

                    switch (cmd)
                    {
                        case fc1:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc1_req));
                            break;
                        case fc2:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc2_req));
                            break;
                        case fc3:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc3_req));
                            break;
                        case fc4:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc4_req));
                            break;
                        case fc5:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc5_req));
                            break;
                        case fc6:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc6_req));
                            break;
                        case fc15:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc15_req));
                            break;
                        case fc16:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_cmd_hanlder(req_json_obj, mbtcp_fc16_req));
                            break;
                        case set_timeout:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_set_response_timeout(tid, json_get_long(req_json_obj, "timeout")));
                            break;
                        case get_timeout:
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                mbtcp_get_response_timeout(tid));
                            break;
                        default: 
                            send_modbus_zmq_resp(zmq_pub, cmd, 
                                set_modbus_fail_resp_str(tid, "unsupport request"));
                            break;
                    }
                }
                // @handle modbus rtu requests
                else if (strcmp(mode, "rtu") == 0)
                {
                    LOG(enable_syslog, "rtu:%d", cmd);
                    // [TODO]
                    // send error response
                }
                // @unkonw mode
                else
                {
                    send_modbus_zmq_resp(zmq_pub, -1, 
                        set_modbus_fail_resp_str(tid, "unsupport mode"));
                }
            }
            else
            {
                send_modbus_zmq_resp(zmq_pub, -2, 
                    set_modbus_fail_resp_str(tid, "Fail to parse command string"));
            }
            
            // @cleanup cJson object (auto mode)
            cJSON_Delete(req_json_obj);
        }
        else
        {
            // @depress this debug message
            //ERR(enable_syslog, "Recv null message");
        }
    }
    
    // @resource clean up
    LOG(enable_syslog, "clean up");
    zctx_destroy(&zmq_context);
    save_config(config_fname, config_json); 
}