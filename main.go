package main

import "github.com/xoreo/basic-p2p/p2p/"

func main() {

	config, err := p2p.ParseFlags()
	if err != nil {
		panic(err)
	}

	err = p2p.NewNode(config)
	if err != nil {
		panic(err)
	}

}
