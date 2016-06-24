# modbus slave simulator in C

## Build
```bash
gcc server1.c -o server -Wall -std=c99 `pkg-config --libs --cflags libmodbus`
```