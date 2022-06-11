package blaze

import (
	"fmt"
	"reflect"

	"github.com/zeebo/errs"
	"storj.io/drpc"
)

var (
	_           drpc.Mux     = (*Router)(nil)
	_           drpc.Handler = (*Router)(nil)
	streamType               = reflect.TypeOf((*drpc.Stream)(nil)).Elem()
	messageType              = reflect.TypeOf((*drpc.Message)(nil)).Elem()
)

type Router struct {
	drpc.Mux
	drpc.Handler
	targets map[string]RpcMuxDataModel
}

type RpcMuxDataModel struct {
	Name     string
	Server   interface{}
	Encoding drpc.Encoding
	Receiver drpc.Receiver
	In1      reflect.Type
	In2      reflect.Type
	Unitary  bool
}

func NewRouter() *Router {
	mux := &Router{}
	mux.targets = make(map[string]RpcMuxDataModel)

	return mux
}

func (d *Router) Register(srv interface{}, desc drpc.Description) error {
	n := desc.NumMethods()
	for i := 0; i < n; i++ {
		rpc, enc, receiver, method, ok := desc.Method(i)
		if !ok {
			return fmt.Errorf("description returned invalid method for index %d", i)
		}
		if err := d.registerOne(srv, rpc, enc, receiver, method); err != nil {
			return err
		}
	}
	return nil
}

func (d *Router) registerOne(srv interface{}, rpc string, enc drpc.Encoding, receiver drpc.Receiver, method interface{}) error {
	data := RpcMuxDataModel{Server: srv, Encoding: enc, Receiver: receiver, Name: rpc}

	switch mt := reflect.TypeOf(method); {
	// unitary input, unitary output
	case mt.NumOut() == 2:
		data.Unitary = true
		data.In1 = mt.In(2)
		if !data.In1.Implements(messageType) {
			return fmt.Errorf("input argument not a drpc message: %v", data.In1)
		}

	// unitary input, stream output
	case mt.NumIn() == 3:
		data.In1 = mt.In(1)
		if !data.In1.Implements(messageType) {
			return fmt.Errorf("input argument not a drpc message: %v", data.In1)
		}
		data.In2 = streamType

	// stream input
	case mt.NumIn() == 2:
		data.In1 = streamType

	// code gen bug?
	default:
		return fmt.Errorf("unknown method type: %v", mt)
	}

	d.targets[rpc] = data

	return nil
}

func (d *Router) HandleRPC(stream drpc.Stream, rpc string) (err error) {

	data, ok := d.targets[rpc]
	if !ok {
		return fmt.Errorf("error finding implementer for %s", rpc)
	}

	in := interface{}(stream)
	if data.In1 != streamType {
		msg, ok := reflect.New(data.In1.Elem()).Interface().(drpc.Message)
		if !ok {
			return drpc.InternalError.New("invalid rpc input type")
		}
		if err := stream.MsgRecv(msg, data.Encoding); err != nil {
			return errs.Wrap(err)
		}
		in = msg
	}

	out, err := data.Receiver(data.Server, stream.Context(), in, stream)
	switch {
	case err != nil:
		return errs.Wrap(err)
	case out != nil && !reflect.ValueOf(out).IsNil():
		return stream.MsgSend(out, data.Encoding)
	default:
		return stream.CloseSend()
	}
}
