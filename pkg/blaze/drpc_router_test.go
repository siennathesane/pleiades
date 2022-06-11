package blaze

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/blaze/testdata"
)

type DrpcMuxerTests struct {
	suite.Suite
}

func TestDrpcMuxer(t *testing.T) {
	suite.Run(t, new(DrpcMuxerTests))
}

func (dmt *DrpcMuxerTests) TestNewMuxerWithRegistration() {
	var mux *DrpcRouter
	var err error
	dmt.Require().NotPanics(func() {
		mux, err = NewDrpcRouter()
	}, "there must not be a panic when building a new muxer")
	dmt.Require().NoError(err, "there must not be an error when creating a new muxer")
	dmt.Require().NotNil(mux, "the muxer must not be nil")
	dmt.Require().NotNil(mux.targets, "the target map must not be nil")

	testServer := &testdata.TestCookieMonsterServer{}

	dmt.Require().NotPanics(func() {
		err = testdata.DRPCRegisterCookieMonster(mux, testServer)
	}, "there must not be a panic registering the test server")
	dmt.Require().NoError(err, "there must not be an error when registering the test server")

	testRpcName := "/testdata.CookieMonster/EatCookie"
	val, ok := mux.targets[testRpcName]
	dmt.Require().Equal(1, len(mux.targets), "there must only be one rpc enabled")
	dmt.Require().True(ok, "the test rpc key must be present")
	dmt.Require().Equal(testRpcName, val.Name, "the stored value must be named properly")
	dmt.Require().Equal(&testdata.TestCookieMonsterServer{}, val.Server, "the server types must match")
}
