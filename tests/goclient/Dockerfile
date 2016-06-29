#
# golang client to modbusd
#
FROM takawang/ubuntu-gozmq
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get github.com/marksalpeter/sugar && go get github.com/taka-wang/psmb

## Load app files
ADD . /go

## Default command
CMD ["go", "test", "-v"]

