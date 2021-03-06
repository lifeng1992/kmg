package kmgHttp

import (
	"testing"

	. "github.com/bronze1man/kmg/kmgTest"
)

func TestAddParameterToUrl(ot *testing.T) {
	u, err := AddParameterToUrl("http://foo.com/", "a", "b")
	Equal(err, nil)
	Equal(u, "http://foo.com/?a=b")

	u, err = AddParameterToUrl("/?n=a&n1=b&n2=c&a=c", "a", "b")
	Equal(err, nil)
	Equal(u, "/?a=c&a=b&n=a&n1=b&n2=c")
}

func TestSetParameterToUrl(ot *testing.T) {
	u, err := SetParameterToUrl("http://foo.com/", "a", "b")
	Equal(err, nil)
	Equal(u, "http://foo.com/?a=b")

	u, err = SetParameterToUrl("/?n=a&n1=b&n2=c&a=c&b=d", "a", "b")
	Equal(err, nil)
	Equal(u, "/?a=b&b=d&n=a&n1=b&n2=c")
}

/*
import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestNewUrlByString(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	url, err := NewUrlByString("http://www.google.com")
	t.Equal(nil, err)
	t.Equal("http://www.google.com", url.String())
}
*/

func TestGetDomainName(t *testing.T) {
	testCaseAList := []string{
		"http://abc.com",
		"https://abc.com",
		"https://abc.com/",
		"http://abc.com/index.html",
		"file://abc.com/upload/index.css",
		"ftp://abc.com/upload/index.css",
	}
	testCaseBList := []string{
		"abc.com/upload/index.css",
		"/upload/index.css",
		"upload/index.css",
	}
	for index, url := range testCaseAList {
		domain, protocol := GetDomainName(url)
		Equal(domain, "abc.com")
		if index == 2 {
			Equal(protocol, "https://")
		}
		if index == 3 {
			Equal(protocol, "http://")
		}
	}
	for _, url := range testCaseBList {
		domain, _ := GetDomainName(url)
		Equal(domain, "")
	}
}
