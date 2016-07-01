package main

import (
	"fmt"
	"strings"
)

// Endian defines byte endianness
type Endian int

// RegDataType defines how to inteprete registers
type RegDataType int

// ScaleRange defines scale range
type ScaleRange struct {
	OriginLow  int `json:"a"`
	OriginHigh int `json:"b"`
	TargetLow  int `json:"c"`
	TargetHigh int `json:"d"`
}

// JSONableByteSlice jsonable uint8 array
type JSONableByteSlice []byte

// MarshalJSON implements the Marshaler interface on JSONableByteSlice (i.e., uint8/byte array).
// Ref: http://stackoverflow.com/questions/14177862/how-to-jsonize-a-uint8-slice-in-go
func (u JSONableByteSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

const (
	// ABCD 32-bit words may be represented in big-endian format
	ABCD Endian = iota
	// DCBA 32-bit words may be represented in little-endian format
	DCBA
	// BADC 32-bit words may be represented in mid-big-endian format
	BADC
	// CDAB 32-bit words may be represented in mid-little-endian format
	CDAB
)
const (
	// AB 16-bit words may be represented in big-endian format
	AB Endian = iota
	// BA 16-bit words may be represented in little-endian format
	BA
)
const (
	// BigEndian 32-bit words may be represented in ABCD format
	BigEndian Endian = iota
	// LittleEndian 32-bit words may be represented in DCBA format
	LittleEndian
	// MidBigEndian 32-bit words may be represented in BADC format
	MidBigEndian
	// MidLittleEndian 32-bit words may be represented in CDAB format
	MidLittleEndian
)

const (
	// RegisterArray register array, ex: [12345, 23456, 5678]
	RegisterArray RegDataType = iota
	// HexString hexadecimal string, ex: "112C004F12345678"
	HexString
	// Scale linearly scale
	Scale
	// UInt16 uint16 array
	UInt16
	// Int16 int16 array
	Int16
	// UInt32 uint32 array
	UInt32
	// Int32 int32 array
	Int32
	// Float32 float32 array
	Float32
)

// ======================= psbm to modbusd structures - downstream =======================

// DMbtcpRes modbus tcp function code generic response
type DMbtcpRes struct {
	Tid    string   `json:"tid"`
	Status string   `json:"status"`
	Data   []uint16 `json:"data,omitempty"`
}

// DMbtcpReadReq modbus tcp read request
type DMbtcpReadReq struct {
	Tid   string `json:"tid"`
	Cmd   int    `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Len   uint16 `json:"len"`
}

// DMbtcpSingleWriteReq modbus tcp write single bit/register request
type DMbtcpSingleWriteReq struct {
	Tid   string `json:"tid"`
	Cmd   int    `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	Data  uint16 `json:"data"`
}

// DMbtcpMultipleWriteReq modbus tcp write multiple bits/registers request
type DMbtcpMultipleWriteReq struct {
	Tid   string   `json:"tid"`
	Cmd   int      `json:"cmd"`
	IP    string   `json:"ip"`
	Port  string   `json:"port"`
	Slave uint8    `json:"slave"`
	Addr  uint16   `json:"addr"`
	Len   uint16   `json:"len"`
	Data  []uint16 `json:"data"`
}

// DMbtcpTimeout modbus tcp set/get timeout request/response
type DMbtcpTimeout struct {
	Tid     string `json:"tid"`
	Cmd     int    `json:"cmd"`
	Status  string `json:"status,omitempty"`
	Timeout int64  `json:"timeout,omitempty"`
}

// ======================= services to psbm structures - upstream =======================

// MbtcpOnceReadReq read coil/register request (1.1).
// Scale range field example:
// Range: &ScaleRange{1,2,3,4},
type MbtcpOnceReadReq struct {
	Tid   int64       `json:"tid"`
	From  string      `json:"from,omitempty"`
	FC    int         `json:"fc"`
	IP    string      `json:"ip"`
	Port  string      `json:"port,omitempty"`
	Slave uint8       `json:"slave"`
	Addr  uint16      `json:"addr"`
	Len   uint16      `json:"len,omitempty"`
	Type  RegDataType `json:"type,omitempty"`
	Order Endian      `json:"order,omitempty"`
	Range *ScaleRange `json:"range,omitempty"`
}

// MbtcpOnceReadRes read coil/register response (1.1).
// `Data interface` supports:
// []uint16, []int16, []uint32, []int32, []float32, string
type MbtcpOnceReadRes struct {
	Tid    int64       `json:"tid"`
	Status string      `json:"status"`
	Type   RegDataType `json:"type,omitempty"`
	// Bytes FC3, FC4 and Type 2~8 only
	Bytes JSONableByteSlice `json:"bytes,omitempty"`
	Data  interface{}       `json:"data,omitempty"`
}

// MbtcpTimeoutReq set/get TCP connection timeout request (1.3, 1.4)
type MbtcpTimeoutReq struct {
	Tid  int64  `json:"tid"`
	From string `json:"from,omitempty"`
	Data int64  `json:"timeout,omitempty"`
}

// MbtcpTimeoutRes set/get TCP connection timeout response (1.3, 1.4)
type MbtcpTimeoutRes struct {
	Tid    int64  `json:"tid"`
	Status string `json:"status"`
	Data   int64  `json:"timeout,omitempty"`
}

// MbtcpSimpleReq generic modbus tcp response
type MbtcpSimpleReq struct {
	Tid  int64  `json:"tid"`
	From string `json:"from,omitempty"`
}

// MbtcpSimpleRes generic modbus tcp response
type MbtcpSimpleRes struct {
	Tid    int64  `json:"tid"`
	Status string `json:"status"`
}
