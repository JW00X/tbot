package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/JW00X/tbot/lib"
)

var verbose = flag.Bool("v", false, "verbose output")
var timeout = flag.Int("t", 10, "timeout in seconds")
var address = flag.String("a", "", "source IP address")
var cookies = flag.Bool("c", false, "use cookies")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] <model ID>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		return
	}
	modelID := flag.Arg(0)
	if !lib.ModelIDRegexp.MatchString(modelID) {
		fmt.Println("invalid model ID")
		return
	}
	client := lib.HTTPClientWithTimeoutAndAddress(*timeout, *address, *cookies, true)
	checker := &lib.BongaCamsChecker{}
	checker.Init(checker, lib.CheckerConfig{Clients: []*lib.Client{client}, Dbg: *verbose})
	fmt.Println(checker.CheckStatusSingle(modelID))
}
