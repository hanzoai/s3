package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/viper"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb"
	"github.com/hanzoai/s3/s3/util"
	util_http_client "github.com/hanzoai/s3/s3/util/http/client"
)

func LoadClientTLSFromFile(configFile string, component string) (pb.DialOption, error) {
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return pb.DialOption{}, fmt.Errorf("failed to read security config %s: %v", configFile, err)
	}
	// Resolve relative PEM paths against the config file's directory.
	configDir := filepath.Dir(configFile)
	for _, key := range []string{"grpc.ca", component + ".cert", component + ".key"} {
		p := v.GetString(key)
		if p != "" && !filepath.IsAbs(p) {
			v.Set(key, filepath.Join(configDir, p))
		}
	}
	return LoadClientTLS(&util.ViperProxy{Viper: v}, component), nil
}

// LoadClientTLS returns the client dial option for the s3 RPC stack: the
// <component> client *tls.Config (presents the component cert, trusts grpc.ca),
// or an empty option (plaintext) when no cert/key is configured. The ZAP
// transport wraps TLS with PQTLSConfig at dial time; the volume/broker grpc
// fallback wraps it with credentials.NewTLS. One config, both transports.
func LoadClientTLS(config *util.ViperProxy, component string) pb.DialOption {
	return pb.DialOption{TLS: pb.ClientTLSConfig(config, component)}
}

// LoadHTTPClientFromFile creates an HTTP client using the https.client TLS
// settings from the given security config file. Returns nil if HTTPS is not
// enabled in the config. This is used by filer.sync to create per-cluster
// HTTP clients when clusters use different certificates.
func LoadHTTPClientFromFile(configFile string) (*util_http_client.HTTPClient, error) {
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read security config %s: %v", configFile, err)
	}

	if !v.GetBool("https.client.enabled") {
		return nil, nil
	}

	configDir := filepath.Dir(configFile)
	resolvePath := func(key string) string {
		p := v.GetString(key)
		if p != "" && !filepath.IsAbs(p) {
			return filepath.Join(configDir, p)
		}
		return p
	}

	return util_http_client.NewHttpClientWithTLS(
		resolvePath("https.client.cert"),
		resolvePath("https.client.key"),
		resolvePath("https.client.ca"),
		v.GetBool("https.client.insecure_skip_verify"),
		util_http_client.AddDialContext,
	)
}

func LoadClientTLSHTTP(clientCertFile string) *tls.Config {
	clientCerts, err := os.ReadFile(clientCertFile)
	if err != nil {
		glog.Fatal(err)
	}
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(clientCerts)
	if !ok {
		glog.Fatalf("Error processing client certificate in %s\n", clientCertFile)
	}

	return &tls.Config{
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
}

func FixTlsConfig(viper *util.ViperProxy, config *tls.Config) error {
	var err error
	config.MinVersion, err = TlsVersionByName(viper.GetString("tls.min_version"))
	if err != nil {
		return err
	}
	config.MaxVersion, err = TlsVersionByName(viper.GetString("tls.max_version"))
	if err != nil {
		return err
	}
	config.CipherSuites, err = TlsCipherSuiteByNames(viper.GetString("tls.cipher_suites"))
	return err
}

func TlsVersionByName(name string) (uint16, error) {
	switch name {
	case "":
		return 0, nil
	case "SSLv3":
		return tls.VersionSSL30, nil
	case "TLS 1.0":
		return tls.VersionTLS10, nil
	case "TLS 1.1":
		return tls.VersionTLS11, nil
	case "TLS 1.2":
		return tls.VersionTLS12, nil
	case "TLS 1.3":
		return tls.VersionTLS13, nil
	default:
		return 0, fmt.Errorf("invalid tls version %s", name)
	}
}

func TlsCipherSuiteByNames(cipherSuiteNames string) ([]uint16, error) {
	cipherSuiteNames = strings.TrimSpace(cipherSuiteNames)
	if cipherSuiteNames == "" {
		return nil, nil
	}
	names := strings.Split(cipherSuiteNames, ",")
	cipherSuites := tls.CipherSuites()
	cipherIds := make([]uint16, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		index := slices.IndexFunc(cipherSuites, func(suite *tls.CipherSuite) bool {
			return name == suite.Name
		})
		if index == -1 {
			return nil, fmt.Errorf("invalid tls cipher suite name %s", name)
		}
		cipherIds = append(cipherIds, cipherSuites[index].ID)
	}
	return cipherIds, nil
}
