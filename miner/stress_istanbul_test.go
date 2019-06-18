package miner

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/galaxy/galaxy-eth/core/types"
	"github.com/galaxy/galaxy-eth/crypto"
	"github.com/galaxy/galaxy-eth/params"

	"github.com/galaxy/galaxy-eth/ethclient"
	"github.com/stretchr/testify/require"
)

func TestMiner_Istanbul(t *testing.T) {
	private, err := crypto.HexToECDSA("155345a175d5ef5cefabb43ac1771d3e981b0eca6081ffa8f448d1939d584eb9")
	require.Nil(t, err)
	from := crypto.PubkeyToAddress(private.PublicKey)
	c, err := ethclient.Dial("http://192.168.0.62:8545")
	require.Nil(t, err)
	txSum := 10000
	for i := 0; i < txSum; i++ {
		nonce, err := c.PendingNonceAt(context.Background(), from)
		require.Nil(t, err)
		to := crypto.CreateAddress(from, nonce)
		tx := types.NewTransaction(nonce, to, big.NewInt(params.GWei), params.TxGas, big.NewInt(params.GWei), nil)
		tx, err = types.SignTx(tx, types.HomesteadSigner{}, private)
		require.Nil(t, err)
		err = c.SendTransaction(context.Background(), tx)
		require.Nil(t, err)
		fmt.Println("hash=", tx.Hash().Hex(), ",i=", i)
	}
}

func TestNew(t *testing.T) {
	c, err := ethclient.Dial("http://192.168.0.62:8545")
	require.Nil(t, err)
	h, err := c.HeaderByNumber(context.Background(), big.NewInt(0))
	require.Nil(t, err)
	t.Log(h.GasLimit)
	//for i := 0; i <= 210; i++ {
	//	block, err := c.BlockByNumber(context.Background(), big.NewInt(int64(i)))
	//	require.Nil(t, err)
	//	for _, v := range block.Transactions() {
	//		t.Log(v.Hash().Hex(), i)
	//	}
	//}
}
