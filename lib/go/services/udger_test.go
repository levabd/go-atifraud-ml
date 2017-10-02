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
	_assert := assert.New(t)

	u, err := GetUdgerInstance()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(-1)
	}

	_assert.Equal(607, len(u.Browsers), "main_table be 607 in len")
	_assert.Equal(154, len(u.OS), "value_table be 154 in len")
	_assert.Equal(8, len(u.Devices), "ordered_table be 7 in len")
}

func TestGetIpClassificationCode(t *testing.T)  {
	assert.Equal(t, "crawler", GetIpClassificationCode("40.77.167.95"))
}

func TestIsCrawler(t *testing.T)  {
	is_crawler := IsCrawler("40.77.167.95", "ua")
	assert.Equal(t, true, is_crawler)

	is_crawler = IsCrawler("62.84.44.222", "ua")
	assert.Equal(t, false, is_crawler)

	is_crawler = IsCrawler("62.84.44.222", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	assert.Equal(t, true, is_crawler)

	is_crawler = IsCrawler("40.77.167.95", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	assert.Equal(t, true, is_crawler)
}

func TestIsCrawlerSql(t *testing.T) {

	ua_class_code,ua_family_code:= IsCrawlerSql("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	assert.Equal(t, "search_engine_bot", ua_class_code)
	assert.Equal(t, "googlebot", ua_family_code)

	ua_class_code, ua_family_code = IsCrawlerSql("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	fmt.Println(fmt.Sprintf("ua_class_code %s, ua_family_code %s", ua_class_code, ua_family_code))

	assert.Equal(t, "", ua_class_code)
	assert.Equal(t, "", ua_family_code)
}