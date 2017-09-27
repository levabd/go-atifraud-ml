package parsers

import (
	"os"
	"fmt"
	"os/exec"
	"github.com/buger/jsonparser"
	"strings"
)

const (
	udgerDBPath = "lib/parsers/data/udgerdb_v3.dat"
)

// private method
func uaIsBot(brawser_family string) bool {
	switch brawser_family {
	case
		// Search engine or antivirus or SEO bots
		"googlebot",
		"siteexplorer",
		"sputnikbot",
		"bingbot",
		"mj12bot",
		"yandexbot",
		"cliqzbot",
		"avast_safezone",
		"megaindex",
		"genieo_web_filter",
		"uptimebot",
		"ahrefsbot",
		"wordpress_pingback",
		"admantx_platform_semantic_analyzer",
		"leikibot",
		"mnogosearch",
		"safednsbot",
		"easybib_autocite",
		"sogou_spider",
		"surveybot",
		"baiduspider",
		"indy_library",
		"mail-ru_bot",
		"pocketparser",
		"virustotal",
		"feedfetcher_google",
		"virusdie_crawler",
		"surdotlybot",
		"yoozbot",
		"facebookbot",
		"linkdexbot",
		"prlog",
		"thinglink_imagebot",
		"obot",
		"spyonweb",
		"avira_crawler",
		"pulsepoint_xt3_web_scraper",
		"comodospider",
		"girafabot",
		"avira_scout",
		"salesintelligent",
		"kaspersky_bot",
		"xenu",
		"maxpointcrawler",
		"seznambot",
		"magpie-crawler",
		"yesupbot",
		"startmebot",
		"brandprotect_bot",
		"ask_jeeves-teoma",
		"duckduckgo_app",
		"linqiabot",
		"flipboardbot",
		"cat_explorador",
		"huaweisymantecspider",
		"coccocbot",
		"powermarks",
		"prism",
		"leechcraft",
		"wkhtmltopdf",

		// I think next is potencial bad bot (framework for apps or bad crowler)
		"java",
		"www::mechanize",
		"grapeshotcrawler",
		"netestate_crawler",
		"ccbot",
		"plukkie",
		"metauri",
		"silk",
		"phantomjs",
		"python-requests",
		"okhttp",
		"python-urllib",
		"netcraft_crawler",
		"go_http_package",
		"google_app",
		"android_httpurlconnection",
		"curl",
		"w3m",
		"wget",
		"getintentcrawler",
		"scrapy",
		"crawler4j",
		"apache-httpclient",
		"feedparser",
		"php",
		"simplepie",
		"lwp::simple",
		"libwww-perl",
		"apache_synapse",
		"scrapy_redis",
		"winhttp",
		"johnhew_crawler",
		"poe-component-client-http",
		"joc_web_spider",

		//Text Browsers
		"elinks",
		"links",
		"lynx":
		return true
	}
	return false
}

// private method - search by udger
func IsCrawlerByUdger(client_ip string, client_ua string) bool {
	ua, err := GetUdgerInstance(udgerDBPath).Lookup(client_ua)
	if err != nil {
		fmt.Println("error GetUdgerInstance: ", err)
		os.Exit(-1)
	}

	if uaIsBot(ua.Browser.Family) {
		return true
	}

	cmd := exec.Command("python", "lib/parsers/python/isCrawler.py", client_ip, client_ua)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error python exec: ", err)
		os.Exit(-1)
	}

	return string(out) == "True"
}

// public method - search by go user agent parser
func UaIsCrawler(ip string, ua string) bool {
	var (
		result bool = false
	)

	if is_crawler_udger_told := IsCrawlerByUdger(ip, ua); is_crawler_udger_told {
		result = true
	}
	fmt.Println("UaIsCrawler: ", result)

	return result
}

func GetUa(client_ua string) map[string]string {
	cmd := exec.Command("python", "lib/parsers/python/getUa.py", client_ua)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("parsers.GetUaerror: ", err)
		os.Exit(-1)
	}

	json_to_parse := strings.Replace(string(out), " ", "", -1)
	json_to_parse = strings.TrimPrefix(strings.TrimSuffix(json_to_parse, "'"), "'")

	data := []byte(json_to_parse)
	i := 0
	map_ro_return := make(map[string]string)
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		map_ro_return[ string(key)] = string(value)
		i = i + 1
		return nil
	})

	return map_ro_return
}
