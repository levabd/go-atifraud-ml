package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"testing"
	"fmt"
	"os"
)

func init() {
	helpers.LoadEnv()
}

func TestUdgerInitiation(t *testing.T) {
	localAssert := assert.New(t)

	u, err := GetUdgerInstance()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(-1)
	}

	localAssert.Equal(607, len(u.Browsers), "main_table be 607 in len")
	localAssert.Equal(154, len(u.OS), "value_table be 154 in len")
	localAssert.Equal(8, len(u.Devices), "ordered_table be 7 in len")
}

func TestGetIpClassificationCode(t *testing.T)  {
	assert.Equal(t, "crawler", GetIpClassificationCode("40.77.167.95"))
}

func TestIsCrawler(t *testing.T)  {
	isCrawler := IsCrawler("40.77.167.95", "ua")
	assert.Equal(t, true, isCrawler)

	isCrawler = IsCrawler("62.84.44.222", "ua")
	assert.Equal(t, false, isCrawler)

	isCrawler = IsCrawler("62.84.44.222", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	assert.Equal(t, true, isCrawler)

	isCrawler = IsCrawler("40.77.167.95", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	assert.Equal(t, true, isCrawler)
}

func TestIsCrawlerSql(t *testing.T) {

	uaClassCode, uaFamilyCode:= IsCrawlerSql("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	assert.Equal(t, "search_engine_bot", uaClassCode)
	assert.Equal(t, "googlebot", uaFamilyCode)

	uaClassCode, uaFamilyCode = IsCrawlerSql("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	fmt.Println(fmt.Sprintf("ua_class_code %s, ua_family_code %s", uaClassCode, uaFamilyCode))

	assert.Equal(t, "", uaClassCode)
	assert.Equal(t, "", uaFamilyCode)
}