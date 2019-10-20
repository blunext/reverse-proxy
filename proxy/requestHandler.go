package proxy

import (
	"Proxy/config"
	"Proxy/domain"
	"Proxy/scanner"
	"Proxy/tools"
	"log"
	"net/http"
)

func processSameHeaders(headers []string, urlAction domain.UrlAction, params config.Params) []string {
	var sameHeaders []string
	for _, h := range headers {
		s := scanner.NewScanner()
		h = string(s.Scan([]byte(h), urlAction, params))
		sameHeaders = append(sameHeaders, h)
		//log.Printf("########: header: %v: %v", name, h)
	}
	return sameHeaders
}

func ProcessRequest(params config.Params) func(req *http.Request) {

	return func(req *http.Request) {

		if params.Config.Debug {
			log.Println("\n################## request oryginal ##################")
			tools.DumpRequest(req, params.Config.DebugBody)
		}

		for name, headers := range req.Header {
			//name = strings.ToLower(name)
			sameHeaders := processSameHeaders(headers, scanner.ReplaceProxyToTarget, params)
			req.Header.Del(name)
			for _, h := range sameHeaders {
				req.Header.Add(name, h)
			}
		}

		req.Header.Del("Accept-Encoding")

		req.Host = params.TargetUrl.Host

		targetQuery := params.TargetUrl.RawQuery
		req.URL.Scheme = params.TargetUrl.Scheme
		req.URL.Host = params.TargetUrl.Host
		req.URL.Path = tools.SingleJoiningSlash(params.TargetUrl.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		if params.Config.Debug {
			log.Println("\n################## request processed ##################")
			tools.DumpRequest(req, params.Config.DebugBody)
		}
	}
}
