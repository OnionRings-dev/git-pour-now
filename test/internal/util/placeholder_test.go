package util_test

import (
	"testing"

	"github.com/OnionRings-dev/git-pour-now/internal/util"
)

func TestExpandMessagePlaceholdersBasic(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		branch   string
		parent   string
		expected string
	}{
		{
			name:     "branch name placeholder",
			message:  "Merge %b",
			branch:   "feature/my-feature",
			parent:   "develop",
			expected: "Merge feature/my-feature",
		},
		{
			name:     "full branch refname placeholder",
			message:  "Merge %B",
			branch:   "feature/my-feature",
			parent:   "develop",
			expected: "Merge refs/heads/feature/my-feature",
		},
		{
			name:     "parent name placeholder",
			message:  "Merge into %p",
			branch:   "feature/my-feature",
			parent:   "develop",
			expected: "Merge into develop",
		},
		{
			name:     "full parent refname placeholder",
			message:  "Merge into %P",
			branch:   "feature/my-feature",
			parent:   "develop",
			expected: "Merge into refs/heads/develop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ExpandMessagePlaceholders(tt.message, tt.branch, tt.parent)
			if result != tt.expected {
				t.Errorf("ExpandMessagePlaceholders(%q, %q, %q) = %q, want %q",
					tt.message, tt.branch, tt.parent, result, tt.expected)
			}
		})
	}
}

func TestExpandMessagePlaceholdersMultiple(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		branch   string
		parent   string
		expected string
	}{
		{
			name:     "branch and parent",
			message:  "Merge %b into %p",
			branch:   "feature/login",
			parent:   "develop",
			expected: "Merge feature/login into develop",
		},
		{
			name:     "all placeholders",
			message:  "Merge %b (%B) into %p (%P)",
			branch:   "hotfix/1.0.1",
			parent:   "main",
			expected: "Merge hotfix/1.0.1 (refs/heads/hotfix/1.0.1) into main (refs/heads/main)",
		},
		{
			name:     "repeated placeholders",
			message:  "%b %b %p %p",
			branch:   "release/2.0",
			parent:   "main",
			expected: "release/2.0 release/2.0 main main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ExpandMessagePlaceholders(tt.message, tt.branch, tt.parent)
			if result != tt.expected {
				t.Errorf("ExpandMessagePlaceholders(%q, %q, %q) = %q, want %q",
					tt.message, tt.branch, tt.parent, result, tt.expected)
			}
		})
	}
}

func TestExpandMessagePlaceholdersEscaped(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		branch   string
		parent   string
		expected string
	}{
		{
			name:     "escaped percent sign",
			message:  "100%% complete",
			branch:   "feature/test",
			parent:   "develop",
			expected: "100% complete",
		},
		{
			name:     "escaped with placeholders",
			message:  "Merge %b (100%% done) into %p",
			branch:   "feature/test",
			parent:   "develop",
			expected: "Merge feature/test (100% done) into develop",
		},
		{
			name:     "multiple escaped",
			message:  "%%b is %b, %%p is %p",
			branch:   "feature/x",
			parent:   "main",
			expected: "%b is feature/x, %p is main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ExpandMessagePlaceholders(tt.message, tt.branch, tt.parent)
			if result != tt.expected {
				t.Errorf("ExpandMessagePlaceholders(%q, %q, %q) = %q, want %q",
					tt.message, tt.branch, tt.parent, result, tt.expected)
			}
		})
	}
}

func TestExpandMessagePlaceholdersNoPlaceholders(t *testing.T) {
	message := "Simple merge message without placeholders"
	result := util.ExpandMessagePlaceholders(message, "feature/test", "develop")
	if result != message {
		t.Errorf("ExpandMessagePlaceholders(%q, ...) = %q, want unchanged", message, result)
	}
}

func TestExpandMessagePlaceholdersEmpty(t *testing.T) {
	result := util.ExpandMessagePlaceholders("", "feature/test", "develop")
	if result != "" {
		t.Errorf("ExpandMessagePlaceholders(\"\", ...) = %q, want empty string", result)
	}
}
