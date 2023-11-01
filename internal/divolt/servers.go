package divolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserStatus struct {
	JoinedAt string `json:"joined_at"`
	Timeout  string `json:"timeout"`
}

const slavArtServerID = "01G96DF05GVMT53VKYH83RMZMN"

func CheckServerStatus(sesh *Session) error {
	resp, err := sesh.AuthenticatedRequest(
		RequestInfo{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("servers/%s/members/%s", slavArtServerID, sesh.Config.Auth.UserID),
		},
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("authenticated response failed")
	}

	var status UserStatus
	
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return err
	}

	if status.JoinedAt == "" {
		return errors.New("user not in slav art server")
	}

	if status.Timeout != "" {
		return errors.New("user is in timeout")
	}

	return nil
}
