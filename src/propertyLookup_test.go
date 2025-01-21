package main

import (
	"reflect"
	"testing"
)

func TestPropertiesLookup_HasProperty(t *testing.T) {
	type fields struct {
		properties []interface{}
	}
	type args struct {
		key string
	}
	map1 := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	map2 := map[string]interface{}{
		"foo":  "value foo",
		"key2": "bar",
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Test HasProperty",
			fields: fields{
				properties: []interface{}{
					map1, map2,
				},
			},
			args: args{
				key: "foo",
			},
			want: true,
		},
		{
			name: "Not found",
			fields: fields{
				properties: []interface{}{
					map1, map2,
				},
			},
			args: args{
				key: "bar",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := PropertiesLookup{
				properties: tt.fields.properties,
			}
			if got := config.HasProperty(tt.args.key); got != tt.want {
				t.Errorf("HasProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPropertiesLookup_GetProperty(t *testing.T) {
	type fields struct {
		properties []interface{}
	}
	type args struct {
		key string
	}
	map1 := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	map2 := map[string]interface{}{
		"foo":  "value foo",
		"key2": "bar",
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "Test GetProperty",
			fields: fields{
				properties: []interface{}{
					map1, map2,
				},
			},
			args: args{
				key: "key2",
			},
			want: "value2",
		},
		{
			name: "Not found",
			fields: fields{
				properties: []interface{}{
					map1, map2,
				},
			},
			args: args{
				key: "bar",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := PropertiesLookup{
				properties: tt.fields.properties,
			}
			if got := config.GetProperty(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
