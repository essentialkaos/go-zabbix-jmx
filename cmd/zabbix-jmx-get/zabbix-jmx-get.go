package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"
	"time"

	"pkg.re/essentialkaos/ek.v10/fmtc"
	"pkg.re/essentialkaos/ek.v10/options"
	"pkg.re/essentialkaos/ek.v10/usage"

	jmx "github.com/essentialkaos/zabbix-jmx"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "zabbix-jmx-get"
	VER  = "1.0.0"
	DESC = "Tool for fetching data from Zabbix Java Gateway"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	OPT_GATEWAY_HOST = "h:host"
	OPT_GATEWAY_PORT = "p:port"
	OPT_SERVER_HOST  = "H:server-host"
	OPT_SERVER_PORT  = "P:server-port"
	OPT_USERNAME     = "user"
	OPT_PASSWORD     = "password"
	OPT_NO_COLOR     = "nc:no-color"
	OPT_HELP         = "help"
	OPT_VER          = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_GATEWAY_HOST: {Type: options.MIXED},
	OPT_GATEWAY_PORT: {Type: options.INT, Min: 1025, Max: 65535},
	OPT_SERVER_HOST:  {},
	OPT_SERVER_PORT:  {Type: options.INT, Min: 1025, Max: 65535},
	OPT_USERNAME:     {},
	OPT_PASSWORD:     {},
	OPT_NO_COLOR:     {Type: options.BOOL},
	OPT_HELP:         {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:          {Type: options.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	keys, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	configureUI()

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.Has(OPT_GATEWAY_HOST) && options.GetS(OPT_GATEWAY_HOST) == "true" {
		showUsage()
		return
	}

	if options.GetB(OPT_HELP) || len(keys) == 0 {
		showUsage()
		return
	}

	process(keys)
}

// checkOptions checks required options
func checkOptions() {
	optionsSet := true

	if !options.Has(OPT_GATEWAY_HOST) {
		printError("Option --%s is required", "host")
		optionsSet = false
	}

	if !options.Has(OPT_GATEWAY_PORT) {
		printError("Option --%s is required", "port")
		optionsSet = false
	}

	if !options.Has(OPT_SERVER_HOST) {
		printError("Option --%s is required", "server-host")
		optionsSet = false
	}

	if !options.Has(OPT_SERVER_PORT) {
		printError("Option --%s is required", "server-port")
		optionsSet = false
	}

	if !optionsSet {
		os.Exit(1)
	}
}

// process starts keys processing
func process(keys []string) {
	client, err := jmx.NewClient(options.GetS(OPT_GATEWAY_HOST) + ":" + options.GetS(OPT_GATEWAY_PORT))

	if err != nil {
		printErrorAndExit(err.Error())
	}

	client.ConnectTimeout = 3 * time.Second
	client.WriteTimeout = 5 * time.Second
	client.ReadTimeout = 5 * time.Second

	resp, err := client.Get(makeRequest(keys))

	if err != nil {
		printErrorAndExit(err.Error())
	}

	renderResponse(resp, keys)
}

// renderResponse renders response data
func renderResponse(resp jmx.Response, keys []string) {
	for index, data := range resp {
		isBeans := isBeansData(keys, index)

		switch isBeans {
		case true:
			renderBeansData(data.Value)
		default:
			fmt.Println(data.Value)
		}
	}
}

// renderBeansData renders beans response
func renderBeansData(data string) {
	beans, err := jmx.ParseBeans(data)

	if err != nil {
		printError(err.Error())
		return
	}

	for _, bean := range beans {
		fmt.Printf(
			"%s %s %s %s\n",
			bean.Domain, bean.Type,
			bean.Object, bean.Name,
		)
	}
}

// isBeansData returns true if key with given index is beans request
func isBeansData(keys []string, index int) bool {
	if len(keys) > index {
		return false
	}

	return strings.HasPrefix(keys[index], "jmx.discovery[beans")
}

// makeRequest creates new request
func makeRequest(keys []string) *jmx.Request {
	r := &jmx.Request{
		Server: options.GetS(OPT_SERVER_HOST),
		Port:   options.GetI(OPT_SERVER_PORT),
		Keys:   keys,
	}

	if options.Has(OPT_USERNAME) {
		r.Username = options.GetS(OPT_USERNAME)
		r.Password = options.GetS(OPT_PASSWORD)
	}

	return r
}

// configureUI configure UI on start
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	info := usage.NewInfo("", "key…")

	info.AddOption(OPT_GATEWAY_HOST, "Java gateway host", "host")
	info.AddOption(OPT_GATEWAY_PORT, "Java gateway port {s-}(1025-65535){!}", "port")
	info.AddOption(OPT_SERVER_HOST, "JMX server host", "port")
	info.AddOption(OPT_SERVER_PORT, "JMX server port {s-}(1025-65535){!}", "port")
	info.AddOption(OPT_USERNAME, "JMX server user", "username")
	info.AddOption(OPT_PASSWORD, "JMX server password", "password")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		`-h 127.0.0.1 -p 10052 -H srv1.domain.com -P 9093 jmx["kafka.server:type=BrokerTopicMetrics,name=BytesInPerSec",OneMinuteRate]`,
		"Request kafka metrics",
	)

	info.AddExample(
		`-h 127.0.0.1 -p 10052 -H srv1.domain.com -P 9093 'jmx.discovery[beans,"*:type=GarbageCollector,name=*"]'`,
		"Request discovery info",
	)

	info.Render()
}

// showAbout shows info about version
func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
	}

	about.Render()
}
