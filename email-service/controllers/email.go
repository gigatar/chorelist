package controllers

import (
	"chorelist/email-service/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// EmailController data type.
type EmailController struct {
	ServerConfig models.EmailConfig
}

// Initialize the email manager
func (e *EmailController) Initialize(configPath string) error {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(data, &e.ServerConfig)
	if err != nil {
		return err
	}
	return nil
}

// SignupEmail Creates and sends an email using the signup email template.
func (e *EmailController) SignupEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var emailRequest models.SignupEmail

	err := json.NewDecoder(r.Body).Decode(&emailRequest)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !emailRequest.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	emailRequest.Code = e.ServerConfig.BaseLink + "?code=" + emailRequest.Code

	var request RequestController
	request.Request = models.Request{
		From:    e.ServerConfig.From,
		To:      []string{emailRequest.Email},
		Subject: "Email Confirm",
	}

	if err := request.parseTemplate("templates/signupEmail.tmpl", emailRequest); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := request.sendEmail(e.ServerConfig); err != nil {
		if strings.Contains(err.Error(), "550 5.1.1") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
