package box

type BoxTask interface {
	Work(chan string)
}
