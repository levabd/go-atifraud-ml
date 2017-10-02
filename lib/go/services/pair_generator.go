package services

import (
	"fmt"
	h "github.com/levabd/go-atifraud-ml/lib/go/helpers"
	"reflect"
	"strconv"
)

func GetPairsDictList(ordered_headers_table []map[string]interface{}) []map[string]int {

	var pairs_dict_list []map[string]int

	for _, ordered_headers := range ordered_headers_table {
		pairs_dict_list = append(pairs_dict_list, GetPairDictForHeaders(ordered_headers))
	}

	return pairs_dict_list
}

func GetPairDictForHeaders(ordered_headers map[string]interface{}) map[string]int {

	pairs_dict := make(map[string]int)
	for combination := range GenerateCombinationsMap(ordered_headers) {
		pairs_dict[DefineKeyFoPair(combination)] = 1
	}
	return pairs_dict
}

func DefineKeyFoPair(two_headers []map[string]interface{}) string {
	var (
		key_to_store string
		first        = two_headers[0]
		second       = two_headers[1]
		f_key        = h.GetRandomKeyFromMap(first)
		s_key        = h.GetRandomKeyFromMap(second)
	)

	if h.StrToInt(typeToStr(first[f_key])) < h.StrToInt(typeToStr(second[s_key])) {
		key_to_store = fmt.Sprintf("%s < %s", f_key, s_key)
	} else {
		key_to_store = fmt.Sprintf("%s < %s", s_key, f_key)
	}

	return key_to_store
}
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func typeToStr(value interface{})  string{
	value_type := reflect.TypeOf(value).Name()
	var _value string
	if value ==0 {
		return "0"
	}
	if value_type == "int" {
		_value = fmt.Sprintf("%d", value)
	} else if value_type == "float64" {
		_value = FloatToString(value.(float64))
	}else if value_type == "string" {
		_value =  value.(string)
	} else {
		println(value.(string))
		panic(fmt.Sprintf("Expect to type string or int, %s provided %g", value_type, value))
	}

	return _value
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

	_length := len(headers)
	if _length == 0 {
		return
	}

	for _key, _value := range headers {

		if key != _key {
			var _return []map[string]interface{}
			var first = make(map[string]interface{})
			var second = make(map[string]interface{})
			first[key] = value
			second[_key] = _value
			_return = append(_return, first)
			_return = append(_return, second)
			c <- _return
		}
	}

	delete(headers, key);
	_length = len(headers)
	if _length == 0 {
		return
	}

	var _key = h.GetRandomKeyFromMap(headers)

	AddHeaderMap(c, _key, typeToStr(headers[_key]), headers)
}
