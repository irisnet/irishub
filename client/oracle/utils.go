package oracle

import (
	"net/url"
	"strings"

	"github.com/irisnet/irishub/app/v3/service/exported"
)

func GetState(state string) exported.RequestContextState {
	state = strings.ToLower(strings.TrimSpace(state))
	if state == "running" {
		return exported.RUNNING
	}
	return exported.PAUSED
}

func GetUrlParam(url *url.URL, arg string) string {
	return url.Query().Get(arg)
}
