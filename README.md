<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-zabbix-jmx.svg"/></a></p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/essentialkaos/go-zabbix-jmx"><img src="https://pkg.go.dev/badge/github.com/essentialkaos/go-zabbix-jmx" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/go-zabbix-jmx"><img src="https://goreportcard.com/badge/github.com/essentialkaos/go-zabbix-jmx"></a>
  <a href="https://github.com/essentialkaos/go-zabbix-jmx/actions"><img src="https://github.com/essentialkaos/go-zabbix-jmx/workflows/CI/badge.svg" alt="GitHub Actions Status" /></a>
  <a href="https://github.com/essentialkaos/go-zabbix-jmx/actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos/go-zabbix-jmx/workflows/CodeQL/badge.svg" /></a>
  <a href='https://coveralls.io/github/essentialkaos/go-zabbix-jmx?branch=develop'><img src='https://coveralls.io/repos/github/essentialkaos/go-zabbix-jmx/badge.svg?branch=develop' alt='Coverage Status' /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-go-zabbix-jmx-master"><img alt="codebeat badge" src="https://codebeat.co/badges/31cf4383-04c5-4ba4-85d2-85835e41d7fc" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#zabbix-jmx-get">zabbix-jmx-get</a> • <a href="#build-status">Build Status</a> • <a href="#license">License</a></p>

<br/>

`zabbix-jmx` is a Go package for retrieving and parsing data from Zabbix Java Gateway.

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.12+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get pkg.re/essentialkaos/go-zabbix-jmx.v1
```

For update to the latest stable release, do:

```
go get -u pkg.re/essentialkaos/go-zabbix-jmx.v1
```

### Usage example

```go
package main

import (
	"fmt"
	jmx "pkg.re/essentialkaos/go-zabbix-jmx.v1"
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
		Keys:     []string{`jmx["kafka.server:type=ReplicaManager,name=PartitionCount",Value]`},
	}

	resp, err := client.Get(r)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Resp value: %s\n", resp[0].Value)
}
```

### `zabbix-jmx-get`

We also provide a command-line tool `zabbix-jmx-get` for retrieving data from Zabbix Java Gateway.

#### Installation

From sources:

```
go get -u github.com/essentialkaos/go-zabbix-jmx/cmd/zabbix-jmx-get
```

Prebuilt binaries:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) zabbix-jmx-get
```

#### Usage examples

```
$ zabbix-jmx-get -h 127.0.0.1 -p 10052 -H kfk-node1.domain.com -P 9093 'jmx.discovery[beans,"*:type=BrokerTopicMetrics,name=*"]'

kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=TotalProduceRequestsPerSec TotalProduceRequestsPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=BytesOutPerSec BytesOutPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=BytesInPerSec BytesInPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=BytesRejectedPerSec BytesRejectedPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=FailedProduceRequestsPerSec FailedProduceRequestsPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=FetchMessageConversionsPerSec FetchMessageConversionsPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=MessagesInPerSec MessagesInPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=FailedFetchRequestsPerSec FailedFetchRequestsPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=ProduceMessageConversionsPerSec ProduceMessageConversionsPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=ReplicationBytesInPerSec ReplicationBytesInPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=TotalFetchRequestsPerSec TotalFetchRequestsPerSec
kafka.server BrokerTopicMetrics kafka.server:type=BrokerTopicMetrics,name=ReplicationBytesOutPerSec ReplicationBytesOutPerSec

$ zabbix-jmx-get -h 127.0.0.1 -p 10052 -H kfk-node1.domain.com -P 9093 'jmx["kafka.server:type=BrokerTopicMetrics,name=BytesInPerSec",OneMinuteRate]'

5668479.780357378

```

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://github.com/essentialkaos/go-zabbix-jmx/workflows/CI/badge.svg?branch=master)](https://github.com/essentialkaos/go-zabbix-jmx/actions) |
| `develop` | [![CI](https://github.com/essentialkaos/go-zabbix-jmx/workflows/CI/badge.svg?branch=develop)](https://github.com/essentialkaos/go-zabbix-jmx/actions) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
