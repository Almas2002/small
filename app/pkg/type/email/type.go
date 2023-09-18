package email

import (
	"errors"
	"regexp"
	"strings"
)

var regexpEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type Email struct {
	value string
}

func New(email string) error {
	if len(email) == 0 {
		return nil
	}

	if !regexpEmail.MatchString(strings.ToLower(email)) {
		return errors.New("invalid email format")
	}

	return nil
}

func (e *Email) Email() Email {
	return *e
}

func (e *Email) String() string {
	return e.value
}

func (e *Email) Equal(email Email) bool {
	return e.value == email.value
}

func (e *Email) IsEmpty() bool {
	return len(strings.TrimSpace(e.String())) == 0
}

func (e *Email) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.value + `"`), nil
}
