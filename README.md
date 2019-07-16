<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-zabbix-jmx.svg"/></a></p>

<p align="center">
  <a href="https://godoc.org/pkg.re/essentialkaos/zabbix-jmx.v1"><img src="https://godoc.org/pkg.re/essentialkaos/zabbix-jmx.v1?status.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/zabbix-jmx"><img src="https://goreportcard.com/badge/github.com/essentialkaos/zabbix-jmx"></a>
  <a href="https://travis-ci.org/essentialkaos/zabbix-jmx"><img src="https://travis-ci.org/essentialkaos/zabbix-jmx.svg"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-zabbix-jmx-master"><img alt="codebeat badge" src="https://codebeat.co/badges/5ea39cff-25ea-4c29-b5f6-97139a584e32" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`zabbix-jmx` is a Go package for retrieving and parsing data from Zabbix Java Gateway.

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.11+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get pkg.re/essentialkaos/zabbix-jmx.v1
```

For update to the latest stable release, do:

```
go get -u pkg.re/essentialkaos/zabbix-jmx.v1
```

### Usage example

```go
package main

import (
	"fmt"
	jmx "pkg.re/essentialkaos/zabbix-jmx.v1"
)

func main() {
	client, err := jmx.NewClient("127.0.0.1:9335")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	r := &jmx.Request{
		Server:   "domain.com",
		Port:     9093,
		Username: "admin",
		Password: "admin",
		Endpoint: `jmx["kafka.server:type=ReplicaManager,name=PartitionCount",Value]`,
	}

	resp, err := client.Get(r)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Resp value: %s\n", resp[0]["value"])
}
```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/zabbix-jmx.svg?branch=master)](https://travis-ci.org/essentialkaos/zabbix-jmx) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/zabbix-jmx.svg?branch=develop)](https://travis-ci.org/essentialkaos/zabbix-jmx) |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
