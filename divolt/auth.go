package divolt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"syscall"

	"github.com/dxbednarczyk/limestone/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

type Session struct {
	Client *http.Client
	Config *config.Config
}

type requestInfo struct {
	method string
	path   string
	body   io.Reader
}

var Login = cli.Command{
	Name:      "login",
	UsageText: "limestone login <email>",
	Action: func(ctx *cli.Context) error {
		email := ctx.Args().First()
		if email == "" {
			return errors.New("no email specified")
		}

		fmt.Printf("Enter the password for %s: ", email)
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		fmt.Print("\nLogging in... ")

		cfg := config.Config{
			Email:    email,
			Password: string(passwordBytes),
		}

		session := NewSession(&cfg)
		err = session.Login()
		if err != nil {
			fmt.Println()
			return err
		}

		fmt.Println("login successful.")

		err = config.CacheLoginDetails(cfg)
		if err != nil {
			return err
		}

		fmt.Println("Login details cached.")

		return nil
	},
}

var Logout = cli.Command{
	Name:      "logout",
	UsageText: "limestone logout",
	Action: func(ctx *cli.Context) error {
		fmt.Print("Logging out... ")

		cfg, err := config.GetLoginDetails()
		if err != nil {
			return err
		}

		// naming seems counterintuitive, but we obviously need
		// to authenticate before we can deauthenticate
		session := NewSession(&cfg)
		err = session.Logout()
		if err != nil {
			return err
		}

		err = config.RemoveConfigDetails()
		if err != nil {
			return err
		}

		fmt.Println("logged out successfully.")

		return nil
	},
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

	var ar config.Authentication

	err = json.NewDecoder(resp.Body).Decode(&ar)
	if err != nil {
		return err
	}

	sesh.Config.Auth = ar

	return nil
}

func (sesh *Session) Logout() error {
	resp, err := sesh.AuthenticatedRequest(
		requestInfo{
			method: http.MethodPost,
			path:   "auth/session/logout",
		},
	)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}

func AuthError(resp *http.Response) error {
	autherr := map[string]string{}

	err := json.NewDecoder(resp.Body).Decode(&autherr)
	if err != nil {
		return err
	}

	if autherr["type"] == "" {
		return nil
	}

	return errors.New(autherr["type"])
}

func (sesh *Session) AuthenticatedRequest(info requestInfo) (*http.Response, error) {
	req, err := http.NewRequest(
		info.method,
		fmt.Sprintf("https://api.divolt.xyz/%s", info.path),
		info.body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-session-token", sesh.Config.Auth.Token)

	return sesh.Client.Do(req)
}
