package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/json"
	"strings"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Bean contains basic bean info
type Bean struct {
	Domain string `json:"{#JMXDOMAIN}"`
	Type   string `json:"{#JMXTYPE}"`
	Object string `json:"{#JMXOBJ}"`
	Name   string `json:"{#JMXNAME}"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

type jmxBeans struct {
	Data []*Bean `json:"data"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ParseBeans parses beans data
func ParseBeans(data string) ([]*Bean, error) {
	data = strings.Replace(data, `\"`, `"`, -1)

	beans := &jmxBeans{}
	err := json.Unmarshal([]byte(data), beans)

	if err != nil {
		return nil, err
	}

	return beans.Data, nil
}
