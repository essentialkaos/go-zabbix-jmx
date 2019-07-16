// Package zabbix-jmx provides methods for working with Zabbix Java Gateway
package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"net"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Client is Zabbix JMX client
type Client struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration

	addr *net.TCPAddr
}

// Request is basic request
type Request struct {
	Server   string
	Port     int
	Username string
	Password string
	Endpoint string
	Keys     []string
}

// Response contains response data
type Response []map[string]string

// ////////////////////////////////////////////////////////////////////////////////// //

type jmxRequest struct {
	Request  string   `json:"request"`
	Conn     string   `json:"conn"`
	Port     int      `json:"port"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Endpoint string   `json:"jmx_endpoint"`
	Keys     []string `json:"keys"`
}

type jmxResponse struct {
	Data   Response `json:"data"`
	Status string   `json:"response"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewClient creates new client
func NewClient(address string) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)

	if err != nil {
		return nil, err
	}

	return &Client{addr: addr}, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Get fetches data from Java Gateway
func (c *Client) Get(r *Request) (Response, error) {
	jr := convertRequest(r)
	conn, err := connectToServer(c)

	if err != nil {
		return nil, err
	}

	defer conn.Close() // Zabbix doesn't support persistent connections

	reqData, err := encodeRequest(jr)

	if err != nil {
		return nil, err
	}

	err = writeToConnection(conn, reqData, c.WriteTimeout)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, 13)
	err = readFromConnection(conn, buf, c.ReadTimeout)

	if err != nil {
		return nil, err
	}

	size, err := decodeMeta(buf)

	if err != nil {
		return nil, err
	}

	buf = make([]byte, size)
	err = readFromConnection(conn, buf, c.ReadTimeout)

	if err != nil {
		return nil, err
	}

	resp, err := decodeResponse(buf)

	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// convertRequest convert request to jmx request
func convertRequest(r *Request) *jmxRequest {
	return &jmxRequest{
		Request:  "java gateway jmx",
		Conn:     fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%d/jmxrmi", r.Server, r.Port),
		Port:     r.Port,
		Username: r.Username,
		Password: r.Password,
		Endpoint: r.Endpoint,
		Keys:     r.Keys,
	}
}

// connectToServer makes connetion to Zabbix server
func connectToServer(c *Client) (*net.TCPConn, error) {
	conn, err := net.DialTCP(c.addr.Network(), nil, c.addr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// readFromConnection reads data fron connection
func readFromConnection(conn *net.TCPConn, buf []byte, timeout time.Duration) error {
	if timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(timeout))
	}

	_, err := conn.Read(buf)

	return err
}

// writeToConnection writes data into connection
func writeToConnection(conn *net.TCPConn, data []byte, timeout time.Duration) error {
	if timeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(timeout))
	}

	_, err := conn.Write(data)

	return err
}
