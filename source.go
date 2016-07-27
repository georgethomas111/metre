// Implements the listener capability of the program.
package metre

import (
	"github.com/pebbe/zmq4"
)

type Source struct {
	Source *zmq4.Socket
}

func NewSource(uri string) (*Listener, error) {
	ctxt, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	source, err := ctxt.NewSocket(zmq4.PUSH)
	if err != nil {
		return nil, err
	}
	source.Connect(uri)
	return &Source{
		Source: source,
	}
}

func (s *Source) Send(data string) error {
	_, err := s.Send(data, zmq4.DONTWAIT)
	return err
}

func (s *Source) Close() {
	defer s.Source.Close()
}
