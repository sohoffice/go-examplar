package main

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestCsvFileSource_Provide(t *testing.T) {
	type fields struct {
		path    string
		headers []string
	}
	type args struct {
		filesystem fs.FS
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantD   interface{}
		wantErr bool
	}{
		{
			name: "Happy path: one record header and record matches",
			fields: fields{
				path:    "test.csv",
				headers: []string{"header1", "header2"},
			},
			args: args{
				filesystem: fstest.MapFS{
					"test.csv": {
						Data: []byte("value1,value2\n"),
					},
				},
			},
			wantD: []map[string]interface{}{
				{"header1": "value1", "header2": "value2"},
			},
			wantErr: false,
		},
		{
			name: "Happy path: multiple records",
			fields: fields{
				path:    "test.csv",
				headers: []string{"header1", "header2"},
			},
			args: args{
				filesystem: fstest.MapFS{
					"test.csv": {
						Data: []byte("value1,value2\nvalue3,value4\n"),
					},
				},
			},
			wantD: []map[string]interface{}{
				{"header1": "value1", "header2": "value2"},
				{"header1": "value3", "header2": "value4"},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CsvFileInputSource{
				path:    tt.fields.path,
				headers: tt.fields.headers,
			}
			gotD, err := config.Provide(tt.args.filesystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("Provide() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestPropertiesInputSource_Provide(t *testing.T) {
	type fields struct {
		path string
	}
	type args struct {
		filesystem fs.FS
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantD   interface{}
		wantErr bool
	}{
		{
			name: "Happy path",
			fields: fields{
				path: "test.properties",
			},
			args: args{
				filesystem: fstest.MapFS{
					"test.properties": {
						Data: []byte("key1=value1\nkey2=value2\n# comment\n"),
					},
				},
			},
			wantD: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &PropertiesInputSource{
				path: tt.fields.path,
			}
			gotD, err := config.Provide(tt.args.filesystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("Provide() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestPlainTextFileInputSource_Provide(t *testing.T) {
	type fields struct {
		path          string
		ignoreComment bool
		trim          bool
	}
	type args struct {
		filesystem fs.FS
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantD   interface{}
		wantErr bool
	}{
		{
			name: "no comment, no trim",
			fields: fields{
				path: "file",
			},
			args: args{
				filesystem: fstest.MapFS{
					"file": {
						Data: []byte("foo\n\nbar"),
					},
				},
			},
			wantD: []string{"foo", "bar"},
		},
		{
			name: "trim=true",
			fields: fields{
				path: "file",
				trim: true,
			},
			args: args{
				filesystem: fstest.MapFS{
					"file": {
						Data: []byte("  foo  \n  bar  "),
					},
				},
			},
			wantD: []string{"foo", "bar"},
		},
		{
			name: "ignoreComment=true",
			fields: fields{
				path:          "file",
				ignoreComment: true,
			},
			args: args{
				filesystem: fstest.MapFS{
					"file": {
						Data: []byte("foo # comment 1\n# line comment\nbar# comment 2"),
					},
				},
			},
			wantD: []string{"foo ", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := PlainTextFileInputSource{
				path:          tt.fields.path,
				ignoreComment: tt.fields.ignoreComment,
				trim:          tt.fields.trim,
			}
			gotD, err := receiver.Provide(tt.args.filesystem)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("Provide() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}
