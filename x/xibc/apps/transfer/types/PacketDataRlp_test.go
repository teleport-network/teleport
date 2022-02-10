package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func Test_RLP_Encode_Struct(t *testing.T) {
	packetData := FungibleTokenPacketData{
		"SrcChain",
		"DestChain",
		"Sender",
		"Receiver",
		[]byte("Amount"),
		"Token",
		"OriToken",
	}
	Bz, err := rlp.EncodeToBytes(&packetData)
	require.NoError(t, err)
	fmt.Println(hex.EncodeToString(Bz))
	var data FungibleTokenPacketData
	err = rlp.Decode(bytes.NewReader(Bz), &data)
	require.NoError(t, err)
	fmt.Println(data.String())
}

func Test_RLP_Encode_Tuple(t *testing.T) {
	type test struct {
		A string
		B common.Address
	}
	testData := test{
		A: "test",
		B: common.BigToAddress(big.NewInt(1)),
	}
	Bz, err := rlp.EncodeToBytes(&testData)
	require.NoError(t, err)
	fmt.Println(hex.EncodeToString(Bz))
	var data test
	err = rlp.Decode(bytes.NewReader(Bz), &data)
	require.NoError(t, err)
	fmt.Println(data.A)
	fmt.Println(data.B.String())
}

func Test_RLP_Encode3(t *testing.T) {
	a := "testdata"
	Bz, err := rlp.EncodeToBytes(&a)
	require.NoError(t, err)
	fmt.Println(hex.EncodeToString(Bz))

	var testdata string
	err = rlp.Decode(bytes.NewReader(Bz), &testdata)
	require.NoError(t, err)
	fmt.Println(testdata)
}

func Test_ABI_Encode_Bool(t *testing.T) {
	Bool, _ := abi.NewType("bool", "", nil)
	pack, _ := abi.Arguments{{Type: Bool}}.Pack(
		true,
	)
	fmt.Println("0x" + hex.EncodeToString(pack))
}

func Test_ABI_Encode_String(t *testing.T) {
	String, _ := abi.NewType("string", "", nil)
	pack, _ := abi.Arguments{{Type: String}}.Pack(
		"testdata",
	)
	fmt.Println("0x" + hex.EncodeToString(pack))
}
