// vsm stands for Versioned State Machine. It contains a server that is run via
// Docker compose and used in the version compatability test.
package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	coretypes "github.com/tendermint/tendermint/types"
)

type BlockResult struct {
	AppVersion int
	AppHash    []byte
	DataHash   []byte
}

type BlockService struct{}

func (t *BlockService) GetBlockResult(block *coretypes.Block, result *BlockResult) error {
	// spin up a new test app for each request that way test cases don't
	// interact with each other in version compatability test.
	fmt.Printf("GetBlockResult %v\n", block)
	*result = BlockResult{
		AppVersion: 1,
		AppHash:    []byte("apphash"),
		DataHash:   []byte("datahash"),
	}
	fmt.Printf("result %v\n", result)
	return nil
}

func main() {
	blockService := new(BlockService)
	rpc.Register(blockService)
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
