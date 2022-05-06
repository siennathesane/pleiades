package fsm

import (
	"encoding/json"
	"fmt"
)

type Item struct {
	Key   []byte
	Value []byte
}

func (i Item) String() string {
	return fmt.Sprintf("{%s: %s}", string(i.Key), json.RawMessage(i.Value))
}
