package requestsystem

import "time"

type Request struct {
	ID        int
	Client    *Client
	Status    string
	CreatedAt time.Time // Время создания заявки
}

// getId returns the ID of the request.
func (r *Request) GetId() int {
	return r.ID
}

// updateStatus updates the status of the request.
func (r *Request) UpdateStatus(status string) {
	r.Status = status
}
