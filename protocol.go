package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
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
func encodeRequest(r *jmxRequest) ([]byte, error) {
	payload, err := json.Marshal(r)

	if err != nil {
		return nil, errors.New("Can't marshal request data: " + err.Error())
	}

	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, uint64(len(payload)))

	var buf bytes.Buffer

	buf.Write(zabbixHeader)
	buf.Write(sizeBuf)
	buf.Write(payload)

	return buf.Bytes(), nil
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
