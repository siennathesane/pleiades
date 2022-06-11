package blaze

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"testing"

	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"r3t.io/pleiades/pkg/blaze/testdata"
)

func TestStreamManager(t *testing.T) {
	suite.Run(t, new(StreamServerTests))
}

type StreamServerTests struct {
	suite.Suite
	caCertPool *x509.CertPool
	keyPair    tls.Certificate
	listener   quic.Listener
	mux        DrpcLayer
	client     *testdata.CookieMonsterClient
}

func (smt *StreamServerTests) BeforeTest(suiteName, testName string) {
	caCert, err := ioutil.ReadFile("testdata/tls.ca")
	if err != nil {
		smt.Require().NoError(err, "there must not be an error reading the ca file")
	}
	smt.caCertPool = x509.NewCertPool()
	smt.caCertPool.AppendCertsFromPEM(caCert)

	smt.keyPair, err = tls.LoadX509KeyPair("testdata/tls.cert", "testdata/tls.key")
	if err != nil {
		smt.Require().NoError(err, "there must not be an error when loading the tls keys")
	}

	testTlsConfig := &tls.Config{
		RootCAs:      smt.caCertPool,
		Certificates: []tls.Certificate{smt.keyPair},
	}

	smt.listener, err = quic.ListenAddr("localhost:8080", testTlsConfig, &quic.Config{})
	smt.Require().NoError(err, "there must not be an error when starting the listener")

	smt.Require().NotPanics(func() {
		smt.mux, err = NewDrpcRouter()
	}, "there must not be a panic when building a new muxer")
	smt.Require().NoError(err, "there must not be an error when creating a new muxer")
	smt.Require().NotNil(smt.mux, "the muxer must not be nil")

	testServer := &testdata.TestCookieMonsterServer{}

	smt.Require().NotPanics(func() {
		err = testdata.DRPCRegisterCookieMonster(smt.mux, testServer)
	}, "there must not be a panic registering the test server")
	smt.Require().NoError(err, "there must not be an error when registering the test server")
}

func (smt *StreamServerTests) AfterTest(suiteName, testName string) {
	smt.Require().NoError(smt.listener.Close(), "there must not be an error when shutting down the listener")
	smt.keyPair = tls.Certificate{}
	smt.caCertPool = nil
}

func (smt *StreamServerTests) TestNewStreamManager() {

	testWriter := zerolog.NewTestWriter(smt.T())
	testLogger := zerolog.New(testWriter)

	var testStreamServer *ConnectionManager
	smt.Require().NotPanics(func() {
		testStreamServer = NewConnectionServer(smt.listener, smt.mux, testLogger)
	}, "there must not be an error when creating the new stream server")
	smt.Require().NotNil(testStreamServer, "the stream server must not be nil")


}
