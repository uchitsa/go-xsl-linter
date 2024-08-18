package main

import (
	"testing"

	"github.com/beevik/etree"
)

func TestLoadRules(t *testing.T) {
	rulesXML := `
	<rules>
		<rule name="testRule">
			<message>Test Message</message>
			<xpath>/test/path</xpath>
			<priority>1</priority>
		</rule>
	</rules>
	`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(rulesXML); err != nil {
		t.Fatal("Failed to parse test XML:", err)
	}

	rules, err := loadRules("testdata/rules.xml")
	if err != nil {
		t.Fatal("Failed to load rules:", err)
	}

	if len(rules) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(rules))
	}

	rule, ok := rules["testRule"]
	if !ok {
		t.Fatal("Expected rule 'testRule' not found")
	}

	if rule.Message != "Test Message" {
		t.Errorf("Expected message 'Test Message', got '%s'", rule.Message)
	}

	if rule.XPath != "/test/path" {
		t.Errorf("Expected XPath '/test/path', got '%s'", rule.XPath)
	}

	if rule.Priority != 1 {
		t.Errorf("Expected priority 1, got %d", rule.Priority)
	}
}
