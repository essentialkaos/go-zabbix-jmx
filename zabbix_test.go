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
	"testing"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type JMXSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&JMXSuite{})

var respData1 = `{
  "data": [
    {
      "value": "112.637"
    }
  ],
  "response": "success"
}`

var respData2 = `{
  "response": "error"
}`

var beansData = `{\"data\":[{\"{#JMXDOMAIN}\":\"kafka.server\",\"{#JMXTYPE}\":\"BrokerTopicMetrics\",\"{#JMXOBJ}\":\"kafka.server:type=BrokerTopicMetrics,name=TotalProduceRequestsPerSec\",\"{#JMXNAME}\":\"TotalProduceRequestsPerSec\"},{\"{#JMXDOMAIN}\":\"kafka.server\",\"{#JMXTYPE}\":\"BrokerTopicMetrics\",\"{#JMXOBJ}\":\"kafka.server:type=BrokerTopicMetrics,name=BytesOutPerSec\",\"{#JMXNAME}\":\"BytesOutPerSec\"}]}`

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *JMXSuite) TestClient(c *C) {
	client, err := NewClient("127.0.")

	c.Assert(client, IsNil)
	c.Assert(err, NotNil)

	client, err = NewClient("127.0.0.0:10051")

	c.Assert(client, NotNil)
	c.Assert(err, IsNil)
}

func (s *JMXSuite) TestEncoder(c *C) {
	r := &Request{
		Server:   "domain.com",
		Port:     9334,
		Username: "admin",
		Password: "admin",
		Keys:     []string{`jmx["kafka.server:type=ReplicaManager,name=PartitionCount",Value]`},
	}

	jr := convertRequest(r)
	payload, err := encodeRequest(jr)

	c.Assert(err, IsNil)
	c.Assert(payload[:5], DeepEquals, zabbixHeader)

	payloadSize := binary.LittleEndian.Uint64(payload[5:13])

	c.Assert(payloadSize, Equals, uint64(249))
}

func (s *JMXSuite) TestDecoder(c *C) {
	r := encodeReponse(respData1)

	size, err := decodeMeta(r)

	c.Assert(size, Equals, 81)
	c.Assert(err, IsNil)

	size, err = decodeMeta([]byte("ABCDEF"))

	c.Assert(size, Equals, -1)
	c.Assert(err, NotNil)

	jr, err := decodeResponse(r[13:])

	c.Assert(err, IsNil)
	c.Assert(jr, NotNil)
	c.Assert(jr.Data, HasLen, 1)

	jr, err = decodeResponse([]byte("ABCDEF"))

	c.Assert(err, NotNil)
	c.Assert(jr, IsNil)

	r = encodeReponse(respData2)
	jr, err = decodeResponse(r[13:])

	c.Assert(err, NotNil)
	c.Assert(jr, IsNil)
}

func (s *JMXSuite) TestBeansDecoder(c *C) {
	beans, err := ParseBeans(beansData)

	c.Assert(beans, HasLen, 2)
	c.Assert(err, IsNil)

	beans, err = ParseBeans("ABCD")

	c.Assert(beans, IsNil)
	c.Assert(err, NotNil)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func encodeReponse(data string) []byte {
	payload := []byte(data)
	sizeBuf := make([]byte, 8)

	binary.LittleEndian.PutUint64(sizeBuf, uint64(len(payload)))

	var buf bytes.Buffer

	buf.Write(zabbixHeader)
	buf.Write(sizeBuf)
	buf.Write(payload)

	return buf.Bytes()
}
