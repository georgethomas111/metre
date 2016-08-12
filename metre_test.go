package metre

import (
	"testing"

	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var test = &Task{
	TimeOut:  time.Second * 5,
	ID:       "Test",
	Interval: "0 * * * * *",
	Schedule: func(t TaskRecord, s Scheduler, c Cache, q Queue) {
		for i := 0; i < 10; i++ {
			t.UID = fmt.Sprintf("%d", i)
			log.Info("Scheduling test " + t.UID)
			time.Sleep(time.Second)
			s.Schedule(t)
		}
		return
	},
	Process: func(t TaskRecord, s Scheduler, c Cache, q Queue) {
		log.Info("Processing Test  " + t.UID)
		return
	},
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func printMsgs(msgChan chan string) {
	go func() {
		for {
			msg := <-msgChan
			log.Info(msg)
		}
	}()
}

func TestLife(t *testing.T) {
	var wg sync.WaitGroup
	met, err := New("", "", "")
	if err != nil {
		t.Errorf("Metre creation error" + err.Error())
	}

	printMsgs(met.MessageChannel)

	go met.Track()
	go met.StartSlave()

	met.Add(test)
	go test.TestTimeOut()
	met.Schedule(test.ID)
	wg.Add(1)
	wg.Wait()
}
