This utility uses the template module of go to convert specific kinds of input together.

The conversion steps are:
  1. Read features from a plain text file
  2. Read feature mapping from a properties file
  3. Convert feature names using the mapping
  4. Read config.yaml from a YAML file
  5. Expand according to configuration in config.yaml
  6. Read properties from property files
  7. For each feature set, filter the features
    a. Filter the feature based on config#featureSet and put it to per-feature set context variable
    b. Sort the features based on config#priority
    c. Prepare the context for rendering
    d. Render the template for feature set
    
Example usage:

    $ go run main.go ./samples/configs \
                     --feature-file features.txt \
                     --feature-mapping-file feature-rename.properties \
                     --feature-set one \
                     --property-file prop1.properties \
                     --template-dir ./samples/templates
