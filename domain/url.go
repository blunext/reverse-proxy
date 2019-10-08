package domain

import (
	"bytes"
)

type Url struct {
	Scheme, SchemeSeparator, Host, Port, Path []byte
	DotSpecial                                [][]byte
}

type UrlAction func(url Url, proxyUrl, targetUrl Url) Url

func (url *Url) CheckValidDomain() bool {
	if len(url.Host) == 0 || len(url.Host) > 253 {
		return false
	}
	if string(url.Host) == "localhost" {
		return true
	}
	i := bytes.LastIndex(url.Host, []byte(".")) // TODO : dodac więcej kropek
	if i == -1 {
		return false
	}
	suffix := url.Host[i+1:] // without the dot
	return chcekValidTLD(suffix)
}
