package goclient

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

// MbSingleWriteReq Modbus tcp write request
type MbSingleWriteReq struct {
	IP    string `json:"ip"`
	Port  string `json:"port"`
	Slave uint8  `json:"slave"`
	Tid   int64  `json:"tid"`
	Cmd   string `json:"cmd"`
	Addr  uint16 `json:"addr"`
	Data  int32  `json:data`
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

type MbTimeoutReq struct {
	Tid  int64  `json:"tid"`
	Cmd  string `json:"cmd"`
	Data int64  `json:data`
}

// MbRes Modbus tcp generic response
type MbRes struct {
	Tid    int64  `json:"tid"`
	Status string `json:status`
}

// MbReadRes Modbus tcp read response
type MbReadRes struct {
	Tid    int64    `json:"tid"`
	Data   []uint16 `json:data` // uint16 for register
	Status string   `json:status`
}
