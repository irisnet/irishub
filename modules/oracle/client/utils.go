package oracle

import (
	"net/url"
)

func GetUrlParam(url *url.URL, arg string) string {
	return url.Query().Get(arg)
}
