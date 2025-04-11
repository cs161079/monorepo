package models

type Notification struct {
	AlertType string `json:"alertType"`
	Title     string `json:"title"`
	Message   string `json:"message"`
}
