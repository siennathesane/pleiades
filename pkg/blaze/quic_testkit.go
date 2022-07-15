
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
MIIDVjCCAj6gAwIBAgIUfkYIfDYgS3yUHW1qX0AQ5fC760gwDQYJKoZIhvcNAQEL
BQAwKzEpMCcGA1UEAxMgcjN0LmlvIEludGVybWVkaWF0ZSBBdXRob3JpdHkgRzEw
HhcNMjIwNzE1MTE1NzE3WhcNMjQwNzE0MTE1NzQ2WjAUMRIwEAYDVQQDEwlsb2Nh
bGhvc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCvwlG7bTAuq+zU
8wHuuLreeIr8wKCt+rcevMeTDcHpaod66Pdw0bwnnnAecHAXy16ZOWr15ix73i8k
K8Urd8QzFCWcrdCplI/7uzQ3O+dVHzfj/d1Xby8UsIISu2UX3pt+4CK8nkfPnH/i
P7eXVnOQFEk/pQPwXrDzW+VuST31/yx1TwpgxYEZ2b1shwPTTF78P/pCnK24Ab2r
kJx1TcPoMyLSBs0+gtIU/nNKq3rdGuJP+Zp+2nJoUiwYGLqgJCHFfMIoSFevedKA
02QVwYud2AqMZ558mc4EbnOK6TVyge8aRgFz+lw4wy7lZyo8UzMeZHIXTcqbvNJc
FCtNWviTAgMBAAGjgYgwgYUwDgYDVR0PAQH/BAQDAgOoMB0GA1UdJQQWMBQGCCsG
AQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUVe+NPLNa8EnkGR9t6NjiMPoU1X8w
HwYDVR0jBBgwFoAUofFBwRk3NkFV3ZnJhQuI/BsZmdcwFAYDVR0RBA0wC4IJbG9j
YWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQC7uBeDvkFB6o+oIEvjOHDDppXObqG6
h1v55krg5soDjgrU9C0lTWv3J6uFgpBSGAn7ELpIqqJRlC/Naino6V0ntlnoZOIX
EnNUt/1uQRdsvcgMElqKiUsPdBKwZiZtqLNVEEnjV1lJg4DdxFQkx87F0CBiPlv7
Aaeci8SEcf+wQbyMq2f/fbSOer8RY0w4BYU45VJcPYyESTMF43rIzWzBZjH78j9o
C2+g6AXRPryP3d1Q3cuCxvgs6EV6v2dHk9CbgK873gPy0rEX83abAHOFg5VtgPeZ
oqCcKFv/lVR9cd4hSIgPmYCOR0Jpb5ftto8Cc3wztRSO3fZ9dZtDL3Jv
-----END CERTIFICATE-----
`
	tlsKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAr8JRu20wLqvs1PMB7ri63niK/MCgrfq3HrzHkw3B6WqHeuj3
cNG8J55wHnBwF8temTlq9eYse94vJCvFK3fEMxQlnK3QqZSP+7s0NzvnVR834/3d
V28vFLCCErtlF96bfuAivJ5Hz5x/4j+3l1ZzkBRJP6UD8F6w81vlbkk99f8sdU8K
YMWBGdm9bIcD00xe/D/6QpytuAG9q5CcdU3D6DMi0gbNPoLSFP5zSqt63RriT/ma
ftpyaFIsGBi6oCQhxXzCKEhXr3nSgNNkFcGLndgKjGeefJnOBG5ziuk1coHvGkYB
c/pcOMMu5WcqPFMzHmRyF03Km7zSXBQrTVr4kwIDAQABAoIBAAw+x579TnwR8mAr
xhixuiNB1r0hVpCIhTWZBXaTYM04ZNQFYyfzFN7VxQ3523Vs79xRCci0Dqlao4Ir
ooMSxxKf9rbhzUXjWdy6ADtQ6x8dG7HAwCkVT/xrB8TyVWrCUacXSpRFms82IhrL
QalDlq2VHD1Y69tVXZK29lcmtzBqRz+60bbLiuY7afwjETRC9tPoFDpMIzd13yuM
OlW8nW7+EgMswgZ/pN2xL+w1w2lCgvmjbVmjTPuXYerR5MzEZCTc95/HsHVezDDC
pIAkRPZQIiInNM6CyOthwftXgvJv/36SjtHiOBi1Q+MXqicqzquYX9cEBIMAM7lJ
QJWaJdECgYEA33cyltfvJXHez11DfeyLRhpiN2mwxABcmXwk2+QZk3Qvd/PZLG/1
uJwM3LZ4LZGA9Z03FnXGU+JchTia88lQaT0RVbRzUO8FeH4uDc9E+RlKDDcks3Vn
Cv2E/r9U61Zxa/kjFYlk7OquI7Nr3Ay5Yg+NBmoCgE1lt3F6ITq00ScCgYEAyVkM
HG2+J4biIuZGsyDxJCiu/nsoYUYa9TvmxNSK6iruav6KoaIqO9KI/VOhGf471Orq
feupyD4sgzQM4UP6Ww4+iwXeNTpF10mxClE1Tk7kwnAZJpYw4zMqpCp46a2Pc2YY
akQB+xVEN+qT9gCs42zc1a1NVcxzNb6z8wmjKLUCgYBMTpy0y7m69J5b+wHv/xUz
9BBz0aBt3Z3BP5YqjEJ7iqIm+NrBBN5IkukFeT2iedwqguvrvH3j6Rkk2MZ41tah
iRvhQ0RZb7VThurdBlkMIqmZcD8VFNMB+r4ua1FpJ1SFxUZItWkEScL7J+p98s5f
AOZsOUjvXP6N3K8Sp8RU4wKBgQCRMiHFrm3d2yrft+dr7Wl3hc8LvIxV+VQfXF8B
ubOjQepEReJ6xJJoKV6YL+KQ+AD1faIzw+nfeNZolvRizb6QQyle35BqGeebZIzC
v+UM31+fx26boNsIPDGXyPkAqiQ0N3+LwhcblS5olES2ta33It3tSNfn81NxgmAJ
9v0tsQKBgDccKQ9I+A9eKlGUFXFPiMHFs9H1VadXlSeajXUkOm90uqXhJUcqByWj
lSCS7P1fg3uFDvWSR54Du0ANSyfuJf7E6H9WhpcL01A2CEWJP0hF7qT2BD+RVaMo
ZYuLm7cH9HKYWharqyPHhQf+MtZSWwgPyIH7YuHwjnMe0RZAa81M
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
