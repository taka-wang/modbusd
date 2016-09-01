/**
 * @file mb.h
 * @author Taka Wang
 * @brief modbus daemon API(Interface)
*/

#ifndef MB_H
#define MB_H

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <time.h>
#include <errno.h>
#include <stdbool.h>

#include <modbus.h>
#include "uthash.h"
#include "log.h"
#include "json.h"

/* ==================================================
 *  struct
================================================== */

/**
 * @brief `structure key` for modbus tcp hash table
 *
 * Hash key. Ip v4/v6 address and port composition.
 */
typedef struct 
{
    char ip[50];        /** IP v4/v6 address or hostname */
    char port[50];      /** service name/port number to connect to */
} mbtcp_key_s;

/**
 * @brief hashable mbtcp handle type
 *
 * Hashable tcp handle strucut for keeping connection.
 */
typedef struct 
{
    mbtcp_key_s key;    /** key */
    bool connected;     /** is connect to modbus slave? */
    modbus_t *ctx;      /** modbus context pointer */
    UT_hash_handle hh;  /** makes this structure hashable */
} mbtcp_handle_s;


/**
 * @brief Function pointer of modbus tcp function code
 *
 * Function pointer to `modbus tcp function code request` for `generic command handle`.
 *
 * @param fc function code
 * @param handle Mbtcp handle.
 * @param req cJSON request object.
 * @return Modbus response string in JSON format.
 */
typedef char * (*mbtcp_fc)(uint8_t fc, mbtcp_handle_s *handle, cJSON *req);


/* ==================================================
 *  api
================================================== */

/**
 * @brief Set modbusd success response string with data (i.e., read func)
 *
 * @param tid Transaction ID.
 * @param json_arr cJSON pointer to data array
 * @return Modbus ok response string in JSON format.
 */ 
char * set_modbus_success_resp_str_with_data (char *tid, cJSON *json_arr);

/**
 * @brief Set modbusd success response string without data (i.e., write func)
 *
 * @param tid Transaction ID.
 * @return Modbus ok response string in JSON format.
 */ 
char * set_modbus_success_resp_str (char *tid);

/**
 * @brief Set modbusd fail response string.
 *
 * @param tid Transaction ID.
 * @param reason Fail reason string.
 * @return Modbus response string in JSON format.
 */
char * set_modbus_fail_resp_str (char *tid, const char *reason);

/**
 * @brief Set modbusd fail response string with error number.
 *
 * @param tid Transaction ID.
 * @param handle Mbtcp handle.
 * @param errnum Error number from modbus tcp handle.
 * @return Modbus error response string in JSON format.
 */ 
char * set_modbus_fail_resp_str_with_errno (char *tid, mbtcp_handle_s *handle, int errnum);

/* ==================================================
 *  modbus tcp (mbtcp)
================================================== */

/**
 * @brief Init mbtcp handle (to hashtable) and try to connect
 *
 * @param ptr_handle Pointer to mbtcp handle.
 * @param ip IP address string.
 * @param port Modbus TCP server port string.
 * @return Success or not.
 */
bool mbtcp_init_handle (mbtcp_handle_s **ptr_handle, char *ip, char *port);

/**
 * @brief Get mbtcp handle from hashtable
 *
 * @param ptr_handle Pointer to mbtcp handle.
 * @param ip IP address string.
 * @param port Modbus TCP server port string or service name.
 * @return Success or not.
 */
bool mbtcp_get_handle (mbtcp_handle_s **ptr_handle, char *ip, char *port);

/**
 * @brief List all handles in mbtcp hash table
 *
 * @return Void.
 */
void mbtcp_list_handles ();

/**
 * @brief Connect to mbtcp slave via mbtcp hashed handle
 *
 * @param handle Mbtcp handle.  
 * @param reason Fail reason string.
 * @return Success or not.
 */
bool mbtcp_do_connect (mbtcp_handle_s *handle, char **reason);

/**
 * @brief Get mbtcp handle's connection status
 *
 * @param handle Mbtcp handle.
 * @return Success or not. 
 */
bool mbtcp_get_connection_status (mbtcp_handle_s *handle);

/**
 * @brief Generic mbtcp command handler
 *
 * @param fc function code
 * @param req cJSON request object.
 * @param ptr_handler Function pointer to modbus tcp fc handler.
 * @return Modbus response string in JSON format.
 */
char * mbtcp_cmd_hanlder (uint8_t fc, cJSON *req, mbtcp_fc ptr_handler);


/**
 * @brief Set mbtcp response timeout in usec
 *
 * @param tid Transaction ID.
 * @param timeout Timeout in usec.
 * @return Modbus response string in JSON format.
 */
char * mbtcp_set_response_timeout (char *tid, long timeout);

/**
 * @brief Get mbtcp response timeout
 *
 * @param tid Transaction ID.
 * @return Modbus response string in JSON format.
 */
char * mbtcp_get_response_timeout (char *tid);

/**
 * @brief Help function. FC1, FC2 request handler
 *
 * @fc Function code 1 and 2 only.
 * @param handle Mbtcp handle.
 * @param req cJSON request object.
 * @return Modbus response string in JSON format.
 */
char * mbtcp_read_bit_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req);

/**
 * @brief Help function. FC3, FC4 request handler
 *
 * @fc Function code 3 and 4 only.
 * @param handle Mbtcp handle.
 * @param req cJSON request object.
 * @return Modbus response string in JSON format.
 */
char * mbtcp_read_reg_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req);

/**
 * @brief Help function. FC5, FC6 request handler
 *
 * @fc Function code 5 and 6 only.
 * @param handle Mbtcp handle.
 * @param req cJSON request object.
 * @return Modbus response string in JSON format.
 */
char * mbtcp_single_write_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req);

/**
 * @brief Help function. FC15, FC16 request handler
 *
 * @fc Function code 15 and 16 only.
 * @param handle Mbtcp handle.
 * @param req cJSON request object.
 * @return Modbus response string in JSON format.
 */
char * mbtcp_multi_write_req (uint8_t fc, mbtcp_handle_s *handle, cJSON *req);
#endif  // MB_H