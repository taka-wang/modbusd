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
    fc1         = 1,
    fc2         = 2,
    fc3         = 3,
    fc4         = 4,
    fc5         = 5,
    fc6         = 6,
    fc15        = 15,
    fc16        = 16,
    set_tcp_timeout = 50,
    get_tcp_timeout = 51
} cmd_t;

#endif  // MAIN_H