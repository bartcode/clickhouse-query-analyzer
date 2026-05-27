package clickhouse

import "testing"

func TestNativeToHTTPPort(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"9000", "8123"},
		{"", "8123"},
		{"9440", "8443"},
		{"19000", "18123"},
		{"invalid", "8123"},
	}

	for _, tt := range tests {
		if got := nativeToHTTPPort(tt.input); got != tt.expected {
			t.Errorf("nativeToHTTPPort(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestColumnsFromNamesTypes(t *testing.T) {
	tests := []struct {
		name     string
		names    []string
		types    []string
		expected []ColumnInfo
	}{
		{
			"names and types match",
			[]string{"id", "name"},
			[]string{"UInt64", "String"},
			[]ColumnInfo{{Name: "id", Type: "UInt64"}, {Name: "name", Type: "String"}},
		},
		{
			"nil types",
			[]string{"id", "name"},
			nil,
			[]ColumnInfo{{Name: "id", Type: ""}, {Name: "name", Type: ""}},
		},
		{
			"fewer types than names",
			[]string{"id", "name", "value"},
			[]string{"UInt64"},
			[]ColumnInfo{{Name: "id", Type: "UInt64"}, {Name: "name", Type: ""}, {Name: "value", Type: ""}},
		},
		{
			"empty",
			[]string{},
			[]string{},
			[]ColumnInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := columnsFromNamesTypes(tt.names, tt.types)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d columns, got %d", len(tt.expected), len(got))
			}
			for i, col := range got {
				if col.Name != tt.expected[i].Name || col.Type != tt.expected[i].Type {
					t.Errorf("column %d: got {%s, %s}, want {%s, %s}", i, col.Name, col.Type, tt.expected[i].Name, tt.expected[i].Type)
				}
			}
		})
	}
}

func TestIsProbablyJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`{"key": "value"}`, true},
		{`[1, 2, 3]`, true},
		{`  {"key": "value"}  `, true},
		{`  [1, 2, 3]  `, true},
		{"plain text", false},
		{"", false},
		{"   ", false},
		{"<xml>", false},
	}

	for _, tt := range tests {
		if got := isProbablyJSON(tt.input); got != tt.expected {
			t.Errorf("isProbablyJSON(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}
