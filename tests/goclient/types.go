package main

import (
	"fmt"
	"strings"
)

// Endian defines byte endianness
type Endian int

// RegValueType defines how to inteprete registers
type RegValueType int

// MbtcpCmdType defines modbus tcp command type
type MbtcpCmdType string

const (
	fc1        MbtcpCmdType = "1"
	fc2        MbtcpCmdType = "2"
	fc3        MbtcpCmdType = "3"
	fc4        MbtcpCmdType = "4"
	fc5        MbtcpCmdType = "5"
	fc6        MbtcpCmdType = "6"
	fc15       MbtcpCmdType = "15"
	fc16       MbtcpCmdType = "16"
	setTimeout MbtcpCmdType = "50"
	getTimeout MbtcpCmdType = "51"
)

// mbtcpReadTask read/poll task request
type mbtcpReadTask struct {
	Name string
	Cmd  string
	Req  interface{}
}

// ScaleRange defines scale range
type ScaleRange struct {
	DomainLow  float64 `json:"a"`
	DomainHigh float64 `json:"b"`
	RangeLow   float64 `json:"c"`
	RangeHigh  float64 `json:"d"`
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
	_ Endian = iota // ignore first value by assigning to blank identifier
	// ABCD 32-bit words may be represented in big-endian format
	ABCD
	// DCBA 32-bit words may be represented in little-endian format
	DCBA
	// BADC 32-bit words may be represented in mid-big-endian format
	BADC
	// CDAB 32-bit words may be represented in mid-little-endian format
	CDAB
)

const (
	_ Endian = iota // ignore first value by assigning to blank identifier
	// AB 16-bit words may be represented in big-endian format
	AB
	// BA 16-bit words may be represented in little-endian format
	BA
)

const (
	_ Endian = iota // ignore first value by assigning to blank identifier
	// BigEndian 32-bit words may be represented in ABCD format
	BigEndian
	// LittleEndian 32-bit words may be represented in DCBA format
	LittleEndian
	// MidBigEndian 32-bit words may be represented in BADC format
	MidBigEndian
	// MidLittleEndian 32-bit words may be represented in CDAB format
	MidLittleEndian
)

const (
	_ RegValueType = iota // ignore first value by assigning to blank identifier
	// RegisterArray register array, ex: [12345, 23456, 5678]
	RegisterArray
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
	Tid string `json:"tid"`
	// Status response status string
	Status string `json:"status"`
	// Data for 'read function code' only.
	Data []uint16 `json:"data,omitempty"`
}

// DMbtcpReadReq downstream modbus tcp read request
type DMbtcpReadReq struct {
	// Tid unique transaction id in `string` format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd int `json:"cmd"`
	// IP ip address or hostname of the modbus tcp slave
	IP string `json:"ip"`
	// Port port number of the modbus tcp slave
	Port string `json:"port"`
	// Slave device id of the modbus tcp slave
	Slave uint8 `json:"slave"`
	// Addr start address for read
	Addr uint16 `json:"addr"`
	// Len the length of registers or bits
	Len uint16 `json:"len"`
}

// DMbtcpWriteReq downstream modbus tcp write single bit/register request
type DMbtcpWriteReq struct {
	// Tid unique transaction id in `string` format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd int `json:"cmd"`
	// IP ip address or hostname of the modbus tcp slave
	IP string `json:"ip"`
	// Port port number of the modbus tcp slave
	Port string `json:"port"`
	// Slave device id of the modbus tcp slave
	Slave uint8 `json:"slave"`
	// Addr start address for write
	Addr uint16 `json:"addr"`
	// Len omit for fc5, fc6
	Len uint16 `json:"len,omitempty"`
	// Data should be []uint16, uint16 (FC5, FC6)
	Data interface{} `json:"data"`
}

// DMbtcpTimeout downstream modbus tcp set/get timeout request/response
type DMbtcpTimeout struct {
	// Tid unique transaction id in `string` format
	Tid string `json:"tid"`
	// Cmd modbusd command type: https://github.com/taka-wang/modbusd#command-mapping-table
	Cmd int `json:"cmd"`
	// Status for response only.
	Status string `json:"status,omitempty"`
	// Timeout set timeout request and get timeout response only.
	Timeout int64 `json:"timeout,omitempty"`
}

// ======================= services to psmb structures - upstream =======================

// MbtcpReadReq read coil/register request (1.1).
// Scale range field example:
// Range: &ScaleRange{1,2,3,4},
type MbtcpReadReq struct {
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

// MbtcpReadRes read coil/register response (1.1).
// `Data interface` supports:
// []uint16, []int16, []uint32, []int32, []float32, string
type MbtcpReadRes struct {
	Tid    int64        `json:"tid"`
	Status string       `json:"status"`
	Type   RegValueType `json:"type,omitempty"`
	// Bytes FC3, FC4 and Type 2~8 only
	Bytes JSONableByteSlice `json:"bytes,omitempty"`
	Data  interface{}       `json:"data,omitempty"` // universal data container
}

// MbtcpWriteReq write coil/register request
type MbtcpWriteReq struct {
	Tid   int64       `json:"tid"`
	From  string      `json:"from,omitempty"`
	FC    int         `json:"fc"`
	IP    string      `json:"ip"`
	Port  string      `json:"port,omitempty"`
	Slave uint8       `json:"slave"`
	Addr  uint16      `json:"addr"`
	Len   uint16      `json:"len,omitempty"`
	Hex   bool        `json:"hex,omitempty"`
	Data  interface{} `json:"data"`
}

// MbtcpWriteRes == MbtcpSimpleRes

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

// MbtcpPollStatus polling coil/register request;
type MbtcpPollStatus struct {
	Tid      int64        `json:"tid"`
	From     string       `json:"from,omitempty"`
	Name     string       `json:"name"`
	Interval uint64       `json:"interval"`
	Enabled  bool         `json:"enabled"`
	FC       int          `json:"fc"`
	IP       string       `json:"ip"`
	Port     string       `json:"port,omitempty"`
	Slave    uint8        `json:"slave"`
	Addr     uint16       `json:"addr"`
	Status   string       `json:"status,omitempty"` // 2.3.2 response only
	Len      uint16       `json:"len,omitempty"`
	Type     RegValueType `json:"type,omitempty"`
	Order    Endian       `json:"order,omitempty"`
	Range    *ScaleRange  `json:"range,omitempty"` // point to struct can be omitted in json encode
}

// MbtcpPollRes == MbtcpSimpleRes

// MbtcpPollData read coil/register response (1.1).
// `Data interface` supports:
// []uint16, []int16, []uint32, []int32, []float32, string
type MbtcpPollData struct {
	TimeStamp int64        `json:"ts"`
	Name      string       `json:"name"`
	Status    string       `json:"status"`
	Type      RegValueType `json:"type,omitempty"`
	// Bytes FC3, FC4 and Type 2~8 only
	Bytes JSONableByteSlice `json:"bytes,omitempty"`
	Data  interface{}       `json:"data,omitempty"` // universal data container
}

// MbtcpPollOpReq generic modbus tcp poll operation request
type MbtcpPollOpReq struct {
	Tid      int64  `json:"tid"`
	From     string `json:"from,omitempty"`
	Name     string `json:"name,omitempty"`
	Interval uint64 `json:"interval,omitempty"`
	Enabled  bool   `json:"enabled,omitempty"`
}

// MbtcpPollsStatus requests status
type MbtcpPollsStatus struct {
	Tid    int64             `json:"tid"`
	Status string            `json:"status"`
	Polls  []MbtcpPollStatus `json:"polls"`
}
