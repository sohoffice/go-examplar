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
		mapping map[string]string
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
				mapping: map[string]string{
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

func TestListExpandTransformer_Transform(t *testing.T) {
	type fields struct {
		dataByKey   map[interface{}]interface{}
		keyMapper   StringMapper
		keepKeyName bool
	}
	type args struct {
		input interface{}
	}
	objectByKey := map[interface{}]interface{}{
		"foo": map[interface{}]interface{}{
			"val": 10,
		},
		"bar": map[interface{}]interface{}{
			"val": 15,
		},
		"foobar": map[interface{}]interface{}{
			"val": 20,
		},
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantRecord []interface{}
		wantErr    bool
	}{
		{
			name: "normal flow",
			fields: fields{
				dataByKey: objectByKey,
				keyMapper: IdentityMapper,
			},
			args: args{
				input: []string{
					"bar", "foo", "baz",
				},
			},
			wantRecord: []interface{}{
				map[interface{}]interface{}{
					"val": 15,
				},
				map[interface{}]interface{}{
					"val": 10,
				},
			},
		},
		{
			name: "keep key name",
			fields: fields{
				dataByKey:   objectByKey,
				keyMapper:   IdentityMapper,
				keepKeyName: true,
			},
			args: args{
				input: []string{
					"bar", "foo", "baz",
				},
			},
			wantRecord: []interface{}{
				map[interface{}]interface{}{
					"val":  15,
					"Name": "bar",
				},
				map[interface{}]interface{}{
					"val":  10,
					"Name": "foo",
				},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ListExpandTransformer{
				dataByKey:   tt.fields.dataByKey,
				keyMapper:   tt.fields.keyMapper,
				keepKeyName: tt.fields.keepKeyName,
			}
			gotRecord, err := config.Transform(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRecord, tt.wantRecord) {
				t.Errorf("Transform() gotRecord = %v, want %v", gotRecord, tt.wantRecord)
			}
		})
	}
}

var testData = []map[string]interface{}{
	{"featureSet": "one", "name": "foo"},
	{"featureSet": "two", "name": "bar"},
	{"featureSet": "one", "name": "baz"},
}

func TestListFilterTransformer_Transform(t *testing.T) {
	type fields struct {
		predicate Predicate
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
				predicate: func(it interface{}) bool {
					return it.(map[string]interface{})["featureSet"] == "one"
				},
			},
			args: args{
				input: testData,
			},
			wantRecords: []interface{}{
				map[string]interface{}{"featureSet": "one", "name": "foo"},
				map[string]interface{}{"featureSet": "one", "name": "baz"},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ListFilterTransformer{
				predicate: tt.fields.predicate,
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

func TestListStringSortTransformer_Transform(t *testing.T) {
	type fields struct {
		mapper StringMapper
	}
	type args struct {
		input interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "normal flow",
			fields: fields{
				mapper: IdentityMapper,
			},
			args: args{
				input: []interface{}{
					"3", "1", "2",
				},
			},
			want: []interface{}{
				"1", "2", "3",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ListStringSortTransformer{
				mapper: tt.fields.mapper,
			}
			got, err := config.Transform(tt.args.input)
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

func TestMapValueStringMapper(t *testing.T) {
	type args struct {
		key  string
		data map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy path",
			args: args{
				key: "priority",
				data: map[string]interface{}{
					"priority": "A01",
				},
			},
			want: "A01",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := MapValueStringMapper(tt.args.key)
			if got := mapper(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapValueStringMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}
