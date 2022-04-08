package keeper_test

func (suite *KeeperTestSuite) TestRegisterRelayer() {
	address := "address0"
	names := []string{"name1", "name2"}
	addresses := []string{"address1", "address2"}

	suite.keeper.RegisterRelayers(
		suite.ctx,
		address,
		names,
		addresses,
	)

	relayers := suite.keeper.GetAllRelayers(suite.ctx)
	relayer, found := suite.keeper.GetRelayer(suite.ctx, address)
	suite.Require().True(found)
	suite.Require().Equal(relayers[0], relayer)

	relayerAddressOnOtherChain, found := suite.keeper.GetRelayerAddressOnOtherChain(suite.ctx, names[0], address)
	suite.Require().True(found)
	suite.Require().Equal(addresses[0], relayerAddressOnOtherChain)

	relayerAddressOnTeleport, found := suite.keeper.GetRelayerAddressOnTeleport(suite.ctx, names[0], addresses[0])
	suite.Require().True(found)
	suite.Require().Equal(address, relayerAddressOnTeleport)
}
