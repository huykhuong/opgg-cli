package prompt

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed *.md
var promptFS embed.FS

func ReadPrompt(templateName string, data any) (string, error) {
	content, err := promptFS.ReadFile(fmt.Sprintf("%s.md", templateName))
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
