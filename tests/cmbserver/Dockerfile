#
# modbus slave simulator in C
#

FROM alpine:latest
MAINTAINER Taka Wang <taka@cmwang.net>

ENV ZMQ_VERSION 3.2.5
ENV CZMQ_VERSION 3.0.2
ENV MB_VERSION 3.1.4

RUN apk update \
    && apk add \
           git autoconf cmake build-base tar libtool zlib musl-dev openssl-dev zlib-dev curl \
    
    && echo " ... adding libmodbus" \
         && curl -L http://libmodbus.org/releases/libmodbus-${MB_VERSION}.tar.gz -o /tmp/libmodbus.tar.gz \
         && cd /tmp/ \
         && tar -xf /tmp/libmodbus.tar.gz \
         && cd /tmp/libmodbus*/ \
         && ./configure --prefix=/usr \
                        --sysconfdir=/etc \
                        --mandir=/usr/share/man \
                        --infodir=/usr/share/info \
         && make && make install \
    
    && echo " ... build mbserver" \
        && cd /tmp/ \
        && git clone https://github.com/taka-wang/modbusd.git \
        && cd /tmp/modbusd/tests/cmbserver/ \
        && gcc server.c -o server -Wall -std=c99 `pkg-config --libs --cflags libmodbus` \
        && chmod +x server \
        && cp server /usr/bin/ \
    
    && echo " ... clean up" \
        && rm -rf /tmp/* \
        && rm /usr/lib/*.a && rm /usr/lib/*.la \
        && apk del \
            git autoconf cmake build-base tar libtool zlib musl-dev openssl-dev zlib-dev curl \
        && rm -rf /var/cache/apk/* 

RUN apk update \
    && apk add libgcc libstdc++

EXPOSE 502

## Default command
CMD /usr/bin/server