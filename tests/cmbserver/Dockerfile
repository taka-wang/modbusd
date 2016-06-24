#
# modbus simulation slave
#
FROM takawang/ubuntu-modbus
MAINTAINER Taka Wang <taka@cmwang.net>

## Load app files
WORKDIR /mbserver
ADD . /mbserver
RUN gcc server.c -o server -Wall -std=c99 `pkg-config --libs --cflags libmodbus`

## Default exported port
EXPOSE 502

## Default command
CMD ["/mbserver/server"]

