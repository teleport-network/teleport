package keeper_test

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	transfer "github.com/teleport-network/teleport/syscontracts/xibc_transfer"
	"github.com/teleport-network/teleport/x/xibc/apps/transfer/types"
)

func (suite *KeeperTestSuite) TestRecvPacket() {
	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"token come in",
			func() {
				srcToken := "token0"
				srcChain := "chain0"
				nativeChain := "teleport"
				sender := "sender"
				amount := big.NewInt(123450)
				scale := uint8(0)

				// deploy ERC20
				tokenAddress := suite.DeployERC20MintableContract(transfer.TransferContractAddress, "kitty", "kit", uint8(18))

				// bind ERC20 trace
				err := suite.app.AggregateKeeper.RegisterERC20Trace(suite.ctx, tokenAddress, srcToken, srcChain, scale)
				suite.Require().NoError(err)

				// check ERC20 trace
				_, _, exist, err := suite.app.AggregateKeeper.QueryERC20Trace(suite.ctx, tokenAddress, srcChain)
				suite.Require().NoError(err)
				suite.Require().True(exist)

				// receive packet and check Ack
				suite.RecvPacket(
					types.FungibleTokenPacketData{
						SrcChain:  srcChain,
						DestChain: nativeChain,
						Sender:    sender,
						Receiver:  suite.address.String(),
						Amount:    amount.Bytes(),
						Token:     srcToken,
						OriToken:  "",
					},
				)

				// check balance and total supply
				totalSupply := suite.TotalSupply(tokenAddress)
				suite.Require().Equal(amount, totalSupply)
				balance := suite.BalanceOf(tokenAddress, suite.address)
				suite.Require().Equal(totalSupply, balance)
			},
		},
		{
			"ERC20 token back to origin",
			func() {
				srcToken := "token0"
				srcChain := "chain0"
				nativeChain := "teleport"
				sender := "sender"
				supply := big.NewInt(2000000)
				amount := big.NewInt(100)

				// deploy ERC20
				tokenAddress := suite.DeployERC20MintableContract(types.ModuleAddress, "kitty", "kit", uint8(18))

				// mint token
				suite.Mint(tokenAddress, types.ModuleAddress, suite.address, supply)
				totalSupply := suite.TotalSupply(tokenAddress)
				suite.Require().Equal(totalSupply, supply)

				// approve
				suite.Approve(tokenAddress, suite.address, transfer.TransferContractAddress, amount)
				allowance := suite.Allowance(tokenAddress, suite.address, transfer.TransferContractAddress)
				suite.Require().Equal(amount, allowance)

				// send transfer ERC20
				suite.SendTransfer(
					suite.address,
					types.TransferData{
						TokenAddress: tokenAddress,
						Receiver:     sender,
						Amount:       amount,
						DestChain:    srcChain,
						RelayChain:   "",
					},
					types.Fee{
						TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Amount:       big.NewInt(0),
					},
				)

				// check balance
				contractBalance := suite.BalanceOf(tokenAddress, transfer.TransferContractAddress)
				suite.Require().Equal(amount, contractBalance)
				senderBalance := suite.BalanceOf(tokenAddress, suite.address)
				suite.Require().Equal(supply.Sub(supply, amount), senderBalance)

				// query token out
				outAmount := suite.OutTokens(tokenAddress, srcChain)
				suite.Require().Equal(amount, outAmount)

				// receive packet
				suite.RecvPacket(
					types.FungibleTokenPacketData{
						SrcChain:  srcChain,
						DestChain: nativeChain,
						Sender:    sender,
						Receiver:  suite.address.String(),
						Amount:    amount.Bytes(),
						Token:     srcToken,
						OriToken:  tokenAddress.String(),
					},
				)

				// query token out
				outAmount = suite.OutTokens(tokenAddress, srcChain)
				suite.Require().Equal(big.NewInt(0).String(), outAmount.String())
			},
		},
		{
			"base token back to origin",
			func() {
				srcToken := "token0"
				srcChain := "chain0"
				nativeChain := "teleport"
				sender := "sender"
				supply := big.NewInt(2000000)
				amount := big.NewInt(100)
				zeroAddress := common.BigToAddress(big.NewInt(0))

				// add balance
				suite.Require().NoError(suite.app.EvmKeeper.SetBalance(suite.ctx, suite.address, supply))
				balance := suite.app.EvmKeeper.GetBalance(suite.ctx, suite.address)
				suite.Require().Equal(supply, balance)

				// send transfer base
				suite.SendTransfer(
					suite.address,
					types.TransferData{
						TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Receiver:     sender,
						Amount:       amount,
						DestChain:    srcChain,
						RelayChain:   "",
					},
					types.Fee{
						TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Amount:       big.NewInt(0),
					},
				)

				// check balance
				contractBalance := suite.app.EvmKeeper.GetBalance(suite.ctx, transfer.TransferContractAddress)
				suite.Require().Equal(amount, contractBalance)
				senderBalance := suite.app.EvmKeeper.GetBalance(suite.ctx, suite.address)
				suite.Require().Equal(supply.Sub(supply, amount), senderBalance)

				// query token out
				outAmount := suite.OutTokens(zeroAddress, srcChain)
				suite.Require().Equal(amount, outAmount)

				// receive packet
				suite.RecvPacket(
					types.FungibleTokenPacketData{
						SrcChain:  srcChain,
						DestChain: nativeChain,
						Sender:    sender,
						Receiver:  suite.address.String(),
						Amount:    amount.Bytes(),
						Token:     srcToken,
						OriToken:  zeroAddress.String(),
					},
				)

				// query token out
				outAmount = suite.OutTokens(zeroAddress, srcChain)
				suite.Require().Equal(big.NewInt(0).String(), outAmount.String())
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}

func (suite *KeeperTestSuite) TestAcknowledgementPacket() {
	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"ack success",
			func() {
				suite.AcknowledgementPacket(
					types.FungibleTokenPacketData{},
					[]byte{byte(1)},
				)
			},
		},
		{
			"refund crossed chain token",
			func() {
				srcToken := "token0"
				srcChain := "chain0"
				nativeChain := "teleport"
				sender := "sender"
				amount := big.NewInt(123450)
				scale := uint8(0)

				// deploy ERC20
				tokenAddress := suite.DeployERC20MintableContract(transfer.TransferContractAddress, "kitty", "kit", uint8(18))

				// bind ERC20 trace
				err := suite.app.AggregateKeeper.RegisterERC20Trace(suite.ctx, tokenAddress, srcToken, srcChain, scale)
				suite.Require().NoError(err)

				// receive packet and check Ack
				suite.RecvPacket(
					types.FungibleTokenPacketData{
						SrcChain:  srcChain,
						DestChain: nativeChain,
						Sender:    sender,
						Receiver:  suite.address.String(),
						Amount:    amount.Bytes(),
						Token:     srcToken,
						OriToken:  "",
					},
				)

				// check balance and total supply
				totalSupply := suite.TotalSupply(tokenAddress)
				suite.Require().Equal(amount, totalSupply)
				balance := suite.BalanceOf(tokenAddress, suite.address)
				suite.Require().Equal(totalSupply, balance)

				// approve
				suite.Approve(tokenAddress, suite.address, transfer.TransferContractAddress, amount)
				allowance := suite.Allowance(tokenAddress, suite.address, transfer.TransferContractAddress)
				suite.Require().Equal(amount, allowance)

				// send transfer ERC20
				suite.SendTransfer(
					suite.address,
					types.TransferData{
						TokenAddress: tokenAddress,
						Receiver:     sender,
						Amount:       amount,
						DestChain:    srcChain,
						RelayChain:   "",
					},
					types.Fee{
						TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Amount:       big.NewInt(0),
					},
				)

				// check balance and total supply
				totalSupply = suite.TotalSupply(tokenAddress)
				suite.Require().Equal(big.NewInt(0).String(), totalSupply.String())
				balance = suite.BalanceOf(tokenAddress, suite.address)
				suite.Require().Equal(big.NewInt(0).String(), balance.String())

				// receive error ack
				suite.AcknowledgementPacket(
					types.FungibleTokenPacketData{
						SrcChain:  nativeChain,
						DestChain: srcChain,
						Sender:    suite.address.String(),
						Receiver:  sender,
						Amount:    amount.Bytes(),
						Token:     tokenAddress.String(),
						OriToken:  srcToken,
					},
					[]byte{},
				)

				// check balance and total supply
				totalSupply = suite.TotalSupply(tokenAddress)
				suite.Require().Equal(amount, totalSupply)
				balance = suite.BalanceOf(tokenAddress, suite.address)
				suite.Require().Equal(totalSupply, balance)
			},
		},
		{
			"refund native ERC20 token",
			func() {
				nativeChain := "teleport"
				dstChain := "chain0"
				receiver := "receiver"
				supply := big.NewInt(2000000)
				amount := big.NewInt(100)

				// deploy ERC20
				tokenAddress := suite.DeployERC20MintableContract(types.ModuleAddress, "kitty", "kit", uint8(18))

				// mint token
				suite.Mint(tokenAddress, types.ModuleAddress, suite.address, supply)
				totalSupply := suite.TotalSupply(tokenAddress)
				suite.Require().Equal(totalSupply, supply)

				// approve
				suite.Approve(tokenAddress, suite.address, transfer.TransferContractAddress, amount)
				allowance := suite.Allowance(tokenAddress, suite.address, transfer.TransferContractAddress)
				suite.Require().Equal(amount, allowance)

				// send transfer ERC20
				suite.SendTransfer(
					suite.address,
					types.TransferData{
						TokenAddress: tokenAddress,
						Receiver:     receiver,
						Amount:       amount,
						DestChain:    dstChain,
						RelayChain:   "",
					},
					types.Fee{
						TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Amount:       big.NewInt(0),
					},
				)

				// check balance
				contractBalance := suite.BalanceOf(tokenAddress, transfer.TransferContractAddress)
				suite.Require().Equal(amount, contractBalance)
				senderBalance := suite.BalanceOf(tokenAddress, suite.address)
				suite.Require().Equal(big.NewInt(0).Sub(supply, amount), senderBalance)

				// query token out
				outAmount := suite.OutTokens(tokenAddress, dstChain)
				suite.Require().Equal(amount, outAmount)

				// receive error ack
				suite.AcknowledgementPacket(
					types.FungibleTokenPacketData{
						SrcChain:  nativeChain,
						DestChain: dstChain,
						Sender:    suite.address.String(),
						Receiver:  suite.address.String(),
						Amount:    amount.Bytes(),
						Token:     tokenAddress.String(),
						OriToken:  "",
					},
					[]byte{},
				)

				senderBalance = suite.BalanceOf(tokenAddress, suite.address)
				suite.Require().Equal(supply, senderBalance)

				// query token out
				outAmount = suite.OutTokens(tokenAddress, dstChain)
				suite.Require().Equal(big.NewInt(0).String(), outAmount.String())
			},
		},
		{
			"refund base token",
			func() {
				dstChain := "chain0"
				nativeChain := "teleport"
				receiver := "sender"
				supply := big.NewInt(2000000)
				amount := big.NewInt(100)
				zeroAddress := common.BigToAddress(big.NewInt(0))

				// add balance
				suite.Require().NoError(suite.app.EvmKeeper.SetBalance(suite.ctx, suite.address, supply))
				balance := suite.app.EvmKeeper.GetBalance(suite.ctx, suite.address)
				suite.Require().Equal(supply, balance)

				// send transfer base
				suite.SendTransfer(
					suite.address,
					types.TransferData{
						TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Receiver:     receiver,
						Amount:       amount,
						DestChain:    dstChain,
						RelayChain:   "",
					},
					types.Fee{
						TokenAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Amount:       big.NewInt(0),
					},
				)

				// check balance
				contractBalance := suite.app.EvmKeeper.GetBalance(suite.ctx, transfer.TransferContractAddress)
				suite.Require().Equal(amount, contractBalance)
				senderBalance := suite.app.EvmKeeper.GetBalance(suite.ctx, suite.address)
				suite.Require().Equal(supply.Sub(supply, amount), senderBalance)

				// query token out
				outAmount := suite.OutTokens(zeroAddress, dstChain)
				suite.Require().Equal(amount, outAmount)

				// receive error ack
				suite.AcknowledgementPacket(
					types.FungibleTokenPacketData{
						SrcChain:  nativeChain,
						DestChain: dstChain,
						Sender:    suite.address.String(),
						Receiver:  receiver,
						Amount:    amount.Bytes(),
						Token:     zeroAddress.String(),
						OriToken:  "",
					},
					[]byte{},
				)

				// query token out
				outAmount = suite.OutTokens(zeroAddress, dstChain)
				suite.Require().Equal(big.NewInt(0).String(), outAmount.String())
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			tc.malleate()
		})
	}
}
