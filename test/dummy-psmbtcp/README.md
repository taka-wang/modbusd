# Dummy proactive service for modbus tcp

Dummy modbusd tester in golang.

## Motivation

I implement this service to test the communication between proactive service and modbusd service.

## From source code

```bash
sudo apt-get install pkg-config
curl -O https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz
tar -xvf go1.6.2.linux-amd64.tar.gz
sudo mv go /usr/local
nano ~/.profile
export PATH=$PATH:/usr/local/go/bin
go get github.com/takawang/zmq3
go test -v
```
