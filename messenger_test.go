package messenger_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mileusna/facebook-messenger"
)

var fs *httptest.Server

var ts *httptest.Server

const (
	verifyToken = "my_secret_token"
)

func TestMain(m *testing.M) {
	// fs will mock up fb messenger server
	fs = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := messenger.FacebookResponse{
			RecipientID: 12123213123,
			MessageID:   "mid00000TEST00000TEST00000TEST",
		}
		b, _ := json.Marshal(rec)
		w.Write(b)
	}))
	defer fs.Close()

	// setup chatbot
	messenger.TestURL = fs.URL + "/"

	msng := &messenger.Messenger{
		AccessToken: "XXXXXXX",
		VerifyToken: verifyToken,
	}

	//chatbot.Messenger.MessageReceived = chatbot.MessageReceived
	//chatbot.Messenger.DeliveryReceived = chatbot.DeliveryReceived
	//chatbot.Messenger.PostbackReceived = chatbot.PostbackReceived

	// ts is our test chatbot
	ts = httptest.NewServer(msng)
	defer ts.Close()

	m.Run()
}

func TestVerify(t *testing.T) {
	challenge := "1122334455"
	verifyReq := ts.URL + "/?test=1&hub.mode=subscribe&hub.challenge=" + challenge + "&hub.verify_token=" + verifyToken
	resp, _ := http.Get(verifyReq)
	defer resp.Body.Close()
	s, _ := ioutil.ReadAll(resp.Body)
	if string(s) != challenge {
		t.Error("Challenge failed, expected", challenge, "returned", string(s))
	}
}
