package tools

import (
	"testing"
)

func TestHtmlToText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "<html><body>Hello, world!</body></html>",
			expected: "Hello, world!",
		},
		{
			name:     "Multiple text nodes",
			input:    "<html><body>Hello <p>World</p></body></html>",
			expected: "Hello World",
		},
		{
			name:     "HTML entities",
			input:    "<html>&lt;Hello&gt;</html>",
			expected: "<Hello>",
		},
		{
			name:     "Script tag content",
			input:    "<html><script>alert('Hello');</script></html>",
			expected: "",
		},
		{
			name:     "Style tag content",
			input:    "<html><style type=\"text/css\">body { color: green; }</style></html>",
			expected: "",
		},
		{
			name:     "Nested tags",
			input:    "<html><div><span>Hello</span> <p>World</p></div></html>",
			expected: "Hello World",
		},
		{
			name:     "Multiple spaces and newlines",
			input:    "<html><body>   Hello   \nWorld  </body></html>",
			expected: "Hello World",
		},
		{
			name:     "Empty HTML",
			input:    "<html></html>",
			expected: "",
		},
		{
			name:     "No text nodes",
			input:    "<html><body></body></html>",
			expected: "",
		},
		{
			name:     "Text with inline tags",
			input:    "<html>Hello <b>bold</b> text.</html>",
			expected: "Hello bold text.",
		},
		{
			name:     "Self-closing tag",
			input:    "<html><img src='image.jpg' /></html>",
			expected: "",
		},
		{
			name:     "Deeply nested structure",
			input:    "<html><body><div><div><p><span>Deeply</span> nested</p></div></div></body></html>",
			expected: "Deeply nested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := htmlToText([]byte(tt.input))
			if string(output) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(output))
			}
		})
	}
}
