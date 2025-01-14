package main

import (
	"errors"
	"log"
	"reflect"
)

type Mapper func(input interface{}) interface{}

// StringMapMapper returns a mapper that extracts the value of the key from the map.
func StringMapMapper(key string) Mapper {
	return func(input interface{}) interface{} {
		if input == nil {
			return nil
		}
		mv := reflect.ValueOf(input)
		kv := reflect.ValueOf(key)
		v := mv.MapIndex(kv)
		if v.Kind() == reflect.Invalid { // if the map value cannot be found
			return nil
		}
		return v.Interface()
	}
}

type StringMapper func(input interface{}) string

// IdentityMapper returns the input as a string.
func IdentityMapper(input interface{}) string {
	return input.(string)
}

type Transformer interface {
	Transform(input interface{}) (interface{}, error)
}

// ListToMapTransformer is a transformer that converts a list to a map using the specified key/value mapper.
type ListToMapTransformer struct {
	keyMapper   Mapper
	valueMapper Mapper
}

func (config *ListToMapTransformer) Transform(input interface{}) (interface{}, error) {
	if input == nil {
		return nil, nil
	}
	if reflect.TypeOf(input).Kind() != reflect.Slice {
		return nil, errors.New("ListToMapTransformer: Input is not a list")
	}
	list := reflect.ValueOf(input)
	m := make(map[interface{}]interface{})
	for i := 0; i < list.Len(); i++ {
		el := list.Index(i)
		k := config.keyMapper(el.Interface())
		v := config.valueMapper(el.Interface())
		m[k] = v
	}
	return m, nil
}

type ListMappingTransformer struct {
	mapping map[string]string
}

// Transform transforms the input list by mapping the elements to the values in the mapping.
// If the element is not found in the mapping, original value is used.
func (config ListMappingTransformer) Transform(input interface{}) (records []interface{}, err error) {
	if input == nil {
		return nil, errors.New("ListMappingTransformer: Input is nil")
	}
	if reflect.TypeOf(input).Kind() != reflect.Slice {
		return nil, errors.New("ListMappingTransformer: Input is not a list")
	}
	listV := reflect.ValueOf(input)
	records = make([]interface{}, 0)
	for i := 0; i < listV.Len(); i++ {
		el := listV.Index(i).String()
		if val, ok := config.mapping[el]; ok {
			records = append(records, val)
		} else {
			records = append(records, el)
		}
	}
	return records, nil
}

type ListExpandTransformer struct {
	dataByKey map[interface{}]interface{}
	keyMapper StringMapper
	// Keep the original key name into the config struct. The config struct must be a map
	keepKeyName bool
}

// Transform transforms the input list by expanding the elements to the values in the data.
func (config ListExpandTransformer) Transform(input interface{}) (record []interface{}, err error) {
	if input == nil {
		return nil, errors.New("ListExpandTransformer: Input is nil")
	}
	if reflect.TypeOf(input).Kind() != reflect.Slice {
		return nil, errors.New("ListExpandTransformer: Input is not a list")
	}
	bookkeeping := make(map[interface{}]bool)
	for k, _ := range config.dataByKey {
		bookkeeping[k] = false
	}
	listV := reflect.ValueOf(input)
	records := make([]interface{}, 0)
	for i := 0; i < listV.Len(); i++ {
		el := listV.Index(i).String()
		keyName := config.keyMapper(el)
		if val, ok := config.dataByKey[keyName]; ok {
			if config.keepKeyName {
				reflect.ValueOf(val).SetMapIndex(reflect.ValueOf("Name"), reflect.ValueOf(keyName))
			}
			records = append(records, val)
			bookkeeping[el] = true
		} else {
			log.Printf("Input key '%s' not found in data", el)
		}
	}

	// report bookkeeping results
	for k, v := range bookkeeping {
		if v == false {
			log.Printf("Config key '%s' not used in expand transformer", k)
		}
	}
	return records, nil
}
