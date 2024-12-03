package models

const (
	INTERNALL_SERVER_ERROR = "INTERNAL SERVER ERROR"
	BAD_SYNTAX             = "REQUEST CONTAINS BAD SYNTAX OR CANNOT BE FULLFILLED"
)

type OasaError struct {
	Error string `json:"error" `
}
