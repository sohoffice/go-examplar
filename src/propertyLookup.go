package main

import "reflect"

type PropertiesLookup struct {
	properties []interface{}
}

// HasProperty returns true if the key is found in any of the properties.
func (config PropertiesLookup) HasProperty(key string) bool {
	for _, prop := range config.properties {
		val := reflect.ValueOf(prop).MapIndex(reflect.ValueOf(key))
		if val.IsValid() {
			return true
		}
	}
	return false
}

// GetProperty returns the value of the key if found in any of the properties.
func (config PropertiesLookup) GetProperty(key string) interface{} {
	for _, prop := range config.properties {
		val := reflect.ValueOf(prop).MapIndex(reflect.ValueOf(key))
		if val.IsValid() {
			return val.Interface()
		}
	}
	return nil
}
