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
	"github.com/jinzhu/gorm"
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
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

func FitUserAgentCodes(userAgentList []string) (map[string]int, map[string]float64) {
	var (
		userAgentIntCodes	= map[string]int {}
		userAgentFloatCodes	= map[string]float64 {}
		uaCodes				= []m.UACode {}
	)

	index := 0

	for _, ua := range userAgentList {
		if _, ok := userAgentIntCodes[ua]; !ok {
			userAgentIntCodes[ua] = index
			userAgentFloatCodes[ua] = float64(index)
			uaCodes = append(uaCodes, m.UACode{
				UserAgent:	ua,
				IntCode:	index,
				FloatCode:	float64(index),
			})
			index++
		}
	}

	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		Logger.Fatalf("user_agent_helpers.go - FitUserAgentCodes: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.UACode{}) {
		db.AutoMigrate(&m.UACode{})
	}

	// Clean last vectoriser
	db.Exec("TRUNCATE TABLE  ua_codes;")

	// Insert new fitted vectoriser
	for _, uaCode := range uaCodes {
		tx := db.Begin()
		tx.Create(&uaCode)
		tx.Commit()
	}

	return userAgentIntCodes, userAgentFloatCodes
}

func LoadFittedUserAgentCodes() (map[string]int, map[string]float64) {
	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		Logger.Fatalf("user_agent_helpers.go - LoadFittedUserAgentCodes: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.UACode{}) {
		db.AutoMigrate(&m.UACode{})
	}

	var (
		userAgentIntCodes	= map[string]int {}
		userAgentFloatCodes	= map[string]float64 {}
		uaCodes				= []m.UACode {}
	)

	db.Find(&uaCodes)
	for _, uaCode := range uaCodes {
		userAgentIntCodes[uaCode.UserAgent]		= uaCode.IntCode
		userAgentFloatCodes[uaCode.UserAgent]	= uaCode.FloatCode
	}

	return userAgentIntCodes, userAgentFloatCodes
}

func GetUAClasses(userAgentList []string, userAgentIntCodes map[string]int, userAgentFloatCodes map[string]float64) ([]int, []float64) {

	var (
		intUAClasses	[]int
		floatUAClasses	[]float64
	)

	for _, userAgent := range userAgentList {
		intUAClasses	= append(intUAClasses, userAgentIntCodes[userAgent])
		floatUAClasses	= append(floatUAClasses, userAgentFloatCodes[userAgent])
	}

	return intUAClasses, floatUAClasses
}

//noinspection GoUnusedExportedFunction
func GetSingleUAClass(userAgent string, userAgentIntCodes map[string]int, userAgentFloatCodes map[string]float64) (int, float64) {

	intUAClass		:= userAgentIntCodes[userAgent]
	floatUAClass	:= userAgentFloatCodes[userAgent]

	return intUAClass, floatUAClass
}

func init() {
	path, err := filepath.Abs(filepath.Join("..", ".env"))
	if err != nil {
		log.Fatal(err)
	}
	godotenv.Load(path)
}
