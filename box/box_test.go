package box

import (
	"fmt"
	"testing"
)

type TestWork struct {
}

func (t *TestWork) Work(updates chan string) {
	for i := 0; i < 100000; i++ {
		if i%100 == 0 {
			updates <- fmt.Sprintf("i = %d", i)
		}
	}
}

func ListenToUpdates(updates chan string) {
	for {
		update := <-updates
		fmt.Println(update)
	}
}

func TestBox(t *testing.T) {
	tw := new(TestWork)
	b := New()
	go b.Run(tw)

	go ListenToUpdates(b.RandomChan)
	b.StartChan <- Start

	message := <-b.EndChan
	if message == Done {
		t.Log("Succesfully completed task")
	}
}
