package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/magiconair/properties"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"strings"
)

// InputSource is an interface that defines the behavior of a source that can be read.
type InputSource interface {
	// Provide reads the source and returns the data.
	Provide(filesystem fs.FS) (d interface{}, err error)
}

type CsvFileInputSource struct {
	// Path is the path to the CSV file.
	path string
	// Headers is the list of headers in the CSV file.
	headers []string
}

// Provide reads the CSV file and returns the data as []map[string]interface{}.
func (config *CsvFileInputSource) Provide(filesystem fs.FS) (d interface{}, err error) {
	f, err := filesystem.Open(config.path)
	if err != nil {
		panic(fmt.Sprintf("Error reading %s: %v", config.path, err))
	}
	defer func() {
		_ = f.Close()
	}()
	records := make([]map[string]interface{}, 0)
	reader := csv.NewReader(f)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			myErr := fmt.Errorf(fmt.Sprintf("Error reading csv line from %s: %v", config.path, err))
			return nil, myErr
		}
		if len(line) != len(config.headers) {
			log.Printf("Record length %d does not match header length %d", len(line), len(config.headers))
		}
		m := make(map[string]interface{})
		for i, header := range config.headers {
			if i < len(line) {
				m[header] = line[i]
			} else {
				m[header] = ""
			}
		}
		records = append(records, m)
	}
	return records, nil
}

type PropertiesInputSource struct {
	// Path is the path to the properties file.
	path string
}

// Provide reads the properties file and returns the data as map[string]interface{}.
func (config *PropertiesInputSource) Provide(filesystem fs.FS) (d interface{}, err error) {
	f, err := filesystem.Open(config.path)
	if err != nil {
		panic(fmt.Sprintf("Error reading %s: %v", config.path, err))
	}
	defer func() {
		_ = f.Close()
	}()
	p := properties.MustLoadReader(f, properties.UTF8)
	m := p.Map()
	return m, nil
}

type PlainTextFileInputSource struct {
	path          string
	ignoreComment bool
	trim          bool
}

// Provide reads the plain text file and returns the data as []string.
func (receiver PlainTextFileInputSource) Provide(filesystem fs.FS) (d interface{}, err error) {
	f, err := filesystem.Open(receiver.path)
	if err != nil {
		panic(fmt.Sprintf("Error reading %s: %v", receiver.path, err))
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	records := make([]string, 0)
	for sc.Scan() {
		line := sc.Text()
		if receiver.ignoreComment {
			idx := strings.Index(line, "#")
			if idx > -1 {
				line = line[0:idx]
			}
		}
		if receiver.trim {
			line = strings.TrimSpace(line)
		}
		if len(line) > 0 {
			records = append(records, line)
		}
	}
	return records, nil
}

type YamlInputSource struct {
	path    string
	flatten int
}

// Provide reads the YAML file and returns the data as map[interface{}]interface{}.
func (config YamlInputSource) Provide(filesystem fs.FS) (data interface{}, err error) {
	f, err := filesystem.Open(config.path)
	if err != nil {
		panic(fmt.Sprintf("Error reading %s: %v", config.path, err))
	}
	defer f.Close()

	m := make(map[interface{}]interface{})
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", m)
	return m, nil
}
