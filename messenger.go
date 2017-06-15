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
	VerifyToken string
	PageID      string

	apiURL  string
	pageURL string

	// MessageReceived event fires when message from Facebook received
	MessageReceived func(msng *Messenger, userID int64, m FacebookMessage)

	// DeliveryReceived event fires when delivery report from Facebook received
	// Omit (nil) if you don't want to manage this events
	DeliveryReceived func(msng *Messenger, userI int64, d FacebookDelivery)

	// PostbackReceived event fires when postback received from Facebook server
	// Omit (nil) if you don't use postbacks and you don't want to manage this events
	PostbackReceived func(msng *Messenger, userID int64, p FacebookPostback)
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
	m := msng.NewTextMessage(receiverID, text)
	return msng.SendMessage(&m)
}

// ServeHTTP is HTTP handler for Messenger so it could be directly used as http.Handler
func (msng *Messenger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msng.VerifyWebhook(w, r)    // verify webhook if needed
	fbRq, _ := DecodeRequest(r) // get FacebookRequest object

	for _, entry := range fbRq.Entry {
		for _, msg := range entry.Messaging {
			userID := msg.Sender.ID
			switch {
			case msg.Message != nil && msng.MessageReceived != nil:
				go msng.MessageReceived(msng, userID, *msg.Message)

			case msg.Delivery != nil && msng.DeliveryReceived != nil:
				go msng.DeliveryReceived(msng, userID, *msg.Delivery)

			case msg.Postback != nil && msng.PostbackReceived != nil:
				go msng.PostbackReceived(msng, userID, *msg.Postback)
			}
		}
	}
}

// VerifyWebhook verifies your webhook by checking VerifyToken and sending challange back to Facebook
func (msng Messenger) VerifyWebhook(w http.ResponseWriter, r *http.Request) {
	// Facebook sends this query for verifying webhooks
	// hub.mode=subscribe&hub.challenge=1085525140&hub.verify_token=moj_token
	if r.FormValue("hub.mode") == "subscribe" {
		if r.FormValue("hub.verify_token") == msng.VerifyToken {
			w.Write([]byte(r.FormValue("hub.challenge")))
			return
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
	return fbRq, err
}

// decodeResponse decodes Facebook response after sending message, usually contains MessageID or Error
func decodeResponse(r *http.Response) (FacebookResponse, error) {
	defer r.Body.Close()
	var fbResp rawFBResponse
	err := json.NewDecoder(r.Body).Decode(&fbResp)
	if err != nil {
		return FacebookResponse{}, err
	}

	if fbResp.Error != nil {
		return FacebookResponse{}, fbResp.Error.Error()
	}

	return FacebookResponse{
		MessageID:   fbResp.MessageID,
		RecipientID: fbResp.RecipientID,
	}, nil
}
