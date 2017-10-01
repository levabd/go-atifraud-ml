package services

import (
	"fmt"
	h "github.com/levabd/go-atifraud-ml/lib/go/helpers"
)

func GetPairsDictList(ordered_headers_table []map[string]interface{}) []map[string]int {

	var pairs_dict_list []map[string]int

	for _, ordered_headers := range ordered_headers_table {
		pairs_dict := make(map[string]int)
		fmt.Println(fmt.Sprintf("ordered_headers %+v", ordered_headers))
		for combination := range GenerateCombinationsMap(ordered_headers) {
			pairs_dict[DefineKeyFoPair(combination)] = 1
		}
		pairs_dict_list = append(pairs_dict_list, pairs_dict)
	}
	return pairs_dict_list
}

func DefineKeyFoPair(two_headers []map[string]interface{}) string {
	var (
		key_to_store string
		first        = two_headers[0]
		second       = two_headers[1]
		f_key        = h.GetRandomKeyFromMap(first)
		s_key        = h.GetRandomKeyFromMap(second)
	)

	if h.StrToInt(first[f_key].(string)) < h.StrToInt(second[s_key].(string)) {
		key_to_store = fmt.Sprintf("%s < %s", f_key, s_key)
	} else {
		key_to_store = fmt.Sprintf("%s < %s", s_key, f_key)
	}

	return key_to_store
}


func GenerateCombinationsMap(headers map[string]interface{}) <-chan []map[string]interface{} {

	c := make(chan []map[string]interface{})

	// Starting a separate goroutine that will create all the combinations,
	// feeding them to the channel c
	go func(c chan []map[string]interface{}) {
		defer close(c) // Once the iteration function is finished, we close the channel
		var key = h.GetRandomKeyFromMap(headers)
		value := headers[key].(string)
		AddHeaderMap(c, key, value, headers) // We start by feeding it an empty string
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

	var _key =  h.GetRandomKeyFromMap(headers)
	_value := headers[_key].(string)

	AddHeaderMap(c, _key, _value, headers)
}