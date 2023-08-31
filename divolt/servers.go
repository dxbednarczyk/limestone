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
		requestInfo{
			method: http.MethodGet,
			path:   fmt.Sprintf("servers/%s/members/%s", slavArtServerID, sesh.Authentication.UserID),
		},
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AuthError(resp)
	}

	var us UserStatus
	err = json.NewDecoder(resp.Body).Decode(&us)
	if err != nil {
		return err
	}

	if us.JoinedAt == "" {
		return errors.New("user not in slav art server")
	}
	if us.Timeout != "" {
		return errors.New("user is in timeout")
	}

	return nil
}
