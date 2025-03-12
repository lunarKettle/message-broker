package broker

type Subscriber struct {
	ID    string
	Queue string
	Ch    chan<- string
}
