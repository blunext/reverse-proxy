package main

import (
	"Proxy/config"
	"Proxy/proxy"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func initParams() config.Params {

	params := config.Params{}

	target := flag.String("target", "https://www.google.com", "target url")
	local := flag.String("proxy", "https://localhost:9013", "proxy url")
	path := "/Users/tom/.ssh/localhost-ssl/"
	flag.StringVar(&params.CerFile, "certFile", fmt.Sprintf("%slocalhost.crt", path), "path to cert file")
	flag.StringVar(&params.KeyFile, "keyFile", fmt.Sprintf("%slocalhost.key", path), "path to key file")

	flag.Parse()

	var err error
	if params.TargetUrl, err = url.ParseRequestURI(*target); err != nil {
		fmt.Println("target url is invalid")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if params.ProxyUrl, err = url.ParseRequestURI(*local); err != nil {
		fmt.Println("proxy url is invalid")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if params.ProxyUrl.Scheme != "http" {
		if !fileExists(params.CerFile) {
			fmt.Println("specify certFile")
			os.Exit(0)
		}
		if !fileExists(params.KeyFile) {
			fmt.Println("specify keyFile")
			os.Exit(0)
		}
	}

	params.ScannerTargetUrl.Scheme = []byte(params.TargetUrl.Scheme)
	params.ScannerTargetUrl.Host = []byte(params.TargetUrl.Host)

	params.ScannerProxyUrl.Scheme = []byte(params.ProxyUrl.Scheme)

	s := strings.Split(params.ProxyUrl.Host, ":")
	if len(s) == 1 {
		params.ScannerProxyUrl.Host = []byte(params.ProxyUrl.Host)
		params.ScannerProxyUrl.Port = []byte(":80")
	} else {
		params.ScannerProxyUrl.Host = []byte(s[0])
		params.ScannerProxyUrl.Port = []byte(":" + s[1])
	}

	params.Config.Debug = debug
	params.Config.DebugBody = debugBody
	params.Config.SaveBodyRequestToFile = saveBodyRequestToFile

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
	log.Fatal(http.ListenAndServeTLS(string(params.ScannerProxyUrl.Port), params.CerFile, params.KeyFile, nil))
	//log.Fatal(http.ListenAndServeTLS(string(params.ScannerProxyUrl.Port), fmt.Sprintf("%slocalhost.crt", path), fmt.Sprintf("%slocalhost.key", path), nil))
	// ./aaa -proxy https://localhost:9013 -target https://wiki.mbank.pl -certFile /Users/tomek/.ssh/localhost-ssl/localhost.crt -keyFile /Users/tomek/.ssh/localhost-ssl/localhost.key
}

const (
	//targetHost            string = "allegro.pl"
	//targetScheme          string = "https"
	//proxyHost             string = "localhost"
	//proxyScheme           string = "https"
	//proxyPort             string = ":9013"
	debugBody             bool = false
	debug                 bool = false
	saveBodyRequestToFile bool = false
)

/*
pomimo pominięcia "Accept-Encoding" w requeście do serwera w sharepoincie *.axd przychodzi zakodowany Gzipem
dlatego dekompresujemy response
*/

//http://localhost:9013/login/form?authorization_uri=https%3a%2f%2flocalhost:9013%2Fauth%2Foauth%2Fauthorize%3Fclient_id%3Dtb5SFf3cRxEyspDN%26redirect_uri%3Dhttp%3a%2f%2flocalhost:9013%2Flogin%2Fauth%3Forigin_url%253D%25252Foferta%25252Fzestaw-mebli-ogrodowych-sofa-2krzesla-stolik-7843444884%26response_type%3Dcode%26state%3DdiznHY&oauth=true
