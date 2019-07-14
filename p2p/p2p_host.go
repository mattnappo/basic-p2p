package p2p

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	maddr "github.com/multiformats/go-multiaddr"
	logging "github.com/whyrusleeping/go-logging"
)

var logger = log.Logger("rendezvous")

// InitLogger initializes the logger.
func InitLogger() {
	log.SetAllLoggers(logging.WARNING)
	log.SetLogLevel("rendezvous", "info")
}

// StartNode starts a new node.
func StartNode(config Config) error {
	logger.Info("starting node")

	ctx := context.Background()

	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs([]maddr.Multiaddr(config.ListenAddresses)...),
	)
	if err != nil {
		return err
	}
	logger.Info("Host created. This host is:", host.ID())
	logger.Info(host.Addrs())

	host.SetStreamHandler(protocol.ID(config.ProtocolID), streamHandler)

	// Initialize a new DHT CLIENT
	dht, err := dht.New(ctx, host)
	if err != nil {
		return err
	}

	err = dht.Bootstrap(ctx)
	if err != nil {
		return err
	}

	logger.Info("bootstrapping nodes")
	// Connect to IPFS bootstrap nodes
	var waitGroup sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
		logger.Info("bootstrap node:", peerAddr)
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		waitGroup.Add(1)
		go func() error {
			defer waitGroup.Done()
			err := host.Connect(ctx, *peerInfo)
			if err != nil {
				logger.Warning(err)
				// return err
			}

			logger.Info("connection established with bootstrap node '", *peerInfo, "'")

			return nil
		}()
	}
	waitGroup.Wait()

	logger.Info("Announcing ourselves to the network")
	routingDiscovery := discovery.NewRoutingDiscovery(dht)
	discovery.Advertise(ctx, routingDiscovery, config.Rendezvous)
	logger.Debug("Successfully announced!")

	// Look for other peers at the rendezvous
	peerChan, err := routingDiscovery.FindPeers(ctx, config.Rendezvous)
	if err != nil {
		return err
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		logger.Debug("found peer:", peer)

		logger.Debug("attempting to connect to peer ", peer)
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))
		if err != nil {
			logger.Warning("could not connect to peer", peer)
			continue
		} else {
			readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			go writeData(readWriter)
			go readData(readWriter)
		}

		logger.Info("successfully connected to peer", peer)

	}

	select {}

	// return nil
}

// streamHandler handles a stream connection to the local p2p node.
func streamHandler(stream network.Stream) {
	reader := bufio.NewReader(stream)
	writer := bufio.NewWriter(stream)
	ioStream := bufio.NewReadWriter(reader, writer)

	go readData(ioStream)
	go writeData(ioStream)

	select {}
}

func readData(stream *bufio.ReadWriter) {
	for {
		str, err := stream.ReadString('\n')
		if err != nil {
			fmt.Println("error reading from buffer in readData")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(stream *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		dataToSend, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading from stdin")
			panic(err)
		}

		_, err = stream.WriteString(dataToSend)
		if err != nil {
			fmt.Println("error writing to stream buffer in writeData")
			panic(err)
		}

		err = stream.Flush()
		if err != nil {
			fmt.Println("error flushing stream buffer")
			panic(err)
		}

	}

}
