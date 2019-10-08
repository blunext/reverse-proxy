package scanner

import (
	"Proxy/domain"
	"bytes"
)

func ReplaceTargetToProxy(url domain.Url, proxyUrl, targetUrl domain.Url) domain.Url {
	//log.Println(fmt.Sprintf("Target: %s%s%s%s%s", url.Scheme, url.SchemeSeparator, url.Host, url.Port, url.Path))
	host := bytes.ToLower(url.Host)
	if bytes.Compare(host, targetUrl.Host) == 0 {
		url.Scheme = proxyUrl.Scheme
		url.Host = proxyUrl.Host
		url.Port = proxyUrl.Port
	} else if bytes.HasPrefix(host, []byte(`.`)) {
		url.Scheme = proxyUrl.Scheme
		url.Host = append([]byte(`.`), proxyUrl.Host...)
		url.Port = proxyUrl.Port
	}
	return url
}

func ReplaceProxyToTarget(url, proxyUrl, targetUrl domain.Url) domain.Url {
	//log.Println(fmt.Sprintf("Proxy: %s%s%s%s%s", url.Scheme, url.SchemeSeparator, url.Host, url.Port, url.Path))
	if bytes.Compare(url.Host, proxyUrl.Host) == 0 && bytes.Compare(url.Port, proxyUrl.Port) == 0 {
		url.Scheme = targetUrl.Scheme
		url.Host = targetUrl.Host
		url.Port = nil
	}
	return url
}
