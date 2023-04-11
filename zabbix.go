// Package zabbix-jmx provides methods for working with Zabbix Java Gateway
package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io"
	"net"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Client is Zabbix JMX client
type Client struct {
	ConnectTimeout time.Duration
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration

	dialer *net.Dialer
	addr   *net.TCPAddr
}

// Request is basic request struct
type Request struct {
	Server   string
	Port     int
	Username string
	Password string
	Keys     []string
}

// Response contains response data
type Response []*ResponseData

// ResponseData contains value for requested key
type ResponseData struct {
	Value string `json:"value"`
}

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
	Error  string   `json:"error"`
	Status string   `json:"response"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewClient creates new client
func NewClient(address string) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp4", address)

	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{Timeout: time.Second * 5}

	return &Client{addr: addr, dialer: dialer}, nil
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

	err = writeToConnection(conn, encodeRequest(jr), c.WriteTimeout)

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
		Conn:     r.Server,
		Port:     r.Port,
		Username: r.Username,
		Password: r.Password,
		Endpoint: fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%d/jmxrmi", r.Server, r.Port),
		Keys:     r.Keys,
	}
}

// connectToServer makes connetion to Zabbix server
func connectToServer(c *Client) (*net.TCPConn, error) {
	if c.ConnectTimeout > 0 && c.dialer.Timeout != c.ConnectTimeout {
		c.dialer.Timeout = c.ConnectTimeout
	}

	conn, err := c.dialer.Dial(c.addr.Network(), c.addr.String())

	if err != nil {
		return nil, err
	}

	return conn.(*net.TCPConn), nil
}

// readFromConnection reads data fron connection
func readFromConnection(conn *net.TCPConn, buf []byte, timeout time.Duration) error {
	if timeout > 0 {
		conn.SetReadDeadline(time.Now().Add(timeout))
	}

	_, err := io.ReadFull(conn, buf)

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
