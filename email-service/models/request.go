package models

// Request struct for emails
type Request struct {
	From    string
	To      []string
	Subject string
	Body    string
}
