package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"

	jmx "github.com/essentialkaos/go-zabbix-jmx"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "zabbix-jmx-get"
	VER  = "2.0.0"
	DESC = "Tool for fetching data from Zabbix Java Gateway"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	OPT_USERNAME = "u:user"
	OPT_PASSWORD = "p:password"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_USERNAME: {},
	OPT_PASSWORD: {},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.BOOL},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error(" - "))
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
	case options.GetB(OPT_HELP) || len(args) < 3:
		genUsage().Print()
		os.Exit(0)
	}

	err := process(args)

	if err != nil {
		terminal.Error(err)
		os.Exit(1)
	}
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

// process starts keys processing
func process(args options.Arguments) error {
	gwHost, gwPort, srvHost, srvPort, keys, err := parseArguments(args)

	if err != nil {
		return err
	}

	client, err := jmx.NewClient(gwHost + ":" + strconv.Itoa(gwPort))

	if err != nil {
		return fmt.Errorf("Can't configure client: %v", err)
	}

	client.ConnectTimeout = 3 * time.Second
	client.WriteTimeout = 5 * time.Second
	client.ReadTimeout = 5 * time.Second

	resp, err := client.Get(makeRequest(srvHost, srvPort, keys))

	if err != nil {
		return fmt.Errorf("Can't send response: %v", err)
	}

	renderResponse(resp, keys)

	return nil
}

// parseArguments parses command arguments
func parseArguments(args options.Arguments) (string, int, string, int, []string, error) {
	gw := args.Get(0).String()
	gwHost, gwPort, ok := strings.Cut(gw, ":")

	if !ok {
		return "", 0, "", 0, nil, fmt.Errorf("Invalid gateway: You must specify the gateway as host:port")
	}

	gwPortInt, err := strconv.Atoi(gwPort)

	if err != nil {
		return "", 0, "", 0, nil, fmt.Errorf("Invalid gateway port: %v", err)
	}

	if gwPortInt < 1025 || gwPortInt > 65535 {
		return "", 0, "", 0, nil, fmt.Errorf("Gateway port must be in range 1024-65535")
	}

	srv := args.Get(1).String()
	srvHost, srvPort, ok := strings.Cut(srv, ":")

	if !ok {
		return "", 0, "", 0, nil, fmt.Errorf("Invalid server: You must specify the server as host:port")
	}

	srvPortInt, err := strconv.Atoi(srvPort)

	if err != nil {
		return "", 0, "", 0, nil, fmt.Errorf("Invalid server port: %v", err)
	}

	if srvPortInt < 1025 || srvPortInt > 65535 {
		return "", 0, "", 0, nil, fmt.Errorf("Server port must be in range 1024-65535")
	}

	keys := args[4:].Strings()

	return gwHost, gwPortInt, srvHost, srvPortInt, keys, err
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
func makeRequest(serverHost string, serverPort int, keys []string) *jmx.Request {
	r := &jmx.Request{
		Server: serverHost,
		Port:   serverPort,
		Keys:   keys,
	}

	if options.Has(OPT_USERNAME) {
		r.Username = options.GetS(OPT_USERNAME)
		r.Password = options.GetS(OPT_PASSWORD)
	}

	return r
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, APP))
	case "fish":
		fmt.Print(fish.Generate(info, APP))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, APP))
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
	info := usage.NewInfo("", "gateway", "server", "key…")

	info.AddOption(OPT_USERNAME, "JMX server user", "username")
	info.AddOption(OPT_PASSWORD, "JMX server password", "password")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		`127.0.0.1:10052 srv1.domain.com:9093 'jmx["kafka.server:type=BrokerTopicMetrics,name=BytesInPerSec",OneMinuteRate]'`,
		"Request kafka metrics",
	)

	info.AddExample(
		`127.0.0.1:10052 srv1.domain.com:9093 'jmx.discovery[beans,"*:type=GarbageCollector,name=*"]'`,
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
