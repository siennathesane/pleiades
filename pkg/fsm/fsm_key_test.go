
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

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
