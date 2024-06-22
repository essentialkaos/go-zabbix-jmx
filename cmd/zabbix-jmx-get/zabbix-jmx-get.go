package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/terminal"
	"github.com/essentialkaos/ek/v12/terminal/tty"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"

	jmx "github.com/essentialkaos/go-zabbix-jmx"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "zabbix-jmx-get"
	VER  = "1.2.0"
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

	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
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
	OPT_HELP:         {Type: options.BOOL},
	OPT_VER:          {Type: options.BOOL},

	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.String())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout().Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_HELP) || options.GetS(OPT_GATEWAY_HOST) == "true" || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	process(args.Strings())
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// validateOptions checks for required options
func validateOptions() {
	optionsSet := true

	if !options.Has(OPT_GATEWAY_HOST) {
		terminal.Error("Option %s is required", options.F(OPT_GATEWAY_HOST))
		optionsSet = false
	}

	if !options.Has(OPT_GATEWAY_PORT) {
		terminal.Error("Option %s is required", options.F(OPT_GATEWAY_PORT))
		optionsSet = false
	}

	if !options.Has(OPT_SERVER_HOST) {
		terminal.Error("Option %s is required", options.F(OPT_SERVER_HOST))
		optionsSet = false
	}

	if !options.Has(OPT_SERVER_PORT) {
		terminal.Error("Option %s is required", options.F(OPT_SERVER_PORT))
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
		terminal.Error(err)
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

// printErrorAndExit print error message and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	terminal.Error(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, APP))
	case "fish":
		fmt.Printf(fish.Generate(info, APP))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout()))
}

// genUsage generates usage info
func genUsage() *usage.Info {
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

	return info
}

// genAbout generates info about version
func genAbout() *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	return about
}
