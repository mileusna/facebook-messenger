package messenger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// Welcome struct used for setting messenger welcome message
type welcome struct {
	SettingType   string         `json:"setting_type"`
	ThreadState   string         `json:"thread_state"`
	CallToActions []callToAction `json:"call_to_actions"`
}

type welcomeResponse struct {
	Result string         `json:"result"`
	Error  *FacebookError `json:"error"`
}

type callToAction struct {
	Message interface{} `json:"message,omitempty"`
}

// SetWelcomeText sets plain text welcome message
func (msng *Messenger) SetWelcomeText(text string) error {
	m := textMessageContent{Text: text}
	return msng.setWelcome(&m)
}

//SetWelcomeGeneric sets generic template welcome message
func (msng *Messenger) SetWelcomeGeneric(m GenericMessage) error {
	return msng.setWelcome(&m.Message)
}

// DeleteWelcome removes welcome message
func (msng *Messenger) DeleteWelcome() error {
	return msng.setWelcome(nil)
}

func (msng *Messenger) setWelcome(m interface{}) error {

	if msng.pageURL == "" {
		msng.pageURL = apiURL + msng.PageID + "/thread_settings?access_token=" + msng.AccessToken
		log.Println(msng.pageURL)
	}

	w := welcome{
		SettingType:   "call_to_actions",
		ThreadState:   "new_thread",
		CallToActions: []callToAction{},
	}

	// if m exists (no delete welcome) add it to callToAction
	if m != nil {
		w.CallToActions = append(w.CallToActions, callToAction{Message: m})
	}

	s, _ := json.Marshal(w)
	log.Println("MESSAGE:", string(s))
	req, err := http.NewRequest("POST", msng.pageURL, bytes.NewBuffer(s))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	reply := welcomeResponse{}
	err = json.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return err
	}

	if reply.Error != nil {
		return reply.Error.Error()
	}

	return nil
}
