package broker

type Command struct {
	ClientID string
	Action   string
	Queue    string
	Message  string
	Response chan string
}
