package blaze

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
	"r3t.io/pleiades/pkg/utils"
)

// TestKit is a pre-defined yet somewhat controllable QUIC transport and connection server
// which is exclusively for testing purposes. It allows you to use a functional QUIC transport
// and connection server for unit testing purposes via sockets.
// Note: currently only one server can run at a time
// todo (sienna): support test server pooling? or just stream provisioning?
type TestKit struct {
	t      *testing.T
	logger zerolog.Logger

	quicConfig *quic.Config
	listener   quic.Listener
	dialConn   quic.Connection

	certPool *x509.CertPool
	keyPair  tls.Certificate
	tlsConf  *tls.Config

	ctx     context.Context
	mux     *Router
	server  *Server
	running bool
}

func NewTestKit(t *testing.T) *TestKit {
	tk := &TestKit{
		logger: utils.NewTestLogger(t),
		ctx:    context.Background(),
	}

	tk.GenerateTlsConfig()

	tk.quicConfig = &quic.Config{MaxIdleTimeout: 300 * time.Second}

	var err error
	tk.listener, err = quic.ListenAddr(testServerAddr, tk.tlsConf, tk.quicConfig)
	if err != nil {
		tk.t.Error(err, "there was an error starting the testkit listener")
	}

	time.Sleep(1 * time.Second)

	tk.dialConn, err = quic.DialAddr(testServerAddr, tk.tlsConf, tk.quicConfig)
	if err != nil {
		tk.t.Error(err, "there was an error dialing the testkit server")
	}

	time.Sleep(1 * time.Second)

	return tk
}

func (tk *TestKit) GetListener() quic.Listener {
	return tk.listener
}

func (tk *TestKit) CloseListener() {
	err := tk.listener.Close()
	if err != nil {
		tk.t.Error(err, "there was an error trying to close the listener")
	}
}

func (tk *TestKit) GetConnection() quic.Connection {
	return tk.dialConn
}

func (tk *TestKit) NewConnectionStream() quic.Stream {
	stream, err := tk.dialConn.OpenStream()
	if err != nil {
		tk.t.Error(err, "there was an error opening a new testkit stream")
	}

	return stream
}

func (tk *TestKit) GenerateTlsConfig() *tls.Config {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		tk.t.Error(err, "error reaching consul for tls config")
	}

	pair, _, err := client.KV().Get(testConsulTlsKey, &api.QueryOptions{})
	if err != nil {
		tk.t.Error(err, "error fetching tls config from consul")
	}

	var tlsPayload *tlsConfig
	if err := json.Unmarshal(pair.Value, &tlsPayload); err != nil {
		tk.t.Error(err, "error unmarshalling tls config")
	}

	tk.certPool = x509.NewCertPool()
	tk.certPool.AppendCertsFromPEM([]byte(tlsPayload.Data.IssuingCa))

	tk.keyPair, err = tls.X509KeyPair([]byte(tlsPayload.Data.Certificate), []byte(tlsPayload.Data.PrivateKey))
	if err != nil {
		tk.t.Error(err, "there must not be an error when loading the tls keys")
	}

	tk.tlsConf = &tls.Config{
		RootCAs:      tk.certPool,
		Certificates: []tls.Certificate{tk.keyPair},
		NextProtos:   []string{"blaze-test-server"},
	}

	return tk.tlsConf
}

type tlsConfig struct {
	Data          tlsData     `json:"data"`
	LeaseDuration int64       `json:"lease_duration"`
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	RequestID     string      `json:"request_id"`
	Warnings      interface{} `json:"warnings"`
}

type tlsData struct {
	CaChain        []string `json:"ca_chain"`
	Certificate    string   `json:"certificate"`
	Expiration     int64    `json:"expiration"`
	IssuingCa      string   `json:"issuing_ca"`
	PrivateKey     string   `json:"private_key"`
	PrivateKeyType string   `json:"private_key_type"`
	SerialNumber   string   `json:"serial_number"`
}

// TestKitServerArgs represents the arguments you can use to configure the TestKit server
type TestKitServerArgs struct {
	// Muxer is required, otherwise RPCs won't be routed properly
	Muxer *Router
	// AutoStart determines whether to start the server
	AutoStart bool
}

func (tk *TestKit) NewServer(args *TestKitServerArgs) {
	if args.Muxer == nil {
		tk.t.Error("the muxer cannot be nil")
	}
	tk.mux = args.Muxer

	tk.server = NewServer(tk.listener, tk.mux, tk.logger)

	if args.AutoStart {
		err := tk.server.Start(tk.ctx)
		if err != nil {
			tk.t.Error(err, "there was an error starting the testkit server")
		}
		tk.running = true
	}
}

func (tk *TestKit) Start() {
	if tk.running {
		return
	}
	err := tk.server.Start(tk.ctx)
	if err != nil {
		tk.t.Error(err, "there was an error starting the testkit server")
	}
}

func (tk *TestKit) Stop() {
	if !tk.running {
		return
	}

	err := tk.server.Stop(tk.ctx)
	if err != nil {
		tk.t.Error(err, "there was an error stopping the testkit server")
	}
}
