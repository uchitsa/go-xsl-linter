package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/beevik/etree"
	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
)

type Rule struct {
	Message  string
	XPath    string
	Priority int
}

func loadRules(rulesFile string) (map[string]Rule, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(rulesFile); err != nil {
		return nil, err
	}

	rules := make(map[string]Rule)
	for _, ruleElem := range doc.FindElements(".//rule") {
		name := ruleElem.SelectAttrValue("name", "")
		message := ruleElem.FindElement("message").Text()
		xpath := ruleElem.FindElement("xpath").Text()
		priority, _ := strconv.Atoi(ruleElem.FindElement("priority").Text())

		rules[name] = Rule{
			Message:  message,
			XPath:    xpath,
			Priority: priority,
		}
	}

	return rules, nil
}

func lintXSLT(xsltFile string, rules map[string]Rule) error {
	doc, err := xmltree.ParseXML(xmltree.ParseXMLWithNS(xsltFile, nil))
	if err != nil {
		return err
	}

	for name, rule := range rules {
		xpathExpr := rule.XPath
		xpath := goxpath.MustParse(xpathExpr)
		res, err := xpath.ExecNode(doc)
		if err != nil {
			return err
		}

		if len(res) > 0 {
			fmt.Printf("Rule '%s' violated: %s\n", name, rule.Message)
			for _, r := range res {
				if node, ok := r.(tree.Node); ok {
					fmt.Printf("  at %s line %d\n", node.GetElement().Tag, node.GetElement().Line)
				}
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: xsl_linter <xslt_file> <rules_file>")
		os.Exit(1)
	}

	xsltFile := os.Args[1]
	rulesFile := os.Args[2]

	rules, err := loadRules(rulesFile)
	if err != nil {
		fmt.Println("Error loading rules:", err)
		os.Exit(1)
	}

	if err := lintXSLT(xsltFile, rules); err != nil {
		fmt.Println("Error linting XSLT:", err)
		os.Exit(1)
	}
}
