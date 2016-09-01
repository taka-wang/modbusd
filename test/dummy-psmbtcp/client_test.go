package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	psmb "github.com/taka-wang/psmb"
	"github.com/takawang/sugar"
	zmq "github.com/takawang/zmq3"
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

// ReadReqBuilder help to build read register/coil command
func ReadReqBuilder(cmd int, addr uint16, len uint16) psmb.DMbtcpReadReq {
	return psmb.DMbtcpReadReq{
		Tid:   strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		Cmd:   cmd,
		IP:    hostName,
		Port:  portNum,
		Slave: 1,
		Addr:  addr,
		Len:   len,
	}
}

// WriteReqBuilder help to build write single register/coil command
func WriteReqBuilder(cmd int, addr uint16, data uint16) psmb.DMbtcpWriteReq {
	return psmb.DMbtcpWriteReq{
		Tid:   strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		Cmd:   cmd,
		IP:    hostName,
		Port:  portNum,
		Slave: 1,
		Addr:  addr,
		Data:  data,
	}
}

// WriteMultiReqBuilder help to build Write multiple register/coil command
func WriteMultiReqBuilder(cmd int, addr uint16, len uint16, data []uint16) psmb.DMbtcpWriteReq {
	return psmb.DMbtcpWriteReq{
		Tid:   strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		Cmd:   cmd,
		IP:    hostName,
		Port:  portNum,
		Slave: 1,
		Addr:  addr,
		Len:   len,
		Data:  data,
	}
}

// init functions
func init() {
	portNum = "502"
	// generalize host reslove for docker/local env
	host, err := net.LookupHost("slave")
	if err != nil {
		fmt.Println("Local run")
		hostName = "127.0.0.1"
	} else {
		fmt.Println("Docker run")
		hostName = host[0] //docker
	}
}

//========= Test cases ==============================

// 4X
func TestHoldingRegisters(t *testing.T) {
	s := sugar.New(t)

	s.Title("4X table test: FC3, FC6, FC16")

	s.Assert("`4X Table: 60000` Read/Write uint16 value test: FC6, FC3", func(logf sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteReqBuilder(6, 10, 60000)
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		logf("req: %s", string(writeReqStr))
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		logf("res: %s", s1)

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
		readReq := ReadReqBuilder(3, 10, 1)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		logf("req: %s", string(readReqStr))
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		logf("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpRes
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

	s.Assert("`4X Table: 30000` Read/Write int16 value test: FC6, FC3", func(logf sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteReqBuilder(6, 10, 30000)
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		logf("req: %s", string(writeReqStr))
		logf("res: %s", s1)

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
		readReq := ReadReqBuilder(3, 10, 1)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		logf("req: %s", string(readReqStr))
		logf("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpRes
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
		s.Assert("`4X Table: -20000` Read/Write int16 value test: FC6, FC3", func(logf sugar.Log) bool {
			// =============== write part ==============
			writeReq := WriteReqBuilder(6, 10, uint16(-20000))
			writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
			go publisher(string(writeReqStr))
			_, s1 := subscriber()
			logf("req: %s", string(writeReqStr))
			logf("res: %s", s1)

			// parse resonse
			var r1 DMbtcpRes
			if err := json.Unmarshal([]byte(s1), &r1); err != nil {
				fmt.Println("json err:", err)
			}
			// check reponse
			if r1.Status != "ok" {
				return false
			}

			// =============== read part ==============
			readReq := ReadReqBuilder(3, 10, 1)
			readReqStr, _ := json.Marshal(readReq) // marshal to json string
			go publisher(string(readReqStr))
			_, s2 := subscriber()
			logf("req: %s", string(readReqStr))
			logf("res: %s", s2)

			// parse resonse
			var r2 DMbtcpRes
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

	s.Assert("`4X Table` Multiple read/write test: FC16, FC3", func(logf sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteMultiReqBuilder(16, 10, 10,
			[]uint16{1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000})

		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		logf("req: %s", string(writeReqStr))
		logf("res: %s", s1)

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
		readReq := ReadReqBuilder(3, 10, 10)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		logf("req: %s", string(readReqStr))
		logf("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}

		var index uint16
		data := writeReq.Data.([]uint16)
		for index = 0; index < readReq.Len; index++ {
			if data[index] != r2.Data[index] {
				return false
			}
		}

		return true
	})

}

// 0x
func TestCoils(t *testing.T) {
	s := sugar.New(t)

	s.Title("0X table test: FC1, FC5, FC15")

	s.Assert("`0X Table` Single read/write test:FC5, FC1", func(logf sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteReqBuilder(5, 400, 1)
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		logf("req: %s", string(writeReqStr))
		logf("res: %s", s1)

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
		readReq := ReadReqBuilder(1, 400, 1)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		logf("req: %s", string(readReqStr))
		logf("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpRes
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

	s.Assert("`0X Table` Multiple read/write test: FC15, FC1", func(logf sugar.Log) bool {
		// =============== write part ==============
		writeReq := WriteMultiReqBuilder(15, 100, 10,
			[]uint16{0, 1, 1, 1, 0, 0, 0, 1, 0, 1})
		writeReqStr, _ := json.Marshal(writeReq) // marshal to json string
		go publisher(string(writeReqStr))
		_, s1 := subscriber()
		logf("req: %s", string(writeReqStr))
		logf("res: %s", s1)

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
		readReq := ReadReqBuilder(1, 100, 10)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		logf("req: %s", string(readReqStr))
		logf("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpRes
		if err := json.Unmarshal([]byte(s2), &r2); err != nil {
			fmt.Println("json err:", err)
		}
		// check reponse
		if r2.Status != "ok" {
			return false
		}

		var index uint16
		data := writeReq.Data.([]uint16)
		for index = 0; index < readReq.Len; index++ {
			if data[index] != r2.Data[index] {
				return false
			}
		}

		return true

	})
}

// 1x
func TestDiscretesInput(t *testing.T) {
	s := sugar.New(t)

	s.Title("1X table test: FC2")

	s.Assert("`1X Table` read test: FC2", func(logf sugar.Log) bool {
		readReq := ReadReqBuilder(2, 0, 12)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		logf("req: %s", string(readReqStr))
		logf("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpRes
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
	s := sugar.New(t)

	s.Title("3X table read test: FC4")
	s.Assert("`3X Table` read test:FC4", func(logf sugar.Log) bool {
		readReq := ReadReqBuilder(4, 0, 12)
		readReqStr, _ := json.Marshal(readReq) // marshal to json string
		go publisher(string(readReqStr))
		_, s2 := subscriber()
		logf("req: %s", string(readReqStr))
		logf("res: %s", s2)

		// parse resonse
		var r2 psmb.DMbtcpRes
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
	s := sugar.New(t)

	s.Assert("`Set timeout` test", func(logf sugar.Log) bool {
		setReq := psmb.DMbtcpTimeout{
			Tid:     strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
			Cmd:     50,
			Timeout: 5100000,
		}
		setReqStr, _ := json.Marshal(setReq) // marshal to json string
		logf("req: %s", string(setReqStr))

		go publisher(string(setReqStr))
		_, s2 := subscriber()
		logf("res: %s", s2)

		return true
	})

	s.Assert("`Get timeout` test", func(logf sugar.Log) bool {

		getReq := psmb.DMbtcpTimeout{
			Tid: strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
			Cmd: 51,
		}
		getReqStr, _ := json.Marshal(getReq) // marshal to json string
		logf("req: %s", string(getReqStr))

		go publisher(string(getReqStr))
		_, s3 := subscriber()
		logf("res: %s", s3)

		return true
	})
}
