package nf

import (
	"strings"
	"testing"
)

func TestGenTasksfromJob(t *testing.T) {
	emailsArray := strings.Split("johuo@yunify.com;danma@yunify.com", ";")
	for _, email := range emailsArray {
		println(email)
	}
}
