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
	Tid    int64    `json:"tid"`
	Data   []uint16 `json:data` // uint16 for register
	Status string   `json:status`
}

// MbMultipleWriteReq Modbus tcp write request
type MbMultipleWriteReq struct {
	IP    string   `json:"ip"`
	Port  string   `json:"port"`
	Slave uint8    `json:"slave"`
	Tid   int64    `json:"tid"`
	Cmd   string   `json:"cmd"`
	Addr  uint16   `json:"addr"`
	Len   uint16   `json:"len"`
	Data  []uint16 `json:data`
}

// MbSingleWriteReq Modbus tcp write request
type MbSingleWriteReq struct {
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

	s.Title("4X table test: FC3, FC6, FC16")

	s.Assert("`4X Table: 60000` Read/Write uint16 value test: FC6, FC3", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbSingleWriteReq{
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
			1,
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

	s.Assert("`4X Table: 30000` Read/Write int16 value test: FC6, FC3", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbSingleWriteReq{
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
			1,
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

	s.Assert("`4X Table: -20000` Read/Write int16 value test: FC6, FC3", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbWriteSingleReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000), // tid
			"fc6",
			10, // addr
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
			1,
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

	s.Assert("`4X Table` Multiple read/write test: FC16, FC3", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbMultipleWriteReq{
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
			if writeReq.Data[index] != r2.Data[index] {
				return false
			}
		}

		return true
	})

	s.Title("0X table test: FC1, FC5, FC15")

	s.Assert("`0X Table` Single read/write test:FC5, FC1", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbSingleWriteReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000), // tid
			"fc5",
			400, // addr
			1,   // should be optional
			1,
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
			"fc1",
			400,
			1,
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
		if r2.Data[0] != 1 {
			return false
		}
		return true
	})

	s.Assert("`0X Table` Multiple read/write test: FC15, FC1", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := MbMultipleWriteReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000), // tid
			"fc15",
			100, // addr
			10,
			[]uint16{0, 1, 1, 1, 0, 0, 0, 1, 0, 1},
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
			"fc1",
			100,
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
			if writeReq.Data[index] != r2.Data[index] {
				return false
			}
		}

		return true

	})

	s.Title("1X table test: FC2")

	s.Assert("`1X Table` read test: FC2", func(log sugar.Log) bool {
		readReq := MbReadReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000),
			"fc2",
			0,
			12,
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
		return true
	})

	s.Title("3X table read test: FC4")
	s.Assert("`3X Table` read test:FC4", func(log sugar.Log) bool {
		readReq := MbReadReq{
			hostName,
			portNum,
			1,
			rand.Int63n(10000000),
			"fc4",
			0,
			12,
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
		return true
	})
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
