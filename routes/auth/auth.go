package auth

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

type authResult struct {
	UniqueId     string `json:"_id"`
	UserId       string `json:"user_id"`
	SessionToken string `json:"token"`
	DisplayName  string `json:"name"`
}

type loginDetails struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	SessionName string `json:"friendly_name"`
}

type Error struct {
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
		return errors.New("failed validation")
	}

	defer resp.Body.Close()

	var ar authResult
	err = util.UnmarshalResponseBody(resp, &ar)
	if err != nil {
		return err
	}

	sesh.Id = ar.UniqueId
	sesh.UserId = ar.UserId
	sesh.DisplayName = ar.DisplayName
	sesh.Token = ar.SessionToken

	util.SessionToken = ar.SessionToken

	return nil
}

func (sesh *Session) Logout() error {
	req, err := util.AuthenticatedRequest(
		http.MethodPost,
		"auth/session/logout",
		nil,
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
		return GetError(resp)
	}

	return nil
}

func GetError(resp *http.Response) error {
	var autherr Error
	err := util.UnmarshalResponseBody(resp, &autherr)
	if err != nil {
		return err
	}

	return errors.New(autherr.Error)
}
