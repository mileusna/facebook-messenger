package main

import (
	"log"
	"net/http"

	"github.com/mileusna/facebook-messenger"
)

// use public messenger for simpler code demonstration
var msng = &messenger.Messenger{
	AccessToken: "YOUR_ACCESS_TOKEN_THAT_YOU_WILL_GENERATE_FOR_YOUR_PAGE_ON_FACEBOOK",
	PageID:      "YOUR_PAGE_ID",
}

// Please check the First example First, it contains more example code for sending messages
// If you don't want to use events receivers like in First example, this is another approach
// that will give you entire FacebookRequest struct received from the Facebook so you can
// use it however you like
func main() {
	// set URL for your webhook and you handler func
	http.HandleFunc("/mychatbot", myHandler)
	http.ListenAndServe(":8008", nil)
}

// myHandler is you regular http Handler
func myHandler(w http.ResponseWriter, r *http.Request) {

	msng.VerifyWebhook(w, r)                   // verify webhook if asked from Facebook
	fbRequest, _ := messenger.DecodeRequest(r) // decode entire request received from Facebook into FacebookRequest struct

	// now you have it all and you can do whatever you want with received request
	// enumerate each entry, and each message in entry
	for _, entry := range fbRequest.Entry {
		// pageID := entry.ID  // here you can find page id that received message
		for _, msg := range entry.Messaging {
			userID := msg.Sender.ID // user that sent you a message

			// but "message" can be text message, delivery report or postback, so check it what it is
			// it can only be one of this, so we use switch
			switch {
			case msg.Message != nil:
				log.Println("Msg received with content:", msg.Message.Text)
				msng.SendTextMessage(userID, "Hello there")
				// check First example for more sending messages examples

			case msg.Delivery != nil:
				// delivery report received, check First example what to do next

			case msg.Postback != nil:
				// postback received, check First example what can you do with that
				log.Println("Postback received with content:", msg.Postback.Payload)
			}
		}
	}
}
