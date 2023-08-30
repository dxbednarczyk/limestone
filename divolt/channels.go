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
	ws "github.com/sacOO7/gowebsocket"
	"github.com/urfave/cli/v2"
)

type Message struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Author  string `json:"author"`

	ID      string `json:"_id"`
	Content string `json:"content"`
	Embeds  []struct {
		Description string `json:"description"`
	} `json:"embeds"`
	Replies []string `json:"replies"`
}

const (
	requestChannelID = "01G9AZ9AMWDV227YA7FQ5RV8WB"
	uploadsChannelID = "01G9AZ9Q2R5VEGVPQ4H99C01YP"
	botUserID        = "01G9824MQPGD7GVYR0F6A6GJ2Q"
)

func SendDownloadMessage(sesh *Session, url string, quality uint) (string, error) {
	var content string

	if quality <= 4 {
		content = fmt.Sprintf(`{"content":"!dl %s %d"}`, url, quality)
	} else {
		content = fmt.Sprintf(`{"content":"!dl %s"}`, url)
	}

	resp, err := sesh.AuthenticatedRequest(
		requestInfo{
			method: http.MethodPost,
			path:   fmt.Sprintf("channels/%s/messages", requestChannelID),
			body:   strings.NewReader(content),
		},
	)
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
	wg := new(sync.WaitGroup)
	wg.Add(1)

	var message Message

	socket := ws.New("wss://ws.divolt.xyz")
	defer socket.Close()

	socket.OnConnected = func(_ ws.Socket) {
		json := fmt.Sprintf(`{"type":"Authenticate","token":"%s"}`, sesh.Authentication.Token)
		socket.SendText(json)

		fmt.Print("Waiting for authentication... ")
	}

	socket.OnTextMessage = func(textMessage string, _ ws.Socket) {
		err := json.Unmarshal([]byte(textMessage), &message)
		if err != nil {
			socket.Close()
			fmt.Fprintln(os.Stderr, err)

			sesh.Logout()
		}

		switch message.Type {
		case "Authenticated":
			fmt.Println("Authenticated. ")

			go func() {
				for {
					time.Sleep(10 * time.Second)
					socket.SendText(`{"type":"Ping","data":0}"`)
				}
			}()

			fmt.Print("Waiting for a response... ")
		case "Message":
			mentionsAuthUser := strings.Contains(message.Content, sesh.Authentication.UserID)

			if message.Channel != uploadsChannelID ||
				message.Author != botUserID ||
				!mentionsAuthUser {
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
