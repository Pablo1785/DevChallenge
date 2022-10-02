package model

type MessageRecipients map[string][]string

type Message struct {
	Text          string   `json:"text" validate:"required"`
	Topics        []string `json:"topics" validate:"required"`
	FromPersonId  string   `json:"from_person_id" validate:"required"`
	MinTrustLevel int      `json:"min_trust_level" validate:"required,gte=1,lte=10"`
}
