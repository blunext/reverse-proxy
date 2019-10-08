package main

import (
	"Proxy/tools"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSomething(t *testing.T) {

	//// assert equality
	//in := []byte("bdsabdasb saiojd sia " + targetScheme + "://" + targetHost + " dsdsds")
	//out := []byte("bdsabdasb saiojd sia " + proxyScheme + "://" + proxyHost + proxyPort + " dsdsds")
	//
	//assert.Equal(t, string(out), string(ReplaceTargetToProxy(in)), "niepoprawna zamiana adresów: replaceTargetToProxy")
	//assert.Equal(t, string(in), string(ReplaceProxyToTarget(out)), "niepoprawna zamiana adresów: replaceProxyToTarget")

	assert.Equal(t, true, tools.TextContentType("text/plain"), "nie rozpoznany content type")
	//assert.Equal(t, false, textContentType("text/plain;ds"), "nie rozpoznany content type")
	//assert.Equal(t, true, textContentType("application/xhtml+xml"), "nie rozpoznany content type")

}
