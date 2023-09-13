package divolt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/recws-org/recws"
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
	textMessage      = 1
)

func SendDownloadMessage(sesh *Session, url string, quality uint) (string, error) {
	var content string

	if quality <= 4 {
		content = fmt.Sprintf(`{"content":"!dl %s %d"}`, url, quality)
	} else {
		content = fmt.Sprintf(`{"content":"!dl %s"}`, url)
	}

	resp, err := sesh.AuthenticatedRequest(
		RequestInfo{
			Method: http.MethodPost,
			Path:   fmt.Sprintf("channels/%s/messages", requestChannelID),
			Body:   strings.NewReader(content),
		},
	)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", errors.New("invalid authentication, is the channel locked?")
	}

	var message Message

	err = json.NewDecoder(resp.Body).Decode(&message)
	if err != nil {
		return "", err
	}

	return message.ID, nil
}

// just over the threshold
//
//nolint:funlen
func GetUploadMessage(sesh *Session) (Message, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	socket, err := authenticateSocket(sesh.Config.Auth.Token)
	if err != nil {
		return Message{}, err
	}

	var message Message

	for socket.IsConnected() {
		select {
		case <-ctx.Done():
			go socket.Close()
		default:
			_, msg, err := socket.ReadMessage()
			if err != nil {
				return Message{}, err
			}

			err = json.Unmarshal(msg, &message)
			if err != nil {
				return Message{}, err
			}

			switch message.Type {
			case "Authenticated":
				slog.Info("Authenticated and awaiting a response")
			case "Message":
				mentionsAuthUser := strings.Contains(message.Content, sesh.Config.Auth.UserID)

				if isMessageInvalid(&message, uploadsChannelID, mentionsAuthUser) {
					break
				}

				cancel()
			}
		}
	}

	slog.Info("Response received")

	return message, nil
}

func authenticateSocket(token string) (*recws.RecConn, error) {
	slog.Info("Authenticating")

	socket := recws.RecConn{
		NonVerbose:       true,
		KeepAliveTimeout: 10 * time.Second,
	}

	socket.Dial("wss://ws.divolt.xyz", http.Header{})

	for start := time.Now(); ; {
		if socket.IsConnected() {
			break
		}

		if time.Since(start) > 10*time.Second {
			return nil, errors.New("dialing websocket timed out")
		}
	}

	authPayload := fmt.Sprintf(`{"type":"Authenticate","token":"%s"}`, token)

	err := socket.WriteMessage(textMessage, []byte(authPayload))
	if err != nil {
		return nil, err
	}

	return &socket, nil
}

func isMessageInvalid(message *Message, wantChannel string, condition bool) bool {
	return message.Channel != wantChannel ||
		message.Author != botUserID ||
		!condition
}
