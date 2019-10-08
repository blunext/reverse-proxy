package scanner

import (
	"Proxy/config"
	"Proxy/domain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// todo: https:\u002f\u002fprofile.mbank.pl:443\u002f_layouts\u002f15\u002fEditProfile.aspx?UserSettingsProvider=234bf0ed\u00252D70db\u00252D4158\u00252Da332\u00252D4dfd683b4148\u0026ReturnUrl=https\u00253A\u00252F\u00252Fintranet\u00252Embank\u00252Epl

func config1(targetHost string) config.Params {
	const (
		//targetHost   string = "abc.com"
		targetScheme string = "https"
		proxyHost    string = "localhost"
		proxyScheme  string = "http"
		proxyPort    string = ":9012"
	)

	var proxyParams config.Params

	//err := errors.NewScanner("")
	//if proxyParams.targetUrl, err = url.Parse(fmt.Sprintf("%s://%s", targetScheme, targetHost)); err != nil {
	//	log.Panic("błędny url target")
	//}
	proxyParams.ScannerTargetUrl.Scheme = []byte(targetScheme)
	proxyParams.ScannerTargetUrl.Host = []byte(targetHost)

	//if proxyParams.proxyUrl, err = url.Parse(fmt.Sprintf("%s://%s%s", proxyScheme, proxyHost, proxyPort)); err != nil {
	//	log.Panic("błędny url proxt")
	//}
	proxyParams.ScannerProxyUrl.Scheme = []byte(proxyScheme)
	proxyParams.ScannerProxyUrl.Host = []byte(proxyHost)
	proxyParams.ScannerProxyUrl.Port = []byte(proxyPort)

	return proxyParams
}

func TestSomething(t *testing.T) {

	proxyParams := config1("abc.com")

	var body, expected, bodyChanged []byte
	//
	body = []byte(`aaa s%2F%2Fs \u003A http endhttps, %3a//, http%3a//google.com:221/dsddsd?Sdad?DA/dsadas/#[]=1221 "http://dupa.com" :\\u002f\\u002f, \\u00253A\\u00252F\\u00252F dsda dsada http%3a//google.com .dsds ds d`)

	bodyChanged = scan(body, nil, false, proxyParams)
	assert.True(t, len(bodyChanged) == 0, "zwrócono nie puste body")
	f := func(url, proxyUrl, targetUrl domain.Url) domain.Url { return url }
	bodyChanged = scan(body, f, true, proxyParams)
	assert.Equal(t, body, bodyChanged, "powinny być takie same:\n"+string(bodyChanged)+"\n"+string(body))

	body = []byte(`} else if( (document.location.href.indexOf('abc.com') !== -1 || document.location.href.indexOf('m.gazeta.pl') !== -1) ||`)
	expected = []byte(`} else if( (document.location.href.indexOf('localhost:9012') !== -1 || document.location.href.indexOf('m.gazeta.pl') !== -1) ||`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	proxyParams = config1("www.abc.com")

	body = []byte(`} else if( (document.location.href.indexOf('www.abc.com') !== -1 || document.location.href.indexOf('m.gazeta.pl') !== -1) ||`)
	expected = []byte(`} else if( (document.location.href.indexOf('localhost:9012') !== -1 || document.location.href.indexOf('m.gazeta.pl') !== -1) ||`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	proxyParams = config1("abc.com")

	body = []byte(`<img src="https:\/\/abc.com\/wp-admin\/admin-ajax.php" class="header__brand__image" title="" alt="naTemat.pl">`)
	expected = []byte(`<img src="http:\/\/localhost:9012\/wp-admin\/admin-ajax.php" class="header__brand__image" title="" alt="naTemat.pl">`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`<img src="//abc.com/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	expected = []byte(`<img src="//localhost:9012/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`<img src="http://abc.com/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	expected = []byte(`<img src="http://localhost:9012/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(` abc.com/logo abc.com/logo abc.com/logo `)
	expected = []byte(` localhost:9012/logo localhost:9012/logo localhost:9012/logo `)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(` abc.com abc.com abc.com `)
	expected = []byte(` localhost:9012 localhost:9012 localhost:9012 `)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`<img src="http://localhost:9012/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	expected = []byte(`<img src="https://abc.com/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	bodyChanged = scan(body, ReplaceProxyToTarget, false, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`Set-Cookie: __cfduid=d4333a86847f959c46d24ed08545d9aac1557651193; expires=Mon, 11-May-20 08:53:13 GMT; path=/; domain=.abc.com; HttpOnly; Secure`)
	expected = []byte(`Set-Cookie: __cfduid=d4333a86847f959c46d24ed08545d9aac1557651193; expires=Mon, 11-May-20 08:53:13 GMT; path=/; domain=.localhost:9012; HttpOnly; Secure`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`Set-Cookie: __cfduid=d4333a86847f959c46d24ed08545d9aac1557651193; expires=Mon, 11-May-20 08:53:13 GMT; path=/; domain=.localhost:9012; HttpOnly; Secure`)
	expected = []byte(`Set-Cookie: __cfduid=d4333a86847f959c46d24ed08545d9aac1557651193; expires=Mon, 11-May-20 08:53:13 GMT; path=/; domain=.abc.com; HttpOnly; Secure`)
	bodyChanged = scan(body, ReplaceProxyToTarget, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`<img src="http%3A\u002f\u002fabc.com/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	expected = []byte(`<img src="http%3A\u002f\u002flocalhost:9012/logo/" class="header__brand__image" title="" alt="naTemat.pl">`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`'ReturnUrl=https\u00253a\u00252f\u00252fabc.com dsds', 'Aktualizacja profilu`)
	expected = []byte(`'ReturnUrl=http\u00253a\u00252f\u00252flocalhost:9012 dsds', 'Aktualizacja profilu`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`bla 'https\u00253A\u00252F\u00252Fabc\u00252Ecom/sasa', bla`)
	expected = []byte(`bla 'http\u00253a\u00252f\u00252flocalhost:9012/sasa', bla`)
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	body = []byte(`'O mnie', 'https:\u002f\u002fabc.com\u002f_layouts\u002f15\u002fEditProfile.aspx?UserSettingsProvider=234bf0ed\u00252D70db\u00252D4158\u00252Da332\u00252D4dfd683b4148\u0026ReturnUrl=https\u00253A\u00252F\u00252Fabc\u00252Ecom/sasa', 'Aktualizacja profilu`)
	expected = []byte(`'O mnie', 'http:\u002f\u002flocalhost:9012\u002f_layouts\u002f15\u002fEditProfile.aspx?UserSettingsProvider=234bf0ed\u00252D70db\u00252D4158\u00252Da332\u00252D4dfd683b4148\u0026ReturnUrl=http\u00253a\u00252f\u00252flocalhost:9012/sasa', 'Aktualizacja profilu`)
	// http\u00253a\u00252f\u00252f -- to się zmienia
	bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	//do analizy

	//body = []byte(`'O mnie', 'https:\u002f\u002fprofile.mbank.pl:443\u002f_layouts\u002f15\u002fEditProfile.aspx?UserSettingsProvider=234bf0ed\u00252D70db\u00252D4158\u00252Da332\u00252D4dfd683b4148\u0026ReturnUrl=https\u00253A\u00252F\u00252Fintranet\u00252Embank\u00252Epl', 'Aktualizacja profilu`)
	//expected = []byte(`'O mnie', 'https:\u002f\u002fprofile.mbank.pl:443\u002f_layouts\u002f15\u002fEditProfile.aspx?UserSettingsProvider=234bf0ed\u00252D70db\u00252D4158\u00252Da332\u00252D4dfd683b4148\u0026ReturnUrl=https\u00253A\u00252F\u00252Fintranet\u00252Embank\u00252Epl', 'Aktualizacja profilu`)
	//bodyChanged = scan(body, ReplaceTargetToProxy, true, proxyParams)
	//assert.Equal(t, expected, bodyChanged, "nie zmieniły się\n"+string(bodyChanged)+"\n"+string(expected))

	//todo: dodac warunek gdy host będzie miał port zrefiniowany np 443 i https
	//	body = []byte(`'O mnie', 'https:\u002f\u002fintranet.mbank.pl:443\u002f_layouts\u002f15\u002fEditProfile.aspx?UserSettingsProvider=234bf0ed\u00252D70db\u00252D4158\u00252Da332\u00252D4dfd683b4148\u0026ReturnUrl=https\u00253A\u00252F\u00252Fintranet\u00252Embank\u00252Epl', 'Aktualizacja profilu`)
	//Location: http://allegro.pl:-1/carts/77404eee-9480-47e6-a703-cf81e2b422f3be07c0e8-d1d2-45c9-98a6-549989f383ba
	//style="background: url(https://s.natemat.pl/gfx/v2/header-footer/footer-photo.jpg?4) no-repeat center center / cover;">
	/*
	   2019/05/04 15:38:26 Target: s s hR.call p p(Q,T)   --- R.call(Q,T)
	   2019/05/04 15:38:26 Target: s s hR.call p p(Q,U,T)
	   2019/05/04 15:38:26 Target: s s hR.call p p(Q,U,T,V)
	   2019/05/04 15:38:26 Target: s s hR.call p p(Q,T,V,U,W)
	   2019/05/04 15:38:26 Target: s s hP.property p p(R)
	   2019/05/04 15:38:26 Target: s s hP.map p p=P.collect=function(V,X,S)
	   2019/05/04 15:38:26 Target: s s hP.select p p=function(T,Q,S)
	   2019/05/04 15:38:26 Target: s s hl.call p p(arguments,2);var
	   2019/05/04 15:38:26 Target: s s hP.map p p(S,function(V)
	   2019/05/04 15:38:26 Target: s s hP.map p p(R,P.property(Q))
	   2019/05/04 15:38:26 Target: s s hP.map p p(R,function(V,T,U)
	   2019/05/04 15:38:26 Target: s s hl.call p p(Q)
	   2019/05/04 15:38:26 Target: s s hP.map p p(Q,P.identity)
	   2019/05/04 15:38:26 Target: s s hl.call p p(S,0,Math.max(0,S.length-(R==null
	   2019/05/04 15:38:26 Target: s s hP.rest p p(S,Math.max(0,S.length-R))
	   2019/05/04 15:38:26 Target: s s hP.rest p p=P.tail=P.drop=function(S,R,Q)
	   2019/05/04 15:38:26 Target: s s hl.call p p(S,R==null
	   2019/05/04 15:38:26 Target: s s hl.call p p(arguments,1))
	   2019/05/04 15:38:26 Target: s s hP.zip p p=function()
	   2019/05/04 15:38:26 Target: s s hl.call p p(X,U,V),P.isNaN);return
	   2019/05/04 15:38:26 Target: s s hl.call p p(arguments,1))
	   2019/05/04 15:38:26 Target: s s hl.call p p(arguments,2);var
	   2019/05/04 15:38:26 Target: s s hl.call p p(arguments)))
	   2019/05/04 15:38:26 Target: s s hl.call p p(arguments,1);var
	   2019/05/04 15:38:26 Target: s s hl.call p p(arguments,2);return
	   2019/05/04 15:38:26 Target: s s hP.now p p();W=null;Y=R.apply(Q,V);if(!W)
	   2019/05/04 15:38:26 Target: s s hP.now p p();if(!U&&X.leading===false)
	   2019/05/04 15:38:26 Target: s s hk.open p p(a,i,true);k.onload=function()
	   2019/05/04 15:38:26 Target: s s hP.now p p()-V;if(Z
	   2019/05/04 15:38:26 Target: s s hP.now p p();var
	   2019/05/04 15:38:26 Target: s s h.call p p(this,S)
	   2019/05/04 15:38:26 Target: s s hP.map p p(x(arguments,false,false,1),String);T=function(V,U)
	   2019/05/04 15:38:26 Target: s s hc.call p p(Y);if(V!==c.call(X))
	   2019/05/04 15:38:26 Target: s s hc.call p p(Q)===
	   2019/05/04 15:38:26 Target: s s hc.call p p(R)===
	   2019/05/04 15:38:26 Target: s s hc.call p p(Q)===
	   2019/05/04 15:38:26 Target: s s hj.call p p(R,Q)
	   2019/05/04 15:38:26 Target: s s hP.property p p=M;P.propertyOf=function(Q)
	   2019/05/04 15:38:26 Target: s s hP.now p p=Date.now
	   2019/05/04 15:38:26 Target: s s hR.call p p(Q):R
	   2019/05/04 15:38:26 Target: s s hj.call p p(arguments,'');
	   2019/05/04 15:38:26 Target: s s hS.call p p(this,aa,P)
	   2019/05/04 15:38:26 Target: s s h.call p p(new
	*/
}

func BenchmarkSearch(b *testing.B) {

	proxyParams := config1("abc.com")

	body := []byte(`aaa s%2F%2Fs \u003A http endhttps, %3a//, http%3a//google.com:221/dsddsd?Sdad?DA/dsadas/#[]=1221 "http://dupa.com" :\\u002f\\u002f, \\u00253A\\u00252F\\u00252F dsda dsada http%3a//google.com .dsds ds d`)
	//f, err := os.Create("bench.out")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//err = trace.Start(f)
	//if err != nil {
	//	panic(err)
	//}
	for n := 0; n < b.N; n++ {
		scan(body, nil, true, proxyParams)
	}
	//trace.Stop()
	//b.StopTimer()
	//// go tool trace bench.out
}

func BenchmarkReplace(b *testing.B) {
	proxyParams := config1("abc.com")

	//body := []byte(`aaa s%2F%2Fs \u003A http endhttps, %3a//, http%3a//google.com:221/dsddsd?Sdad?DA/dsadas/#[]=1221 "http://dupa.com" :\\u002f\\u002f, \\u00253A\\u00252F\\u00252F dsda dsada http%3a//google.com .dsds ds d`)
	body := getFileBytes()
	for n := 0; n < b.N; n++ {
		scan(body, ReplaceProxyToTarget, true, proxyParams)
	}
}

func BenchmarkRange(b *testing.B) {
	body := []byte(`aaas%2F%2Fs\u003Ahttpendhttps,%3a//, http%3a//google.com:221/dsddsd?Sdad?DA/dsadas/#[]=1221 "http://dupa.com" :\\u002f\\u002f, \\u00253A\\u00252F\\u00252F dsda dsada http%3a//google.com .dsds ds d`)
	s := scanner{body: body}

	for n := 0; n < b.N; n++ {
		s.getRange(pathAndQueryAllowedChars)
	}
}

func getFileBytes() []byte {
	file, err := os.Open("/Users/tomek/Dysk Google/Projects/Proxy/files/_browse_SDJ-442423.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return b
}
