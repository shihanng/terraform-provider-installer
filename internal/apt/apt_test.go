package apt_test

import (
	"testing"

	"github.com/shihanng/terraform-provider-installer/internal/apt"
	"gotest.tools/v3/assert"
)

func TestGetInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected apt.Info
	}{
		{
			input: "",
			expected: apt.Info{
				Name:    "",
				Version: "",
			},
		},
		{
			input: "nginx",
			expected: apt.Info{
				Name:    "nginx",
				Version: "",
			},
		},
		{
			input: "nginx=1.18.0-6ubuntu14.3",
			expected: apt.Info{
				Name:    "nginx",
				Version: "1.18.0-6ubuntu14.3",
			},
		},
		{
			input: "nginx=1.18.0-6ubuntu14.3=abc",
			expected: apt.Info{
				Name:    "nginx",
				Version: "1.18.0-6ubuntu14.3=abc",
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			actual := apt.GetInfo(tc.input)
			assert.DeepEqual(t, actual, tc.expected)
		})
	}
}

func TestExtractVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input: `
Package: nginx
Status: install ok installed
Priority: optional
Section: httpd
Installed-Size: 49
Maintainer: Ubuntu Developers <ubuntu-devel-discuss@lists.ubuntu.com>
Architecture: arm64
Version: 1.18.0-6ubuntu14.3
Depends: nginx-core (<< 1.18.0-6ubuntu14.3.1~) | nginx-full (<< 1.18.0-6ubuntu14.3.1~) | nginx-light (<< 1.18.0-6ubuntu14.3.1~) | nginx-extras (<< 1.18.0-6ubuntu14.3.1~), nginx-core (>= 1.18.0-6ubuntu14.3) | nginx-full (>= 1.18.0-6ubuntu14.3) | nginx-light (>= 1.18.0-6ubuntu14.3) | nginx-extras (>= 1.18.0-6ubuntu14.3)
Breaks: libnginx-mod-http-lua (<< 1.18.0-6ubuntu5)
Description: small, powerful, scalable web/proxy server
 Nginx ("engine X") is a high-performance web and reverse proxy server
 created by Igor Sysoev. It can be used both as a standalone web server
 and as a proxy to reduce the load on back-end HTTP or mail servers.
 .
 This is a dependency package to install either nginx-core (by default),
 nginx-full, nginx-light or nginx-extras.
Homepage: https://nginx.net
Original-Maintainer: Debian Nginx Maintainers <pkg-nginx-maintainers@alioth-lists.debian.net>
`,
			expected: "1.18.0-6ubuntu14.3",
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			actual := apt.ExtractVersion(tc.input)
			assert.Equal(t, actual, tc.expected)
		})
	}
}
