package channels

import (
	"encoding/json"
	"errors"
	"fmt"
	"limestone/routes/auth/session"
	"limestone/util"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sacOO7/gowebsocket"
)

type Message struct {
	Id      string `json:"_id"`
	Content string `json:"content"`
	Embeds  []struct {
		Description string `json:"description"`
	} `json:"embeds"`
}

const music_dl_request_channel_id = "01G9AZ9AMWDV227YA7FQ5RV8WB"

func SendDownloadMessage(sesh *session.Session, url string) error {
	content := fmt.Sprintf(`{"content":"!dl %s"}`, url)

	req, err := util.RequestWithSessionToken(
		http.MethodPost,
		fmt.Sprintf("channels/%s/messages", music_dl_request_channel_id),
		strings.NewReader(content),
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
	wg := new(sync.WaitGroup)
	wg.Add(1)

	socket := gowebsocket.New("wss://ws.divolt.xyz")
	defer socket.Close()

	socket.OnConnected = func(_ gowebsocket.Socket) {
		json := fmt.Sprintf(`{"type":"Authenticate","token":"%s"}`, sesh.Token)
		socket.SendText(json)

		fmt.Print("Waiting for authentication... ")
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

			go func() {
				for {
					time.Sleep(10 * time.Second)
					socket.SendText("{\"type\":\"Ping\",\"data\":0}")
				}
			}()
		case "Message":
			err := json.Unmarshal([]byte(textMessage), &message)
			if err != nil {
				socket.Close()
				sesh.Logout()

				log.Fatal(err)
			}

			if strings.Contains(message.Content, sesh.UserId) {
				socket.Close()
				wg.Done()
			}
		}
	}

	socket.Connect()
	wg.Wait()

	fmt.Println("Response recieved.")

	return message, nil
}
