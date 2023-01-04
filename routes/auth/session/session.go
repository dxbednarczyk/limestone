package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"limestone/util"
	"log"
	"net/http"
)

type Session struct {
	Client      *http.Client
	login       loginDetails
	Id          string
	UserId      string
	Token       string
	DisplayName string
}

type loginDetails struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	SessionName string `json:"friendly_name"`
}

type DefaultError struct {
	Error string `json:"type"`
}

func NewSession(email string, password string, sessionName string) Session {
	return Session{
		Client: &http.Client{},
		login: loginDetails{
			email,
			password,
			sessionName,
		},
	}
}

func (sesh *Session) Login() error {
	body, err := json.Marshal(sesh.login)
	if err != nil {
		return err
	}

	resp, err := sesh.Client.Post(
		"https://api.divolt.xyz/auth/session/login",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		log.Fatal("Failed to login.")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed validation")
	}

	defer resp.Body.Close()

	result := struct {
		UniqueId     string `json:"_id"`
		UserId       string `json:"user_id"`
		SessionToken string `json:"token"`
		DisplayName  string `json:"name"`
	}{}
	err = util.UnmarshalResponseBody(resp, &result)
	if err != nil {
		return err
	}

	sesh.Id = result.UniqueId
	sesh.UserId = result.UserId
	sesh.Token = result.SessionToken
	sesh.DisplayName = result.DisplayName

	return nil
}

func (sesh *Session) Logout() error {
	req, err := util.RequestWithSessionToken(
		http.MethodPost,
		"auth/session/logout",
		nil,
		sesh.Token,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := sesh.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var derr DefaultError
		err = util.UnmarshalResponseBody(resp, &derr)
		if err != nil {
			return err
		}

		return errors.New(derr.Error)
	}

	return nil
}
