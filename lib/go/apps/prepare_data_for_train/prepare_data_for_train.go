package main

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/levabd/go-atifraud-ml/lib/go/models"
	"gopkg.in/cheggaaa/pb.v1"
	"sync"
	"github.com/levabd/go-atifraud-ml/lib/go/services"
)

func main() {

	auVersionIntCodes, _, _, intFullFeatures, familyCodesList, logIds := services.PrepareUaFamilyCodes(90000)

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
