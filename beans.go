package jmx

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
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
