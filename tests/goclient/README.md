#goclient

Modbusd test cases in golang

## Docker

### From the scratch
```bash
# build docker image 
docker build -t takawang/modbus-goclient .

# build arm version image 
#docker build -t takawang/arm-modbus-goclient -f Dockerfile.arm .

# run the image (host_port:container_port)
docker run -p 502:502 -d --name slave takawang/modbus-cserver

# mount file system
docker run -v /tmp:/tmp -it takawang/modbus-goclient /bin/bash

# run go test
go test -v

# Print app output
docker logs <container id>
# Enter the container
docker exec -it <container id> /bin/bash
```

### Pull pre-built docker image
```bash
docker pull takawang/modbus-goclient

# arm version
#docker pull takawang/arm-modbus-goclient
```
