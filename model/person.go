package model

import (
	"errors"
	"fmt"
)

const MinTrustLevel = 1
const MaxTrustLevel = 10

type TrustConnections map[string]int

type Person struct {
	Id               string   `json:"id"`
	Topics           []string `json:"topics"`
	TrustConnections TrustConnections
}

func (tc TrustConnections) Validate() error {
	for _, trustLevel := range tc {
		if trustLevel > MaxTrustLevel || trustLevel < MinTrustLevel {
			return errors.New(fmt.Sprintf("Trust levels should be between %d and %d", MinTrustLevel, MaxTrustLevel))
		}
	}
	return nil
}
