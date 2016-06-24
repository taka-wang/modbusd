package goclient

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/marksalpeter/sugar"
	zmq "github.com/taka-wang/zmq3"
)

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

// MbReadReq Modbus tcp read request
type MbReadReq struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	Addr  uint16 `json:"addr"`
	Len   uint16 `json:"len"`
}

// MbReadRes Modbus tcp read response
type MbReadRes struct {
	Tid    int64   `json:"tid"`
	Data   []int32 `json:data` // uint16 for register
	Status string  `json:status`
}

// MbWriteReq Modbus tcp write request
type MbWriteReq struct {
	IP    string   `json:"ip"`
	Port  string   `json:"port"`
	Slave uint8    `json:"slave"`
	Tid   int64    `json:"tid"`
	Cmd   string   `json:"cmd"`
	Addr  uint16   `json:"addr"`
	Len   uint16   `json:"len"`
	Data  []uint16 `json:data`
}

// MbWriteReq Modbus tcp write request
type MbWriteSingleReq struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	Addr  uint16 `json:"addr"`
	Len   uint16 `json:"len"`
	Data  int32  `json:data`
}

func TestModbus(t *testing.T) {
	var hostName string
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println(err)
		hostName = "127.0.0.1"
	} else {
		hostName = host[0] //docker
	}
	portNum := "502"

	s := sugar.New(nil)

	s.Title("4x table read/write test")

	s.Assert("`4X Table: 60000` Read/Write uint16 value test", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbWriteSingleReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000), // tid
			"fc6",
			10, // addr
			1,  // should be optional
			60000,
		}

		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 MbRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := MbReadReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000),
			"fc3",
			10,
			1, // should be optional
		}

		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 MbReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}
		if r2.Data[0] != 60000 {
			return false
		}
		return true
	})

	s.Assert("`4X Table: 30000` Read/Write int16 value test", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbWriteSingleReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000), //tid
			"fc6",
			10, // addr
			1,  // should be optional
			30000,
		}

		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 MbRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := MbReadReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000),
			"fc3",
			10,
			1, // should be optional
		}

		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 MbReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}
		if r2.Data[0] != 30000 {
			return false
		}
		return true
	})

	s.Assert("`4X Table: -20000` Read/Write int16 value test", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbWriteSingleReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000), // tid
			"fc6",
			10, // addr
			1,  // should be optional
			-20000,
		}

		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 MbRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := MbReadReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000),
			"fc3",
			10,
			1, // should be optional
		}

		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 MbReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}
		if r2.Data[0] != 0xB1E0 {
			return false
		}
		return true
	})

	s.Assert("`4X Table` Write multiple registers", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbWriteReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000), // tid
			"fc16",
			10, // addr
			10,
			[]uint16{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000},
		}

		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 MbRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := MbReadReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000),
			"fc3",
			10,
			10,
		}

		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 MbReadRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}

		var index uint16
		for index = 0; index < readReq.Len; index++ {
			if writeReq.Data[index] != uint16(r2.Data[index]) {
				return false
			}
		}

		return true
	})

	s.Title("0x table read/write test")

	// TODO: single read/write

	// TODO: multiple read/write

	s.Assert("`Function 1` should work", func(log sugar.Log) bool {
		log("Hello")
		go publisher(gen())
		a, b := subscriber()
		log("Get method:%s", a)
		log("Get json:%s", b)

		var s MbReadRes
		if err := json.Unmarshal([]byte(b), &s); err != nil {
			fmt.Println("json err:", err)
		}
		log("Get status %s", s.Status)
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
		"502",
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

// generic tcp publisher
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
		//fmt.Println(msg[0]) // frame 1: method
		//fmt.Println(msg[1]) // frame 2: command
		return msg[0], msg[1]
	}
}
