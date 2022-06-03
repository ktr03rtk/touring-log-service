package model

import (
	"errors"
	"regexp"
)

// example key: touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-51.dat.
var keyValidater = regexp.MustCompile(`^touring-log/raw/thing=.+/year=\d{4}/month=[01][0-9]/day=[01][0-9]/\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}.dat$`)

type Payload struct {
	message []byte
	key     string
}

func NewPayload(message []byte, key string) (*Payload, error) {
	if !keyValidater.MatchString(key) {
		return nil, errors.New("failed to validate key")
	}

	return &Payload{
		message,
		key,
	}, nil
}

func (p *Payload) GetMessage() []byte {
	return p.message
}

func (p *Payload) GetKey() string {
	return p.key
}
