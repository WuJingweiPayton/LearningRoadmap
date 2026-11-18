package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// PrintBlock writes hash, timestamp, and tx count for the given block (nil number = latest).
func PrintBlock(ctx context.Context, client *ethclient.Client, number *big.Int) error {
	block, err := client.BlockByNumber(ctx, number)
	if err != nil {
		return err
	}
	label := "latest"
	if number != nil {
		label = number.String()
	}
	fmt.Printf("block number: %s\n", label)
	fmt.Printf("hash:         %s\n", block.Hash().Hex())
	fmt.Printf("timestamp:    %d\n", block.Time())
	fmt.Printf("tx count:     %d\n", len(block.Transactions()))
	return nil
}