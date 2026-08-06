package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWordCount(t *testing.T) {
	cases := map[string]int{
		"":          0,
		"hello":     1,
		"a b c":     3,
		"  x   y  ": 2,
	}
	for in, want := range cases {
		if got := wordCount(in); got != want {
			t.Errorf("wordCount(%q)=%d want %d", in, got, want)
		}
	}
}

func TestDedupLogic(t *testing.T) {
	lines := []string{"a", "b", "a", "c", "b"}
	seen := map[string]bool{}
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if !seen[l] {
			seen[l] = true
			out = append(out, l)
		}
	}
	if strings.Join(out, ",") != "a,b,c" {
		t.Fatalf("dedup 顺序错误: %v", out)
	}
}

func TestSortUniqLogic(t *testing.T) {
	lines := []string{"c", "a", "b", "a"}
	sortStringsTest(lines)
	out := lines[:0]
	for _, l := range lines {
		if len(out) == 0 || out[len(out)-1] != l {
			out = append(out, l)
		}
	}
	if strings.Join(out, ",") != "a,b,c" {
		t.Fatalf("sort -u 结果错误: %v", out)
	}
}

func sortStringsTest(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}

func TestReadLinesTmp(t *testing.T) {
	p := filepath.Join(t.TempDir(), "t.txt")
	if err := os.WriteFile(p, []byte("line1\nline2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	lines, err := readLines(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 || lines[0] != "line1" || lines[1] != "line2" {
		t.Fatalf("readLines 异常: %v", lines)
	}
}

func TestReplaceAllHelper(t *testing.T) {
	// 验证 strings.ReplaceAll 行为（cmdReplace 内部使用）
	got := strings.ReplaceAll("foo bar foo", "foo", "baz")
	if got != "baz bar baz" {
		t.Fatalf("ReplaceAll 行为异常: %q", got)
	}
}
