package goclient

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/marksalpeter/sugar"
	"github.com/taka-wang/psmb"
	zmq "github.com/taka-wang/zmq3"
)

var hostName, portNum string

// generic tcp publisher
func publisher(cmd string) {

	sender, _ := zmq.NewSocket(zmq.PUB)
	defer sender.Close()
	sender.Connect("ipc:///tmp/to.modbus")

	for {
		time.Sleep(time.Duration(1) * time.Second)
		sender.Send("tcp", zmq.SNDMORE) // frame 1
		sender.Send(cmd, 0)             // convert to string; frame 2
		// send the exit loop
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
		// recv then exit loop
		return msg[0], msg[1]
	}
}

// ReadReqBuilder Read register/coil command builder
func ReadReqBuilder(cmd string, addr uint16, len uint16) psmb.DMbtcpReadReq {
	return psmb.DMbtcpReadReq{
		Tid:   uint64(rand.Int63n(10000000)),
		Cmd:   cmd,
		IP:    hostName,
		Port:  portNum,
		Slave: 1,
		Addr:  addr,
		Len:   len,
	}
}

// WriteReqBuilder Write single register/coil command builder
func WriteReqBuilder(cmd string, addr uint16, data uint16) psmb.DMbtcpSingleWriteReq {
	return psmb.DMbtcpSingleWriteReq{
		Tid:   uint64(rand.Int63n(10000000)),
		Cmd:   cmd,
		IP:    hostName,
		Port:  portNum,
		Slave: 1,
		Addr:  addr,
		Data:  data,
	}
}

// WriteMultiReqBuilder Write multiple register/coil command builder
func WriteMultiReqBuilder(cmd string, addr uint16, len uint16, data []uint16) psmb.DMbtcpMultipleWriteReq {
	return psmb.DMbtcpMultipleWriteReq{
		Tid:   uint64(rand.Int63n(10000000)),
		Cmd:   cmd,
		IP:    hostName,
		Port:  portNum,
		Slave: 1,
		Addr:  addr,
		Len:   len,
		Data:  data,
	}
}

//========= Test cases ==============================

// 4X
func TestHoldingRegisters(t *testing.T) {
	s := sugar.New(nil)
	portNum = "502"

	// generalize host reslove for docker/local env
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}

	s.Title("4X table test: FC3, FC6, FC16")

	s.Assert("`4X Table: 60000` Read/Write uint16 value test: FC6, FC3", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteReqBuilder("fc6", 10, 60000)
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 psmb.DMbtcpRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := ReadReqBuilder("fc3", 10, 1)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpReadRes
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
		writeReq := WriteReqBuilder("fc6", 10, 30000)
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 psmb.DMbtcpRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := ReadReqBuilder("fc3", 10, 1)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpReadRes
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

	/*
		s.Assert("`4X Table: -20000` Read/Write int16 value test: FC6, FC3", func(log sugar.Log) bool {
			// =============== write part ==============
			writeReq := WriteReqBuilder("fc6", 10, uint16(-20000))
			writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
			go publisher(string(writeReqStr))
			_, s1 := subscriber()
			log("req: %s", string(writeReqStr))
			log("res: %s", s1)

			// parse resonse
			var r1 psmb.DMbtcpRes
			if err := json.Unmarshal([]byte(s1), &r1); err != nil {
				fmt.Println("json err:", err)
			}
			// check reponse
			if r1.Status != "ok" {
				return false
			}

			// =============== read part ==============
			readReq := ReadReqBuilder("fc3", 10, 1)
			readReqStr, _ := json.Marshal(readReq) // marshal to json string
			go publisher(string(readReqStr))
			_, s2 := subscriber()
			log("req: %s", string(readReqStr))
			log("res: %s", s2)

			// parse resonse
			var r2 psmb.DMbtcpReadRes
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
	*/

	s.Assert("`4X Table` Multiple read/write test: FC16, FC3", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteMultiReqBuilder("fc16", 10, 10,
			[]uint16{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000})

		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 psmb.DMbtcpRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := ReadReqBuilder("fc3", 10, 10)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpReadRes
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

}

// 0x
func TestCoils(t *testing.T) {
	s := sugar.New(nil)
	portNum = "502"

	// generalize host reslove for docker/local env
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}

	s.Title("0X table test: FC1, FC5, FC15")

	s.Assert("`0X Table` Single read/write test:FC5, FC1", func(log sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteReqBuilder("fc5", 400, 1)
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 psmb.DMbtcpRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := ReadReqBuilder("fc1", 400, 1)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpReadRes
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
		writeReq := WriteMultiReqBuilder("fc15", 100, 10,
			[]uint16{0, 1, 1, 1, 0, 0, 0, 1, 0, 1})
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		log("req: %s", string(writeReqStr))
		log("res: %s", s1)

		// parse resonse
		var r1 psmb.DMbtcpRes
		if err := json.Unmarshal([]byte(s1), &r1); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r1.Status != "ok" {
			return false
		}

		// =============== read part ==============
		readReq := ReadReqBuilder("fc1", 100, 10)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpReadRes
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
}

// 1x
func TestDiscretesInput(t *testing.T) {
	s := sugar.New(nil)
	portNum = "502"

	// generalize host reslove for docker/local env
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}

	s.Title("1X table test: FC2")

	s.Assert("`1X Table` read test: FC2", func(log sugar.Log) bool {
		readReq := ReadReqBuilder("fc2", 0, 12)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpReadRes
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

// 3X
func TestInputRegisters(t *testing.T) {
	s := sugar.New(nil)
	portNum = "502"

	// generalize host reslove for docker/local env
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("docker run")
		hostName = host[0] //docker
	}

	s.Title("3X table read test: FC4")
	s.Assert("`3X Table` read test:FC4", func(log sugar.Log) bool {
		readReq := ReadReqBuilder("fc4", 0, 12)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		log("req: %s", string(readReqStr))
		log("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpReadRes
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

// timeout
func TestTimeout(t *testing.T) {
	s := sugar.New(nil)

	s.Assert("`Set timeout` test", func(log sugar.Log) bool {
		setReq := psmb.DMbtcpTimeoutReq{
			Tid:     uint64(rand.Int63n(10000000)),
			Cmd:     "timeout.set",
			Timeout: 5100000,
		}
		setReqStr, _ := json.Marshal(setReq) // marshal to json string
		go publisher(string(setReqStr))
		_, s2 := subscriber()
		log("req: %s", string(setReqStr))
		log("res: %s", s2)

		getReq := psmb.DMbtcpTimeoutReq{
			Tid: uint64(rand.Int63n(10000000)),
			Cmd: "timeout.get",
		}
		getReqStr, _ := json.Marshal(getReq) // marshal to json string
		go publisher(string(getReqStr))
		_, s3 := subscriber()
		log("req: %s", string(setReqStr))
		log("res: %s", s3)

		return true
	})
}
