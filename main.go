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

	target := flag.String("target", targetUrl, "target url")
	local := flag.String("proxy", proxyUrl, "proxy url")
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
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)

	})

	if params.ProxyUrl.Scheme == "http" {
		log.Fatal(http.ListenAndServe(string(params.ScannerProxyUrl.Port), nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(string(params.ScannerProxyUrl.Port), params.CerFile, params.KeyFile, nil))
	}

}

const (
	targetUrl             string = "https://abc.xyz"
	proxyUrl              string = "https://localhost:9013"
	debugBody             bool   = false
	debug                 bool   = false
	saveBodyRequestToFile bool   = false
)
