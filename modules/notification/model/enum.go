package model

import (
	"encoding/json"
	"fmt"
)

type NotificationEnum string

const (
	NotificationPlaceOrder   NotificationEnum = "confirmation"
	NotificationPreparing    NotificationEnum = "preparing"
	NotificationWorking      NotificationEnum = "working"
	NotificationComplete     NotificationEnum = "complete"
	NotificationUserCancel   NotificationEnum = "user_cancel"
	NotificationWorkerCancel NotificationEnum = "worker_cancel"
	NotificationRating       NotificationEnum = "user_rating"
	NotificationSystemCancel NotificationEnum = "system_cancel"
)

func (n *NotificationEnum) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	Notification := NotificationEnum(val)
	if Notification != NotificationPlaceOrder && Notification != NotificationPreparing && Notification != NotificationWorking && Notification != NotificationComplete && Notification != NotificationUserCancel && Notification != NotificationWorkerCancel && Notification != NotificationRating {
		return fmt.Errorf("invalid Application enum value: %s", Notification)
	}

	*n = Notification

	return nil
}
