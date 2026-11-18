package chain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Dial(ctx context.Context, rpcURL string) (*ethclient.Client, error) {
	if rpcURL == "" {
		return nil, fmt.Errorf("RPC URL is empty")
	}
	return ethclient.DialContext(ctx, rpcURL)
}