package blaze
//
//import (
//	"context"
//	"crypto/tls"
//	"crypto/x509"
//	"io/ioutil"
//	"os"
//	"testing"
//	"time"
//
//	"github.com/lucas-clemente/quic-go"
//	"github.com/rs/zerolog"
//	"github.com/stretchr/testify/suite"
//	"r3t.io/pleiades/pkg/blaze/testdata"
//)
//
//func TestFunctional(t *testing.T) {
//	suite.Run(t, new(FunctionalTests))
//}
//
//type FunctionalTests struct {
//	suite.Suite
//	caCertPool *x509.CertPool
//	keyPair    tls.Certificate
//	listener   quic.Listener
//	tls        *tls.Config
//	logger     zerolog.Logger
//	mux        *Router
//	client     *testdata.CookieMonsterClient
//	config *quic.Config
//}
//
//func (ft *FunctionalTests) SetupSuite() {
//	testWriter := zerolog.NewTestWriter(ft.T())
//	ft.logger = zerolog.New(testWriter)
//}
//
//func (ft *FunctionalTests) BeforeTest(suiteName, testName string) {
//	caCert, err := ioutil.ReadFile("testdata/tls.ca")
//	if err != nil {
//		ft.Require().NoError(err, "there must not be an error reading the ca file")
//	}
//	ft.caCertPool = x509.NewCertPool()
//	ft.caCertPool.AppendCertsFromPEM(caCert)
//
//	ft.keyPair, err = tls.LoadX509KeyPair("testdata/tls.cert", "testdata/tls.key")
//	if err != nil {
//		ft.Require().NoError(err, "there must not be an error when loading the tls keys")
//	}
//
//	ft.tls = &tls.Config{
//		RootCAs:      ft.caCertPool,
//		Certificates: []tls.Certificate{ft.keyPair},
//		NextProtos:   []string{"multiplexed-string-tests"},
//	}
//
//	ft.config = &quic.Config{MaxIdleTimeout: 300 * time.Second}
//
//	ft.listener, err = quic.ListenAddr("localhost:8080", ft.tls, ft.config)
//	ft.Require().NoError(err, "there must not be an error when starting the listener")
//
//	ft.Require().NotPanics(func() {
//		ft.mux = NewRouter()
//	}, "there must not be a panic when building a new muxer")
//	ft.Require().NotNil(ft.mux, "the muxer must not be nil")
//
//	testServer := &testdata.TestCookieMonsterServer{}
//
//	ft.Require().NotPanics(func() {
//		err = testdata.DRPCRegisterCookieMonster(ft.mux, testServer)
//	}, "there must not be a panic registering the test server")
//	ft.Require().NoError(err, "there must not be an error when registering the test server")
//}
//
//func (ft *FunctionalTests) AfterTest(suiteName, testName string) {
//	ft.Require().NoError(ft.listener.Close(), "there must not be an error when shutting down the listener")
//	ft.keyPair = tls.Certificate{}
//	ft.caCertPool = nil
//}
//
//func (ft *FunctionalTests) TestMultiplexedStream() {
//	// get environment variable CI
//	ci := os.Getenv("CI")
//	if ci == "true" {
//		ft.T().Skip("this test is not run in CI")
//	}
//
//	var testStreamServer *Server
//	ft.Require().NotPanics(func() {
//		testStreamServer = NewTestKitServer(ft.listener, ft.mux, ft.logger)
//	}, "there must not be an error when creating the new stream server")
//	ft.Require().NotNil(testStreamServer, "the stream server must not be nil")
//
//	ctx := context.Background()
//	err := testStreamServer.Start(ctx)
//	defer testStreamServer.Stop(ctx)
//
//	// gotta let the network stack open the connections
//	time.Sleep(5 * time.Second)
//
//	ft.Require().NoError(err, "there must not be an error when starting the stream server")
//
//	dialConn, err := quic.DialAddr(testServerAddr, ft.tls, ft.config)
//	ft.Require().NoError(err, "there must not be an error when dialing the test server")
//
//	stream, err := dialConn.OpenStream()
//	ft.Require().NoError(err, "there must not be an error when creating a new stream")
//
//	clientStream := NewConnectionStream(stream, ft.mux, ft.logger)
//	ft.Require().NotNil(clientStream, "the client stream must not be null")
//
//	client := testdata.NewDRPCCookieMonsterClient(clientStream)
//	ft.Require().NotNil(client, "the cookie monster client must not be null")
//
//	ctx, cancel := context.WithTimeout(context.Background(), ft.config.MaxIdleTimeout)
//	resp, err := client.EatCookie(ctx, &testdata.Cookie{Type: testdata.Cookie_Oatmeal})
//	ft.Require().NoError(err, "there must not be an error when trying to eat a cookie")
//	ft.Assert().Equal(testdata.Cookie_Oatmeal, resp.Cookie.Type, "the cookie types should match")
//	cancel()
//}
