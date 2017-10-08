package services

func GetFullFeatures(orderedHeadersTable []map[string]interface{}, valuesHeadersTable []map[string]interface{}, valuesFeaturesOrder map[string]int) [][]float64 {

	var fullFeatures [][]float64

	for index, valuesHeaders := range valuesHeadersTable {
		valueFeatures := GetSingleValueFeatures(valuesHeaders, valuesFeaturesOrder)
		orderFeatures := GetSingleOrderFeatures(orderedHeadersTable[index])
		fullLineFeatures := append(orderFeatures, valueFeatures...)
		fullFeatures = append(fullFeatures, fullLineFeatures)
	}

	return fullFeatures
}

func GetSingleFullFeatures(orderedHeaders map[string]interface{}, valueHeaders map[string]interface{}, valuesFeaturesOrder map[string]int) []float64 {

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