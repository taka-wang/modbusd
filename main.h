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

#endif  // MAIN_H