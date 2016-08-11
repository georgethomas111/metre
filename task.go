package metre

type Task struct {
	ID       string // Type Type of task (user as class prefix in cache)
	Interval string // Schedule String in cron notation
	Schedule func(t TaskRecord, s Scheduler, c Cache, q Queue)
	Process  func(t TaskRecord, s Scheduler, c Cache, q Queue)
}

func (t Task) GetID() string {
	return t.ID
}

func (t Task) GetInterval() string {
	return t.Interval
}

type TaskInt interface {
	GetID() string
	GetInterval() string
}
