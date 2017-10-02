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
		path_to_udger_db := filepath.Join(os.Getenv("APP_ROOT_DIR"), "data", "db", "udgerdb_v3.dat")
		println("path_to_udger_db: ", path_to_udger_db)

		u, err := udger.New(path_to_udger_db)
		if err != nil {
			Logger.Fatalln(err)
			return nil, err;
		}
		instantiated = u;
	}
	return instantiated, nil;
}

func IsCrawler(client_ip string, client_ua string) bool {

	ua_class_code,ua_family_code:= IsCrawlerSql(client_ua)
	is_crawler_by_ua := (IsInBotsUaFamily(strings.ToLower(ua_family_code)) || IsInClassCode(strings.ToLower(ua_class_code)))

	if is_crawler_by_ua || GetIpClassificationCode(client_ip) == "crawler"  {
		return true
	}

	return false
}

func GetIpClassificationCode(client_ip string) string {

	db, err := gorm.Open("sqlite3", os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		Logger.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
	}
	defer db.Close()
	var ip_classification_code string
	row:= db.Raw(fmt.Sprintf(`
	SELECT ip_classification_code
	FROM udger_ip_list
	JOIN udger_ip_class ON udger_ip_class.id=udger_ip_list.class_id
	LEFT JOIN udger_crawler_list ON udger_crawler_list.id=udger_ip_list.crawler_id
	LEFT JOIN udger_crawler_class ON udger_crawler_class.id=udger_crawler_list.class_id
	WHERE ip = '%s' ORDER BY sequence`, client_ip)).Row()

	row.Scan(&ip_classification_code)
	return ip_classification_code
}

func IsCrawlerSql(ua string) (string, string){

	db, err := gorm.Open("sqlite3", os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		Logger.Fatalf("parse_gz_logs.go - main: Failed to connect database: %s", err)
	}
	defer db.Close()

	crawler_classification_code :=""
	family_code :=""

	row:= db.Raw(fmt.Sprintf(`
	SELECT
	 crawler_classification_code, family_code
	FROM
	  udger_crawler_list
	  LEFT JOIN
	  udger_crawler_class ON udger_crawler_class.id = udger_crawler_list.class_id
	WHERE
	  ua_string = '%s' `, ua)).Row()

	row.Scan(&crawler_classification_code, &family_code)

	return crawler_classification_code, family_code
}

func GetUa( client_ua string) map[string]interface{} {
	udger, err := GetUdgerInstance()
	if err != nil{
		panic(err)
	}

	ua, err:=udger.Lookup(client_ua)
	if err != nil{
		panic(err)
	}

	main_data :=make(map[string]interface{})
	main_data["ua_family_code"] = strings.ToLower(ua.Browser.Family)
	main_data["ua_class_code"] = strings.ToLower(ua.Browser.Type)
	main_data["ua_version"] = strings.ToLower(ua.Browser.Version)
	main_data["os_family_code"] = strings.ToLower(ua.OS.Family)
	main_data["os_code"] = strings.Replace(strings.ToLower(ua.OS.Name), " ", "_", -1)

	if ua.Device.Name =="Personal computer"{
		main_data["device_class_code"] = "desktop"
	} else{
		main_data["device_class_code"] = strings.Replace(strings.ToLower( ua.Device.Name), " ", "_", -1)
	}

	return main_data
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
