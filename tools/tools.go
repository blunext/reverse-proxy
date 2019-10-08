package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func DumpRequest(req *http.Request, debugBody bool) {
	debug := debugBody && TextContentType(req.Header.Get("Content-Type"))
	requestDump, err := httputil.DumpRequest(req, debug)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(requestDump))
}

func DumpResponse(response *http.Response, debugBody bool) {
	debug := debugBody && TextContentType(response.Header.Get("Content-Type"))
	requestDump, err := httputil.DumpResponse(response, debug)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(requestDump))
}

func SingleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func TextContentType(ct string) bool {
	//	see:
	//	https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
	//  https://www.iana.org/assignments/media-types/media-types.xhtml

	// TODO: get out of CSS

	str := strings.ToLower(ct)
	switch {
	// np html, plain, javascript, css, xml
	case strings.Contains(str, "text/html"):
	case strings.Contains(str, "text/plain"):
	case strings.Contains(str, "text/javascript"):
	case strings.Contains(str, "text/x-javascript"):
	case strings.Contains(str, "text/xml"):
	case strings.Contains(str, "application/json"):
	//case strings.Contains(str, "application/xml"):
	//case strings.Contains(str, "application/xhtml+xml"):
	case strings.Contains(str, "application/javascript"):
	case strings.Contains(str, "application/x-javascript"):
	default:
		return false
	}
	return true
}

func URIToFileName(name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(name, "/", "_"), "?", "QQQ")
}

func WriteTestFile(body []byte, name string) {
	fileName := "files/" + name + ".html"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// file not exists so write it
		err := ioutil.WriteFile(fileName, body, 0644)
		if err != nil {
			panic(err)
		}
	}
}
