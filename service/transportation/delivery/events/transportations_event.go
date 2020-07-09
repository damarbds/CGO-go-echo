package events

import "github.com/models"

var UserCreated userCreated

// UserCreatedPayload is the data for when a user is created
type UserCreatedPayload struct {
	Schedule models.NewCommandSchedule
}

type userCreated struct {
	handlers []interface{ Handle(UserCreatedPayload) }
}

// Register adds an event handler for this event
func (u *userCreated) Register(handler interface{ Handle(UserCreatedPayload) }) {
	u.handlers = append(u.handlers, handler)
}

// Trigger sends out an event with the payload
func (u userCreated) Trigger(payload UserCreatedPayload) {
	for _, handler := range u.handlers {
		go handler.Handle(payload)
	}
}
