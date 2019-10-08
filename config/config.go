package config

import (
	"Proxy/domain"
	"net/url"
)

type Config struct {
	DebugBody             bool
	Debug                 bool
	SaveBodyRequestToFile bool
}

type Params struct {
	TargetUrl        *url.URL
	ProxyUrl         *url.URL
	ScannerTargetUrl domain.Url
	ScannerProxyUrl  domain.Url
	Config
}
