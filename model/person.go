package model

type TrustConnection struct {
	Id         string
	TrustLevel int
}

type Person struct {
	Id               string   `json:"id"`
	Topics           []string `json:"topics"`
	TrustConnections []TrustConnection
}
