#
# golang client to modbusd
#
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get github.com/marksalpeter/sugar


## Load app files
ADD . /go
RUN go get

## Default command
CMD ["go", "test", "-v"]

