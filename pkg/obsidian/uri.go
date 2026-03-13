package obsidian

import (
	"errors"
	"net/url"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

type Uri struct {
}

type UriManager interface {
	Construct(baseUri string, params map[string]string) string
	Execute(uri string) error
}

func (u *Uri) Construct(baseUri string, params map[string]string) string {
	uri := baseUri
	for key, value := range params {
		if value != "" && value != "false" {
			encoded := strings.ReplaceAll(url.QueryEscape(value), "+", "%20")
			if uri == baseUri {
				uri += "?" + key + "=" + encoded
			} else {
				uri += "&" + key + "=" + encoded
			}
		}
	}
	return uri
}

var Run = open.Run

func (u *Uri) Execute(uri string) error {
	//fmt.Println("Opening URI: ", uri)
	err := Run(uri)
	if err != nil {
		return errors.New(ExecuteUriError)

	}
	return nil
}
