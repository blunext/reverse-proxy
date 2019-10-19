package proxy

import (
	"Proxy/config"
	"Proxy/scanner"
	"Proxy/tools"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func ProcessResponse(params config.Params) func(response *http.Response) error {

	return func(response *http.Response) error {
		if response.StatusCode == http.StatusSwitchingProtocols {
			return nil
		}

		if params.Config.Debug {
			log.Println("\n################## response oryginal ##################")
			tools.DumpResponse(response, params.Config.DebugBody)
		}

		if tools.TextContentType(response.Header.Get("Content-Type")) {

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				return err
			}

			err = response.Body.Close()
			if err != nil {
				log.Printf("błąd zamykania body: %s", err)
				return err
			}

			if params.Config.SaveBodyRequestToFile {
				tools.WriteTestFile(body, tools.URIToFileName(response.Request.RequestURI))
			}
			ch := make(chan []byte)
			s := scanner.NewScanner()
			go s.ScanAsynch(body, scanner.ReplaceTargetToProxy, false, params, ch)
			body = <-ch
			response.ContentLength = int64(len(body))
			response.Header.Set("Content-Length", strconv.Itoa(len(body)))
			response.Header.Set("x-proxy-body", "processed")

			response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		} else {
			response.Header.Set("x-proxy-body", "not processed")
		}

		response.Header.Del("Access-Control-Allow-Origin")
		for name, headers := range response.Header {
			//name = strings.ToLower(name)
			sameHeaders := processSameHeaders(headers, scanner.ReplaceTargetToProxy, params)
			response.Header.Del(name)
			for _, h := range sameHeaders {
				response.Header.Add(name, h)
			}
		}

		if params.Config.Debug {
			log.Println("\n################## response modyfied ##################")
			tools.DumpResponse(response, params.Config.DebugBody)

		}

		return nil
	}

}
