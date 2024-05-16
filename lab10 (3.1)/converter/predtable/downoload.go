package predtable

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
)

var tableView = `
package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
%s
	}
}

`

var axiomView = `
func newGenAxiom() string {
	return "%s"
}

`

func UploadTableToFile(filepath string, genTable map[string][]string, rules semantic.Rules) error {
	newFileView := makeTableView(genTable) + makeAxiomView(rules.Axiom)
	return uploadNewFile(filepath, newFileView)
}

func makeTableView(genTable map[string][]string) string {
	keys := make([]string, 0, len(genTable))
	for key := range genTable {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var lines []string
	for _, key := range keys {
		values := genTable[key]
		var quotedValues []string
		for _, value := range values {
			if value != "" {
				quotedValues = append(quotedValues, fmt.Sprintf(`"%s"`, value))
			}
		}
		line := fmt.Sprintf("\t\t\"%s\": {%s},\n", key, strings.Join(quotedValues, ", "))
		lines = append(lines, line)
	}

	return fmt.Sprintf(tableView, strings.Join(lines, ""))
}

func makeAxiomView(axiom string) string {
	return fmt.Sprintf(axiomView, axiom)
}

func uploadNewFile(filepath, view string) error {
	err := os.WriteFile(filepath, []byte(view), 0644)
	if err != nil {
		return fmt.Errorf("failed uploud view: %w", err)
	}
	return nil
}
