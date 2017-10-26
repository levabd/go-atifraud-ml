package services


func GetFullFeatures(
	orderedHeadersTable []map[string]interface{},
	valuesHeadersTable []map[string]interface{},
	valuesFeaturesOrder map[string]int) ([][]float64, [][]int) {

	var floatFullFeatures [][]float64
	var intFullFeatures [][]int

	for index, valuesHeaders := range valuesHeadersTable {
		//floatValueFeatures := GetSingleValueFeatures(valuesHeaders, valuesFeaturesOrder)
		//floatOrderFeatures := GetSingleOrderFeatures(orderedHeadersTable[index])
		//floatFullLineFeatures := append(floatOrderFeatures, floatValueFeatures...)
		//floatFullFeatures = append(floatFullFeatures, floatFullLineFeatures)

		intValueFeatures := GetSingleValueFeaturesInt(valuesHeaders, valuesFeaturesOrder)
		intOrderFeatures := GetSingleOrderFeaturesInt(orderedHeadersTable[index])
		intFullLineFeatures := append(intOrderFeatures, intValueFeatures...)
		intFullFeatures = append(intFullFeatures, intFullLineFeatures)
	}

	return floatFullFeatures, intFullFeatures
}

func GetSingleFullFeatures(orderedHeaders map[string]interface{},
	valueHeaders map[string]interface{},
	valuesFeaturesOrder map[string]int) []float64 {

	valueFeatures := make([]float64, len(valuesFeaturesOrder))
	for header, value := range valueHeaders {
		valueFeatures[valuesFeaturesOrder[header + "=" + typeToStr(value)]] = 1.0
	}

	orderFeatures := make([]float64, len(OrderPairsFeaturesOrder))
	for combination := range GenerateCombinationsMap(orderedHeaders) {
		orderFeatures[OrderPairsFeaturesOrder[DefineKeyFoPair(combination)]] = 1.0
	}

	return append(orderFeatures, valueFeatures...)
}