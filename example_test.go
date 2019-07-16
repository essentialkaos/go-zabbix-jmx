package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func ExampleNewClient() {
	client, err := NewClient("127.0.0.1:9335")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	client.WriteTimeout = time.Second * 3
	client.ReadTimeout = time.Second * 3
}

func ExampleClient_Get() {
	client, err := NewClient("127.0.0.1:9335")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	r := &Request{
		Server:   "domain.com",
		Port:     9093,
		Username: "admin",
		Password: "admin",
		Keys:     []string{`jmx["kafka.server:type=ReplicaManager,name=PartitionCount",Value]`},
	}

	resp, err := client.Get(r)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Resp value: %s\n", resp[0].Value)
}
