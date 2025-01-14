package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"os"
)

type Args struct {
	ConfigDir          string   `arg:"positional"`
	FeatureFile        string   `arg:"--feature-file,required"`
	FeatureMappingFile string   `arg:"--feature-mapping-file,required"`
	ConfigFile         string   `arg:"--config-file" default:"config.yaml"`
	FeatureSet         []string `arg:"--feature-set"`
}

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
	// 1. Read features from a plain text file
	context["raw-features"] = readFeatures(context)
	// 2. Read feature mapping from a properties file
	context["feature-mapping"] = readFeatureMapping(context)
	// 3. Convert feature names using the mapping
	context["features"] = convertFeatureNames(context)
	// 4. Read config.yaml from a YAML file
	context["config"] = readConfigFile(context)
	// 5. Expand according to configuration in config.yaml
	context["features"] = expandFeatureByConfig(context)
	//  6. Filter the feature based on config#feature-set
	//  7. Sort the features based on config#priority

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

func readConfigFile(context map[string]interface{}) interface{} {
	step := YamlInputSource{
		path: context["args"].(Args).ConfigFile,
	}
	value, err := step.Provide(os.DirFS(context["args"].(Args).ConfigDir))
	if err != nil {
		panic(err)
	}
	return value
}

func expandFeatureByConfig(context map[string]interface{}) interface{} {
	step := ListExpandTransformer{
		dataByKey: (context["config"]).(map[interface{}]interface{}),
		keyMapper: IdentityMapper,
	}
	value, err := step.Transform(context["features"])
	if err != nil {
		panic(err)
	}
	return value
}
