package scanner

//var beforePatterns = [...]string{` `, `'`, `"`}
var schemaPatterns = [][]byte{
	[]byte(`https`), []byte(`http`),
	//[]byte(`wss`), []byte(`wss`)
}

var slashSlashPatterns = [][]byte{
	[]byte(`//`),
	[]byte(`\/\/`),
	[]byte(`\\/\\/`),
	[]byte(`\\u002f\\u002f`),
	[]byte(`\u002f\u002f`),
	[]byte(`%2f%2f`),
	[]byte(`&#x2f;&#x2f;`),
	//[]byte(`%2F%2F`), ---> co to jest?
	[]byte(`\\u00252f\\u00252f`),
	[]byte(`\u00252f\u00252f`),
}

var colonPatterns = [][]byte{[]byte(`:`), []byte(`%3A`), []byte(`\u00253a`), []byte(`\\u00253a`)}

var dotNonStandardPatterns = [][]byte{[]byte(`\uff61'`), []byte(`\uff0e'`), []byte(`\uff61'`), []byte(`\\u00252e`), []byte(`\u00252e`)}

var hostAllowedChars = []byte(`abcdefghijklmnopqrstuvwxyz-.1234567890`)

// https://tools.ietf.org/html/rfc3986#section-2.2
var portAllowedChars = []byte(`1234567890`)

// https://tools.ietf.org/html/rfc3986
// i dont like ' sign in this
var pathAndQueryAllowedChars = []byte(`abcdefghijklmnopqrstuvwxyz1234567890:/?#[]@!$&'()*+,;=-._~^\`)
