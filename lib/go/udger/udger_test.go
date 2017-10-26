package udger

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"log"
)

var udger *Udger

func init() {
	_udger,err := New("/home/vagrant/go/src/github.com/levabd/go-atifraud-ml/data/db/udgerdb_v3.dat")
	if err!=nil{
		log.Fatal(err)
	}

	udger= _udger
}

//func TestParseUa(t *testing.T) {
//	assert := assert.New(t)
//	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36";
//	data , _:= udger.ParseUa(userAgent)
//
//	assert.Equal("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36", data["user_agent"]["ua_string"])
//	assert.Equal("Browser", data["user_agent"]["ua_class"])
//	assert.Equal("browser", data["user_agent"]["ua_class_code"])
//	assert.Equal("Chrome 39.0.2171.71", data["user_agent"]["ua"])
//	assert.Equal("Chrome 39.0.2171.71", data["user_agent"]["ua"])
//	assert.Equal("39.0.2171.71", data["user_agent"]["ua_version"])
//	assert.Equal("39", data["user_agent"]["ua_version_major"])
//	assert.Equal("Chrome", data["user_agent"]["ua_family"])
//	assert.Equal("chrome", data["user_agent"]["ua_family_code"])
//	assert.Equal("http://www.google.com/chrome/", data["user_agent"]["ua_family_homepage"])
//	assert.Equal("Google Inc.", data["user_agent"]["ua_family_vendor"])
//	assert.Equal("google_inc", data["user_agent"]["ua_family_vendor_code"])
//	assert.Equal("https://www.google.com/about/company/", data["user_agent"]["ua_family_vendor_homepage"])
//	assert.Equal("chrome.png", data["user_agent"]["ua_family_icon"])
//	assert.Equal("chrome_big.png", data["user_agent"]["ua_family_icon_big"])
//	assert.Equal("https://udger.com/resources/ua-list/browser-detail?browser=Chrome", data["user_agent"]["ua_family_info_url"])
//	assert.Equal("WebKit/Blink", data["user_agent"]["ua_engine"])
//
//	assert.Equal("OS X 10.9 Mavericks", data["user_agent"]["os"])
//	assert.Equal("osx_10_9", data["user_agent"]["os_code"])
//	assert.Equal("https://en.wikipedia.org/wiki/OS_X_Mavericks", data["user_agent"]["os_homepage"])
//	assert.Equal("macosx.png", data["user_agent"]["os_icon"])
//	assert.Equal("macosx_big.png", data["user_agent"]["os_icon_big"])
//	assert.Equal("https://udger.com/resources/ua-list/client-detail?client=OS X 10.9 Mavericks", data["user_agent"]["os_info_url"])
//	assert.Equal("OS X", data["user_agent"]["os_family"])
//	assert.Equal("osx", data["user_agent"]["os_family_code"])
//	assert.Equal("Apple Computer, Inc.", data["user_agent"]["os_family_vendor"])
//	assert.Equal("apple_inc", data["user_agent"]["os_family_vendor_code"])
//	assert.Equal("http://www.apple.com/", data["user_agent"]["os_family_vendor_homepage"])
//
//	// test crawlers
//	data , _ = udger.ParseUa("BabalooSpider/1.3 (BabalooSpider; http://www.babaloo.si; spider@babaloo.si)")
//	assert.Equal("crawler", data["user_agent"]["ua_class_code"])
//	assert.Equal("uncategorised", data["user_agent"]["crawler_category_code"])
//	assert.Equal("unknown", data["user_agent"]["crawler_respect_robotstxt"])
//	assert.Equal("2010-05-31 08:30:59", data["user_agent"]["crawler_last_seen"])
//}
//
//func TestParseIp(t *testing.T) {
//
//	assert := assert.New(t)
//
//	data := udger.ParseIp(`69.162.124.227`)
//
//	assert.Equal("69.162.124.227", data["ip_address"]["ip"])
//	assert.Equal("v4", data["ip_address"]["ip_ver"])
//	assert.Equal("crawler", data["ip_address"]["ip_classification_code"])
//}

func TestIsCrawler(t *testing.T) {

	assert := assert.New(t)

	assert.Equal(true, udger.IsCrawler( `69.162.124.227`, `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36`, false))
	assert.Equal(false, udger.IsCrawler(`62.84.44.222`, `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36`, false))
	assert.Equal(true, udger.IsCrawler(`52.34.47.242`, `Abrave v5.5 (http://robot.abrave.co.uk)`, true))
	assert.Equal(true, udger.IsCrawler(`52.34.47.242`, `Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)`, true))
}
