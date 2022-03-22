package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	_PORT_OK          = "50001"
	_PORT_META_ERR    = "50002"
	_PORT_PAYLOAD_ERR = "50003"
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

func (s *JMXSuite) SetUpSuite(c *C) {
	go runServer(c, _PORT_OK)
	go runServer(c, _PORT_META_ERR)
	go runServer(c, _PORT_PAYLOAD_ERR)

	time.Sleep(time.Second)
}

func (s *JMXSuite) TestClient(c *C) {
	client, err := NewClient("127.0.")

	c.Assert(client, IsNil)
	c.Assert(err, NotNil)

	client, err = NewClient("127.0.0.1:10051")

	c.Assert(client, NotNil)
	c.Assert(err, IsNil)
}

func (s *JMXSuite) TestClientGet(c *C) {
	client, err := NewClient("127.0.0.1:" + _PORT_OK)

	c.Assert(client, NotNil)
	c.Assert(err, IsNil)

	client.ConnectTimeout = time.Second * 3
	client.WriteTimeout = time.Second * 3
	client.ReadTimeout = time.Second * 3

	r := &Request{
		Server: "domain.com",
		Port:   9334,
		Keys:   []string{`jmx["kafka.server:type=ReplicaManager,name=PartitionCount",Value]`},
	}

	resp, err := client.Get(r)

	c.Assert(err, IsNil)
	c.Assert(resp, NotNil)
	c.Assert(resp[0].Value, Equals, "112.637")

	// -------

	client, err = NewClient("127.0.0.1:" + _PORT_META_ERR)

	c.Assert(client, NotNil)
	c.Assert(err, IsNil)

	resp, err = client.Get(r)

	c.Assert(err, NotNil)
	c.Assert(resp, IsNil)

	// -------

	client, err = NewClient("127.0.0.1:" + _PORT_PAYLOAD_ERR)

	c.Assert(client, NotNil)
	c.Assert(err, IsNil)

	resp, err = client.Get(r)

	c.Assert(err, NotNil)
	c.Assert(resp, IsNil)

	// -------

	client, err = NewClient("127.0.0.0:10000")

	c.Assert(client, NotNil)
	c.Assert(err, IsNil)

	resp, err = client.Get(r)

	c.Assert(err, NotNil)
	c.Assert(resp, IsNil)
}

func (s *JMXSuite) TestEncoder(c *C) {
	r := &Request{
		Server:   "domain.com",
		Port:     9334,
		Username: "admin",
		Password: "admin",
		Keys:     []string{`jmx["kafka.server:type=ReplicaManager,name=PartitionCount",Value]`},
	}

	payload := encodeRequest(convertRequest(r))

	c.Assert(payload[:5], DeepEquals, zabbixHeader)

	payloadSize := binary.LittleEndian.Uint64(payload[5:13])

	c.Assert(payloadSize, Equals, uint64(249))
}

func (s *JMXSuite) TestDecoder(c *C) {
	r := encodePayload([]byte(respData1))

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

	r = encodePayload([]byte(respData2))
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

func runServer(c *C, port string) {
	server, err := net.Listen("tcp4", "127.0.0.1:"+port)

	if err != nil {
		c.Fatal(err.Error())
	}

	defer server.Close()

	fmt.Printf("Fake server started on %s\n", port)

	for {
		conn, err := server.Accept()

		if err != nil {
			c.Fatal(err.Error())
		}

		handleRequest(conn, port)
	}
}

func handleRequest(conn net.Conn, port string) {
	switch port {
	case _PORT_OK:
		conn.Write(encodePayload([]byte(respData1)))
	case _PORT_META_ERR:
		conn.Write([]byte(`PAYLOAD12345678`))
	case _PORT_PAYLOAD_ERR:
		conn.Write(encodePayload([]byte(`PAYLOAD12345678`)))
	}

	conn.Close()
}
