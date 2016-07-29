#
# golang client to modbusd
#
FROM takawang/psmb:x86
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go

## Load app files
ADD . /go

## Default command
CMD ["go", "test", "-v"]

