package udger

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"strings"
	_ "github.com/mattn/go-sqlite3"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"errors"
	"os"
	"net"
	"time"
	"log"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (udger *Udger) GetNewParsedData() (parsedData map[string]map[string]string){
	parsedData = make(map[string]map[string]string)

	parsedData["user_agent"] = make(map[string]string)
	parsedData["user_agent"]["ua_string"] = ""
	parsedData["user_agent"]["ua_class"] = ""
	parsedData["user_agent"]["ua_class_code"] = ""
	parsedData["user_agent"]["ua"] = ""
	parsedData["user_agent"]["ua_version"] = ""
	parsedData["user_agent"]["ua_version_major"] = ""
	parsedData["user_agent"]["ua_uptodate_current_version"] = ""
	parsedData["user_agent"]["ua_family"] = ""
	parsedData["user_agent"]["ua_family_code"] = ""
	parsedData["user_agent"]["ua_family_homepage"] = ""
	parsedData["user_agent"]["ua_family_vendor"] = ""
	parsedData["user_agent"]["ua_family_vendor_code"] = ""
	parsedData["user_agent"]["ua_family_vendor_homepage"] = ""
	parsedData["user_agent"]["ua_family_icon"] = ""
	parsedData["user_agent"]["ua_family_icon_big"] = ""
	parsedData["user_agent"]["ua_family_info_url"] = ""
	parsedData["user_agent"]["ua_engine"] = ""
	parsedData["user_agent"]["os"] = ""
	parsedData["user_agent"]["os_code"] = ""
	parsedData["user_agent"]["os_homepage"] = ""
	parsedData["user_agent"]["os_icon"] = ""
	parsedData["user_agent"]["os_icon_big"] = ""
	parsedData["user_agent"]["os_info_url"] = ""
	parsedData["user_agent"]["os_family"] = ""
	parsedData["user_agent"]["os_family_code"] = ""
	parsedData["user_agent"]["os_family_vendor"] = ""
	parsedData["user_agent"]["os_family_vendor_code"] = ""
	parsedData["user_agent"]["os_family_vendor_homepage"] = ""
	parsedData["user_agent"]["device_class"] = ""
	parsedData["user_agent"]["device_class_code"] = ""
	parsedData["user_agent"]["device_class_icon"] = ""
	parsedData["user_agent"]["device_class_icon_big"] = ""
	parsedData["user_agent"]["device_class_info_url"] = ""
	parsedData["user_agent"]["device_marketname"] = ""
	parsedData["user_agent"]["device_brand"] = ""
	parsedData["user_agent"]["device_brand_code"] = ""
	parsedData["user_agent"]["device_brand_homepage"] = ""
	parsedData["user_agent"]["device_brand_icon"] = ""
	parsedData["user_agent"]["device_class_code"] = ""
	parsedData["user_agent"]["device_brand_icon_big"] = ""
	parsedData["user_agent"]["device_brand_info_url"] = ""
	parsedData["user_agent"]["crawler_last_seen"] = ""
	parsedData["user_agent"]["crawler_category"] = ""
	parsedData["user_agent"]["crawler_category_code"] = ""
	parsedData["user_agent"]["crawler_respect_robotstxt"] = ""
	parsedData["ip_address"] = make(map[string]string)
	parsedData["ip_address"]["ip"] = ""
	parsedData["ip_address"]["ip_ver"] = ""
	parsedData["ip_address"]["ip_classification"] = ""
	parsedData["ip_address"]["ip_classification_code"] = ""
	parsedData["ip_address"]["ip_hostname"] = ""
	parsedData["ip_address"]["ip_last_seen"] = ""
	parsedData["ip_address"]["ip_country"] = ""
	parsedData["ip_address"]["ip_country_code"] = ""
	parsedData["ip_address"]["ip_city"] = ""
	parsedData["ip_address"]["crawler_name"] = ""
	parsedData["ip_address"]["crawler_ver"] = ""
	parsedData["ip_address"]["crawler_ver_major"] = ""
	parsedData["ip_address"]["crawler_family"] = ""
	parsedData["ip_address"]["crawler_family_code"] = ""
	parsedData["ip_address"]["crawler_family_homepage"] = ""
	parsedData["ip_address"]["crawler_family_vendor"] = ""
	parsedData["ip_address"]["crawler_family_vendor_code"] = ""
	parsedData["ip_address"]["crawler_family_vendor_homepage"] = ""
	parsedData["ip_address"]["crawler_family_icon"] = ""
	parsedData["ip_address"]["crawler_family_info_url"] = ""
	parsedData["ip_address"]["crawler_last_seen"] = ""
	parsedData["ip_address"]["crawler_category"] = ""
	parsedData["ip_address"]["crawler_category_code"] = ""
	parsedData["ip_address"]["crawler_respect_robotstxt"] = ""
	parsedData["ip_address"]["datacenter_name"] = ""
	parsedData["ip_address"]["datacenter_name_code"] = ""
	parsedData["ip_address"]["datacenter_homepage"] = ""

	return
}

func (udger *Udger) init() error {
	defer timeTrack(time.Now(), "udger init")

	// init data keys
	udger.ParseData = udger.GetNewParsedData()

	// load dictionaries
	var clients []Client
	udger.DB.Raw(ClientsSql).Scan(&clients);
	udger.Clients = make(map[int64]Client)

	for _, client := range clients {
		var d rexData
		d.ID = client.ClientId
		d.Regex = udger.cleanRegex(client.Regstring)
		r, err := pcre.Compile(d.Regex, pcre.CASELESS)
		if err != nil {
			return errors.New(err.String())
		}
		d.RegexCompiled = r
		udger.rexClients = append(udger.rexClients, d)
		udger.Clients[client.ClientId] = client
	}

	var operationSystems []Os

	udger.DB.Raw(OSystemsSql).Scan(&operationSystems);
	udger.OS = make(map[int64]Os)

	for _, operationSystem := range operationSystems {
		var d rexData
		d.ID = operationSystem.OsId
		d.Regex = udger.cleanRegex(operationSystem.Regstring)
		r, err := pcre.Compile(d.Regex, pcre.CASELESS)
		if err != nil {
			return errors.New(err.String())
		}
		d.RegexCompiled = r
		udger.rexOS = append(udger.rexOS, d)
		udger.OS[operationSystem.OsId] = operationSystem
	}

	var ips []Ip
	udger.IPS = make(map[string]string)

	udger.DB.Raw(fmt.Sprintf(`SELECT
			  udger_ip_list.ip as ip,
			  ip_classification_code
			FROM udger_ip_list
			  JOIN udger_ip_class ON udger_ip_class.id = udger_ip_list.class_id
			  LEFT JOIN udger_crawler_list ON udger_crawler_list.id = udger_ip_list.crawler_id
			  LEFT JOIN udger_crawler_class ON udger_crawler_class.id = udger_crawler_list.class_id
			ORDER BY sequence`, )).Scan(&ips);
	for _, ip := range ips {
		udger.IPS[ip.Ip] = ip.IpClassificationCode
	}

	//println("ips", len(ips))
	//println("IPS", len(udger.IPS))

	var crawlers []Crawler
	udger.Crawlers = make(map[string]Crawler)

	// crawler
	udger.DB.Raw(fmt.Sprintf(`SELECT
			  udger_crawler_list.id AS botid,
			  name,
			  ver,
			  ver_major,
			  last_seen,
			  respect_robotstxt,
			  family,
			  family_code,
			  family_homepage,
			  family_icon,
			  vendor,
			  vendor_code,
			  vendor_homepage,
			  crawler_classification,
			  crawler_classification_code,
			  ua_string
			FROM udger_crawler_list
			  LEFT JOIN udger_crawler_class ON udger_crawler_class.id=udger_crawler_list.class_id`)).Scan(&crawlers);
	for _, crawler := range crawlers {
		udger.Crawlers[crawler.UaString] = crawler
	}

	return nil
}

func New(dbPath string) (*Udger, error) {
	u := &Udger{}

	var err error

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, err
	}
	u.DbPath = dbPath
	u.DB, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer u.DB.Close()

	err = u.init()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (udger *Udger) findDataWithVersion(ua string, data []rexData, withVersion bool) (idx int64, value string, err error) {
	defer func() {
		if r := recover(); r != nil {
			idx, value, err = udger.findData(ua, data, false)
		}
	}()

	idx, value, err = udger.findData(ua, data, withVersion)

	return idx, value, err
}

func (udger *Udger) cleanRegex(r string) string {
	if strings.HasSuffix(r, "/si") {
		r = r[:len(r)-3]
	}
	if strings.HasPrefix(r, "/") {
		r = r[1:]
	}

	return r
}

func (udger *Udger) findData(ua string, data []rexData, withVersion bool) (idx int64, value string, err error) {
	for i := 0; i < len(data); i++ {
		r := data[i].RegexCompiled
		matcher := r.MatcherString(ua, 0)
		if !matcher.MatchString(ua, 0) {
			continue
		}

		if withVersion && matcher.Present(1) {
			return data[i].ID, matcher.GroupString(1), nil
		}

		return data[i].ID, "", nil
	}

	return -1, "", nil
}

func (udger *Udger) ParseUa(ua string, parsedData map[string]map[string]string) (map[string]map[string]string, error) {

	if ua != "" {

		parsedData["user_agent"]["ua_string"] = ua;
		parsedData["user_agent"]["ua_class"] = "Unrecognized";
		parsedData["user_agent"]["ua_class_code"] = "unrecognized";

		if crawler, ok := udger.Crawlers[ua]; ok {

			//client_class_id = 99;
			parsedData["user_agent"]["ua_class"] = "Crawler";
			parsedData["user_agent"]["ua_class_code"] = "crawler";

			ver := fmt.Sprintf("%.2f", crawler.Ver)
			if ver != "" {
				parsedData["user_agent"]["ua"] = crawler.Name + " " + ver
				parsedData["user_agent"]["ua_version"] = fmt.Sprintf("%.2f", crawler.Ver);
				parsedData["user_agent"]["ua_version_major"] = string(strings.Split(string(ver), ".")[0])
			} else {
				parsedData["user_agent"]["ua"] = crawler.Name
				parsedData["user_agent"]["ua_version"] = ""
				parsedData["user_agent"]["ua_version_major"] = ""
			}

			parsedData["user_agent"]["ua_family"] = crawler.Family;
			parsedData["user_agent"]["ua_family_code"] = crawler.FamilyCode;
			parsedData["user_agent"]["ua_family_homepage"] = crawler.FamilyHomepage;
			parsedData["user_agent"]["ua_family_vendor"] = crawler.Vendor;
			parsedData["user_agent"]["ua_family_vendor_code"] = crawler.VendorCode;
			parsedData["user_agent"]["ua_family_vendor_homepage"] = crawler.VendorHomepage;
			parsedData["user_agent"]["ua_family_icon"] = crawler.FamilyIcon;
			parsedData["user_agent"]["ua_family_info_url"] = "https://udger.com/resources/ua-list/bot-detail?bot=" + crawler.Family + "#id" + string(crawler.Botid);
			parsedData["user_agent"]["crawler_last_seen"] = string(crawler.LastSeen);
			parsedData["user_agent"]["crawler_category"] = crawler.CrawlerClassification;
			parsedData["user_agent"]["crawler_category_code"] = crawler.CrawlerClassificationCode;
			parsedData["user_agent"]["crawler_respect_robotstxt"] = crawler.RespectRobotstxt;

		} else {

			// client
			clientID, version, err := udger.findDataWithVersion(ua, udger.rexClients, true)
			if err != nil {
				return nil, err
			}

			if client, ok := udger.Clients[clientID]; ok {
				parsedData["user_agent"]["ua_class"] = client.ClientClassification
				parsedData["user_agent"]["ua_class_code"] = client.ClientClassificationCode

				if version != "" {
					parsedData["user_agent"]["ua"] = client.Name + " " + version
					parsedData["user_agent"]["ua_version"] = version
					parsedData["user_agent"]["ua_version_major"] = string(strings.Split(string(version), ".")[0])
				} else {
					parsedData["user_agent"]["ua"] = client.Name
					parsedData["user_agent"]["ua_version"] = ""
					parsedData["user_agent"]["ua_version_major"] = ""
				}

				parsedData["user_agent"]["ua_uptodate_current_version"] = client.UptodateCurrentVersion
				parsedData["user_agent"]["ua_family"] = client.Name
				parsedData["user_agent"]["ua_family_code"] = client.NameCode
				parsedData["user_agent"]["ua_family_homepage"] = client.Homepage
				parsedData["user_agent"]["ua_family_vendor"] = client.Vendor
				parsedData["user_agent"]["ua_family_vendor_code"] = client.VendorCode
				parsedData["user_agent"]["ua_family_vendor_homepage"] = client.VendorHomepage
				parsedData["user_agent"]["ua_family_icon"] = client.Icon
				parsedData["user_agent"]["ua_family_icon_big"] = client.IconBig
				parsedData["user_agent"]["ua_family_info_url"] = "https://udger.com/resources/ua-list/browser-detail?browser=" + client.Name;
				parsedData["user_agent"]["ua_engine"] = client.Engine
			}
		}
	}

	return parsedData, nil
}

func (udger *Udger) ParseIp(ip string, parsedData map[string]map[string]string ) map[string]map[string]string {


	if ip == "" {
		return parsedData
	}

	parsedData["ip_address"]["ip"] = ip;
	ipver := getIpVersion(ip);
	if ipver == "v4" {
		parsedData["ip_address"]["ip_ver"] = ipver
		if IpClassificationCode, ok := udger.IPS[ip]; ok {
			parsedData["ip_address"]["ip_classification_code"] = IpClassificationCode
			return parsedData
		}
	}

	parsedData["ip_address"]["ip_classification_code"] = ""

	return parsedData
}

func getIpVersion(ipAddress string) string {

	testInput := net.ParseIP(ipAddress)

	if testInput.To4() != nil {
		return "v4"
	}

	if testInput.To16() != nil {
		return "v6"
	}

	return ""
}

func (udger *Udger) IsCrawler(ip string, ua string, parsedData  map[string]map[string]string) bool {
	udger.ParseUa(ua, parsedData)
	udger.ParseIp(ip, parsedData)

	if  udger.ParseData["ip_address"]["ip_classification_code"] == "crawler" ||
		( udger.ParseData["user_agent"]["ua_class_code"] == "crawler" ||
			IsInBotsUaFamily( udger.ParseData["user_agent"]["ua_family_code"]) ||
			UaContainsCrawler(ua)) {
		return true
	}

	return false
}

func IsInBotsUaFamily(category string) bool {
	switch category {
	case
		"begunadvertising",
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
		"lynx",
		"demandbase-bot",
		"dotbot",
		"duckduckbot",
		"exabot",
		"extlinksbot",
		"istellabot",
		"linkpadbot",
		"nettrack_info-bot",
		"seokicks-robot",
		"smtbot",
		"spbot",
		"telegrambot",
		"twitterbot":
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
	"crawler",
	"Crawler",
	"Magic Browser",
	"Microsoft Office Protocol Discovery",
	"Microsoft Office Word 2014",
	"Spider",
	"spider",
	"scraper",
	"Scraper",
	"semrushbot",
}

func UaContainsCrawler(s string) bool {
	for _, b := range crawlerWords {
		if b == s || strings.Contains(s, b) || strings.HasPrefix(s, b) || strings.HasSuffix(s, b) {
			return true
		}
	}
	return false
}
