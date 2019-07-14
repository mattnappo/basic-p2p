package main

import (
	"flag"
	"fmt"

	"github.com/xoreo/basic-p2p/p2p"
)

func main() {
	help := flag.Bool("h", false, "Display Help")
	if *help {
		fmt.Println("This program demonstrates a simple p2p chat application using libp2p")
		fmt.Println()
		fmt.Println("Usage: Run './chat in two different terminals. Let them connect to the bootstrap nodes, announce themselves and connect to the peers")
		flag.PrintDefaults()
		return
	}
	// flag.Parse()
	p2p.InitLogger()

	config, err := p2p.ParseFlags()
	if err != nil {
		panic(err)
	}

	err = p2p.StartNode(config)
	if err != nil {
		panic(err)
	}

}
