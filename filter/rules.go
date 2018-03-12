package filter

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func LoadRules(filename string) string {
	rules, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(rules)
}

func MustParseRules(rules string) RuleSet {
	r := RuleSet{}
	err := yaml.Unmarshal([]byte(rules), &r)
	if err != nil {
		panic(err)
	}
	return r
}

func FilterRulesByTag(rules []Rule, tags []string) []Rule {
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		if tags == nil || listsIntersect(tags, rule.Tags) {
			filteredRules = append(filteredRules, rule)
		}
	}
	return filteredRules
}

func FilterRulesById(rules []Rule, ruleIds []string) []Rule {
	if len(ruleIds) == 0 {
		return rules
	}
	filteredRules := make([]Rule, 0)
	for _, rule := range rules {
		for _, id := range ruleIds {
			if id == rule.Id {
				filteredRules = append(filteredRules, rule)
			}
		}
	}
	return filteredRules
}