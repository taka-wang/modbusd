# modbus slave simulator in C

## Build
```bash
gcc server1.c -o server `pkg-config --libs --cflags libmodbus`
```