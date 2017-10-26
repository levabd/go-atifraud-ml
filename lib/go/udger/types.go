package udger

import (
	"github.com/jinzhu/gorm"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type Udger struct {
	ParseData map[string]map[string]string

	Clients    map[int64]Client
	OS         map[int64]Os
	IPS        map[string]string
	Crawlers   map[string]Crawler
	Devices    []Device
	rexClients []rexData
	rexDevices []rexData
	rexOS      []rexData
	DB         *gorm.DB
	DbPath     string
}

type rexData struct {
	ID            int64
	Regex         string
	RegexCompiled pcre.Regexp
}

type Crawler struct {
	Botid                     int64
	Name                      string
	Ver                       []uint8
	VerMajor                  []uint8
	LastSeen                  []uint8
	RespectRobotstxt          string
	Family                    string
	FamilyCode                string
	FamilyHomepage            string
	FamilyIcon                string
	Vendor                    string
	VendorCode                string
	VendorHomepage            string
	CrawlerClassification     string
	CrawlerClassificationCode string
	UaString                  string
}

type Client struct {
	ClassId                  int64
	ClientId                 int64
	Regstring                string
	Name                     string
	NameCode                 string
	Homepage                 string
	Icon                     string
	IconBig                  string
	Engine                   string
	Vendor                   string
	VendorCode               string
	VendorHomepage           string
	UptodateCurrentVersion   string
	ClientClassification     string
	ClientClassificationCode string
}

type Os struct {
	OsId           int64
	Regstring      string
	Family         string
	FamilyCode     string
	Name           string
	NameCode       string
	Homepage       string
	Icon           string
	IconBig        string
	Vendor         string
	VendorCode     string
	VendorHomepage string
}

type ClientOsRelation struct {
	OsId           int64
	Family         string
	FamilyCode     string
	Name           string
	NameCode       string
	Homepage       string
	Icon           string
	IconBig        string
	Vendor         string
	VendorCode     string
	VendorHomepage string
}

type Device struct {
	DeviceclassId int64
	Regstring     string
	Name          string
	NameCode      string
	Icon          string
	IconBig       string
}

type DeviceMarketName struct {
	Id        int64
	Regstring string
}

type DeviceBrand struct {
	Marketname   string
	Brand        string
	BrandCode    string
	brandUrl     string
	Icon         string
	Icon_big     string
	BrandInfoUrl string
}

type Ip struct {
	Ip                   string
	IpClassificationCode string
}
