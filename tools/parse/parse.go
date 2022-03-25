package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/teleport-network/teleport/app"
	"github.com/tharsis/ethermint/encoding"
)

func main() {
	testKey := "021493354845030274cd4bf1686abd60ab28ec52e1a76174656c65"
	keyResult, teleAddress := parseAddress(testKey)
	testValue := "0a056174656c6512173431333632343936353930323030303030303030303030"
	valueResult := parseBalance(testValue)

	testBase64 := "MTA1MDA5MDAwMDAwMDAwMGF0ZWxl"
	base64Result := base64ToHex(testBase64)

	fmt.Println(keyResult)
	fmt.Println(teleAddress)
	fmt.Println(valueResult)
	fmt.Println(base64Result)
}

func parseAddress(key string) (string, string) {
	addrBytes, _ := hex.DecodeString(key[4:44])
	return key[4:44], newBech32Output(addrBytes).Formats[0]
}

func parseBalance(value string) sdk.Coin {
	bz, _ := hex.DecodeString(value)

	var balance sdk.Coin
	encoding.MakeConfig(app.ModuleBasics).Marshaler.MustUnmarshal(bz, &balance)

	return balance
}

type bech32Output struct {
	Formats []string `json:"formats"`
}

func newBech32Output(bs []byte) bech32Output {
	bech32Prefixes := "teleport"
	out := bech32Output{Formats: make([]string, len(bech32Prefixes))}

	bech32Addr, err := bech32.ConvertAndEncode(bech32Prefixes, bs)
	if err != nil {
		panic(err)
	}

	out.Formats[0] = bech32Addr
	return out
}

func base64ToHex(base64Str string) string {
	var decodeBytes []byte
	var err error
	if decodeBytes, err = base64.StdEncoding.DecodeString(base64Str); err != nil {
		panic(err.Error())
	}

	return string(decodeBytes)
}
