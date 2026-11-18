// Task 2: deploy / read / increment Counter using abigen bindings.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"homework05/bindings"
	"homework05/internal/chain"
)

func main() {
	rpc := os.Getenv("SEPOLIA_RPC_URL")
	key := os.Getenv("PRIVATE_KEY")
	action := strings.ToLower(strings.TrimSpace(os.Getenv("COUNTER_ACTION")))
	if action == "" {
		action = "read"
	}
	if rpc == "" || key == "" {
		fmt.Fprintln(os.Stderr, "set SEPOLIA_RPC_URL and PRIVATE_KEY")
		os.Exit(1)
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
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		fatal(err)
	}
	auth, err := chain.TransactOpts(chainID, pk)
	if err != nil {
		fatal(err)
	}

	switch action {
	case "deploy":
		addr, tx, _, err := bindings.DeployCounter(auth, client)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("Counter deployed at %s\n", addr.Hex())
		fmt.Printf("deploy tx: %s\n", tx.Hash().Hex())
	case "increment", "read":
		addrHex := strings.TrimSpace(os.Getenv("COUNTER_ADDRESS"))
		if addrHex == "" {
			fatal(fmt.Errorf("set COUNTER_ADDRESS for %s", action))
		}
		if !common.IsHexAddress(addrHex) {
			fatal(fmt.Errorf("invalid COUNTER_ADDRESS"))
		}
		addr := common.HexToAddress(addrHex)
		contract, err := bindings.NewCounter(addr, client)
		if err != nil {
			fatal(err)
		}
		if action == "increment" {
			tx, err := contract.Increment(auth)
			if err != nil {
				fatal(err)
			}
			fmt.Printf("increment tx: %s\n", tx.Hash().Hex())
		}
		count, err := contract.Count(&bind.CallOpts{Context: ctx})
		if err != nil {
			fatal(err)
		}
		fmt.Printf("count = %s\n", count.String())
	default:
		fatal(fmt.Errorf("unknown COUNTER_ACTION %q (deploy|increment|read)", action))
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}