// Package metre is used to schedule end execute cron jobs in a simplified fashion
package metre

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

const LOCALHOST string = "127.0.0.1" // Default host for cache and queue
const QUEUEPORT string = "5555"      // Default port for queue
const TRACKQUEUEPORT string = "5556" // Default port for queue
const CACHEPORT string = "6379"      // Default port for cache

type Metre struct {
	Cron       cron.Cron
	Queue      Queue
	TrackQueue Queue
	Cache      Cache
	Scheduler  Scheduler
	TaskMap    map[string]Task
	// Add func(t Task)
}

// New creates a new scheduler to manage task scheduling and states
func New(queueUri string, trackQueueUri string, cacheUri string) (Metre, error) {
	if cacheUri == "" {
		cacheUri = LOCALHOST + ":" + CACHEPORT
	} else if strings.Index(cacheUri, ":") == 0 {
		cacheUri = LOCALHOST + ":" + cacheUri
	}

	if queueUri == "" {
		queueUri = LOCALHOST + ":" + QUEUEPORT
	} else if strings.Index(queueUri, ":") == 0 {
		queueUri = LOCALHOST + ":" + queueUri
	}

	if trackQueueUri == "" {
		trackQueueUri = LOCALHOST + ":" + TRACKQUEUEPORT
	} else if strings.Index(trackQueueUri, ":") == 0 {
		trackQueueUri = LOCALHOST + ":" + trackQueueUri
	}

	cron := *cron.New()
	c, cErr := NewCache(cacheUri)
	if cErr != nil {
		return Metre{}, cErr
	}
	q, qErr := NewQueue(queueUri)
	if qErr != nil {
		return Metre{}, qErr
	}

	t, tErr := NewQueue(trackQueueUri)
	if qErr != nil {
		return Metre{}, tErr
	}

	err := t.BindPush()
	if err != nil {
		return Metre{}, err
	}
	s := NewScheduler(q, c)
	m := make(map[string]Task)
	return Metre{cron, q, t, c, s, m}, nil
}

// Add adds a cron job task to schedule and process
func (m *Metre) Add(t Task) {
	id := t.GetID()
	if _, exists := m.TaskMap[id]; exists {
		panic("attempted to add two tasks with the same ID [" + t.ID + "]")
	}

	m.TaskMap[id] = t
	m.Cron.AddFunc(t.Interval, func() {
		t.Schedule(NewTaskRecord(id), m.Scheduler, m.Cache, m.Queue)
	})
}

// Schedule schedules a singular cron task
func (m *Metre) Schedule(ID string) (string, error) {
	e := m.Queue.BindPush()
	if e != nil {
		return "", nil
	}
	t, ok := m.TaskMap[ID]
	if ok == false {
		return "", errors.New("task [" + ID + "] not recognized")
	}

	tr := NewTaskRecord(t.GetID())
	t.Schedule(tr, m.Scheduler, m.Cache, m.Queue)
	return buildTaskKey(tr), nil
}

// Scheduler processes a singular cron task
func (m *Metre) Process(ID string) (string, error) {
	t, ok := m.TaskMap[ID]
	if ok == false {
		return "", errors.New("task [" + ID + "] not recognized")
	}

	tr := NewTaskRecord(t.GetID())
	t.Process(tr, m.Scheduler, m.Cache, m.Queue)
	return buildTaskKey(tr), nil
}

func (m *Metre) StartMaster() {
	e := m.Queue.BindPush()
	if e != nil {
		panic(e)
	}
	m.Cron.Start()
}

// This function tracks if the schedules get completed.
func (m *Metre) Track() {
	e := m.TrackQueue.ConnectPull()
	if e != nil {
		panic(e)
	}
	log.Info("Waiting for messages from track queue")
	for {
		msg := m.TrackQueue.Pop()
		// Handle different types of messages
		// FIXME log the message in temporarily
		log.Info(msg)
	}
}

func (m *Metre) runAndSendComplete(tr TaskRecord) {
	tsk := m.TaskMap[tr.ID]
	tsk.Process(tr, m.Scheduler, m.Cache, m.Queue)
	log.Info("Sending Completed")
	_, err := m.TrackQueue.Push("Status:Completed")
	if err != nil {
		log.Warn("Error while pushing completed status for a process")
	}

}

func (m *Metre) StartSlave() {
	e := m.Queue.ConnectPull()
	if e != nil {
		panic(e)
	}
	for {
		msg := m.Queue.Pop()
		tr, _ := ParseTask(msg)
		if tr.ID == "" || tr.UID == "" {
			log.Warn("Failed to parse task from message: " + msg)
			continue
		}

		m.Cache.Delete(buildTaskKey(tr))
		go m.runAndSendComplete(tr)
	}
}
