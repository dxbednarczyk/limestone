package divolt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dxbednarczyk/limestone/util"
	"github.com/google/uuid"
	ws "github.com/sacOO7/gowebsocket"
	"github.com/urfave/cli/v2"
)

type Message struct {
	ID      string `json:"_id"`
	Content string `json:"content"`
	Embeds  []struct {
		Description string `json:"description"`
	} `json:"embeds"`
	Replies []string `json:"replies"`
}

type messageInfo struct {
	Author  string `json:"author"`
	Channel string `json:"channel"`
}

type messageType struct {
	Type string `json:"type"`
}

const (
	requestChannelID = "01G9AZ9AMWDV227YA7FQ5RV8WB"
	uploadsChannelID = "01G9AZ9Q2R5VEGVPQ4H99C01YP"
	botUserID        = "01G9824MQPGD7GVYR0F6A6GJ2Q"
)

func SendDownloadMessage(sesh *Session, url string) (string, error) {
	content := fmt.Sprintf(`{"content":"!dl %s"}`, url)

	req, err := sesh.AuthenticatedRequest(
		http.MethodPost,
		fmt.Sprintf("channels/%s/messages", requestChannelID),
		strings.NewReader(content),
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("Idempotency-Key", uuid.NewString())

	resp, err := sesh.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", AuthError(resp)
	}

	var message Message
	err = util.UnmarshalResponseBody(resp, &message)
	if err != nil {
		return "", err
	}

	return message.ID, nil
}

func GetUploadMessage(ctx *cli.Context, sesh *Session, sentId string) (Message, error) {
	var message Message

	wg := new(sync.WaitGroup)
	wg.Add(1)

	socket := ws.New("wss://ws.divolt.xyz")
	defer socket.Close()

	socket.OnConnected = func(_ ws.Socket) {
		json := fmt.Sprintf(`{"type":"Authenticate","token":"%s"}`, sesh.Token)
		socket.SendText(json)

		fmt.Print("Waiting for authentication... ")
	}

	socket.OnTextMessage = func(textMessage string, _ ws.Socket) {
		var mt messageType
		err := json.Unmarshal([]byte(textMessage), &mt)
		if err != nil {
			abort(&socket, sesh, err)
		}

		switch mt.Type {
		case "Authenticated":
			fmt.Print("Authenticated.\nWaiting for a response... ")

			go func() {
				for {
					time.Sleep(10 * time.Second)
					socket.SendText(`{"type":"Ping","data":0}"`)
				}
			}()
		case "Message":
			var mi messageInfo
			err := json.Unmarshal([]byte(textMessage), &mi)
			if err != nil {
				abort(&socket, sesh, err)
			}

			if mi.Channel != uploadsChannelID {
				break
			}
			if mi.Author != botUserID {
				break
			}

			err = json.Unmarshal([]byte(textMessage), &message)
			if err != nil {
				abort(&socket, sesh, err)
			}

			if !strings.Contains(message.Content, sesh.UserID) {
				break
			}

			wg.Done()
		}
	}

	socket.Connect()
	wg.Wait()

	fmt.Println("Response received.")

	return message, nil
}

//nolint:errcheck
func abort(socket *ws.Socket, sesh *Session, err error) {
	socket.Close()
	sesh.Logout()

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
