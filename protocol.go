package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// zabbixHeader is Zabbix header
var zabbixHeader = []byte("ZBXD\x01")

// ////////////////////////////////////////////////////////////////////////////////// //

// encodeRequest encodes request
func encodeRequest(r *jmxRequest) []byte {
	payload, _ := json.Marshal(r)
	return encodePayload(payload)
}

// encodePayload encodes payload
func encodePayload(payload []byte) []byte {
	size := uint64(len(payload))

	var buf bytes.Buffer

	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, size)

	buf.Write(zabbixHeader)
	buf.Write(sizeBuf)
	buf.Write(payload)

	return buf.Bytes()
}

// decodeMeta decodes response meta
func decodeMeta(data []byte) (int, error) {
	if len(data) < 5 || !bytes.Equal(data[:5], zabbixHeader) {
		return -1, errors.New("Wrong header format")
	}

	return int(binary.LittleEndian.Uint64(data[5:])), nil
}

// decodeResponse decodes response
func decodeResponse(data []byte) (*jmxResponse, error) {
	resp := &jmxResponse{}
	err := json.Unmarshal(data, resp)

	if err != nil {
		return nil, errors.New("Can't unmarshal response data: " + err.Error())
	}

	if resp.Status != "success" {
		return nil, errors.New(resp.Error)
	}

	return resp, nil
}
