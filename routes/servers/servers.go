package servers

import (
	"errors"
	"fmt"
	"limestone/routes/auth/session"
	"limestone/util"
	"net/http"
)

type user struct {
	JoinedAt string `json:"joined_at"`
	Timeout  string `json:"timeout"`
}

const SLAV_ART_SERVER_ID = "01G96DF05GVMT53VKYH83RMZMN"

func CheckServerStatus(sesh *session.Session) error {
	req, err := util.RequestWithSessionToken(
		http.MethodGet,
		fmt.Sprintf("servers/%s/members/%s", SLAV_ART_SERVER_ID, sesh.UserId),
		nil,
		sesh.Token,
	)
	if err != nil {
		return err
	}

	resp, err := sesh.Client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var derr session.DefaultError
		err = util.UnmarshalResponseBody(resp, &derr)
		if err != nil {
			return err
		}

		return errors.New(derr.Error)
	}

	var user user
	err = util.UnmarshalResponseBody(resp, &user)
	if err != nil {
		return err
	}

	if user.JoinedAt == "" {
		return errors.New("user not in slav art server")
	}

	if user.Timeout != "" {
		return errors.New("user is in timeout")
	}

	return nil
}
