#
# golang client to modbusd
#
FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb

## Load app files
ADD . /go/src/github.com/taka-wang/dummy-psmbtcp
RUN cd /go/src/github.com/taka-wang/dummy-psmbtcp && go get -t ./... && go test -v

## Default command
#CMD ["go", "test", "-v"]
