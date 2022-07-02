package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AggregateTestSuite struct {
	suite.Suite
}

func TestAggregateTestSuite(t *testing.T) {
	suite.Run(t, new(AggregateTestSuite))
}

func (suite *AggregateTestSuite) SetupTest() {
}
