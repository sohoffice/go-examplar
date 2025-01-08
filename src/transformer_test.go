package main

import (
	"reflect"
	"testing"
)

func TestListToSetTransformer_Transform(t *testing.T) {
	type fields struct {
		keyMapper   Mapper
		valueMapper Mapper
	}
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Happy path: simple list",
			fields: fields{
				keyMapper:   StringMapMapper("foo"),
				valueMapper: StringMapMapper("bar"),
			},
			args: args{
				input: []map[string]interface{}{
					{"foo": 1, "bar": "BAR1"},
					{"foo": 2, "bar": "BAR2"},
				},
			},
			want: map[interface{}]interface{}{
				1: "BAR1",
				2: "BAR2",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := &ListToMapTransformer{
				keyMapper:   tt.fields.keyMapper,
				valueMapper: tt.fields.valueMapper,
			}
			got, err := transformer.Transform(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transform() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringMapMapper(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Happy path",
			args: args{
				key: "foo",
			},
			want: "FOO",
		},
		{
			name: "Not found",
			args: args{
				key: "not-found",
			},
			want: nil,
		},
	}
	testData := make(map[string]interface{})
	testData["foo"] = "FOO"
	testData["bar"] = "BAR"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := StringMapMapper(tt.args.key)
			got := mapper(testData)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringMapMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListMappingTransformer_Transform(t *testing.T) {
	type fields struct {
		mapping map[interface{}]interface{}
	}
	type args struct {
		input interface{}
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRecords []interface{}
		wantErr     bool
	}{
		{
			name: "normal flow",
			fields: fields{
				mapping: map[interface{}]interface{}{
					"1": "one",
					"9": "nine",
				},
			},
			args: args{
				input: []string{
					"0", "1", "2",
				},
			},
			wantRecords: []interface{}{
				"0", "one", "2",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ListMappingTransformer{
				mapping: tt.fields.mapping,
			}
			gotRecords, err := config.Transform(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecords, tt.wantRecords) {
				t.Errorf("Transform() gotRecords = %v, want %v", gotRecords, tt.wantRecords)
			}
		})
	}
}
