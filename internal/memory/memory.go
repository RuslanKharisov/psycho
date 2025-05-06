package memory

type ChatMessage struct {
	Role    string
	Content string
}

type Memory interface {
	Append(userID int64, msg ChatMessage) error
	Get(userID int64) ([]ChatMessage, error)
	Truncate(userID int64, limit int) error
}
