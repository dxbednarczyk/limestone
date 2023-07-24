package divolt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dxbednarczyk/limestone/util"
)

type Session struct {
	Client         *http.Client
	Authentication *Authentication
	loginDetails   loginDetails
}

type Authentication struct {
	ID          string `json:"_id"`
	UserID      string `json:"user_id"`
	Token       string `json:"token"`
	DisplayName string `json:"name"`
}

type loginDetails struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	SessionName string `json:"friendly_name"`
}

type AuthenticationError struct {
	Error string `json:"type"`
}

func NewSession(email string, password string, sessionName string) Session {
	return Session{
		Client: http.DefaultClient,
		loginDetails: loginDetails{
			email,
			password,
			sessionName,
		},
	}
}

func (sesh *Session) Login() error {
	body, err := json.Marshal(sesh.loginDetails)
	if err != nil {
		return err
	}

	resp, err := sesh.Client.Post(
		"https://api.divolt.xyz/auth/session/login",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return errors.New("failed to send login request")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed validation")
	}

	defer resp.Body.Close()

	var ar Authentication
	err = util.UnmarshalResponseBody(resp, &ar)
	if err != nil {
		return err
	}

	sesh.Authentication = &ar

	return nil
}

func (sesh *Session) Logout() {
	var req *http.Request
	var resp *http.Response

	req, err := sesh.AuthenticatedRequest(
		http.MethodPost,
		"auth/session/logout",
		nil,
	)
	if err != nil {
		goto logouterr
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err = sesh.Client.Do(req)
	if err != nil {
		goto logouterr
	}

	resp.Body.Close()
	return

logouterr:
	fmt.Fprintln(os.Stderr, "failed to logout current session")
	os.Exit(1)
}

func AuthError(resp *http.Response) error {
	var autherr AuthenticationError
	err := util.UnmarshalResponseBody(resp, &autherr)
	if err != nil {
		return err
	}

	return errors.New(autherr.Error)
}

func (sesh *Session) AuthenticatedRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("https://api.divolt.xyz/%s", path),
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-session-token", sesh.Authentication.Token)

	return req, nil
}
