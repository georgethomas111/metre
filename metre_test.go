package metre

import (
	"testing"

	"fmt"
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

var test = Task{
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
	m.Run()
}

func TestLife(t *testing.T) {
	var wg sync.WaitGroup
	met, err := New("", "", "")
	if err != nil {
		t.Errorf("Metre creation error" + err.Error())
	}

	go met.Track()
	go met.StartSlave()
	met.Add(test)
	test.TestTimeOut()
	met.TaskMap[test.ID] = &test
	met.Schedule(test.ID)
	log.Info("Waiting ...")
	wg.Add(1)
	wg.Wait()
}
