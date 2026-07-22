package graph

import "time"

type MessagesResponse struct {
	Value []Message `json:"value"`
}

type Recipient struct {
	EmailAddress EmailAddress `json:"emailAddress"`
}

type EmailAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Message struct {
	ID               string    `json:"id"`
	Subject          string    `json:"subject"`
	HasAttachments   bool      `json:"hasAttachments"`
	ReceivedDateTime time.Time `json:"receivedDateTime"`
	From             Recipient `json:"from"`
}

type AttachmentsResponse struct {
	Value []Attachment `json:"value"`
}

type Attachment struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ContentType  string `json:"contentType"`
	Size         int64  `json:"size"`
	ContentBytes []byte `json:"contentBytes"`
}

type MailMetadata struct {
	ID               string    `json:"id"`
	Subject          string    `json:"subject"`
	ReceivedDateTime time.Time `json:"receivedDateTime"`
	From             string    `json:"from"`
}
