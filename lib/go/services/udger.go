package services

import (
	"github.com/udger/udger"
	"path/filepath"
	"os"
	"github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"fmt"
	"strings"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func init() {
	err := helpers.LoadEnv()
	if err != nil {
		Logger.Fatalln(err)
	}
}

var instantiated *udger.Udger = nil

func GetUdgerInstance() (*udger.Udger, error) {
	if instantiated == nil {
		pathToUdgerDb := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "db", "udgerdb_v3.dat")
		log.Println("path_to_udger_db: ", pathToUdgerDb)

		u, err := udger.New(pathToUdgerDb)
		if err != nil {
			Logger.Fatalln(err)
			return nil, err
		}
		instantiated = u
	}
	return instantiated, nil
}

func IsCrawler(clientIp string, clientUa string) bool {

	uaClassCode, uaFamilyCode := IsCrawlerSql(clientUa)

	isBotByUaString := UaContainsCrawler(clientUa)
	isCrawlerByUa := IsInBotsUaFamily(strings.ToLower(uaFamilyCode)) || IsInClassCode(strings.ToLower(uaClassCode))

	if isBotByUaString || isCrawlerByUa || GetIpClassificationCode(clientIp) == "crawler"  {
		return true
	}

	return false
}

func GetIpClassificationCode(clientIp string) string {

	db, err := gorm.Open("sqlite3", os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		Logger.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
	}
	defer db.Close()
	var ipClassificationCode string
	row := db.Raw(fmt.Sprintf(`
	SELECT ip_classification_code
	FROM udger_ip_list
	JOIN udger_ip_class ON udger_ip_class.id=udger_ip_list.class_id
	LEFT JOIN udger_crawler_list ON udger_crawler_list.id=udger_ip_list.crawler_id
	LEFT JOIN udger_crawler_class ON udger_crawler_class.id=udger_crawler_list.class_id
	WHERE ip = '%s' ORDER BY sequence`, clientIp)).Row()
	row.Scan(&ipClassificationCode)

	return ipClassificationCode
}

func IsCrawlerSql(ua string) (string, string){

	db, err := gorm.Open("sqlite3", os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		Logger.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
	}
	defer db.Close()

	crawlerClassificationCode :=""
	familyCode :=""

	row:= db.Raw(fmt.Sprintf(`
	SELECT
	 crawler_classification_code, family_code,
	FROM
	  udger_crawler_list
	  LEFT JOIN
	  udger_crawler_class ON udger_crawler_class.id = udger_crawler_list.class_id
	WHERE
	  ua_string = '%s' `, ua)).Row()

	row.Scan(&crawlerClassificationCode, &familyCode)

	return crawlerClassificationCode, familyCode
}

func GetUa( clientUa string) map[string]interface{} {
	udg, err := GetUdgerInstance()
	if err != nil{
		panic(err)
	}

	ua, err := udg.Lookup(clientUa)
	if err != nil{
		panic(err)
	}

	mainData := make(map[string]interface{})
	mainData["ua_family_code"] = strings.ToLower(ua.Browser.Family)
	mainData["ua_class_code"] = strings.ToLower(ua.Browser.Type)
	mainData["ua_version"] = strings.ToLower(ua.Browser.Version)
	mainData["os_family_code"] = strings.ToLower(ua.OS.Family)
	mainData["os_code"] = strings.Replace(strings.ToLower(ua.OS.Name), " ", "_", -1)

	if ua.Device.Name =="Personal computer"{
		mainData["device_class_code"] = "desktop"
	} else{
		mainData["device_class_code"] = strings.Replace(strings.ToLower( ua.Device.Name), " ", "_", -1)
	}

	return mainData
}

func IsInClassCode(category string) bool {
	switch category {
	case
		"crawler":
		return true
	}
	return false
}

func IsInBotsUaFamily(category string) bool {
	switch category {
	case
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
		"easybib_autocite",
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

		"elinks",
		"links",
		"lynx":
		return true
	}
	return false
}



var crawlerWords = []string{
	"CodeGator",
	"spbot",
	"Barkrowler",
	"HybridBot",
	"MoodleBot",
	"www.ru",
	"Java",
	"Googlebot",
	"ia_archiver",
	"Mediapartners-Google",
	"OpenLinkProfiler.org/bot",
	"www.proximic.com/info/spider",
	"top100.rambler.ru crawler",
	"YandexBot",
	"Bot",
	"bot",
	"crawler",
	"Crawler",
	"Magic Browser",
	"Microsoft Office Protocol Discovery",
	"Microsoft Office Word 2014",
	"Spider",
	"spider",
	"scraper",
	"Scraper",
}

func UaContainsCrawler(s string)bool  {
	for _, b := range crawlerWords {
		if b==s || strings.Contains( s, b) || strings.HasPrefix( s, b)|| strings.HasSuffix( s, b){
			//log.Println("Contain crawler: ", s)
			return true
		}
	}
	return false
}


