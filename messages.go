package messenger

// ButtonType for buttons, it can be ButtonTypeWebURL or ButtonTypePostback
type ButtonType string

// AttachmentType describes attachment type in GenericMessage
type AttachmentType string

// TemplateType of template in GenericMessage
type TemplateType string

// NotificationType for sent messages
type NotificationType string

// Message interface that represents all type of messages that we can send to Facebook Messenger
type Message interface {
	foo()
}

func (m TextMessage) foo()    {} // Message interface
func (m GenericMessage) foo() {} // Message interface

const (
	// ButtonTypeWebURL is type for web links
	ButtonTypeWebURL = ButtonType("web_url")

	//ButtonTypePostback is type for postback buttons that sends data back to webhook
	ButtonTypePostback = ButtonType("postback")

	// AttachmentTypeTemplate for template attachments
	AttachmentTypeTemplate = AttachmentType("template")

	// TemplateTypeGeneric for generic message templates
	TemplateTypeGeneric = TemplateType("generic")

	// NotificationTypeRegular for regular notification type
	NotificationTypeRegular = NotificationType("REGULAR")

	// NotificationTypeSilentPush for silent push
	NotificationTypeSilentPush = NotificationType("SILENT_PUSH")

	// NotificationTypeNoPush for no push
	NotificationTypeNoPush = NotificationType("NO_PUSH")
)

// TextMessage struct used for sending text messages to messenger
type TextMessage struct {
	Message          textMessageContent `json:"message"`
	Recipient        recipient          `json:"recipient"`
	NotificationType NotificationType   `json:"notification_type,omitempty"`
}

// GenericMessage struct used for sending structural messages to messenger (messages with images, links, and buttons)
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
