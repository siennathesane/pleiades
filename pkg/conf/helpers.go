package conf

import (
	"net"
	"time"
)

// portChecker just verifies that specific ports are open, used for tests
// ref: https://stackoverflow.com/a/56336811/4949938
func portChecker(host string, ports ...string) error {
	for _, port := range ports {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			return err
		}
		if conn != nil {
			err := conn.Close()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
