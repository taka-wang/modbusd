#
# golang client to modbusd
#
FROM takawang/gozmq:pull
MAINTAINER Taka Wang <taka@cmwang.net>

WORKDIR /go
RUN go get -u github.com/taka-wang/psmb

## Load app files
ADD . /go/src/github.com/taka-wang/dummy-psmbtcp
RUN cd /go/src/github.com/taka-wang/dummy-psmbtcp && go get && go test -v

## Default command
#CMD ["go", "test", "-v"]

ADD . /go/src/github.com/taka-wang/psmb
RUN go get github.com/takawang/sugar \
    && go get github.com/taka-wang/gocron \
    && go get github.com/takawang/logrus \ 
    && cd /go/src/github.com/taka-wang/psmb \
    && go test -v