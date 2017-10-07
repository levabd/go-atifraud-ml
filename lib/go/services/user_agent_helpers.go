package services

import (
	"os"
	"fmt"
	"os/exec"
	"strings"
	"path/filepath"
	"log"
	"github.com/joho/godotenv"
	"github.com/buger/jsonparser"
)

// private method
//noinspection GoUnusedFunction
func uaIsBot(browserFamily string) bool {
	switch browserFamily {
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
func IsCrawlerByUdger(clientIp string, clientUa string) bool {
	cmd := exec.Command("python3", filepath.Join(os.Getenv("APP_ROOT_DIR"),"lib", "python", "isCrawler.py"), clientIp, clientUa)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("error python exec: ", err)
		os.Exit(-1)
	}

	return strings.TrimRight(string(out) , "\n")  == "True"
}

// public method - search by go user agent parser
//noinspection GoUnusedExportedFunction
func PyIsCrawler(ip string, ua string) bool {
	var (
		result bool = false
	)

	if isCrawlerUdgerTold := IsCrawlerByUdger(ip, ua); isCrawlerUdgerTold {
		result = true
	}
	return result
}

//noinspection GoUnusedExportedFunction
func PyGetUa(clientUa string) map[string]interface{} {
	cmd := 	exec.Command("python3", filepath.Join(os.Getenv("APP_ROOT_DIR"),"lib", "python", "getUa.py"), clientUa)
	cmd.Wait()
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("parsers.GetUaerror: after Wait", err)
		os.Exit(-1)
	}

	jsonToParse := strings.Replace(string(out), " ", "", -1)
	jsonToParse = strings.Replace(string(out), "'", "\"", -1)
	jsonToParse = strings.TrimPrefix(strings.TrimSuffix(jsonToParse, "'"), "'")

	data := []byte(jsonToParse)
	i := 0
	mapToReturn := make(map[string]interface{})

	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		mapToReturn[ string(key)] = string(value)
		i = i + 1
		return nil
	})

	return mapToReturn
}

func init() {
	path, err := filepath.Abs(filepath.Join("..", ".env"))
	if err != nil {
		log.Fatal(err)
	}
	godotenv.Load(path)
}
