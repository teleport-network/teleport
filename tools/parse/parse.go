package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/bitdao-io/bitnetwork/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/evmos/ethermint/encoding"
)

func main() {

	parseAddr := flag.Bool("isParseAddress", false, "whether to parse address")
	parseBala := flag.Bool("isParseBalances", false, "whether to parse balances")
	parsebase64 := flag.Bool("isParseBase64", false, "whether to parse base64")
	inputData := flag.String("input", "", "input data waiting to parse")
	flag.Parse() // 解析参数

	if *parseAddr {
		keyResult, bitAddress := parseAddress(*inputData)
		fmt.Printf("eth address: %s\n", keyResult)
		fmt.Printf("bit address: %s\n", bitAddress)
	} else if *parseBala {
		valueResult := parseBalance(*inputData)
		fmt.Println(valueResult)
	} else if *parsebase64 {
		base64Result := base64ToHex(*inputData)
		fmt.Println(base64Result)
	}
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
	bech32Prefixes := "bit"
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
