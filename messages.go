package messenger

// messenger := &messenger.Messenger {
//     VerifyToken: "VERIFY_TOKEN/optional",
//     AppSecret: "APP_SECRET/optional",
//     AccessToken: "PAGE_ACCESS_TOKEN",
//     PageID: "PAGE_ID/optional",
// }

// ButtonType for buttons, it can be ButtonTypeWebURL or ButtonTypePostback
type ButtonType string

type AttachmentType string

type TemplateType string

type NotificationType string

type Message interface {
	foo()
}

func (m TextMessage) foo()    {}
func (m GenericMessage) foo() {}

const (
	// ButtonTypeWebURL is type for web links
	ButtonTypeWebURL = ButtonType("web_url")
	//ButtonTypePostback is type for postback buttons that sends data back to webhook
	ButtonTypePostback = ButtonType("postback")

	AttachmentTypeTemplate = AttachmentType("template")

	TemplateTypeGeneric = TemplateType("generic")

	NotificationTypeRegular    = NotificationType("REGULAR")
	NotificationTypeSilentPusg = NotificationType("SILENT_PUSH")
	NotificationTypeNoPush     = NotificationType("NO_PUSH")
)

type TextMessage struct {
	Message          textMessageContent `json:"message"`
	Recipient        recipient          `json:"recipient"`
	NotificationType NotificationType   `json:"notification_type,omitempty"`
}

type GenericMessage struct {
	Message          genericMessageContent `json:"message"`
	Recipient        recipient             `json:"recipient"`
	NotificationType NotificationType      `json:"notification_type,omitempty"`
}

type recipient struct {
	ID int64 `json:"id,string"`
}

type textMessageContent struct {
	Text string `json:"text,omitempty"`
}

type genericMessageContent struct {
	Attachment *attachment `json:"attachment,omitempty"`
}

// type MessageData struct {
// 	Text       string      `json:"text,omitempty"`
// 	Attachment *Attachment `json:"attachment,omitempty"`
// }

type attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload payload `json:"payload,omitempty"`
}

type payload struct {
	TemplateType string    `json:"template_type,omitempty"`
	Elements     []Element `json:"elements,omitempty"`
}

// Element in Generic Message template attachment
type Element struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle,omitempty"`
	ItemURL  string   `json:"item_url,omitempty"`
	ImageURL string   `json:"image_url,omitempty"`
	Buttons  []Button `json:"buttons,omitempty"`
}

// Button on Generic Message template element
type Button struct {
	Type    ButtonType `json:"type"`
	URL     string     `json:"url,omitempty"`
	Title   string     `json:"title"`
	Payload string     `json:"payload,omitempty"`
}

// NewTextMessage creates new text message for receiverID
// This function is here for convenient reason, you will
// probably use shorthand version SentTextMessage which sends message immediatly
func NewTextMessage(receiverID int64, text string) TextMessage {
	return TextMessage{
		Recipient: recipient{ID: receiverID},
		Message:   textMessageContent{Text: text},
	}
}

// NewGenericMessage creates new Generic Template message for receiverID
// Generic template messages are used for structured messages with images, links, buttons and postbacks
func NewGenericMessage(receiverID int64) GenericMessage {
	return GenericMessage{
		Recipient: recipient{ID: receiverID},
		Message: genericMessageContent{
			Attachment: &attachment{
				Type:    "template",
				Payload: payload{TemplateType: "generic"},
			},
		},
	}
}

// AddNewElement adds element to Generic template message with defined title, subtitle, link url and image url
// Only title is mandatory, other params can be empty string
// Generic messages can have up to 10 elements which are scolled horizontaly in users messenger
func (m *GenericMessage) AddNewElement(title, subtitle, itemURL, imageURL string) {
	m.AddElement(NewElement(title, subtitle, itemURL, imageURL))
}

// AddElement adds element e to Generic Message
// Generic messages can have up to 10 elements which are scolled horizontaly in users messenger
// If element contain buttons, you can create element with NewElement, than add some buttons with
// AddWebURLButton and AddPostbackButton and add it to message using this method
func (m *GenericMessage) AddElement(e Element) {
	m.Message.Attachment.Payload.Elements = append(m.Message.Attachment.Payload.Elements, e)
}

// NewElement creates new element with defined title, subtitle, link url and image url
// Only title is mandatory, other params can be empty string
// Instead of calling this function you can also initialize Element struct, depends what you prefere
func NewElement(title, subtitle, itemURL, imageURL string) Element {
	e := Element{
		Title:    title,
		Subtitle: subtitle,
		ItemURL:  itemURL,
		ImageURL: imageURL,
	}
	return e
}

// AddWebURLButton adds web link URL button to the element
func (e *Element) AddWebURLButton(title, URL string) {
	b := Button{
		Type:  ButtonTypeWebURL,
		Title: title,
		URL:   URL,
	}
	e.Buttons = append(e.Buttons, b)
}

// AddPostbackButton adds button that sends payload string back to webhook when pressed
func (e *Element) AddPostbackButton(title, payload string) {
	b := Button{
		Type:    ButtonTypePostback,
		Title:   title,
		Payload: payload,
	}
	e.Buttons = append(e.Buttons, b)
}
