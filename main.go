package main

import (
	"Proxy/config"
	"Proxy/proxy"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func initParams() config.Params {

	params := config.Params{}

	params.ScannerTargetUrl.Scheme = []byte(targetScheme)
	params.ScannerTargetUrl.Host = []byte(targetHost)

	params.ScannerProxyUrl.Scheme = []byte(proxyScheme)
	params.ScannerProxyUrl.Host = []byte(proxyHost)
	params.ScannerProxyUrl.Port = []byte(proxyPort)

	params.Config.Debug = debug
	params.Config.DebugBody = debugBody
	params.Config.SaveBodyRequestToFile = saveBodyRequestToFile

	err := errors.New("")
	if params.TargetUrl, err = url.Parse(fmt.Sprintf("%s://%s", targetScheme, targetHost)); err != nil {
		log.Panic("błędny url target")
	}
	if params.ProxyUrl, err = url.Parse(fmt.Sprintf("%s://%s%s", proxyScheme, proxyHost, proxyPort)); err != nil {
		log.Panic("błędny url proxt")
	}
	return params
}

func main() {

	params := initParams()

	//f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()
	//wrt := io.MultiWriter(os.Stdout, f)
	//log.SetOutput(wrt)
	//log.Println(" main")

	reverseProxy := &httputil.ReverseProxy{
		Director:       proxy.ProcessRequest(params),
		ModifyResponse: proxy.ProcessResponse(params),
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			//MaxIdleConnsPerHost:   numCoroutines,
			//MaxIdleConns:          100,
			//IdleConnTimeout:       90 * time.Second,
			//TLSHandshakeTimeout:   10 * time.Second,
			//ExpectContinueTimeout: 1 * time.Second,
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)

	})
	//log.Fatal(http.ListenAndServe(proxyPort, nil))
	path := "/Users/tom/.ssh/localhost-ssl/"
	log.Fatal(http.ListenAndServeTLS(proxyPort, fmt.Sprintf("%slocalhost.crt", path), fmt.Sprintf("%slocalhost.key", path), nil))
}

const (
	targetHost            string = "olamundo.pl"
	targetScheme          string = "https"
	proxyHost             string = "localhost"
	proxyScheme           string = "https"
	proxyPort             string = ":9013"
	debugBody             bool   = false
	debug                 bool   = false
	saveBodyRequestToFile bool   = false
)

/*
pomimo pominięcia "Accept-Encoding" w requeście do serwera w sharepoincie *.axd przychodzi zakodowany Gzipem
dlatego dekompresujemy response
*/

//http://localhost:9013/login/form?authorization_uri=https%3a%2f%2flocalhost:9013%2Fauth%2Foauth%2Fauthorize%3Fclient_id%3Dtb5SFf3cRxEyspDN%26redirect_uri%3Dhttp%3a%2f%2flocalhost:9013%2Flogin%2Fauth%3Forigin_url%253D%25252Foferta%25252Fzestaw-mebli-ogrodowych-sofa-2krzesla-stolik-7843444884%26response_type%3Dcode%26state%3DdiznHY&oauth=true
