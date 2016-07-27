// Implements the listener capability of the program.
package metre

import (
	"github.com/pebbe/zmq4"
)

type Sink struct {
	Sink *zmq4.Socket
	// Have task map to arbitrate messages
}

func NewSink(bindURI string) (*Sink, error) {
	ctxt, err := zmq4.NewContext()
	if err != nil {
		return nil, err
	}
	sink, err := ctxt.NewSocket(zmq4.PULL)
	if err != nil {
		return nil, err
	}
	source.Bind(bindURI)
	return &Sink{
		Sink: sink,
	}, nil
}

func (s *Sink) Recv(data string) error {
	for {
		msg, err := s.Sink.Recv(0)
		if err != nil {
			fmt.Println(err.Error())
		}
		// Arbitrate the message to appropriate task channel.
	}
}

func (s *Sink) Close() {
	defer s.Sink.Close()
}
