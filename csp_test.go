package csp

import (
	"net/url"
	"strings"
	"testing"
)

func checkErr(t *testing.T, i int, got error, want string) {
	if got == nil && want == "" {
		return
	} else if got != nil && want == "" {
		t.Fatalf("%d. unexpected error %+v", i, got)
	} else if got == nil && want != "" {
		t.Fatalf("%d. expected error matching %q", i, want)
	}

	if !strings.Contains(got.Error(), want) {
		t.Fatalf("%d. expected error matching %q; got %+v", i, want, got)
	}
}

func TestCSP(t *testing.T) {
	cases := []struct {
		policy                 string
		page                   string
		html                   string
		valid                  bool
		policyErr, validateErr string
	}{
		{
			policy: "default-src 'self'",
			page:   "https://google.com",
			html:   `<script src="https://google.com"></script>`,
			valid:  true,
		},
		{
			policy: "default-src 'none'",
			page:   "https://google.com",
			html:   `<script src="https://google.com"></script>`,
			valid:  false,
		},
		{
			policy: "script-src *; default-src 'none'",
			page:   "https://google.com",
			html:   `<script src="https://google.com"></script>`,
			valid:  true,
		},
		{
			policy: "default-src http:",
			page:   "https://google.com",
			html:   `<script src="https://google.com"></script>`,
			valid:  true,
		},
		{
			policy: "default-src google.com",
			page:   "https://google.com",
			html:   `<script src="https://google.com"></script>`,
			valid:  true,
		},
		{
			policy: "default-src *.google.com",
			page:   "https://google.com",
			html:   `<script src="https://www.google.com"></script>`,
			valid:  true,
		},
		{
			policy: "default-src 'nonce-foo'",
			page:   "https://google.com",
			html:   `<script src="https://www.google.com"></script>`,
			valid:  false,
		},
		{
			policy: "default-src 'nonce-foo'",
			page:   "https://google.com",
			html:   `<script nonce="foo" src="https://www.google.com"></script>`,
			valid:  true,
		},
		{
			policy: "default-src 'sha256-LCa0a2j/xo/5m0U8HTBBNBNCLXBkg7+g+YpeiGJm564='",
			page:   "https://google.com",
			html:   `<script>foo</script>`,
			valid:  true,
		},
		{
			policy: "default-src 'sha256-LCa0a2j/xo/5m0U8HTBBNBNCLXBkg7+g+YpeiGJm564='",
			page:   "https://google.com",
			html:   `<script>bar</script>`,
			valid:  false,
		},
		{
			policy: "default-src 'unsafe-inline'",
			page:   "https://google.com",
			html:   `<script>bar</script>`,
			valid:  true,
		},
	}

	for i, c := range cases {
		p, err := ParsePolicy(c.policy)
		checkErr(t, i, err, c.policyErr)
		page, err := url.Parse(c.page)
		if err != nil {
			t.Fatal(err)
		}
		valid, reports, err := ValidatePage(p, *page, strings.NewReader(c.html))
		checkErr(t, i, err, c.validateErr)
		if valid != c.valid {
			t.Errorf("%d. ValidatePage(...) = %v; not %v; reports = %+v", i, valid, c.valid, reports)
		}
	}
}