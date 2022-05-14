package conf

import (
	"context"
	"encoding/json"
	"fmt"

	clog "cloud.google.com/go/logging"
	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog"
	"google.golang.org/api/option"
)

const (
	googleCloudCredentialsConfigPath string = "hosts/%s/config/google"
)

type gcpLogger struct {
	env              *EnvironmentConfig
	clogger          *clog.Logger
	cloudClient      *clog.Client
	initialWriteSent bool
}

// DefaultSeverityMap contains the default zerolog.Level -> clog.Severity mappings.
var defaultSeverityMap = map[zerolog.Level]clog.Severity{
	zerolog.DebugLevel: clog.Debug,
	zerolog.InfoLevel:  clog.Info,
	zerolog.WarnLevel:  clog.Warning,
	zerolog.ErrorLevel: clog.Error,
	zerolog.PanicLevel: clog.Critical,
	zerolog.FatalLevel: clog.Critical,
}

func newGcpLogger(l Logger) (gcpLogger, error) {
	// Sets your Google Cloud Platform project ID.
	projectID := l.env.GCPProjectId

	// grab the google auth from consul
	pair, _, err := l.client.KV().Get(fmt.Sprintf(googleCloudCredentialsConfigPath, l.env.Hostname), &api.QueryOptions{})
	if err != nil {
		return gcpLogger{}, err
	}

	// Creates a client.
	cloudClient, err := clog.NewClient(context.TODO(), projectID, option.WithCredentialsJSON(pair.Value))
	if err != nil {
		return gcpLogger{}, err
	}

	clogger := cloudClient.Logger(fmt.Sprintf("%s", l.env.Hostname), clog.CommonLabels(map[string]string{
		"hostname": l.env.Hostname,
		"env":      string(l.env.Environment),
	}))

	return gcpLogger{clogger: clogger, cloudClient: cloudClient}, nil
}

func (l gcpLogger) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	entry := clog.Entry{
		Payload:  json.RawMessage(p),
		Severity: defaultSeverityMap[level],
	}

	if !l.initialWriteSent {
		if err := l.clogger.LogSync(context.TODO(), entry); err != nil {
			return 0, err
		}
		l.initialWriteSent = true
	} else {
		l.clogger.Log(entry)
	}

	if level == zerolog.FatalLevel {
		if err := l.clogger.Flush(); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func (l gcpLogger) Write(p []byte) (n int, err error) {
	entry := clog.Entry{Payload: json.RawMessage(p)}

	if !l.initialWriteSent {
		if err := l.clogger.LogSync(context.TODO(), entry); err != nil {
			return 0, err
		}
		l.initialWriteSent = true
	} else {
		l.clogger.Log(entry)
	}
	return len(p), nil
}

func (l gcpLogger) Close() error {
	if err := l.clogger.Flush(); err != nil {
		return err
	}
	return l.cloudClient.Close()
}
