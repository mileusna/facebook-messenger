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

// DecodeRequest decodes http request from FB messagner to FacebookRequest struct
// DecodeRequest will close the Body reader
func DecodeRequest(r *http.Request) FacebookRequest {
	defer r.Body.Close()
	var fbRq FacebookRequest
	err := json.NewDecoder(r.Body).Decode(&fbRq)
	if err != nil {
		log.Println(err)
	}

	return fbRq
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
