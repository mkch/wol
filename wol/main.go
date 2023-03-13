/*
wol is a command line tool to send Wake-on-LAN(WoL) magic packet.
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"

	"github.com/mkch/wol"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(),
			`Usage of wol:

wol sends Wake-on-LAN magic packet.

wol [flags] MacAddress

All flags available:`)
		flag.PrintDefaults()
	}

	var addr = flag.String("addr", "255.255.255.255:6", "The UDP address to send the packet")
	var pass = flag.String("pass", "", "The SecureOn password")
	flag.Parse()

	var password *[6]byte
	if *pass != "" {
		var err error
		if password, err = parsePassword(*pass); err != nil {
			fmt.Fprintf(os.Stderr, "%v: invalid password: %v", *pass, err)
			os.Exit(1)
		}
	}

	if flag.NArg() != 1 {
		var errMsg string
		if flag.NArg() == 0 {
			errMsg = "no argument"
		} else {
			errMsg = "only one argument is allowed"
		}
		fmt.Fprintln(os.Stderr, errMsg)
		os.Exit(1)
	}

	destMac, err := net.ParseMAC(flag.Args()[0])
	if err == nil {
		err = wol.SendUDP(*addr, destMac, password)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

var hex = regexp.MustCompile("^([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$")

func parsePassword(pass string) (*[6]byte, error) {
	bytes := hex.FindStringSubmatch(pass)
	if bytes == nil {
		return nil, errors.New("not a 6-byte hexadecimal")
	}
	var parsedPass [6]byte
	for i := 0; i < len(parsedPass); i++ {
		// bitSize 9: unsigned byte
		if b, err := strconv.ParseInt(bytes[i+1], 16, 9); err != nil {
			log.Panic(err)
		} else {
			parsedPass[i] = byte(b)
		}
	}
	return &parsedPass, nil
}
