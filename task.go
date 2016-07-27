package metre

type Task struct {
	// All the messages and synchronization info happens here.
	MessageChan chan string

	// Type Type of task (user as class prefix in cache).
	ID string

	// Schedule String in cron notation.
	Interval string

	// Schedules a task for execution by slave.
	Schedule func(t TaskRecord, s Scheduler, c Cache, q Queue)

	// The process action where slave picks up the part that is given to it.
	Process func(t TaskRecord, s Scheduler, c Cache, q Queue)

	// Track the completion and sends updates for a task.
	Track func() chan string
}
