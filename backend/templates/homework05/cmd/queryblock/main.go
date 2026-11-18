// Task 1: query a block on Sepolia via ethclient.
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"homework05/internal/chain"
)

func main() {
	rpc := os.Getenv("SEPOLIA_RPC_URL")
	if rpc == "" {
		fmt.Fprintln(os.Stderr, "set SEPOLIA_RPC_URL")
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := chain.Dial(ctx, rpc)
	if err != nil {
		fatal(err)
	}
	defer client.Close()

	var number *big.Int
	if raw := strings.TrimSpace(os.Getenv("BLOCK_NUMBER")); raw != "" {
		n, ok := new(big.Int).SetString(raw, 10)
		if !ok {
			fatal(fmt.Errorf("invalid BLOCK_NUMBER %q", raw))
		}
		number = n
	}

	if err := chain.PrintBlock(ctx, client, number); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}