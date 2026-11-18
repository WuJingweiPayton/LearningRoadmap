// Task 1: send a simple ETH transfer on Sepolia.
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"homework05/internal/chain"
)

func main() {
	rpc := os.Getenv("SEPOLIA_RPC_URL")
	key := os.Getenv("PRIVATE_KEY")
	toHex := os.Getenv("TRANSFER_TO")
	if rpc == "" || key == "" || toHex == "" {
		fmt.Fprintln(os.Stderr, "set SEPOLIA_RPC_URL, PRIVATE_KEY, TRANSFER_TO")
		os.Exit(1)
	}

	amountEth := strings.TrimSpace(os.Getenv("TRANSFER_ETH"))
	if amountEth == "" {
		amountEth = "0.0001"
	}
	wei, ok := ethToWei(amountEth)
	if !ok {
		fatal(fmt.Errorf("invalid TRANSFER_ETH %q", amountEth))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	client, err := chain.Dial(ctx, rpc)
	if err != nil {
		fatal(err)
	}
	defer client.Close()

	pk, err := chain.PrivateKeyFromHex(key)
	if err != nil {
		fatal(err)
	}
	if !common.IsHexAddress(toHex) {
		fatal(fmt.Errorf("invalid TRANSFER_TO"))
	}
	to := common.HexToAddress(toHex)

	if _, err := chain.SendEther(ctx, client, pk, to, wei); err != nil {
		fatal(err)
	}
}

func ethToWei(eth string) (*big.Int, bool) {
	// minimal parse: "0.001" style
	if strings.Contains(eth, ".") {
		parts := strings.SplitN(eth, ".", 2)
		whole, ok := new(big.Int).SetString(parts[0], 10)
		if !ok {
			return nil, false
		}
		frac := parts[1]
		if len(frac) > 18 {
			return nil, false
		}
		for len(frac) < 18 {
			frac += "0"
		}
		f, ok := new(big.Int).SetString(frac, 10)
		if !ok {
			return nil, false
		}
		wei := new(big.Int).Mul(whole, big.NewInt(params.Ether))
		return wei.Add(wei, f), true
	}
	whole, ok := new(big.Int).SetString(eth, 10)
	if !ok {
		return nil, false
	}
	return new(big.Int).Mul(whole, big.NewInt(params.Ether)), true
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}