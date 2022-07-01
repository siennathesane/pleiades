package blaze

import (
	"fmt"
	"math/rand"
)

const (
	testServerPortStart int = 8080
	testServerPortStop  int = 8100
)

func testServerAddr() string {
	testPort := rand.Intn(testServerPortStop-testServerPortStart) + testServerPortStart
	return fmt.Sprintf("localhost:%d", testPort)
}
