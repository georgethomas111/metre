// A box is a container to run the tasks.
// It can be either a master slave or
package box

type Message int

const (
	Start Message = iota
	Done
)

type Box struct {
	StartChan  chan Message
	EndChan    chan Message
	RandomChan chan string
}

func New() *Box {
	return &Box{
		StartChan:  make(chan Message),
		EndChan:    make(chan Message),
		RandomChan: make(chan string),
	}
}

func (b *Box) Run(task BoxTask) {
	for {
		// Wait for the start signal
		<-b.StartChan
		task.Work(b.RandomChan)
		// Convey that work is done and wait for more
		b.EndChan <- Done
	}
}
