#
# golang client to modbusd
#
FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb

## Load app files
ADD . /go
RUN export GOPATH=/go && go get && go test -v

## Default command
#CMD ["go", "test", "-v"]

