// messenger package encapsulate Facebook Messenger API, essential API used fro Facebook chat bots
// It can send and receive messages from Facebook Messenger
package messenger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const apiURL = "https://graph.facebook.com/v2.6/"

// TestURL to mock FB server, used for testing
var TestURL = ""

// Messenger struct
type Messenger struct {
	AccessToken string
	PageID      string
	apiURL      string
	pageURL     string

	// MessageReceived event fires when message from Facebook received
	MessageReceived func(userID, pageID int64, m FacebookMessage)

	// DeliveryReceived event fires when delivery report from Facebook received
	// Omit (nil) if you don't want to manage this events
	DeliveryReceived func(userID, pageID int64, d FacebookDelivery)

	// PostbackReceived event fires when postback received from Facebook server
	// Omit (nil) if you don't use postbacks and you don't want to manage this events
	PostbackReceived func(userID, pageID int64, p FacebookPostback)
}

// New creates new messenger instance
func New(accessToken, pageID string) Messenger {
	return Messenger{
		AccessToken: accessToken,
		PageID:      pageID,
	}
}

// SendMessage sends chat message
func (msng *Messenger) SendMessage(m Message) (FacebookResponse, error) {
	if msng.apiURL == "" {
		if TestURL != "" {
			msng.apiURL = TestURL + "me/messages?access_token=" + msng.AccessToken // testing, mock FB URL
		} else {
			msng.apiURL = apiURL + "me/messages?access_token=" + msng.AccessToken
		}
	}

	s, _ := json.Marshal(m)
	log.Println("MESSAGE:", string(s))
	req, err := http.NewRequest("POST", msng.apiURL, bytes.NewBuffer(s))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return FacebookResponse{}, err
	}

	return decodeResponse(resp)
}

// SendTextMessage sends text messate to receiverID
// it is shorthand instead of crating new text message and then sending it
func (msng Messenger) SendTextMessage(receiverID int64, text string) (FacebookResponse, error) {
	m := NewTextMessage(receiverID, text)
	return msng.SendMessage(&m)
}

// ServeHTTP is HTTP handler for Messenger so it could be directly used as http.Handler
func (msng Messenger) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fbRq, _ := DecodeRequest(r)

	for _, entry := range fbRq.Entry {
		pageID := entry.ID
		for _, msg := range entry.Messaging {
			userID := msg.Sender.ID

			switch {
			case msg.Message != nil && msng.MessageReceived != nil:
				go msng.MessageReceived(userID, pageID, *msg.Message)

			case msg.Delivery != nil && msng.DeliveryReceived != nil:
				go msng.DeliveryReceived(userID, pageID, *msg.Delivery)

			case msg.Postback != nil && msng.PostbackReceived != nil:
				go msng.PostbackReceived(userID, pageID, *msg.Postback)
			}
		}
	}
}

// DecodeRequest decodes http request from FB messagner to FacebookRequest struct
// DecodeRequest will close the Body reader
// Usually you don't have to use DecodeRequest if you setup events for specific types
func DecodeRequest(r *http.Request) (FacebookRequest, error) {
	defer r.Body.Close()
	var fbRq FacebookRequest
	err := json.NewDecoder(r.Body).Decode(&fbRq)
	if err != nil {
		return fbRq, err //log.Println("!!!!", err)
	}

	return fbRq, nil
}

// decodeResponse decodes Facebook response after sending message, usually contains MessageID or Error
func decodeResponse(r *http.Response) (FacebookResponse, error) {
	defer r.Body.Close()
	var fbResp rawFBResponse
	err := json.NewDecoder(r.Body).Decode(&fbResp)
	if err != nil {
		return FacebookResponse{}, nil
	}

	if fbResp.Error != nil {
		return FacebookResponse{}, fbResp.Error.Error()
	}

	return FacebookResponse{
		MessageID:   fbResp.MessageID,
		RecipientID: fbResp.RecipientID,
	}, nil
}
