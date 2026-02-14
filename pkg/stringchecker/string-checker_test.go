package stringchecker_test

import (
	"testing"

	"github.com/Grisha1Kadetov/BeautyLogs/pkg/stringchecker"
)

func TestCheckFirstLowercase(t *testing.T) {
	tests := []struct {
		in     string
		want   string
		wantOK bool
	}{
		{"", "", true},
		{"hello", "hello", true},
		{"Hello", "hello", false},
		{"  Hello", "  hello", false},
		{"123", "123", true},
		{"😊Hello", "😊hello", false},
	}

	for _, test := range tests {
		got, ok := stringchecker.CheckFirstLowercase(test.in)
		if got != test.want || ok != test.wantOK {
			t.Fatalf("CheckFirstLowercase(%q) = (%q,%v), want (%q,%v)",
				test.in, got, ok, test.want, test.wantOK)
		}
	}
}

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"Hello world", true},
		{"Hello123", true},
		{"Привет", false},
		{"HelloПривет", false},
		{"é", true},
	}

	for _, test := range tests {
		got := stringchecker.CheckEnglish(test.in)
		if got != test.want {
			t.Fatalf("CheckEnglish(%q) = %v, want %v", test.in, got, test.want)
		}
	}
}

func TestCheckSpecial(t *testing.T) {
	ignore := map[rune]any{
		':': true,
		'!': true,
	}

	tests := []struct {
		name   string
		in     string
		ignore map[rune]any
		want   bool
	}{
		{"base", "Hello 123", nil, true},
		{"punct_without_ignore", "Hello!", nil, false},
		{"punct_with_ignore", "Hello!", ignore, true},
		{"colon_with_ignore", "key: value", ignore, true},
		{"emoji", "Hello 😊", ignore, false},
	}

	for _, test := range tests {
		got := stringchecker.CheckSpecial(test.in, test.ignore)
		if got != test.want {
			t.Fatalf("%s: CheckSpecial(%q, ignore) = %v, want %v", test.name, test.in, got, test.want)
		}
	}
}

func TestCheckSensitive(t *testing.T) {
	keys := []string{"password", "api_key", "username"}

	tests := []struct {
		name string
		in   string
		keys []string
		want bool
	}{
		{"nil_keys", "password", nil, true},
		{"match", "password", keys, false},
		{"match_case", "PASSWORD", keys, false},
		{"match_snake_case", "my api_key is 123", keys, false},
		{"match_camel_case", "my Username", keys, false},
		{"no_match", "hello world", keys, true},
		{"solid", "mypassword", keys, true},
	}

	for _, test := range tests {
		got := stringchecker.CheckSensitive(test.in, test.keys)
		if got != test.want {
			t.Fatalf("%s: CheckSensitive(%q, keys) = %v, want %v", test.name, test.in, got, test.want)
		}
	}
}
