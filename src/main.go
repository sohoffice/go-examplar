package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"os"
)

type Args struct {
	ConfigDir          string `arg:"positional"`
	FeatureFile        string `arg:"--feature-file,required"`
	FeatureMappingFile string `arg:"--feature-mapping-file,required"`
}

// Processing steps
// 1. Read features from a plain text file
// 2. Read feature mapping from a properties file
// 3. Convert feature names using the mapping
func main() {
	// use struct embedding to create a anonymous struct while still using a declared interface
	// so it can be referred to later
	var args struct {
		Args
	}
	arg.MustParse(&args)
	fmt.Printf("args: %v\n", args)

	context := make(map[string]interface{})
	context["args"] = args.Args
	// read the source feature names
	context["raw-features"] = readFeatures(context)
	// read the feature mapping to new names
	context["feature-mapping"] = readFeatureMapping(context)
	// convert feature to new names
	context["features"] = convertFeatureNames(context)

	fmt.Printf("Output: %+v\n", context)
}

// Read the raw feature file
func readFeatures(context map[string]interface{}) interface{} {
	step := PlainTextFileInputSource{
		path:          context["args"].(Args).FeatureFile,
		ignoreComment: true,
		trim:          true,
	}
	value1, err := step.Provide(os.DirFS(context["args"].(Args).ConfigDir))
	if err != nil {
		panic(err)
	}
	return value1
}

func readFeatureMapping(context map[string]interface{}) interface{} {
	step := PropertiesInputSource{
		path: context["args"].(Args).FeatureMappingFile,
	}
	value, err := step.Provide(os.DirFS(context["args"].(Args).ConfigDir))
	if err != nil {
		panic(err)
	}
	return value
}

func convertFeatureNames(context map[string]interface{}) interface{} {
	step := ListMappingTransformer{
		mapping: context["feature-mapping"].(map[string]string),
	}
	v, err := step.Transform((context)["raw-features"])
	if err != nil {
		panic(err)
	}
	return v
}
