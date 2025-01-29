package main

import (
	"bytes"
	"fmt"
	"github.com/alexflint/go-arg"
	"html/template"
	"io"
	"os"
	"path"
)

type Args struct {
	ConfigDir          string   `arg:"positional"`
	FeatureFile        string   `arg:"--feature-file,required"`
	FeatureMappingFile string   `arg:"--feature-mapping-file,required"`
	ConfigFile         string   `arg:"--config-file" default:"config.yaml"`
	FeatureSet         []string `arg:"--feature-set"`
	PropertyFiles      []string `arg:"--property-file"`
	TemplateDir        string   `arg:"--template-dir,required"`
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
	// 6. Read properties from property files
	properties := readProperties(context)
	// 7. For each feature set, filter the features
	for _, featureSet := range args.FeatureSet {
		// The per-feature set results will be stored at "feature-<featureSet>"
		contextVarName := fmt.Sprintf("feature-%s", featureSet)
		// a. Filter the feature based on config#featureSet and put it to per-feature set context variable
		context[contextVarName] = filterFeatureByFeatureSet(context, featureSet)
		// b. Sort the features based on config#priority
		context[contextVarName] = sortFeatureSetByPriority(context, contextVarName)
		fmt.Printf("Context: %+v\n", context)
		// c. Prepare the context for rendering
		tc := make(map[string]interface{})
		tc["features"] = context[contextVarName]
		// d. Render the template for feature set
		tmpl := prepareTemplateForFeature(args.TemplateDir, featureSet, properties)
		buffer := bytes.Buffer{}
		err := tmpl.Execute(&buffer, tc)
		if err != nil {
			panic(err)
		}
		fmt.Printf("== BEGIN OUTPUT ==\n%s", buffer.String())
	}

}

func prepareTemplateForFeature(templateDir string, featureSet string, properties []interface{}) *template.Template {
	tmplFile := path.Join(templateDir, fmt.Sprintf("%s.tmpl", featureSet))
	f, err := os.Open(tmplFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	functions := prepareTemplateFunctions(properties)
	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New(featureSet).Funcs(functions).Parse(string(bytes))
	if err != nil {
		panic(err)
	}
	return tmpl
}

func prepareTemplateFunctions(properties []interface{}) template.FuncMap {
	lookup := PropertiesLookup{properties: properties}
	functions := template.FuncMap{
		"hasProperty": lookup.HasProperty,
		"getProperty": lookup.GetProperty,
		"add": func(a, b int) int {
			return a + b
		},
	}
	return functions
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
		dataByKey:   (context["config"]).(map[interface{}]interface{}),
		keyMapper:   IdentityMapper,
		keepKeyName: true,
	}
	value, err := step.Transform(context["features"])
	if err != nil {
		panic(err)
	}
	return value
}

// Read properties specified in the input arguments
func readProperties(context map[string]interface{}) []interface{} {
	files := context["args"].(Args).PropertyFiles
	list := make([]interface{}, len(files))
	for i, f := range files {
		step := PropertiesInputSource{
			path: f,
		}
		properties, err := step.Provide(os.DirFS(context["args"].(Args).ConfigDir))
		if err != nil {
			panic(err)
		}
		list[i] = properties
	}
	return list
}

func filterFeatureByFeatureSet(context map[string]interface{}, featureSet string) interface{} {
	step := ListFilterTransformer{
		predicate: MapValuePredicate("feature-set", featureSet),
	}
	value, err := step.Transform(context["features"])
	if err != nil {
		panic(err)
	}
	return value
}

func sortFeatureSetByPriority(context map[string]interface{}, contextVarName string) interface{} {
	step := ListStringSortTransformer{
		mapper: MapValueStringMapper("priority"),
	}
	value, err := step.Transform(context[contextVarName])
	if err != nil {
		panic(err)
	}
	return value
}
