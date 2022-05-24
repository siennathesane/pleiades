package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type FsmKeyTests struct {
	suite.Suite
}

func TestFsmKey(t *testing.T) {
	suite.Run(t, new(FsmKeyTests))
}

func (fk *FsmKeyTests) TestParsePrn() {
	validTestPrn1 := "prn:global:pleiades:*:123456:bucket/test-bucket"
	validTestPrn1Struct := &PleiadesResourceName{
		Partition:    GlobalPartition,
		Service:      Pleiades,
		Region:       GlobalRegion,
		AccountId:    testAccountKey,
		ResourceType: Bucket,
		ResourceId:   "test-bucket",
	}

	var convertedPrn1 *PleiadesResourceName
	var err error
	require.NotPanics(fk.T(), func() {
		convertedPrn1, err = ParsePrn(validTestPrn1)
	})
	require.NoError(fk.T(), err, "there must not be an error parsing a valid prn")
	require.Equal(fk.T(), validTestPrn1Struct, convertedPrn1, "the converted prn must be parsed correctly")
}

func (fk *FsmKeyTests) TestToFsmRootKeyTest() {
	prn := &PleiadesResourceName{
		Partition:    GlobalPartition,
		Service:      Pleiades,
		Region:       GlobalRegion,
		AccountId:    testAccountKey,
		ResourceType: Bucket,
		ResourceId:   "noop-bucket",
	}
	preRendered := fmt.Sprintf("/global/pleiades/*/%d/test-bucket", testAccountKey)

	target := prn.ToFsmRootPath("test-bucket")
	require.Equal(fk.T(), preRendered, target, "the bucket prefix keys must match")
}
