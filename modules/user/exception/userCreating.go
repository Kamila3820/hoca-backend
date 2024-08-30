package exception

import "fmt"

type UserCreating struct {
	UserID string
}

func (e *UserCreating) Error() string {
	return fmt.Sprintf("creating userID: %s failed", e.UserID)
}
