package model

import (
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// example key: touring-log/raw/thing=thingName/year=2022/month=01/day=12/2022-01-12-12-51.dat.
var keyValidater = regexp.MustCompile(`^touring-log/raw/thing=.+/year=\d{4}/month=[01][0-9]/day=[0-2][0-9]/\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}.dat$`)

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

func (p *Payload) GetDate() (*time.Time, error) {
	words := strings.Split(p.GetKey(), "/")

	year := strings.TrimPrefix(words[3], "year=")
	month := strings.TrimPrefix(words[4], "month=")
	day := strings.TrimPrefix(words[5], "day=")

	const shortForm = "2006-01-02"
	date, err := time.Parse(shortForm, year+"-"+month+"-"+day)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse date")
	}

	return &date, nil
}

func (p *Payload) GetUnit() string {
	words := strings.Split(p.GetKey(), "/")

	return strings.TrimPrefix(words[2], "thing=")
}
