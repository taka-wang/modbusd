/*
 * Copyright © 2008-2014 Stéphane Raimbault <stephane.raimbault@gmail.com>
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Modified by Taka
 */

/*
uint8_t *tab_bits;              // 0x
uint8_t *tab_input_bits;        // 1x
uint16_t *tab_input_registers;  // 3x
uint16_t *tab_registers;        // 4x
*/


#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <modbus.h>

int main(int argc, char *argv[])
{
    printf("init server\n");
    char * ip = NULL;
    int port = 502;
    if (argc > 1) {
        port = atoi(argv[1]);
    }
    
    int s = -1;
    modbus_t *ctx;
    modbus_mapping_t *mb_mapping;

    ctx = modbus_new_tcp(ip, port);
    modbus_set_debug(ctx, TRUE);

    // allocate memory map
    //mb_mapping = modbus_mapping_new(10000, 10000, 10000, 10000); // max
     mb_mapping = modbus_mapping_new(500, 500, 500, 500); // max
    if (mb_mapping == NULL)
    {
        fprintf(stderr, "Failed to allocate the mapping: %s\n", modbus_strerror(errno));
        modbus_free(ctx);
        return -1;
    }

    // initalize bits
    const uint8_t UT_BITS_TAB[] = { 0xAC, 0xDB, 0x35 };
    const uint16_t UT_BITS_NB = 0x16;

    // 0x; little endian set
    modbus_set_bits_from_bytes(mb_mapping->tab_bits, 
                               3, 
                               UT_BITS_NB,
                               UT_BITS_TAB);

    // 1x; little endian set
    modbus_set_bits_from_bytes(mb_mapping->tab_input_bits, 
                               0, 
                               UT_BITS_NB,
                               UT_BITS_TAB);

    const uint16_t UT_INPUT_REGISTERS_NB = 0x9;
    const uint16_t UT_INPUT_REGISTERS_TAB[] = { 0x000A, 0x000B, 0x000C, 0x000D, 0x000E, 0x000F, 0x0001, 0x0002, 0x0003 };
    
    // Initialize INPUT REGISTERS: 3x
    for (int i=0; i < UT_INPUT_REGISTERS_NB; i++) {
        mb_mapping->tab_input_registers[i] = UT_INPUT_REGISTERS_TAB[i];;
    }

    // Initialize  REGISTERS: 4x
    for (int i=0; i < UT_INPUT_REGISTERS_NB; i++) {
        mb_mapping->tab_registers[i] = UT_INPUT_REGISTERS_TAB[i];;
    }


    printf("start listening at: port:%d\n", port);

    s = modbus_tcp_listen(ctx, 1); // only one connection allow
    modbus_tcp_accept(ctx, &s);

    
    for (;;) {
        uint8_t query[MODBUS_TCP_MAX_ADU_LENGTH];
        int rc;

        rc = modbus_receive(ctx, query);
        if (rc > 0) {
            /* rc is the query size */
            modbus_reply(ctx, query, rc, mb_mapping);
        } else if (rc == -1) {
            /* Connection closed by the client or error */
            break;
        }
    }

    printf("Quit the loop: %s\n", modbus_strerror(errno));

    if (s != -1) {
        close(s);
    }
    modbus_mapping_free(mb_mapping);
    modbus_close(ctx);
    modbus_free(ctx);

    return 0;
}
