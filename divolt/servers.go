package divolt

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dxbednarczyk/limestone/util"
)

type UserStatus struct {
	JoinedAt string `json:"joined_at"`
	Timeout  string `json:"timeout"`
}

const slavArtServerID = "01G96DF05GVMT53VKYH83RMZMN"

func CheckServerStatus(sesh *Session) error {
	req, err := sesh.AuthenticatedRequest(
		http.MethodGet,
		fmt.Sprintf("servers/%s/members/%s", slavArtServerID, sesh.UserId),
		nil,
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
		return AuthError(resp)
	}

	var us UserStatus
	err = util.UnmarshalResponseBody(resp, &us)
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
