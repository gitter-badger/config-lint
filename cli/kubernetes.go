package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/stelligent/config-lint/assertion"
	"io/ioutil"
	"os"
	"path/filepath"
)

// KubernetesLinter lints resources in Kubernets YAML files
type KubernetesLinter struct {
	BaseLinter
	Log assertion.LoggingFunction
}

// KubernetesResourceLoader converts Terraform configuration files into a collection of Resource objects
type KubernetesResourceLoader struct {
	Log assertion.LoggingFunction
}

func loadYAML(filename string, log assertion.LoggingFunction) ([]interface{}, error) {
	empty := []interface{}{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, filename, err.Error())
		return empty, err
	}

	var yamlData interface{}
	err = yaml.Unmarshal(content, &yamlData)
	if err != nil {
		fmt.Fprintln(os.Stderr, filename, err.Error())
		return empty, err
	}
	m := yamlData.(map[string]interface{})
	return []interface{}{m}, nil
}

func getResourceIDFromMetadata(m map[string]interface{}) (string, bool) {
	if metadata, ok := m["metadata"].(map[string]interface{}); ok {
		if name, ok := metadata["name"].(string); ok {
			return name, true
		}
	}
	return "", false
}

func getResourceIDFromFilename(filename string) string {
	_, resourceID := filepath.Split(filename)
	return resourceID
}

// Load converts a text file into a collection of Resource objects
func (l KubernetesResourceLoader) Load(filename string) []assertion.Resource {
	resources := make([]assertion.Resource, 0)
	yamlResources, _ := loadYAML(filename, l.Log)
	for _, resource := range yamlResources {
		m := resource.(map[string]interface{})
		var resourceID string
		if name, ok := getResourceIDFromMetadata(m); ok {
			resourceID = name
		} else {
			resourceID = getResourceIDFromFilename(filename)
		}
		kr := assertion.Resource{
			ID:         resourceID,
			Type:       m["kind"].(string),
			Properties: m,
			Filename:   filename,
		}
		resources = append(resources, kr)
	}
	return resources
}

// Validate runs validate on a collection of filenames using a RuleSet
func (l KubernetesLinter) Validate(filenames []string, ruleSet assertion.RuleSet, tags []string, ruleIDs []string) ([]string, []assertion.Violation) {
	loader := KubernetesResourceLoader{Log: l.Log}
	return l.ValidateFiles(filenames, ruleSet, tags, ruleIDs, loader, l.Log)
}

// Search evaluates a JMESPath expression against the resources in a collection of filenames
func (l KubernetesLinter) Search(filenames []string, ruleSet assertion.RuleSet, searchExpression string) {
	loader := KubernetesResourceLoader{Log: l.Log}
	l.SearchFiles(filenames, ruleSet, searchExpression, loader)
}
