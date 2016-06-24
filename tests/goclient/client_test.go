package goclient

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/marksalpeter/sugar"
	zmq "github.com/taka-wang/zmq3"
)

// MbReadReq Modbus tcp read request
type MbReadReq struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave int    `json:"slave"`
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	Addr  int    `json:"addr"`
	Len   int    `json:"len"`
}

// MbReadRes Modbus tcp read response
type MbReadRes struct {
	Tid    int64    `json:"tid"`
	Data   []string `json:data`
	Status string   `json:status`
}

// MbWriteReq Modbus tcp write request
type MbWriteReq struct {
	IP    string   `json:"ip"`
	Port  string   `json:"port"`
	Slave int      `json:"slave"`
	Tid   int64    `json:"tid"`
	Cmd   string   `json:"cmd"`
	Addr  int      `json:"addr"`
	Len   int      `json:"len"`
	Data  []string `json:data`
}

// MbRes Modbus tcp generic response
type MbRes struct {
	Tid    int64  `json:"tid"`
	Status string `json:status`
}

type MbTimeoutReq struct {
	Tid  int64  `json:"tid"`
	Cmd  string `json:"cmd"`
	Data int64  `json:data`
}

func TestModbus(t *testing.T) {
	s := sugar.New(nil)
	s.Title("modbus test")

	s.Assert("`Function 1` should work", func(log sugar.Log) bool {
		log("Hello")
		go publisher(gen())
		a, b := subscriber()
		log("Get method:%s", a)
		log("Get json:%s", b)
		return true
	})
	s.Assert("`Function 2` should work", func(log sugar.Log) bool {
		log("World")
		go publisher(gen())
		a, b := subscriber()
		log("Get method:%s", a)
		log("Get json:%s", b)
		return true
	})
}

func gen() string {
	command := MbReadReq{
		"127.0.0.1",
		"1502",
		1,
		12,
		"fc1",
		10,
		10,
	}

	cmd, err := json.Marshal(command) // marshal to json string
	if err != nil {
		fmt.Println("json err:", err)
		return ""
	}
	return string(cmd)
}

func publisher(cmd string) {

	sender, _ := zmq.NewSocket(zmq.PUB)
	defer sender.Close()
	sender.Connect("ipc:///tmp/to.modbus")

	for {
		time.Sleep(time.Duration(1) * time.Second)
		sender.Send("tcp", zmq.SNDMORE) // frame 1
		sender.Send(cmd, 0)             // convert to string; frame 2
		break
	}
}

// generic subscribe
func subscriber() (string, string) {
	receiver, _ := zmq.NewSocket(zmq.SUB)
	defer receiver.Close()
	receiver.Connect("ipc:///tmp/from.modbus")
	filter := ""
	receiver.SetSubscribe(filter) // filter frame 1
	for {
		msg, _ := receiver.RecvMessage(0)
		fmt.Println(msg[0]) // frame 1: method
		fmt.Println(msg[1]) // frame 2: command
		return msg[0], msg[1]
	}
}
