package pb

import (
	"reflect"
	"testing"
)

func TestServerAddresses_ToAddressMapOrSrv_shouldRemovePrefix(t *testing.T) {
	str := ServerAddresses("dnssrv+hello.srv.consul")

	d := str.ToServiceDiscovery()

	expected := ServerSrvAddress("hello.srv.consul")
	if *d.srvRecord != expected {
		t.Fatalf(`ServerAddresses("dnssrv+hello.srv.consul") = %s, expected %s`, *d.srvRecord, expected)
	}
}

func TestServerAddresses_ToAddressMapOrSrv_shouldHandleIPPortList(t *testing.T) {
	str := ServerAddresses("10.0.0.1:23,10.0.0.2:24")

	d := str.ToServiceDiscovery()

	if d.srvRecord != nil {
		t.Fatalf(`ServerAddresses("dnssrv+hello.srv.consul") = %s, expected nil`, *d.srvRecord)
	}

	expected := []ServerAddress{
		ServerAddress("10.0.0.1:23"),
		ServerAddress("10.0.0.2:24"),
	}

	if !reflect.DeepEqual(d.list, expected) {
		t.Fatalf(`Expected %q, got %q`, expected, d.list)
	}
}

func TestServerAddress_ToHost(t *testing.T) {
	testCases := []struct {
		name     string
		address  ServerAddress
		expected string
	}{
		{
			name:     "hostname with port",
			address:  ServerAddress("master.example.com:9333"),
			expected: "master.example.com",
		},
		{
			name:     "IPv4 with port",
			address:  ServerAddress("192.168.1.1:9333"),
			expected: "192.168.1.1",
		},
		{
			name:     "IPv6 with port",
			address:  ServerAddress("[2001:db8::1]:9333"),
			expected: "2001:db8::1",
		},
		{
			name:     "hostname without port",
			address:  ServerAddress("master.example.com"),
			expected: "master.example.com",
		},
		{
			name:     "hostname with port.grpcPort",
			address:  ServerAddress("master.example.com:443.10443"),
			expected: "master.example.com",
		},
		{
			name:     "IPv4 with port.grpcPort",
			address:  ServerAddress("192.168.1.1:8080.18080"),
			expected: "192.168.1.1",
		},
		{
			name:     "IPv6 with port.grpcPort",
			address:  ServerAddress("[2001:db8::1]:8080.18080"),
			expected: "2001:db8::1",
		},
		{
			name:     "bracketed IPv6 without port",
			address:  ServerAddress("[2001:db8::1]"),
			expected: "2001:db8::1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.address.ToHost()
			if got != tc.expected {
				t.Errorf("ServerAddress(%q).ToHost() = %q, want %q", tc.address, got, tc.expected)
			}
		})
	}
}

func TestIPv6ServerAddressFormatting(t *testing.T) {
	testCases := []struct {
		name         string
		sa           ServerAddress
		expectedHttp string
		expectedGrpc string
	}{
		{
			name:         "unbracketed IPv6",
			sa:           NewServerAddress("2001:db8::1", 8080, 18080),
			expectedHttp: "[2001:db8::1]:8080",
			expectedGrpc: "[2001:db8::1]:18080",
		},
		{
			name:         "bracketed IPv6",
			sa:           NewServerAddressWithGrpcPort("[2001:db8::1]:8080", 18080),
			expectedHttp: "[2001:db8::1]:8080",
			expectedGrpc: "[2001:db8::1]:18080",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if httpAddr := tc.sa.ToHttpAddress(); httpAddr != tc.expectedHttp {
				t.Errorf("%s: ToHttpAddress() = %s, want %s", tc.name, httpAddr, tc.expectedHttp)
			}
			if grpcAddr := tc.sa.ToGrpcAddress(); grpcAddr != tc.expectedGrpc {
				t.Errorf("%s: ToGrpcAddress() = %s, want %s", tc.name, grpcAddr, tc.expectedGrpc)
			}
		})
	}
}

// TestZapPort pins the overflow-safe ZAP-port derivation. The nominal offset is
// +10000, but a high grpcPort (an ephemeral test port near 55535, as
// AllocateMiniPorts can hand out) must not push the result past 65535 — that
// produced an invalid port and the master Fatalf'd on its ZAP listener, hanging
// cluster readiness. ZapPort folds the offset into [1,65535].
func TestZapPort(t *testing.T) {
	cases := []struct{ grpc, want int }{
		{9333, 19333},  // common case: plain +10000
		{19333, 29333}, // common case (filer grpc derived from http+10000)
		{45535, 55535}, // still in range
		{55535, 65535}, // exactly the top of the range
		{55536, 1},     // first overflow: wraps to a low port instead of 65536
		{60000, 4465},  // mid overflow
		{65535, 10000}, // top grpc port wraps
	}
	for _, c := range cases {
		if got := ZapPort(c.grpc); got != c.want {
			t.Errorf("ZapPort(%d) = %d, want %d", c.grpc, got, c.want)
		}
	}

	// Over the whole valid port space the result is always a legal port and the
	// map is injective, so client and server always agree on the same port and no
	// two grpc ports collide on one ZAP port.
	seen := make(map[int]int, 65535)
	for g := 1; g <= 65535; g++ {
		z := ZapPort(g)
		if z < 1 || z > 65535 {
			t.Fatalf("ZapPort(%d) = %d out of [1,65535]", g, z)
		}
		if prev, dup := seen[z]; dup {
			t.Fatalf("ZapPort not injective: %d and %d both map to %d", prev, g, z)
		}
		seen[z] = g
	}
}

// TestToMasterZapAddress_OverflowSafe proves the master client derives a valid,
// agreed ZAP address even when the grpc port is high enough that grpcPort+10000
// would overflow — the exact condition that broke cluster bring-up.
func TestToMasterZapAddress_OverflowSafe(t *testing.T) {
	// httpPort 50000 -> grpcPort 60000 -> ZapPort(60000) = 4465.
	sa := NewServerAddress("127.0.0.1", 50000, 60000)
	got := sa.ToMasterZapAddress()
	want := "127.0.0.1:4465"
	if got != want {
		t.Fatalf("ToMasterZapAddress() = %s, want %s", got, want)
	}
	// IAM rides the same offset, so it agrees.
	if iam := sa.ToIamZapAddress(); iam != want {
		t.Fatalf("ToIamZapAddress() = %s, want %s", iam, want)
	}
}
