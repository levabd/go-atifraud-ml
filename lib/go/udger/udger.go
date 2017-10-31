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
func (udger *Udger) init() error {
	defer timeTrack(time.Now(), "udger init")

	// init data keys
	udger.ParseData = make(map[string]map[string]string)
	udger.ParseData["user_agent"] = make(map[string]string)
	udger.ParseData["user_agent"]["ua_string"] = ""
	udger.ParseData["user_agent"]["ua_class"] = ""
	udger.ParseData["user_agent"]["ua_class_code"] = ""
	udger.ParseData["user_agent"]["ua"] = ""
	udger.ParseData["user_agent"]["ua_version"] = ""
	udger.ParseData["user_agent"]["ua_version_major"] = ""
	udger.ParseData["user_agent"]["ua_uptodate_current_version"] = ""
	udger.ParseData["user_agent"]["ua_family"] = ""
	udger.ParseData["user_agent"]["ua_family_code"] = ""
	udger.ParseData["user_agent"]["ua_family_homepage"] = ""
	udger.ParseData["user_agent"]["ua_family_vendor"] = ""
	udger.ParseData["user_agent"]["ua_family_vendor_code"] = ""
	udger.ParseData["user_agent"]["ua_family_vendor_homepage"] = ""
	udger.ParseData["user_agent"]["ua_family_icon"] = ""
	udger.ParseData["user_agent"]["ua_family_icon_big"] = ""
	udger.ParseData["user_agent"]["ua_family_info_url"] = ""
	udger.ParseData["user_agent"]["ua_engine"] = ""
	udger.ParseData["user_agent"]["os"] = ""
	udger.ParseData["user_agent"]["os_code"] = ""
	udger.ParseData["user_agent"]["os_homepage"] = ""
	udger.ParseData["user_agent"]["os_icon"] = ""
	udger.ParseData["user_agent"]["os_icon_big"] = ""
	udger.ParseData["user_agent"]["os_info_url"] = ""
	udger.ParseData["user_agent"]["os_family"] = ""
	udger.ParseData["user_agent"]["os_family_code"] = ""
	udger.ParseData["user_agent"]["os_family_vendor"] = ""
	udger.ParseData["user_agent"]["os_family_vendor_code"] = ""
	udger.ParseData["user_agent"]["os_family_vendor_homepage"] = ""
	udger.ParseData["user_agent"]["device_class"] = ""
	udger.ParseData["user_agent"]["device_class_code"] = ""
	udger.ParseData["user_agent"]["device_class_icon"] = ""
	udger.ParseData["user_agent"]["device_class_icon_big"] = ""
	udger.ParseData["user_agent"]["device_class_info_url"] = ""
	udger.ParseData["user_agent"]["device_marketname"] = ""
	udger.ParseData["user_agent"]["device_brand"] = ""
	udger.ParseData["user_agent"]["device_brand_code"] = ""
	udger.ParseData["user_agent"]["device_brand_homepage"] = ""
	udger.ParseData["user_agent"]["device_brand_icon"] = ""
	udger.ParseData["user_agent"]["device_class_code"] = ""
	udger.ParseData["user_agent"]["device_brand_icon_big"] = ""
	udger.ParseData["user_agent"]["device_brand_info_url"] = ""
	udger.ParseData["user_agent"]["crawler_last_seen"] = ""
	udger.ParseData["user_agent"]["crawler_category"] = ""
	udger.ParseData["user_agent"]["crawler_category_code"] = ""
	udger.ParseData["user_agent"]["crawler_respect_robotstxt"] = ""
	udger.ParseData["ip_address"] = make(map[string]string)
	udger.ParseData["ip_address"]["ip"] = ""
	udger.ParseData["ip_address"]["ip_ver"] = ""
	udger.ParseData["ip_address"]["ip_classification"] = ""
	udger.ParseData["ip_address"]["ip_classification_code"] = ""
	udger.ParseData["ip_address"]["ip_hostname"] = ""
	udger.ParseData["ip_address"]["ip_last_seen"] = ""
	udger.ParseData["ip_address"]["ip_country"] = ""
	udger.ParseData["ip_address"]["ip_country_code"] = ""
	udger.ParseData["ip_address"]["ip_city"] = ""
	udger.ParseData["ip_address"]["crawler_name"] = ""
	udger.ParseData["ip_address"]["crawler_ver"] = ""
	udger.ParseData["ip_address"]["crawler_ver_major"] = ""
	udger.ParseData["ip_address"]["crawler_family"] = ""
	udger.ParseData["ip_address"]["crawler_family_code"] = ""
	udger.ParseData["ip_address"]["crawler_family_homepage"] = ""
	udger.ParseData["ip_address"]["crawler_family_vendor"] = ""
	udger.ParseData["ip_address"]["crawler_family_vendor_code"] = ""
	udger.ParseData["ip_address"]["crawler_family_vendor_homepage"] = ""
	udger.ParseData["ip_address"]["crawler_family_icon"] = ""
	udger.ParseData["ip_address"]["crawler_family_info_url"] = ""
	udger.ParseData["ip_address"]["crawler_last_seen"] = ""
	udger.ParseData["ip_address"]["crawler_category"] = ""
	udger.ParseData["ip_address"]["crawler_category_code"] = ""
	udger.ParseData["ip_address"]["crawler_respect_robotstxt"] = ""
	udger.ParseData["ip_address"]["datacenter_name"] = ""
	udger.ParseData["ip_address"]["datacenter_name_code"] = ""
	udger.ParseData["ip_address"]["datacenter_homepage"] = ""

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

func (udger *Udger) ParseUa(ua string) (map[string]map[string]string, error) {

	udger.mux.Lock()

	if ua != "" {

		udger.ParseData["user_agent"]["ua_string"] = ua;
		udger.ParseData["user_agent"]["ua_class"] = "Unrecognized";
		udger.ParseData["user_agent"]["ua_class_code"] = "unrecognized";

		if crawler, ok := udger.Crawlers[ua]; ok {

			//client_class_id = 99;
			udger.ParseData["user_agent"]["ua_class"] = "Crawler";
			udger.ParseData["user_agent"]["ua_class_code"] = "crawler";

			ver := fmt.Sprintf("%.2f", crawler.Ver)
			if ver != "" {
				udger.ParseData["user_agent"]["ua"] = crawler.Name + " " + ver
				udger.ParseData["user_agent"]["ua_version"] = fmt.Sprintf("%.2f", crawler.Ver);
				udger.ParseData["user_agent"]["ua_version_major"] = string(strings.Split(string(ver), ".")[0])
			} else {
				udger.ParseData["user_agent"]["ua"] = crawler.Name
				udger.ParseData["user_agent"]["ua_version"] = ""
				udger.ParseData["user_agent"]["ua_version_major"] = ""
			}

			udger.ParseData["user_agent"]["ua_family"] = crawler.Family;
			udger.ParseData["user_agent"]["ua_family_code"] = crawler.FamilyCode;
			udger.ParseData["user_agent"]["ua_family_homepage"] = crawler.FamilyHomepage;
			udger.ParseData["user_agent"]["ua_family_vendor"] = crawler.Vendor;
			udger.ParseData["user_agent"]["ua_family_vendor_code"] = crawler.VendorCode;
			udger.ParseData["user_agent"]["ua_family_vendor_homepage"] = crawler.VendorHomepage;
			udger.ParseData["user_agent"]["ua_family_icon"] = crawler.FamilyIcon;
			udger.ParseData["user_agent"]["ua_family_info_url"] = "https://udger.com/resources/ua-list/bot-detail?bot=" + crawler.Family + "#id" + string(crawler.Botid);
			udger.ParseData["user_agent"]["crawler_last_seen"] = string(crawler.LastSeen);
			udger.ParseData["user_agent"]["crawler_category"] = crawler.CrawlerClassification;
			udger.ParseData["user_agent"]["crawler_category_code"] = crawler.CrawlerClassificationCode;
			udger.ParseData["user_agent"]["crawler_respect_robotstxt"] = crawler.RespectRobotstxt;

		} else {

			// client
			clientID, version, err := udger.findDataWithVersion(ua, udger.rexClients, true)
			if err != nil {
				return nil, err
			}

			if client, ok := udger.Clients[clientID]; ok {
				udger.ParseData["user_agent"]["ua_class"] = client.ClientClassification
				udger.ParseData["user_agent"]["ua_class_code"] = client.ClientClassificationCode

				if version != "" {
					udger.ParseData["user_agent"]["ua"] = client.Name + " " + version
					udger.ParseData["user_agent"]["ua_version"] = version
					udger.ParseData["user_agent"]["ua_version_major"] = string(strings.Split(string(version), ".")[0])
				} else {
					udger.ParseData["user_agent"]["ua"] = client.Name
					udger.ParseData["user_agent"]["ua_version"] = ""
					udger.ParseData["user_agent"]["ua_version_major"] = ""
				}

				udger.ParseData["user_agent"]["ua_uptodate_current_version"] = client.UptodateCurrentVersion
				udger.ParseData["user_agent"]["ua_family"] = client.Name
				udger.ParseData["user_agent"]["ua_family_code"] = client.NameCode
				udger.ParseData["user_agent"]["ua_family_homepage"] = client.Homepage
				udger.ParseData["user_agent"]["ua_family_vendor"] = client.Vendor
				udger.ParseData["user_agent"]["ua_family_vendor_code"] = client.VendorCode
				udger.ParseData["user_agent"]["ua_family_vendor_homepage"] = client.VendorHomepage
				udger.ParseData["user_agent"]["ua_family_icon"] = client.Icon
				udger.ParseData["user_agent"]["ua_family_icon_big"] = client.IconBig
				udger.ParseData["user_agent"]["ua_family_info_url"] = "https://udger.com/resources/ua-list/browser-detail?browser=" + client.Name;
				udger.ParseData["user_agent"]["ua_engine"] = client.Engine
			}
		}
	}
	udger.mux.Unlock()
	return udger.ParseData, nil
}

func (udger *Udger) ParseIp(ip string) map[string]map[string]string {

	udger.mux.Lock()
	defer udger.mux.Unlock()

	if ip == "" {
		return udger.ParseData
	}

	udger.ParseData["ip_address"]["ip"] = ip;
	ipver := getIpVersion(ip);
	if ipver == "v4" {
		udger.ParseData["ip_address"]["ip_ver"] = ipver
		if IpClassificationCode, ok := udger.IPS[ip]; ok {
			udger.ParseData["ip_address"]["ip_classification_code"] = IpClassificationCode
			return udger.ParseData
		}
	}

	udger.ParseData["ip_address"]["ip_classification_code"] = ""

	return udger.ParseData
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

func (udger *Udger) IsCrawler(ip string, ua string, onlyUaParsing bool) bool {
	udger.ParseUa(ua)

	if onlyUaParsing == false{
		udger.ParseIp(ip)
	}

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
