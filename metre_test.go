package metre

import (
	"github.com/hart/fred/tasks"

	log "github.com/Sirupsen/logrus"
)

var met Metre

func viewUpdate(updateChan chan string) {
	for {
		update := <-updateChan
		if update == TaskCompletd {
			log.Info("Job completed. Time the job to calculate total run time.")
			break
		}
		log.Info("Received update :" + update)
	}
}

func TestMain(m *testing.M) {
	m.Run()
	met = metre.New()
}

func TestLife(t *testing.T) {
	var wg sync.WaitGroup
	met.Schedule(tasks.Test)
	wg.Add(1)
	go viewUpdate(tasks.MessageChan)
	met.Process()
	wg.Wait()
}
