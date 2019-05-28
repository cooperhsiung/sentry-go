package notice

import (
	"testing"
)

func TestSendWx(t *testing.T) {

	SendWx("hello")
}

func TestSendEmail(t *testing.T) {
	SendEmail()
}
