package messenger

import "fmt"

// FacebookRequest received from Facebook server on webhook, contains messages, delivery reports and/or postbacks
type FacebookRequest struct {
	Entry []struct {
		ID        int64 `json:"id"`
		Messaging []struct {
			Recipient struct {
				ID int64 `json:"id,string"`
			} `json:"recipient"`
			Sender struct {
				ID int64 `json:"id,string"`
			} `json:"sender"`
			Timestamp int               `json:"timestamp"`
			Message   *FacebookMessage  `json:"message,omitempty"`
			Delivery  *FacebookDelivery `json:"delivery"`
			Postback  *FacebookPostback `json:"postback"`
		} `json:"messaging"`
		Time int `json:"time"`
	} `json:"entry"`
	Object string `json:"object"`
}

// FacebookMessage struct for text messaged received from facebook server as part of FacebookRequest struct
type FacebookMessage struct {
	Mid  string `json:"mid"`
	Seq  int    `json:"seq"`
	Text string `json:"text"`
}

// FacebookDelivery struct for delivery reports received from Facebook server as part of FacebookRequest struct
type FacebookDelivery struct {
	Mids      []string `json:"mids"`
	Seq       int      `json:"seq"`
	Watermark int      `json:"watermark"`
}

// FacebookPostback struct for postbacks received from Facebook server  as part of FacebookRequest struct
type FacebookPostback struct {
	Payload string `json:"payload"`
}

// rawFBResponse received from Facebook server after sending the message
// if Error is null we copy this into FacebookResponse object
type rawFBResponse struct {
	MessageID   string         `json:"message_id"`
	RecipientID int64          `json:"recipient_id,string"`
	Error       *FacebookError `json:"error"`
}

// FacebookResponse received from Facebook server after sending the message
type FacebookResponse struct {
	MessageID   string `json:"message_id"`
	RecipientID int64  `json:"recipient_id,string"`
}

// FacebookError received form Facebook server if sending messages failed
type FacebookError struct {
	Code      int    `json:"code"`
	FbtraceID string `json:"fbtrace_id"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

// Error returns Go error object constructed from FacebookError data
func (err *FacebookError) Error() error {
	return fmt.Errorf("FB Error: Type %s: %s; FB trace ID: %s", err.Type, err.Message, err.FbtraceID)
}
