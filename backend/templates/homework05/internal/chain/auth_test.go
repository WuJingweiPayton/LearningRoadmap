package chain

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestPrivateKeyFromHex(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	raw := "0x" + hex.EncodeToString(crypto.FromECDSA(key))
	parsed, err := PrivateKeyFromHex(raw)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.D.Cmp(key.D) != 0 {
		t.Fatal("key mismatch")
	}
}