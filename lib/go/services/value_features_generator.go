package services

import  (
	m "github.com/levabd/go-atifraud-ml/lib/go/models"
	"github.com/jinzhu/gorm"
)

func GetValueFeatures(valuesHeadersTable []map[string]interface{}, valuesFeaturesOrder map[string]int) [][]float64 {

	var valueFeatures [][]float64

	for _, valuesHeaders := range valuesHeadersTable {
		valueFeatures = append(valueFeatures, GetSingleValueFeatures(valuesHeaders, valuesFeaturesOrder))
	}

	return valueFeatures
}

func GetSingleValueFeatures(valueHeaders map[string]interface{}, valuesFeaturesOrder map[string]int) []float64 {

	valueFeatures := make([]float64, len(valuesFeaturesOrder))
	for header, value := range valueHeaders {
	//	fmt.Println(fmt.Sprintf(`valuesFeaturesOrder[header + "=" + typeToStr(value) %s %s `, header, value))
		valueFeatures[valuesFeaturesOrder[header + "=" + typeToStr(value)]] = 1.0
		break
	}
	return valueFeatures
}

func GetSingleValueFeaturesInt(valueHeaders map[string]interface{}, valuesFeaturesOrder map[string]int) []int {

	valueFeatures := make([]int, len(valuesFeaturesOrder))
	for header, value := range valueHeaders {
		valueFeatures[valuesFeaturesOrder[header + "=" + typeToStr(value)]] = 1
	}
	return valueFeatures
}

func FitValuesFeaturesOrder(valuesHeadersTable []map[string]interface{}) map[string]int {
	var valuesFeaturesOrder = map[string]int {}

	index := 0

	for _, valuesHeader := range valuesHeadersTable {
		for header, value := range valuesHeader {
			potentialKey := header + "=" + typeToStr(value)
			if _, ok := valuesFeaturesOrder[potentialKey]; !ok {
				valuesFeaturesOrder[potentialKey] = index
				index++
			}
		}
	}


	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		Logger.Fatalf("value_features_generator.go - FitValuesFeaturesOrder: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.ValueFeatureOrder{}) {
		db.AutoMigrate(&m.ValueFeatureOrder{})
	}

	// Clean last vectoriser
	db.Exec("TRUNCATE TABLE  value_features_order;")

	// Insert new fitted vectoriser
	for featureName, featureOrder := range valuesFeaturesOrder {
		tx := db.Begin()
		dbFeature := m.ValueFeatureOrder {
			FeatureName: featureName,
			Order: featureOrder,
		}
		tx.Create(&dbFeature)
		tx.Commit()
	}

	return valuesFeaturesOrder
}

func LoadFittedValuesFeaturesOrder() map[string]int {
	db, err := gorm.Open("postgres", m.GetDBConnectionStr())
	if err != nil {
		Logger.Fatalf("value_features_generator.go - LoadFittedValuesFeaturesOrder: Failed to connect database: %s", err)
	}
	defer db.Close()
	if !db.HasTable(&m.ValueFeatureOrder{}) {
		db.AutoMigrate(&m.ValueFeatureOrder{})
	}

	var(
		dbFeatures = []m.ValueFeatureOrder {}
		valuesFeaturesOrder = map[string]int {}
	)

	db.Find(&dbFeatures)
	for _, dbFeature := range dbFeatures {
		valuesFeaturesOrder[dbFeature.FeatureName] = dbFeature.Order
	}

	return valuesFeaturesOrder
}