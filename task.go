package metre

import (
	log "github.com/Sirupsen/logrus"

	"time"
)

type Task struct {
	TimeOut       time.Duration
	MessageCount  int
	ScheduleCount int
	ScheduleDone  bool
	ID            string // Type Type of task (user as class prefix in cache)
	Interval      string // Schedule String in cron notation
	Schedule      func(t TaskRecord, s Scheduler, c Cache, q Queue)
	Process       func(t TaskRecord, s Scheduler, c Cache, q Queue)
}

func (t Task) GetID() string {
	return t.ID
}

func (t Task) GetInterval() string {
	return t.Interval
}

func (t *Task) checkComplete() bool {
	if t.MessageCount == t.ScheduleCount && t.ScheduleDone {
		return true
	}
	return false
}

func (t *Task) TestTimeOut() {
	go func() {
		if t.TimeOut != 0 {
			time.Sleep(t.TimeOut)
			if !t.checkComplete() {
				log.Info("Task hit timeout at time ", time.Now())
			}
		}
	}()
}

func (t *Task) Track(trackMsg *trackMessage) {
	switch trackMsg.MessageType {
	case Status:
		t.MessageCount++
		if t.checkComplete() {
			log.Info(trackMsg.TaskId + ": Complete")
			log.Info(trackMsg.TaskId, ": Completed at ", time.Now())
		}
	case Debug:
		log.Debug(trackMsg.TaskId + ":" + trackMsg.Message)
	case Error:
		log.Warn(trackMsg.TaskId + ":" + trackMsg.Message)
	default:
		log.Warn("Unknown message type")
	}
}
