package p2p

import (
	"encoding/json"
	"flag"
	"fmt"

	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// Config is a struct storing important node configuration data.
type Config struct {
	Rendezvous      string      `json:"rendezvous"`      // The rendezvous point
	BootstrapPeers  AddressList `json:"bootstrapPeers"`  // A list of addresses of bootstrap peers from the DHT
	ListenAddresses AddressList `json:"listenAddresses"` // The addresses to listen on
	// ProtocolID      protocol.ID `json:"protocolID"`      // The protol ID
	ProtocolID string `json:"protocolID"` // The protol ID
}

// NewConfig generates a new configuration.
func NewConfig(rendezvous string, bootstrapPeers, listenAddresses AddressList, protocolID string) (Config, error) {
	return Config{}, nil
}

// ParseFlags parses the input flags and generates a Config struct.
func ParseFlags() (Config, error) {
	config := Config{}

	flag.StringVar(&config.Rendezvous, "rendezvous", "meeting point",
		"Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.Var(&config.BootstrapPeers, "peer", "Adds a peer to the bootstrap list")
	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")
	flag.StringVar(&config.ProtocolID, "pid", "/chat/1.1.0", "Sets the protocol ID to be used in the stream headers")
	flag.Parse()

	if len(config.BootstrapPeers) == 0 {
		config.BootstrapPeers = dht.DefaultBootstrapPeers
	}
	fmt.Println("there are not none")
	// fmt.Println(config.String())

	return config, nil

}

// String converts a Config struct to a string.
func (config Config) String() string {
	json, _ := json.MarshalIndent(config, "", "  ")
	return string(json)
}
