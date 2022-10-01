package model

type Message struct {
	Text          string   `json:"string"`
	Topics        []string `json:"topics"`
	FromPersonId  string   `json:"from_person_id"`
	MinTrustLevel int      `json:"min_trust_level"`
}
