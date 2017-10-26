package services

import (
	"fmt"
	h "github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"reflect"
	"strconv"
)

var (OrderPairsFeaturesOrder = map[string]int {
	"Accept < Accept-Encoding": 0,
	"Accept < Connection": 1,
	"Accept < From": 2,
	"Accept < Host": 3,
	"Accept < If-Modified-Since": 4,
	"Accept < Upgrade-Insecure-Requests": 5,
	"Accept < User-Agent": 6,
	"Accept-Encoding < Accept": 7,
	"Accept-Encoding < Connection": 8,
	"Accept-Encoding < Host": 9,
	"Accept-Encoding < If-Modified-Since": 10,
	"Accept-Encoding < Upgrade-Insecure-Requests": 11,
	"Accept-Encoding < User-Agent": 12,
	"Connection < Accept": 13,
	"Connection < Accept-Encoding": 14,
	"Connection < From": 15,
	"Connection < Host": 16,
	"Connection < If-Modified-Since": 17,
	"Connection < Upgrade-Insecure-Requests": 18,
	"Connection < User-Agent": 19,
	"From < Accept": 20,
	"From < Accept-Encoding": 21,
	"From < User-Agent": 22,
	"Host < Accept": 23,
	"Host < Accept-Encoding": 24,
	"Host < Connection": 25,
	"Host < From": 26,
	"Host < If-Modified-Since": 27,
	"Host < Upgrade-Insecure-Requests": 28,
	"Host < User-Agent": 29,
	"If-Modified-Since < Accept": 30,
	"If-Modified-Since < Accept-Encoding": 31,
	"If-Modified-Since < Connection": 32,
	"If-Modified-Since < Host": 33,
	"If-Modified-Since < Upgrade-Insecure-Requests": 34,
	"If-Modified-Since < User-Agent": 35,
	"Upgrade-Insecure-Requests < Accept": 36,
	"Upgrade-Insecure-Requests < Accept-Encoding": 37,
	"Upgrade-Insecure-Requests < Connection": 38,
	"Upgrade-Insecure-Requests < Host": 39,
	"Upgrade-Insecure-Requests < If-Modified-Since": 40,
	"Upgrade-Insecure-Requests < User-Agent": 41,
	"User-Agent < Accept": 42,
	"User-Agent < Accept-Encoding": 43,
	"User-Agent < Connection": 44,
	"User-Agent < From": 45,
	"User-Agent < Host": 46,
	"User-Agent < If-Modified-Since": 47,
	"User-Agent < Upgrade-Insecure-Requests": 48,
})

func GetOrderFeatures(orderedHeadersTable []map[string]interface{}) [][]float64 {

	var orderFeatures [][]float64

	for _, orderedHeaders := range orderedHeadersTable {
		orderFeatures = append(orderFeatures, GetSingleOrderFeatures(orderedHeaders))
	}

	return orderFeatures
}

func GetSingleOrderFeatures(orderedHeaders map[string]interface{}) (orderFeatures []float64) {

	orderFeatures = make([]float64, len(OrderPairsFeaturesOrder))
	for combination := range GenerateCombinationsMap(orderedHeaders) {
		orderFeatures[OrderPairsFeaturesOrder[DefineKeyFoPair(combination)]] = 1.0
	}
	return
}
func GetSingleOrderFeaturesInt(orderedHeaders map[string]interface{}) (orderFeatures []int) {
	orderFeatures = make([]int, len(OrderPairsFeaturesOrder))
	for combination := range GenerateCombinationsMap(orderedHeaders) {
		orderFeatures[OrderPairsFeaturesOrder[DefineKeyFoPair(combination)]] = 1
	}
	return
}
func DefineKeyFoPair(twoHeaders []map[string]interface{}) string {
	var (
		key2Store string
		first        = twoHeaders[0]
		second       = twoHeaders[1]
		fKey        = h.GetRandomKeyFromMap(first)
		sKey        = h.GetRandomKeyFromMap(second)
	)

	//println("fKey", fKey, fmt.Sprintf("%+v", first), typeToStr(first[fKey]), "sKey", sKey, fmt.Sprintf("%+v", second), typeToStr(second[sKey]))

	if h.StrToInt(typeToStr(first[fKey])) < h.StrToInt(typeToStr(second[sKey])) {
		key2Store = fmt.Sprintf("%s < %s", fKey, sKey)
	} else {
		key2Store = fmt.Sprintf("%s < %s", sKey, fKey)
	}

	return key2Store
}
func FloatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func typeToStr(value interface{})  string{
	if value ==nil {
		return ""
	}
	valueType := reflect.TypeOf(value).Name()
	var tmpValue string

	if value ==0 {
		return "0"
	}
	if valueType == "int" {
		tmpValue = fmt.Sprintf("%d", value)
	} else if valueType == "float64" {
		tmpValue = FloatToString(value.(float64))
	}else if valueType == "string" {
		tmpValue =  value.(string)
	} else {
		println(value.(string))
		panic(fmt.Sprintf("Expect to type string or int, %s provided %g", valueType, value))
	}

	return tmpValue
}

func GenerateCombinationsMap(headers map[string]interface{}) <-chan []map[string]interface{} {

	c := make(chan []map[string]interface{})

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan []map[string]interface{}) {
		defer close(c) // Once the iteration function is finished, we close the channel
		var key = h.GetRandomKeyFromMap(headers)
		AddHeaderMap(c, key, typeToStr( headers[key]), headers) // We start by feeding it an empty string
	}(c)

	return c // Return the channel to the calling function
}

func AddHeaderMap(c chan []map[string]interface{}, key string, value string, headers map[string]interface{}) {

	length := len(headers)
	if length == 0 {
		return
	}

	for tmpKey, tmpValue := range headers {

		if key != tmpKey {
			var result []map[string]interface{}
			var first = make(map[string]interface{})
			var second = make(map[string]interface{})
			first[key] = value
			second[tmpKey] = tmpValue
			result = append(result, first)
			result = append(result, second)
			c <- result
		}
	}

	delete(headers, key)
	length = len(headers)
	if length == 0 {
		return
	}

	var tmpKey = h.GetRandomKeyFromMap(headers)

	AddHeaderMap(c, tmpKey, typeToStr(headers[tmpKey]), headers)
}
