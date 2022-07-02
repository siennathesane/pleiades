
/*
 * Copyright (c) 2022 Sienna Lloyd
 *
 * Licensed under the PolyForm Strict License 1.0.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License here:
 *  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
 */

package blaze

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"github.com/mxplusb/pleiades/pkg/utils"
	"github.com/lucas-clemente/quic-go"
	"github.com/rs/zerolog"
)

const (
	tlsCa = `-----BEGIN CERTIFICATE-----
MIIDlTCCAn2gAwIBAgIUVC411J3YD6VlJAEXkCVCgVgyYo4wDQYJKoZIhvcNAQEL
BQAwETEPMA0GA1UEAxMGcjN0LmlvMB4XDTIyMDYwNDAyMTcyNloXDTI3MDYwMzAy
MTc1NlowKzEpMCcGA1UEAxMgcjN0LmlvIEludGVybWVkaWF0ZSBBdXRob3JpdHkg
RzEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDF58VTT6IiloCbyodP
bkNFb3PgAgWvCRxYqNmMXxeCDFlz4BAhwkHKFCpn3sIVsLfnz7UJzVT9u4yjzwuf
CDdtrGjKWF1yQ6FX2wOxLU0LeSZCaSQV8qosANg/J4mx6tYWhVmLZvpnw6/fMn78
Wwb3ZQoblpX29MSaqTJIrWfuKnJTs9NL+3a+/EpW+Rl8itGtrDgPfCd361/H403H
7dbasbdoHV/uFxI+SXDjaLXt9Cc/uHtG7h0UYYso1ZPv0zZbggdFpm55KErbZHMP
PSN5By0OldK6oyiS2idxhLjBg6HTCDUO4UwXlQzmifXoclhY1rCyX56tFx2Sxi3m
WyNPAgMBAAGjgcowgccwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8w
HQYDVR0OBBYEFKHxQcEZNzZBVd2ZyYULiPwbGZnXMB8GA1UdIwQYMBaAFHOVcRxX
YXXj6jvgOp0Nv89aaSI3MDYGCCsGAQUFBwEBBCowKDAmBggrBgEFBQcwAoYaaHR0
cHM6Ly92LnIzdC5pby92MS9wa2kvY2EwLAYDVR0fBCUwIzAhoB+gHYYbaHR0cHM6
Ly92LnIzdC5pby92MS9wa2kvY3JsMA0GCSqGSIb3DQEBCwUAA4IBAQAwmTEVjaVg
OeYTGXinRGPX68yHPqdppJcO0rpY1REVKamvRSyGu9/2bVWkttJ/WZW7ZlstaVln
vh3s1NBBDvi6cik04MrnIGqlOzbKjm+9JQeKao6Gs/dSiI2V7CMAohsrJugovpN4
ODXZDzILNIp6wOfMoICg+aKgWbd1tzfi5eV/6OvW7T4LEsEUdELS5Z4M/uYI8hna
rrEqcTu6lPt+WsXcRBRr4YW0wCVtGIESjDihSVDkryAkRNFOUaQrNmWPFh2dIulI
gm6OgQ/Jkzytt+0NqfCQ1uEdrTKO4vIPvmWcDFIq/LAh9GBkV56LyID+ZIZMGZ9f
DdGYKdthyS0n
-----END CERTIFICATE-----
`
	tlsCert = `-----BEGIN CERTIFICATE-----
MIIDVjCCAj6gAwIBAgIUVBNbHLqmRSye3nQy0bzLEpKPT0owDQYJKoZIhvcNAQEL
BQAwKzEpMCcGA1UEAxMgcjN0LmlvIEludGVybWVkaWF0ZSBBdXRob3JpdHkgRzEw
HhcNMjIwNjEwMjM0NzU2WhcNMjIwNzEyMjM0ODI2WjAUMRIwEAYDVQQDEwlsb2Nh
bGhvc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDXMh3o+eQ50bXq
KTcgToFeDRDAwKu47KRQgsM69mqkontBFwDbK8AaerVchOznl/66vah2HdzEN7Hb
3y66NtrP6rNg0zcWz+iEAfOdd+JGH5vZhsn9agG/eJZE516mtJU4QWC5iGDcRo2T
q368cIHa3UDE1vXX+Xq0f64g3o1p3xN3iJxsv1KMpO3qixxefDprMhazYJ+Jc0wa
grk1sSKPlYFfL0DbnbtN+1ZJVwtlX6gcq9bibJTvww+YE9Yes2pwa3jktYr+mW3P
GRGhqEQmBIgbe0CHJy2YAiJIw33I2XFQF6fRI7ubklVvTDChQgCzirX2WyGGeq0a
6moBwLlRAgMBAAGjgYgwgYUwDgYDVR0PAQH/BAQDAgOoMB0GA1UdJQQWMBQGCCsG
AQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQU7ctSf8Q2PTWHEIXLGXSVUw+8zDkw
HwYDVR0jBBgwFoAUofFBwRk3NkFV3ZnJhQuI/BsZmdcwFAYDVR0RBA0wC4IJbG9j
YWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQAf9Ndcmaw7rwItQQO7wkbpPGKDvhDM
ZxrMzVRUUIBJTxg5BqfM1KHetZcHycmDkb850C1Pgxqp58J1YC1Va9yc9tQKk8qb
vb1kDpZqD6j5cbZStbl8JYPjIDZac7NlLvsCCHQBMFbcYN01EWj2qdnB2W0ATB/K
z5OIEHg56EV8LPAVEtutTVyW6Jj1a3g3tHuq88MVxqasuNYwKj23X1ivs8TH6TSQ
jGlxYrJeyCmqXF/Pf2nQAA9L/yZ2HFhneo6OEqjx9b4fcsLeOE4xTgGJGK/ZohRZ
CJOI/8ctCtyAjxNTSmTqFqMnbkA/NeGdKFTfEDV5PpoctxmUx0PZNJIF
-----END CERTIFICATE-----
`
	tlsKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA1zId6PnkOdG16ik3IE6BXg0QwMCruOykUILDOvZqpKJ7QRcA
2yvAGnq1XITs55f+ur2odh3cxDex298uujbaz+qzYNM3Fs/ohAHznXfiRh+b2YbJ
/WoBv3iWROdeprSVOEFguYhg3EaNk6t+vHCB2t1AxNb11/l6tH+uIN6Nad8Td4ic
bL9SjKTt6oscXnw6azIWs2CfiXNMGoK5NbEij5WBXy9A2527TftWSVcLZV+oHKvW
4myU78MPmBPWHrNqcGt45LWK/pltzxkRoahEJgSIG3tAhyctmAIiSMN9yNlxUBen
0SO7m5JVb0wwoUIAs4q19lshhnqtGupqAcC5UQIDAQABAoIBAQCymoTajTSPjG84
dparOJ0Ea0GhSQf9RmKl87GHaWdfVv+HKUlrnmclUvzdWfGp0av4X1rHFcfaDwOO
IjENGmQHNptEXGkXhN1NWrVP12U0oB1gsA9LRUVIHhRSAdm78Jr8gr59niQODnyI
uEhKq/IKraGI/YQziXr+/g2OeEfUOmu7OZchzORyqmr+qTuXqLDKIw1oNAKE0DXL
fqCj44c0zTxJEgjvVJ2odXLLI4D7VcZHKnAPRa+qFD/jMHmV5H3G2nRfUadKA+0s
o536lSVdpdI52949nRrxsG8xiwdfqL4N78eE1zxqryvyuyhvOQwzMZqoZlTgTB5J
wS7enXVBAoGBANqRHqQ2zm44Y0naWH2cb34tnFzW8IUhYqX4+4k/vX212Ey+bQjf
NIijnFpTcjNqbZmE7RZU8wwtgfKkwJzEMlC3gcNAQFE9CvtveKpyHyP38RyUP/a8
ZFTxiZhWG2kgu4ToxyyZZGx8BjsD+4LGTrRe56vrkMNgYHu3Os5UKRT5AoGBAPwN
MaLM2+7IyPQ5SzEIDX7jOcqoR1UIBpVnzt3UgbMNnKaAAAAwWCyyhHhfr1liymbm
n1vb0ykAb9bZmQSB4hRHra6TB5gysday9yBpzSdNR2wwTDBT4g1pNZz3htgHeXh/
l1qrYwmDwBaNCWE3vcFSx1rtho4uJ7d0TWeAv1UZAoGBAJGFIXOQEe0Mmf6n41bu
esT1tS+S49yfl8CNf1uoFo/GLNcbyhioE6AN3qG9AUH+UC5wdDH0KUYoXmahDqTR
c/aN11WaR7hO/ird0ucYyGb4Q44Vnmi2kc6EamoEmodqBa++FC47isM36CYOxrwR
MIGi1nh+hImwd0ynd/27xwZxAoGBAKk9nreGwKJ8FVrPYaqhkpZBspteFM+GnQ0S
7/dJanE00ZuG1PlLfNk+YO6GqTHmwKsJbbV7TDT6wx3LbBB3ubsOShOvS+kpGPpl
nsQX5pXeMPf3EiFdIasJmuMz3UoO8sQzQAi0jcJkwcUinEq35+T4VT27wZ6UZTys
jhDShSZ5AoGAYvLIHoNpx/fWruvTsV+UH8GkyhToDjjlWYfH0hgMG3fPVFk5RkN9
Q9whW/iQa5xkLOr5+nmCRyYUvnTyJQ5lMNN3MlU+zh2Rpl4M1WLv29nD173CYw1g
BLlzHCDBvh3zBlkMWivluxbG5XZzQPQN+Y+WaC+ldJqE9oQeaOw3GDI=
-----END RSA PRIVATE KEY-----
`
)

type QuicTestKit struct {
	t              *testing.T
	logger         zerolog.Logger
	testServerAddr string

	quicConfig *quic.Config
	listener   quic.Listener
	dialConn   quic.Connection

	certPool *x509.CertPool
	keyPair  tls.Certificate
	tlsConf  *tls.Config

	ctx context.Context
}

func NewQuicTestKit(t *testing.T) *QuicTestKit {
	tk := &QuicTestKit{
		logger:         utils.NewTestLogger(t),
		ctx:            context.Background(),
		t:              t,
		testServerAddr: testServerAddr(),
	}

	var err error
	tk.tlsConf, err = tk.GenerateTlsConfig()
	if err != nil {
		t.Fatalf("failed to generate tls host: %v", err)
	}

	tk.quicConfig = &quic.Config{MaxIdleTimeout: 300 * time.Second}

	tk.listener, err = quic.ListenAddr(tk.testServerAddr, tk.tlsConf, tk.quicConfig)
	if err != nil {
		t.Fatalf("failed to listen server: %v", err)
	}

	tk.dialConn, err = quic.DialAddr(tk.testServerAddr, tk.tlsConf, tk.quicConfig)
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}

	return tk
}

func (stk *QuicTestKit) Start() {
	var err error
	stk.listener, err = quic.ListenAddr(testServerAddr(), stk.tlsConf, stk.quicConfig)
	if err != nil {
		stk.t.Fatalf("failed to listen server: %v", err)
	}
}

func (stk *QuicTestKit) Stop() {
	stk.listener.Close()
	//if err != nil {
	//	stk.t.Fatalf("failed to close listener: %v", err)
	//}
}

func (stk *QuicTestKit) GenerateTlsConfig() (*tls.Config, error) {
	stk.certPool = x509.NewCertPool()
	stk.certPool.AppendCertsFromPEM([]byte(tlsCa))

	var err error
	stk.keyPair, err = tls.X509KeyPair([]byte(tlsCert), []byte(tlsKey))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tls key pair: %v", err)
	}

	stk.tlsConf = &tls.Config{
		Certificates: []tls.Certificate{stk.keyPair},
		RootCAs:      stk.certPool,
		NextProtos:   []string{"blaze-test-server"},
	}

	return stk.tlsConf, nil
}

func (stk *QuicTestKit) GetListener() quic.Listener {
	return stk.listener
}

func (stk *QuicTestKit) CloseListener() {
	err := stk.listener.Close()
	if err != nil {
		stk.t.Error(err, "there was an error trying to close the listener")
	}
}

func (stk *QuicTestKit) GetConnection() quic.Connection {
	return stk.dialConn
}

func (stk *QuicTestKit) NewConnectionStream() quic.Stream {
	stream, err := stk.dialConn.OpenStream()
	if err != nil {
		stk.t.Error(err, "there was an error opening a new testkit stream")
	}

	return stream
}
