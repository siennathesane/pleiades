package conf

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/nats-io/nats-server/v2/server"
)

type logLevel string

const (
	natsConfigTemplate string   = "hosts/%s/config/nats"
	trace              logLevel = "trace"
	debug              logLevel = "debug"
	traceVerbose       logLevel = "trace-verbose"
)

// https://github.com/nats-io/natscli/blob/13775fd198e89a6358bb66dc59bad27694e33c08/cli/server_run_command.go#L325
// https://pkg.go.dev/github.com/nats-io/nats-server/v2/server#Options
// https://github.com/nats-io/nats-server/blob/v2.8.2/server/opts.go#L4561
// https://docs.nats.io/using-nats/nats-tools
// https://dev.to/karanpratapsingh/embedding-nats-in-go-19o
// https://github.com/nats-io/nats-top/blob/220bc613e4a5d26c8a2a24131e75be19e588645b/util/toputils.go#L94
type natsServerOpts struct {
	HostAddress        string   `json:"host-address,omitempty"`
	Port               int      `json:"port,omitempty"`
	ServerName         string   `json:"server-name,omitempty"`
	HttpMetricsPort    int      `json:"http-metrics-port,omitempty"`
	ClientAdvertiseUrl string   `json:"client-advertise-url,omitempty"`
	EnableJetStream    bool     `json:"enable-jet-stream,omitempty"`
	Username           string   `json:"username,omitempty"`
	Password           string   `json:"password,omitempty"`
	Routes             []string `json:"routes,omitempty"`
	ClusterUrl         string   `json:"cluster-url,omitempty"`
	ClusterName        string   `json:"cluster-name,omitempty"`
	ClusterAdvertise   bool     `json:"cluster-advertise,omitempty"`
	ConnectRetryCount  int      `json:"connect-retry-count,omitempty"`
	LogLevel           logLevel `json:"log-level,omitempty"`
}

func NewNatsConfig(client *api.Client, env *EnvironmentConfig) (*server.Options, error) {
	pair, _, err := client.KV().Get(fmt.Sprintf(natsConfigTemplate, env.Hostname), &api.QueryOptions{})
	if err != nil {
		return nil, err
	}

	var config *natsServerOpts
	if err := json.Unmarshal(pair.Value, &config); err != nil {
		return nil, err
	}

	serverOpts := &server.Options{
		ServerName:                 config.ServerName,
		Host:                       config.HostAddress,
		Port:                       config.Port,
		ClientAdvertise:            config.ClientAdvertiseUrl,
		Logtime:                    true,
		MaxConn:                    50_000,
		MaxSubs:                    10_000,
		MaxSubTokens:               100,
		Nkeys:                      nil,
		Users:                      nil,
		Accounts:                   nil,
		NoAuthUser:                 "",
		SystemAccount:              "",
		NoSystemAccount:            false,
		Username:                   "",
		Password:                   "",
		Authorization:              "",
		PingInterval:               0,
		MaxPingsOut:                0,
		HTTPHost:                   "",
		HTTPPort:                   0,
		HTTPBasePath:               "",
		HTTPSPort:                  0,
		AuthTimeout:                0,
		MaxControlLine:             0,
		MaxPayload:                 0,
		MaxPending:                 0,
		Cluster:                    server.ClusterOpts{},
		Gateway:                    server.GatewayOpts{},
		LeafNode:                   server.LeafNodeOpts{},
		JetStream:                  false,
		JetStreamMaxMemory:         0,
		JetStreamMaxStore:          0,
		JetStreamDomain:            "",
		JetStreamExtHint:           "",
		JetStreamKey:               "",
		JetStreamUniqueTag:         "",
		JetStreamLimits:            server.JSLimitOpts{},
		StoreDir:                   "",
		JsAccDefaultDomain:         nil,
		Websocket:                  server.WebsocketOpts{},
		MQTT:                       server.MQTTOpts{},
		ProfPort:                   0,
		PidFile:                    "",
		PortsFileDir:               "",
		LogFile:                    "",
		LogSizeLimit:               0,
		Syslog:                     false,
		RemoteSyslog:               "",
		Routes:                     nil,
		RoutesStr:                  "",
		TLSTimeout:                 0,
		TLS:                        false,
		TLSVerify:                  false,
		TLSMap:                     false,
		TLSCert:                    "",
		TLSKey:                     "",
		TLSCaCert:                  "",
		TLSConfig:                  nil,
		TLSPinnedCerts:             nil,
		TLSRateLimit:               0,
		AllowNonTLS:                false,
		WriteDeadline:              0,
		MaxClosedClients:           0,
		LameDuckDuration:           0,
		LameDuckGracePeriod:        0,
		MaxTracedMsgLen:            0,
		TrustedKeys:                nil,
		TrustedOperators:           nil,
		AccountResolver:            nil,
		AccountResolverTLSConfig:   nil,
		AlwaysEnableNonce:          false,
		CustomClientAuthentication: nil,
		CustomRouterAuthentication: nil,
		CheckConfig:                false,
		ConnectErrorReports:        0,
		ReconnectErrorReports:      0,
		Tags:                       nil,
		OCSPConfig:                 nil,
	}

	switch config.LogLevel {
	case trace:
		serverOpts.Trace = true
		break
	case traceVerbose:
		serverOpts.TraceVerbose = true
		break
	case debug:
		serverOpts.Debug = true
		break
	default:
	}

	return nil, nil
}
