package main

import (
	"time"
	"log"
	"fmt"
	//"github.com/levabd/go-atifraud-ml/lib/go/services"
	"runtime"
	"github.com/jinzhu/gorm"
	"github.com/levabd/go-atifraud-ml/lib/go/models"
	"gopkg.in/cheggaaa/pb.v1"
	"sync"
	//"os"
	"github.com/levabd/go-atifraud-ml/lib/go/udger"
	"github.com/levabd/go-atifraud-ml/lib/go/services"
	"os"
	"github.com/uniplaces/carbon"
)

var (
	udger_instance *udger.Udger = nil
)

func main() {
	u, err := udger.New(os.Getenv("DB_FILE_PATH_UDGER"))
	if err != nil {
		log.Fatalf("server.go - init: Failed to initiate udger helper: %s", err)
	}

	udger_instance = u

	numCPU := runtime.NumCPU()
	fmt.Println("NumCPU", numCPU)
	runtime.GOMAXPROCS(numCPU)

	//fmt.Println("ddd %+v ", fmt.Sprintf(strings.Join(strings.SplitN("56.0.1750.154",".", 2)[:2], ".")))
	//fmt.Println("ddd %+v ", strings.SplitN("56.0.1750.154",".", 4)[:3])

	//println(startTime.Minute(), startTime.Second(), startNanosecond)

	// ver all
	//from 1505399120 to 1502236800 - 425493 records (shape - 3564)
	//from 1505399120 to 1507802544 - 178439 records (shape - 2421)
	//from 1505399120 to 1505808000 - 120074 records (shape - 2073)

	// ver 00.00
	//from 1505399120 to 1505808000 - 119992 records (shape - 1274)
	//from 1502150400 to 1507802544 ~ 240000 records (shape - 1274)
	//from 1502150400 to 1507075200 - 481843 records (shape - 2073)

	//PrepareData: after GetTrimmedLodMapsForPeriod. len(uaVersionList):  54057
	//PrepareData: after FitValuesFeaturesOrder. len(valuesFeaturesOrder):  165
	// 1502150400 - 1503082800 - 60 000
	// 1502150400 - 1503360000 - 80 000
	// 1502150400 - 1504051200 - 100 000
	// 1502150400 - 1504328400 - 120 000

	// 1502150400 - 1505304000 - 240 000
	// 1502150400 - 1507093200 - 480 000
	start, _ := carbon.Create(2017, 8,28,0,0,0,0,"Asia/Almaty")
	end, _ := carbon.Create(2017, 9,7,23,0,0,0,"Asia/Almaty")

	auVersionIntCodes, _, _, intFullFeatures, familyCodesList, logIds := services.PrepareDataUaVersion(
		start.Unix(),
		end.Unix(),
		udger_instance,
	)

	storePreparedData(intFullFeatures, auVersionIntCodes, familyCodesList, logIds)
}

func storePreparedData(fullFeatures [][]int, yIntCodes [][]int, list []string, logIds []uint) {

	db, err := gorm.Open("postgres", models.GetDBConnectionStr())
	if err != nil {
		panic("user_agent_helpers.go - LoadFittedUserAgentCodes: Failed to connect to database")
	}
	defer db.Close()
	if !db.HasTable(&models.Features{}) {
		db.AutoMigrate(&models.Features{})
	}
	if !db.HasTable(&models.Browsers{}) {
		db.AutoMigrate(&models.Browsers{})
	}
	db.Exec("TRUNCATE TABLE features;")
	db.Exec("TRUNCATE TABLE browsers;")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		tx := db.Begin()
		bar := pb.StartNew(len(fullFeatures))
		bar.SetRefreshRate(time.Second)
		for index_row, featureValues := range fullFeatures {
			for index_column, value := range featureValues {
				if value == 1 {
					cacheFeatures := models.Features{
						LogId:  logIds[index_row],
						Row:    index_row,
						Column: index_column,
					}
					tx.Create(&cacheFeatures)
				}
			}
			bar.Increment()
		}
		tx.Commit()
	}()

	go func() {
		defer wg.Done()
		tx := db.Begin()
		bar := pb.StartNew(len(yIntCodes))
		bar.SetRefreshRate(time.Second)
		for i, name := range list {
			bar.Increment()
			cacheFeatures := models.Browsers{Name: name, LogId: logIds[i]}
			tx.Create(&cacheFeatures)
		}
		tx.Commit()
	}()
	wg.Wait()
}
