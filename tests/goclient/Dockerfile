#
# golang client to modbusd
#
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

## Load app files
ADD . /go
RUN go build

## Default command
CMD ["go", "test", "-v"]

