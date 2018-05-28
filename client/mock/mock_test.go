package mock_test

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
	"github.com/j7b/mailgun/ips"
)

func TestClient(t *testing.T) {
	c := mock.Client(t)
	api := ips.API(c)
	if _, err := api.List(false); err != nil {
		t.Fatal(err)
	}
}
