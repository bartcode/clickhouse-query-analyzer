package clickhouse

import "testing"

func TestShortenSymbol_NoPrefix(t *testing.T) {
	if got := shortenSymbol("SomeClass::method"); got != "SomeClass::method" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestShortenSymbol_WithParens(t *testing.T) {
	if got := shortenSymbol("DB::MergeTreeData::fetchPart(std::shared_ptr)"); got != "MergeTreeData::fetchPart" {
		t.Errorf("expected trimmed at parens, got %q", got)
	}
}

func TestShortenSymbol_WithConst(t *testing.T) {
	if got := shortenSymbol("DB::MergeTreeData const"); got != "MergeTreeData" {
		t.Errorf("expected trimmed at const, got %q", got)
	}
}
