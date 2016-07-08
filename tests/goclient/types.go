package main

import (
	"fmt"
	"strings"
)

// Endian defines byte endianness
type Endian int

// RegValueType defines how to inteprete registers
type RegValueType int

// ScaleRange defines scale range
type ScaleRange struct {
	DomainLow  int `json:"a"`
	DomainHigh int `json:"b"`
	RangeLow   int `json:"c"`
	RangeHigh  int `json:"d"`
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
	RegisterArray RegValueType = iota
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

// ======================= psmb to modbusd structures - downstream =======================

// DMbtcpRes downstream modbus tcp read/write response
type DMbtcpRes struct {
	// Tid unique transaction id in string format
	Tid    string `json:"tid"`
	Status string `json:"status"`
	// Data for read function code only.
	Data []uint16 `json:"data,omitempty"`
}

// DMbtcpReadReq downstream modbus tcp read request
type DMbtcpReadReq struct {
	// Tid unique transaction id in string format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd   int    `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	// Len the length of register or bit
	Len uint16 `json:"len"`
}

// DMbtcpWriteReq downstream modbus tcp write single bit/register request
type DMbtcpWriteReq struct {
	// Tid unique transaction id in string format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd   int    `json:"cmd"`
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Addr  uint16 `json:"addr"`
	// Len omit for fc5, fc6
	Len uint16 `json:"len,omitempty"`
	// Data should be []uint16, uint16 (FC5, FC6)
	Data interface{} `json:"data"`
}

// DMbtcpTimeout downstream modbus tcp set/get timeout request/response
type DMbtcpTimeout struct {
	// Tid unique transaction id in string format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd int `json:"cmd"`
	// Status for response only.
	Status string `json:"status,omitempty"`
	// Timeout set timeout request and get timeout response only.
	Timeout int64 `json:"timeout,omitempty"`
}

// ======================= services to psmb structures - upstream =======================

// MbtcpOnceReadReq read coil/register request (1.1).
// Scale range field example:
// Range: &ScaleRange{1,2,3,4},
type MbtcpOnceReadReq struct {
	Tid   int64        `json:"tid"`
	From  string       `json:"from,omitempty"`
	FC    int          `json:"fc"`
	IP    string       `json:"ip"`
	Port  string       `json:"port,omitempty"`
	Slave uint8        `json:"slave"`
	Addr  uint16       `json:"addr"`
	Len   uint16       `json:"len,omitempty"`
	Type  RegValueType `json:"type,omitempty"`
	Order Endian       `json:"order,omitempty"`
	Range *ScaleRange  `json:"range,omitempty"` // point to struct can be omitted in json encode
}

// MbtcpOnceReadRes read coil/register response (1.1).
// `Data interface` supports:
// []uint16, []int16, []uint32, []int32, []float32, string
type MbtcpOnceReadRes struct {
	Tid    int64        `json:"tid"`
	Status string       `json:"status"`
	Type   RegValueType `json:"type,omitempty"`
	// Bytes FC3, FC4 and Type 2~8 only
	Bytes JSONableByteSlice `json:"bytes,omitempty"`
	Data  interface{}       `json:"data,omitempty"` // universal data container
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
