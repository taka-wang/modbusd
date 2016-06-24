/**
 * @file mb.c
 * @author taka-wang
 * @brief modbus common functions
*/

#include "mb.h"

extern int enable_syslog; // syslog flag

/* ==================================================
 *  public functions
================================================== */

char * set_modbus_success_resp_str(int tid)
{
    BEGIN(enable_syslog);

    cJSON *resp_root;
    resp_root = cJSON_CreateObject();
    cJSON_AddNumberToObject(resp_root, "tid", tid);
    cJSON_AddStringToObject(resp_root, "status", "ok");
    char * resp_json_str = cJSON_PrintUnformatted(resp_root);
    LOG(enable_syslog, "resp: %s", resp_json_str);
    // clean up
    cJSON_Delete(resp_root);
    return resp_json_str;     
}

char * set_modbus_success_resp_str_with_data(int tid, cJSON * json_arr)
{
    BEGIN(enable_syslog);

    cJSON *resp_root;
    resp_root = cJSON_CreateObject();
    cJSON_AddNumberToObject(resp_root, "tid", tid);
    cJSON_AddItemToObject(resp_root, "data", json_arr);
    cJSON_AddStringToObject(resp_root, "status", "ok");
    char * resp_json_str = cJSON_PrintUnformatted(resp_root);
    LOG(enable_syslog, "resp: %s", resp_json_str);
    // clean up
    cJSON_Delete(resp_root);
    return resp_json_str;
}

char * set_modbus_fail_resp_str(int tid, const char *reason)
{
    BEGIN(enable_syslog);
    
    cJSON *resp_root;
    resp_root = cJSON_CreateObject();
    cJSON_AddNumberToObject(resp_root, "tid", tid);
    cJSON_AddStringToObject(resp_root, "status", reason);
    char * resp_json_string = cJSON_PrintUnformatted(resp_root);
    LOG(enable_syslog, "resp: %s", resp_json_string);
    
    // clean up
    cJSON_Delete(resp_root);
    return resp_json_string;
}

char * set_modbus_fail_resp_str_with_errno(int tid, mbtcp_handle_s *handle, int errnum)
{
    BEGIN(enable_syslog);
    // [todo][enhance] reconnect proactively?
    // ... if the request interval is very large, 
    // we should try to reconnect automatically
    
    if (errnum == 104) // Connection reset by peer (i.e, tcp connection timeout)
    {
        handle->connected = false;
    }
    ERR(enable_syslog, "%s:%d", modbus_strerror(errnum), errnum);
    return set_modbus_fail_resp_str(tid, modbus_strerror(errnum));
}
