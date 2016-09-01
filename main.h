/**
 * @file main.h
 * @author Taka Wang
* @brief main header with version string and command mapping table
*/

#ifndef MAIN_H
#define MAIN_H

/* ==================================================
 *  marco
================================================== */

// define version string from cmake
#ifndef VERSION
    #define VERSION @MODBUSD_VERSION@
#endif

// command mapping table
typedef enum
{
    FC1             = 1,
    FC2             = 2,
    FC3             = 3,
    FC4             = 4,
    FC5             = 5,
    FC6             = 6,
    FC15            = 15,
    FC16            = 16,
    SET_TCP_TIMEOUT = 50,
    GET_TCP_TIMEOUT = 51
} cmd_s;

/* ==================================================
 *  configuration keys
================================================== */

const char * KEY_LOGGER                 = "syslog";
const char * KEY_IPC_TYPE               = "zmq";
const char * KEY_IPC_PUB                = "pub";
const char * KEY_IPC_SUB                = "sub";
const char * KEY_MB_TYPE                = "mbtcp";
const char * KEY_MBTCP_CONN_TIMEOUT     = "connect_timeout";

#endif  // MAIN_H