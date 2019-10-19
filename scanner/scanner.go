package scanner

import (
	"Proxy/config"
	"Proxy/domain"
	"bytes"
	"unicode"
)

type scanner struct {
	body    []byte
	newBody []byte
	index   int
	url     domain.Url
}

// buffer to minimize memory allocation occurrence
const extendNewBodyByteSize = 200

func (s *scanner) getPort() {
	if s.index >= len(s.body) {
		return
	}
	x := s.getRange(portAllowedChars)
	s.url.Port = s.body[s.index-1 : x] // -1 because we take : too
	s.index = x
}

func (s *scanner) getPathAndQuery(urlAction domain.UrlAction, parseHostAlone bool, params config.Params) {
	if s.index >= len(s.body) {
		return
	}
	x := s.getRange(pathAndQueryAllowedChars)
	if urlAction == nil {
		s.url.Path = s.body[s.index:x]
	} else {
		s.url.Path = scan(s.body[s.index:x], urlAction, parseHostAlone, params)
	}
	s.index = x

}

func (s *scanner) findDomainIndex() int {
	xCache := s.index
	x := s.getRange(hostAllowedChars)
	if xCache != x {
		s.index = x
		for _, dot := range dotNonStandardPatterns {
			if s.match(dot) {
				s.url.DotSpecial = append(s.url.DotSpecial, dot)
				s.index += len(dot)
				x = s.findDomainIndex()
				break
			}
		}
		s.index = xCache
	}
	return x
}

func (s *scanner) copyAndProcessHost(host []byte) {
	for _, dot := range s.url.DotSpecial {
		for i := 0; i < len(s.url.Host); i++ {
			if match(s.url.Host, dot, i) {
				cleanHost := append([]byte{}, s.url.Host[:i]...)
				cleanHost = append(cleanHost, `.`...)
				s.url.Host = append(cleanHost, s.url.Host[i+len(dot):]...)
				break
			}
		}
	}
}

func (s *scanner) findDomain() bool {
	x := s.findDomainIndex()
	s.url.Host = s.body[s.index:x]

	if len(s.url.DotSpecial) > 0 {
		s.copyAndProcessHost(s.url.Host)
	}

	if s.url.CheckValidDomain() {
		//fmt.Println("#################### domain ok: ", string(s.url.Host))
		s.index = x
		return true
	}
	s.url.DotSpecial = nil
	return false
}

func (s *scanner) getRange(allowedChars []byte) int {
	// TODO: można zoptymalizowac zwracając odrazu odppwiedni zakres
	var x int
	for x = s.index; x < len(s.body); x++ {
		char := byte(unicode.ToLower(rune(s.body[x])))
		if bytes.IndexByte(allowedChars, char) == -1 {
			break
		}
	}
	return x
}

func (s *scanner) match(pattern []byte) bool {
	length := len(pattern)
	if s.index+length >= len(s.body) {
		return false
	}
	return bytes.Equal(bytes.ToLower(s.body[s.index:s.index+length]), pattern)
}

func (s *scanner) findAndSetPatterns(patterns [][]byte, part *[]byte) bool {
	for _, pattern := range patterns {
		if s.match(pattern) {
			*part = pattern
			s.index += len(pattern)
			return true
		}
	}
	return false
}

func (s *scanner) findSchemeSeparator() bool {
	var colon []byte
	for _, pattern := range colonPatterns {
		if s.match(pattern) {
			colon = pattern
			s.index += len(colon)
			break
		}
	}
	if s.findAndSetPatterns(slashSlashPatterns, &s.url.SchemeSeparator) {
		if colon == nil {
			s.url.SchemeSeparator = append([]byte(""), s.url.SchemeSeparator...)
		} else {
			s.url.SchemeSeparator = append(colon, s.url.SchemeSeparator...)
		}
		return true
	}
	return false
}

func (s *scanner) scanForPortAndPath(urlAction domain.UrlAction, parseHostAlone bool, params config.Params) {
	if s.findAndSetPatterns(colonPatterns, &s.url.Port) {
		s.getPort()
	} else {
		s.url.Port = nil
	}
	s.getPathAndQuery(urlAction, parseHostAlone, params)
}

func (s *scanner) scan(urlAction domain.UrlAction, parseHostAlone bool, params config.Params) {
	lastIndex := 0
	for s.index = 0; s.index < len(s.body); s.index++ {
		lastIndex = s.index
		if s.findAndSetPatterns(schemaPatterns, &s.url.Scheme) && s.findSchemeSeparator() && s.findDomain() {
			s.scanForPortAndPath(urlAction, parseHostAlone, params)

			if urlAction != nil {
				s.url = urlAction(s.url, params.ScannerProxyUrl, params.ScannerTargetUrl)
				s.newBody = append(s.newBody, s.url.Scheme...)
				s.newBody = append(s.newBody, s.url.SchemeSeparator...)
				s.newBody = append(s.newBody, s.url.Host...)
				s.newBody = append(s.newBody, s.url.Port...)
				s.newBody = append(s.newBody, s.url.Path...)
			}
			s.url.Path = nil
			s.url.DotSpecial = nil
			s.index -= 1
		} else if s.findSchemeSeparator() && s.findDomain() {
			s.url.Scheme = nil
			s.scanForPortAndPath(urlAction, parseHostAlone, params)
			if urlAction != nil {
				s.url = urlAction(s.url, params.ScannerProxyUrl, params.ScannerTargetUrl)
				s.newBody = append(s.newBody, s.url.SchemeSeparator...)
				s.newBody = append(s.newBody, s.url.Host...)
				s.newBody = append(s.newBody, s.url.Port...)
				s.newBody = append(s.newBody, s.url.Path...)
			}
			s.url.Path = nil
			s.url.DotSpecial = nil
			s.index -= 1
		} else if parseHostAlone && s.findDomain() {

			s.url.Scheme = nil
			s.url.SchemeSeparator = nil
			s.scanForPortAndPath(urlAction, parseHostAlone, params)
			if urlAction != nil {
				s.url = urlAction(s.url, params.ScannerProxyUrl, params.ScannerTargetUrl)
				s.newBody = append(s.newBody, s.url.Host...)
				s.newBody = append(s.newBody, s.url.Port...)
				s.newBody = append(s.newBody, s.url.Path...)
			}
			s.url.Path = nil
			s.url.DotSpecial = nil
			s.index -= 1
		} else {
			if urlAction != nil {
				s.newBody = append(s.newBody, s.body[lastIndex:s.index+1]...)
			}
		}
	}
}

func init() {
	domain.TldPrepare()
}

func NewScanner() *scanner {
	s := scanner{}
	return &s
}

func (s *scanner) Scan(body []byte, urlAction domain.UrlAction, parseHostAlone bool, params config.Params) []byte {
	s.body = body
	s.newBody = make([]byte, 0, len(body)+extendNewBodyByteSize)
	s.scan(urlAction, parseHostAlone, params)
	if urlAction != nil {
		return s.newBody
	}
	return nil
}

func (s *scanner) ScanAsynch(body []byte, urlAction domain.UrlAction, parseHostAlone bool, params config.Params, newBody chan []byte) {
	newBody <- s.Scan(body, urlAction, parseHostAlone, params)
}

func scan(body []byte, urlAction domain.UrlAction, parseHostAlone bool, params config.Params) []byte {
	s := NewScanner()
	return s.Scan(body, urlAction, parseHostAlone, params)
}
