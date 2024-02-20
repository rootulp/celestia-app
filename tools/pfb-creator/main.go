package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tendermint/tendermint/rpc/client/http"
)

func main() {
	if err := Run(); err != nil {
		fmt.Printf("ERROR: %s", err.Error())
	}
}

func Run() error {
	if len(os.Args) < 2 {
		fmt.Println("Usage: NODE_RPC path/to/ascii.json\n", os.Args[0])
		return nil
	}

	url := os.Args[1]
	c, err := http.New(url, "/websocket")
	if err != nil {
		return err
	}
	resp, err := c.Status(context.Background())
	if err != nil {
		return err
	}

	lastHeight := resp.SyncInfo.LatestBlockHeight
	chainID := resp.NodeInfo.Network

	fmt.Printf("lastHeight %v\n", lastHeight)
	fmt.Printf("chainID %v\n", chainID)

	asciiJsonPath := os.Args[2]
	content, err := os.ReadFile(asciiJsonPath)
	if err != nil {
		return err
	}
	fmt.Printf("content %v", content)
	return nil
}
