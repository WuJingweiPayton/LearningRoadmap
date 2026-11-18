package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendEther signs and broadcasts an EIP-1559 transfer to Sepolia or any EVM network.
func SendEther(ctx context.Context, client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, amountWei *big.Int) (common.Hash, error) {
	from := crypto.PubkeyToAddress(privateKey.PublicKey)
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		return common.Hash{}, err
	}
	gasTip, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	gasFeeCap, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: gasTip,
		GasFeeCap: gasFeeCap,
		Gas:       21000,
		To:        &to,
		Value:     amountWei,
	})
	signed, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return common.Hash{}, err
	}
	if err := client.SendTransaction(ctx, signed); err != nil {
		return common.Hash{}, err
	}
	fmt.Printf("tx hash: %s\n", signed.Hash().Hex())
	if err := waitMined(ctx, client, signed.Hash()); err != nil {
		return signed.Hash(), err
	}
	return signed.Hash(), nil
}

func waitMined(ctx context.Context, client *ethclient.Client, hash common.Hash) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, hash)
			if err != nil {
				continue
			}
			fmt.Printf("mined in block %d (status %d)\n", receipt.BlockNumber.Uint64(), receipt.Status)
			return nil
		}
	}
}