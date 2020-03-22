package controllers

import (
	"bytes"
	"chorelist/email-service/models"
	"html/template"
	"net/smtp"
	"strconv"
)

// RequestController data type.
type RequestController struct {
	Request models.Request
}

// SendEmail sends our email.
func (r *RequestController) sendEmail(config models.EmailConfig) error {
	auth := smtp.PlainAuth("", config.User, config.Password, config.Server)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.Request.Subject + "\n"
	msg := []byte("from: " + r.Request.From + "\n" + subject + mime + "\n" + r.Request.Body)
	addr := config.Server + ":" + strconv.Itoa(config.ServerPort)

	if err := smtp.SendMail(addr, auth, r.Request.From, r.Request.To, msg); err != nil {
		return err
	}

	return nil
}

// ParseTemplate parses our email template and sets the proper email body.
func (r *RequestController) parseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.Request.Body = buf.String()
	return nil
}
