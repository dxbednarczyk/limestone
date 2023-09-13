package divolt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dxbednarczyk/limestone/internal/config"
)

type Session struct {
	Client *http.Client
	Config *config.Config
}

type RequestInfo struct {
	Method string
	Path   string
	Body   io.Reader
}

func NewSession(cfg *config.Config) Session {
	return Session{
		Client: http.DefaultClient,
		Config: cfg,
	}
}

func (sesh *Session) Login() error {
	body, err := json.Marshal(sesh.Config)
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

	err = json.NewDecoder(resp.Body).Decode(&sesh.Config.Auth)
	if err != nil {
		return err
	}

	return nil
}

func (sesh *Session) Logout() error {
	resp, err := sesh.AuthenticatedRequest(
		RequestInfo{
			Method: http.MethodPost,
			Path:   "auth/session/logout",
		},
	)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

func (sesh *Session) AuthenticatedRequest(info RequestInfo) (*http.Response, error) {
	req, err := http.NewRequest(
		info.Method,
		fmt.Sprintf("https://api.divolt.xyz/%s", info.Path),
		info.Body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-session-token", sesh.Config.Auth.Token)

	return sesh.Client.Do(req)
}
