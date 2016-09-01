/**
 * @file mbtcp.c
 * @author Taka Wang
 * @brief modbus tcp functions
*/

#include "mb.h"

/* ==================================================
 *  global variable
================================================== */

extern int enable_syslog;                       // syslog flag
long tcp_conn_timeout_usec = 200000;            // tcp connection timeout in usec
static mbtcp_handle_s *mbtcp_htable = NULL;     // hashtable header

/* ==================================================
 *  static functions
================================================== */

/**
 * @brief Combo function: get or init mbtcp handle.
 *
 * @param ptr_handle Pointer to mbtcp handle.
 * @param req cJSON request object.
 * @return Success or not.
 */
static bool lazy_init_mbtcp_handle (mbtcp_handle_s **ptr_handle, cJSON *req)
{
    BEGIN (enable_syslog);
    
    char *ip   = json_get_char (req, "ip");
    char *port = json_get_char (req, "port");
    
    // check before add to hash table, uthash requirement
    if (mbtcp_get_handle (ptr_handle, ip, port))
    {
       return true;
    }
    else
    {
        if (mbtcp_init_handle (ptr_handle, ip, port))
        {
            return true;
        }
        else
        {
            // Fatal! Unable to allocate mbtcp context,
            // maybe system resourse issue!
            return false;
        }
    }
}

/**
 * @brief Combo function: check mbtcp connection status, 
 *        if not connected, try to connect to slave.
 *
 * @param handle Mbtcp handle.
 * @param reason Pointer to fail reason string.
 * @return Success or not.
 */
static bool lazy_mbtcp_connect (mbtcp_handle_s *handle, char **reason)
{
    BEGIN (enable_syslog);
        
    if (mbtcp_get_connection_status (handle))
    {
        return true;
    }
    else
    {
        if (mbtcp_do_connect (handle, reason))
        {
             return true;
        }
        else
        {
            // could get fail reason from '*reason'
            return false;
        }
    }   
}

/* ==================================================
 *  public functions
================================================== */

bool mbtcp_init_handle (mbtcp_handle_s **ptr_handle, char *ip, char *port)
{
    BEGIN (enable_syslog);

    // create a mbtcp context
    modbus_t *ctx = modbus_new_tcp_pi (ip, port);
    
    if (ctx == NULL)
    {
        ERR (enable_syslog, "Unable to allocate mbtcp context");
        return false; // fail to allocate context
    }
        
    // set tcp connection timeout
    modbus_set_response_timeout (ctx, 0, tcp_conn_timeout_usec);
    LOG (enable_syslog, "set response timeout: %ld", tcp_conn_timeout_usec);
    
    // @add context to mbtcp hashtable
    mbtcp_handle_s *handle = (mbtcp_handle_s*) malloc (sizeof (mbtcp_handle_s));
    if (handle != NULL)
    {
        // let alignment bytes being set to zero-value!!
        memset (handle, 0, sizeof (mbtcp_handle_s));
        handle->connected = false;
        strcpy (handle->key.ip, ip);
        strcpy (handle->key.port, port);
        handle->ctx = ctx;

        HASH_ADD(hh, mbtcp_htable, key, sizeof (mbtcp_key_s), handle);
        LOG (enable_syslog, "Add %s: %s to mbtcp hashtable", handle->key.ip, mbtcp_htable->key.port);

        // call by reference to `mbtcp handle address`
        *ptr_handle = handle;

        // @connect to server
        char * reason = NULL;
        mbtcp_do_connect (handle, &reason);
        return true;
    }
    else
    {
        ERR (enable_syslog, "Unable to allocate mbtcp handle");
        return false;
    }
}

bool mbtcp_get_handle (mbtcp_handle_s **ptr_handle, char *ip, char *port)
{
    BEGIN(enable_syslog);
    
    mbtcp_handle_s query, *hash_ctx;
    memset (&query, 0, sizeof (mbtcp_handle_s));
    strcpy (query.key.ip, ip);
    strcpy (query.key.port, port);
    // get handle from hash table
    HASH_FIND (hh, mbtcp_htable, &query.key, sizeof (mbtcp_key_s), hash_ctx);
    
    if (hash_ctx != NULL)
    {
        LOG (enable_syslog, "tcp server %s:%s found", hash_ctx->key.ip, hash_ctx->key.port);
        // call by reference to `mbtcp handle address`
        *ptr_handle = hash_ctx; 
        return true;
    }
    else
    {
        ERR (enable_syslog, "tcp server %s:%s not found", query.key.ip, query.key.port);
        *ptr_handle = NULL; 
        return false; // not found
    }
}

void mbtcp_list_handles () 
{
    BEGIN (enable_syslog);

    mbtcp_handle_s * handle;
    for (handle = mbtcp_htable; handle != NULL; handle = handle->hh.next)
    {
        LOG (enable_syslog, "ip:%s, port:%s", handle->key.ip, handle->key.port);
    }

    END (enable_syslog);
}

bool mbtcp_do_connect (mbtcp_handle_s *handle, char ** reason)
{
    BEGIN (enable_syslog);
    
    if (handle != NULL)
    {
        if (modbus_connect (handle->ctx) == -1) 
        {
            ERR (enable_syslog, "Connection failed: %s", modbus_strerror (errno));
            handle->connected = false;
            *reason = (char *) modbus_strerror (errno);
            return false;
        }
        else
        {
            LOG (enable_syslog, "%s:%s connected", handle->key.ip, handle->key.port);
            handle->connected = true;
            return true;
        }
    }
    else
    {
        ERR (enable_syslog, "NULL handle");
        *reason = "NULL handle";
        return false;
    }
}

bool mbtcp_get_connection_status (mbtcp_handle_s *handle)
{
    BEGIN (enable_syslog);
    
    if (handle != NULL)
    {
        LOG (enable_syslog, "%s:%s connected: %s",  handle->key.ip, 
                                                    handle->key.port, 
                                                    handle->connected ? "true" : "false");
        return handle->connected;
    }
    else
    {
        ERR (enable_syslog, "NULL handle");
        return false;
    }
}

char * mbtcp_cmd_hanlder (uint8_t fc, cJSON *req, mbtcp_fc ptr_handler)
{
    BEGIN (enable_syslog);

    char * tid = json_get_char (req, "tid");
    mbtcp_handle_s *handle = NULL;
    
    if (lazy_init_mbtcp_handle (&handle, req))
    {
        char * reason = NULL;
        if (lazy_mbtcp_connect (handle, &reason))
        {
            // set slave id
            int slave = json_get_int (req, "slave");
            LOG (enable_syslog, "slave id: %d", slave);
            modbus_set_slave (handle->ctx, slave);
            return ptr_handler (fc, handle, req);
        }
        else
        {
            // [enhance]: get reason from modbus response
            return set_modbus_fail_resp_str (tid, reason);
        }
    }
    else
    {
        return set_modbus_fail_resp_str (tid, "Fail to init modbus tcp handle");
    }
}

char * mbtcp_set_response_timeout (char *tid, long timeout)
{
    BEGIN (enable_syslog);

    // set timeout
    tcp_conn_timeout_usec = timeout;
    
    cJSON *resp_root;
    resp_root = cJSON_CreateObject ();
    cJSON_AddStringToObject (resp_root, "tid", tid);
    cJSON_AddStringToObject (resp_root, "status", "ok");
    char * resp_json_string = cJSON_PrintUnformatted (resp_root);
    LOG (enable_syslog, "resp: %s", resp_json_string);
    
    // clean up
    cJSON_Delete (resp_root);
    END (enable_syslog);
    return resp_json_string;
}

char * mbtcp_get_response_timeout (char *tid)
{
    BEGIN (enable_syslog);

    cJSON *resp_root;
    resp_root = cJSON_CreateObject();
    cJSON_AddStringToObject (resp_root, "tid", tid);
    cJSON_AddNumberToObject (resp_root, "timeout", tcp_conn_timeout_usec);
    cJSON_AddStringToObject (resp_root, "status", "ok");
    char * resp_json_string = cJSON_PrintUnformatted (resp_root);
    LOG (enable_syslog, "resp: %s", resp_json_string);
    
    // clean up
    cJSON_Delete (resp_root);
    END (enable_syslog);
    return resp_json_string;
}

char * mbtcp_read_bit_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req)
{
    BEGIN (enable_syslog);

    int addr = json_get_int (req, "addr");
    int len  = json_get_int (req, "len");
    char *tid = json_get_char (req, "tid");
    uint8_t bits[len];
    int ret = 0;

    if (len > MODBUS_MAX_READ_BITS) // 2000
    {
        return set_modbus_fail_resp_str (tid, "Too many bits requested");
    }
    else
    {
        switch (fc)
        {
            case 1:
                ret = modbus_read_bits (handle->ctx, addr, len, bits);
                break;
            case 2:
                ret = modbus_read_input_bits (handle->ctx, addr, len, bits);
                break;
            default:
                return set_modbus_fail_resp_str (tid, "Invalid function code");
        }

        if (ret < 0) 
        {
            return set_modbus_fail_resp_str_with_errno (tid, handle, errno);
        } 
        else 
        {
            LOG (enable_syslog, "fc:%d, desired length: %d, read length:%d", fc, len, ret);
            
            /* debug only
            for (int ii = 0; ii < ret; ii++) 
            {
                LOG (enable_syslog, "[%d]=%d", ii, bits[ii]);
            }
            */

            // uint8_t array
            return set_modbus_success_resp_str_with_data (tid, cJSON_CreateUInt8Array (bits, len));
        }
    }    
}

char * mbtcp_read_reg_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req)
{
    BEGIN (enable_syslog);

    int addr = json_get_int (req, "addr");
    int len  = json_get_int (req, "len");
    char * tid = json_get_char (req, "tid");
    uint16_t regs[len];
    int ret = 0;

    if (len > MODBUS_MAX_READ_REGISTERS) // 125
    {
        return set_modbus_fail_resp_str (tid, "Too many registers requested");
    }
    else
    {
        switch (fc)
        {
            case 3:
                ret = modbus_read_registers (handle->ctx, addr, len, regs);
                break;
            case 4:
                ret = modbus_read_input_registers (handle->ctx, addr, len, regs);
                break;
            default:
                return set_modbus_fail_resp_str (tid, "Invalid function code");
        }
        
        if (ret < 0) 
        {
            return set_modbus_fail_resp_str_with_errno (tid, handle, errno);
        } 
        else 
        {
            LOG (enable_syslog, "fc:%d, desired length: %d, read length:%d", fc, len, ret);
            
            /* debug only
            for (int ii = 0; ii < ret; ii++) 
            {
                LOG (enable_syslog, "[%d]=%d", ii, regs[ii]);
            }
            */
            
            // uint16_t array
            return set_modbus_success_resp_str_with_data (tid, cJSON_CreateUInt16Array (regs, len));
        }
    }
}

char * mbtcp_single_write_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req)
{
    BEGIN (enable_syslog);

    int addr = json_get_int (req, "addr");
    char * tid = json_get_char (req, "tid");
    int data = json_get_int (req, "data");
    int ret  = 0;

    switch (fc)
    {
        case 5:
            ret = modbus_write_bit (handle->ctx, addr, data);
            break;
        case 6:
            ret = modbus_write_register (handle->ctx, addr, data);
            break;
        default:
            return set_modbus_fail_resp_str (tid, "Invalid function code");
    }

    if (ret < 0) 
    {
        return set_modbus_fail_resp_str_with_errno (tid, handle, errno);
    }
    else
    {
        return set_modbus_success_resp_str (tid);   
    }
}

char * mbtcp_multi_write_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req)
{
    BEGIN (enable_syslog);

    int addr = json_get_int (req, "addr");
    int len  = json_get_int (req, "len");
    char * tid = json_get_char (req, "tid");
    cJSON * data = NULL;
    int ret = 0;
    uint8_t bits[len];  // FC15, VLAs
    uint16_t regs[len]; // FC16, VLAs

    switch (fc)
    {
        case 15:
            data = cJSON_GetObjectItem (req, "data");
            for (int i = 0 ; i < cJSON_GetArraySize (data) ; i++)
            {
                uint8_t subitem = cJSON_GetArrayItem (data, i)->valueint;
                bits[i] = subitem;
                LOG(enable_syslog, "[%d]=%d", i, bits[i]);
            }
            ret = modbus_write_bits (handle->ctx, addr, len, bits);
            break;

        case 16:
            data = cJSON_GetObjectItem (req, "data");
            for (int i = 0 ; i < cJSON_GetArraySize (data) ; i++)
            {
                uint16_t subitem = cJSON_GetArrayItem (data, i)->valueint;
                regs[i] = subitem;
                LOG(enable_syslog, "[%d]=%d", i, regs[i]);
            }
            ret = modbus_write_registers (handle->ctx, addr, len, regs);
            break;

        default:
            return set_modbus_fail_resp_str (tid, "Invalid function code");
    }

    if (ret < 0) 
    {
        return set_modbus_fail_resp_str_with_errno (tid, handle, errno);
    } 
    else
    {
        return set_modbus_success_resp_str (tid);         
    }
}