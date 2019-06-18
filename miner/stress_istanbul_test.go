package miner

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/galaxy/galaxy-eth/common"
	"github.com/galaxy/galaxy-eth/common/hexutil"

	"github.com/galaxy/galaxy-eth/rpc"

	"github.com/galaxy/galaxy-eth/core/types"
	"github.com/galaxy/galaxy-eth/crypto"
	"github.com/galaxy/galaxy-eth/params"

	"github.com/galaxy/galaxy-eth/ethclient"
	"github.com/stretchr/testify/require"
)

func client(t *testing.T) *ethclient.Client {
	c, err := ethclient.Dial("http://192.168.0.62:8545")
	require.Nil(t, err)
	return c
}

func TestMiner_Istanbul(t *testing.T) {
	private, err := crypto.HexToECDSA("155345a175d5ef5cefabb43ac1771d3e981b0eca6081ffa8f448d1939d584eb9")
	require.Nil(t, err)
	from := crypto.PubkeyToAddress(private.PublicKey)
	c := client(t)
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
	c := client(t)
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

type addrs []common.Address

func (a addrs) equal(l []common.Address) bool {
	if len(a) != len(l) {
		return false
	}
	var b bool
	for _, v := range a {
		b = false
		for _, val := range l {
			if v == val {
				b = true
			}
		}
		if !b {
			return false
		}
	}
	return true
}

func TestMiner_Close(t *testing.T) {
	client, err := rpc.Dial("http://192.168.0.62:8545")
	require.Nil(t, err)
	var numberResult hexutil.Uint64
	err = client.Call(&numberResult, "eth_blockNumber")
	require.Nil(t, err)
	t.Log(uint64(numberResult))
}

func TestMiner_HashRate(t *testing.T) {
	host := []string{"192.168.0.62", "192.168.0.184", "192.168.0.171"}
	var clients []*rpc.Client
	for _, v := range host {
		client, err := rpc.Dial(fmt.Sprintf("http://%s:8545", v))
		require.Nil(t, err)
		clients = append(clients, client)
	}
	for _, v := range clients {
		var result []common.Address
		err := v.Call(&result, "istanbul_getValidators")
		require.Nil(t, err)

	}
}

func Test_Istanbul_addProposal(t *testing.T) {
	host := []string{"192.168.0.62", "192.168.0.184", "192.168.0.171"}
	header := types.Header{
		Extra: hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82faf86df86994dab21cf3587db83c9c5a871baf1a09006df90d6d946fe95084fd21d672e4db06e1ba92e9e3d2dc0bfc94955d6e2d9ec391b31600a7dded61b7024e7184aa944ba8f974d2d8355c1981999c2edf7399d39a9b92946b1fbb95c1d2758bf637d5683c30eaa4eb05497980c0"),
	}
	istanbul, err := types.ExtractIstanbulExtra(&header)
	require.Nil(t, err)
	var clients []*rpc.Client
	for _, v := range host {
		client, err := rpc.Dial(fmt.Sprintf("http://%s:8545", v))
		require.Nil(t, err)
		clients = append(clients, client)
	}
	require.NotNil(t, clients)
	var numberResult hexutil.Uint64
	err = clients[0].Call(&numberResult, "eth_blockNumber")
	require.Nil(t, err)
	numberFun := func() uint64 {
		var numberResult hexutil.Uint64
		err = clients[0].Call(&numberResult, "eth_blockNumber")
		require.Nil(t, err)
		return uint64(numberResult)
	}
	proposalAddr := common.BytesToAddress([]byte("proposal addr"))
	firstNumber := numberFun()
	for _, v := range clients {
		err = v.Call(nil, "istanbul_propose", proposalAddr, true)
		require.Nil(t, err)
	}
	for {
		time.Sleep(time.Second * 10)
		currentNumber := numberFun()
		if currentNumber >= firstNumber+1024 {
			break
		}
	}
	nowAddr := append(istanbul.Validators, proposalAddr)
	var result []common.Address
	err = clients[0].Call(&result, "istanbul_getValidators")
	require.Nil(t, err)
	require.True(t, addrs(nowAddr).equal(result))
}

func Test_Istanbul_dropProposal(t *testing.T) {
	host := []string{"192.168.0.62", "192.168.0.184", "192.168.0.171"}
	header := types.Header{
		Extra: hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82faf86df86994dab21cf3587db83c9c5a871baf1a09006df90d6d946fe95084fd21d672e4db06e1ba92e9e3d2dc0bfc94955d6e2d9ec391b31600a7dded61b7024e7184aa944ba8f974d2d8355c1981999c2edf7399d39a9b92946b1fbb95c1d2758bf637d5683c30eaa4eb05497980c0"),
	}
	istanbul, err := types.ExtractIstanbulExtra(&header)
	require.Nil(t, err)
	require.Equal(t, len(istanbul.Validators), 5)
	var clients []*rpc.Client
	for _, v := range host {
		client, err := rpc.Dial(fmt.Sprintf("http://%s:8545", v))
		require.Nil(t, err)
		clients = append(clients, client)
	}
	require.NotNil(t, clients)
	numberFun := func() uint64 {
		var numberResult hexutil.Uint64
		err = clients[0].Call(&numberResult, "eth_blockNumber")
		require.Nil(t, err)
		return uint64(numberResult)
	}

	proposalAddr := istanbul.Validators[0]
	var nowAddr addrs
	for _, v := range istanbul.Validators {
		if v != proposalAddr {
			nowAddr = append(nowAddr, v)
		}
	}
	require.Equal(t, len(nowAddr), 4)
	firstNumber := numberFun()
	for _, v := range clients {
		err = v.Call(nil, "istanbul_propose", proposalAddr, false)
		require.Nil(t, err)
	}
	for {
		time.Sleep(time.Second * 10)
		currentNumber := numberFun()
		if currentNumber >= firstNumber+1024 {
			break
		}
	}
	var result []common.Address
	err = clients[0].Call(&result, "istanbul_getValidators")
	require.Nil(t, err)
	require.True(t, nowAddr.equal(result))
}

func TestMiner_Mining(t *testing.T) {
	host := []string{"192.168.0.62", "192.168.0.184", "192.168.0.171"}
	//nodes := []string{
	//	"enode://6319f32c27f6cc58529672e40219881966e199bcc68ba4f098e11a74931ba46f0e7ab463f9a4f9d93bf4623105b81bb17c06a75ed37cb2a78ae0a19c8d6173ae@192.168.0.62:30303",
	//	"enode://fd2d1885263d663910c95e923fc44c4138dd1de35e7321fc9d8b5d91bb781f1064ec810dac0979ad7f32e38e68233ecd83d7a312cdba2ebb77350e272d15f372@192.168.0.184:30303",
	//	"enode://4f45efb895880bc9321eb39f82b3b407bf20f176e4e8d2ee68081287c348c8390ea43ebd481ff754829b9639478779602efe3421b9051aa5607d2cd74c2a908f@192.168.0.171:30303",
	//}
	var clients []*rpc.Client
	for _, v := range host {
		client, err := rpc.Dial(fmt.Sprintf("http://%s:8545", v))
		require.Nil(t, err)
		clients = append(clients, client)
	}
	for _, v := range clients {
		//for _, node := range nodes {
		//	err := v.Call(nil, "admin_addPeer", node)
		//	require.Nil(t, err)
		//}
		//var infos []*p2p.PeerInfo
		//err := v.Call(&infos, "admin_peers")
		//require.Nil(t, err)
		//require.Len(t, infos, 2)
		err := v.Call(nil, "miner_start")
		require.Nil(t, err)
	}
}
