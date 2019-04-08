package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

var (
	commentsEnabled = flag.Bool("comments", false, "Enable comments after post")

	commentsLimit = flag.Int("limit", 10, "Comments number per page")

	commentsServices = &AddrList{}
)

// go run flag.go --comments=true --servers="127.0.0.1:8081,127.0.0.1:8082"

func init() {
	flag.Var(commentsServices, "servers", "Comments number per page")
}

func main() {
	flag.Parse()

	if *commentsEnabled {
		fmt.Println("Comments per page", *commentsLimit)
		fmt.Println("Comments services", *commentsServices)
	} else {
		fmt.Println("Comments disabled")
	}
}

type AddrList []string

func (v *AddrList) String() string {
	return fmt.Sprint(*v)
}

func (v *AddrList) Set(in string) error {
	for _, addr := range strings.Split(in, ",") {
		ipRaw, _, err := net.SplitHostPort(addr)
		if err != nil {
			return fmt.Errorf("bad addr %v", addr)
		}
		ip := net.ParseIP(ipRaw)
		if ip.To4() == nil {
			return fmt.Errorf("invalid ipv4 addr %v", addr)
		}
		*v = append(*v, addr)
	}
	return nil
}
