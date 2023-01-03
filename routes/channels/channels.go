package channels

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"limestone/routes/auth/session"
	"limestone/util"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/sacOO7/gowebsocket"
)

type downloadMessage struct {
	Content string `json:"content"`
}

type Message struct {
	Id      string  `json:"_id"`
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

type Embed struct {
	Description string `json:"description"`
}

const music_dl_request_channel_id = "01G9AZ9AMWDV227YA7FQ5RV8WB"
const music_dl_uploads_channel_id = "01G9AZ9Q2R5VEGVPQ4H99C01YP"

func SendDownloadMessage(sesh *session.Session, url string) error {
	jsoned, err := json.Marshal(
		downloadMessage{
			Content: "!dl " + url,
		},
	)
	if err != nil {
		return err
	}

	req, err := util.RequestWithSessionToken(
		http.MethodPost,
		fmt.Sprintf("channels/%s/messages", music_dl_request_channel_id),
		bytes.NewReader(jsoned),
		sesh.Token,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Idempotency-Key", uuid.NewString())

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

	return nil
}

func GetUploadMessage(sesh *session.Session) (Message, error) {
	var message Message

	fmt.Print("Waiting for authentication... ")

	socket := gowebsocket.New("wss://ws.divolt.xyz")
	defer socket.Close()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	socket.OnConnected = func(_ gowebsocket.Socket) {
		json := fmt.Sprintf(`{"type":"Authenticate","token":"%s"}`, sesh.Token)
		socket.SendText(json)
	}

	socket.OnTextMessage = func(textMessage string, _ gowebsocket.Socket) {
		rawData := struct {
			Type string `json:"type"`
		}{}

		err := json.Unmarshal([]byte(textMessage), &rawData)
		if err != nil {
			socket.Close()
			sesh.Logout()
			log.Fatal(err)
		}

		switch rawData.Type {
		case "Authenticated":
			fmt.Println("Authenticated.")
			fmt.Print("Waiting for a response... ")
		case "Message":
			err := json.Unmarshal([]byte(textMessage), &message)
			if err != nil {
				socket.Close()
				sesh.Logout()
				log.Fatal(err)
			}

			if strings.Contains(message.Content, sesh.UserId) {
				fmt.Println("Response recieved.")

				socket.Close()
				wg.Done()
			}
		}
	}

	socket.Connect()
	wg.Wait()

	return message, nil
}
