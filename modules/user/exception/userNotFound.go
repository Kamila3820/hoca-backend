package exception

import "fmt"

type UserNotFound struct {
	UserID string
}

func (e *UserNotFound) Error() string {
	return fmt.Sprintf("userID: %s not found", e.UserID)
}
