package client_test

import (
	"testing"

	"github.com/orov-io/lbasket/client"
	. "github.com/smartystreets/goconvey/convey"
)

// convey phrases
const (
	givenAClient         = "Given a lbasket client"
	callHandlerByService = "When call is handler by the service"
	responseShouldBeOK   = "Then response should be OK"
)

func TestPing(t *testing.T) {
	Convey(givenAClient, t, func() {

		Convey(callHandlerByService, func() {
			pong, err := client.Ping()
			Convey(responseShouldBeOK, func() {

				So(err, ShouldBeNil)
				So(pong.Message, ShouldNotBeEmpty)
			})
		})
	})
}
