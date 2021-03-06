package options

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/zalando/planb-tokeninfo/processor"
)

// The Settings type contains the application configurable options
type Settings struct {
	ListenAddress                     string
	MetricsListenAddress              string
	UpstreamTokenInfoURL              *url.URL
	UpstreamTimeout                   time.Duration
	UpstreamCacheMaxSize              int64
	UpstreamCacheTTL                  time.Duration
	OpenIDProviderConfigurationURL    *url.URL
	OpenIDProviderRefreshInterval     time.Duration
	HTTPClientTimeout                 time.Duration
	HTTPClientTLSTimeout              time.Duration
	RevocationCacheTTL                time.Duration
	RevocationProviderRefreshInterval time.Duration
	RevocationRefreshTolerance        time.Duration
	RevocationProviderUrl             *url.URL
	HashingSalt                       string
	JwtProcessors                     map[string]processor.JwtProcessor
}

const (
	defaultListenAddress                 = ":9021"
	defaultMetricsListenAddress          = ":9020"
	defaultUpstreamCacheMaxSize          = 10000
	defaultUpstreamCacheTTL              = 60 * time.Second
	defaultUpstreamTimeout               = 1 * time.Second
	defaultOpenIDRefreshInterval         = 30 * time.Second
	defaultHTTPClientTimeout             = 10 * time.Second
	defaultHTTPClientTLSTimeout          = 10 * time.Second
	defaultRevocationCacheTTL            = 30 * 24 * time.Hour
	defaultRevokeProviderRefreshInterval = 10 * time.Second
	defaultRevocationRereshTolerance     = 60 * time.Second
	defaultHashingSalt                   = "seasaltisthebest"
)

var (
	// AppSettings is a global variable that holds the application settings
	AppSettings = defaultSettings()
)

func defaultSettings() *Settings {
	return &Settings{
		ListenAddress:                     defaultListenAddress,
		MetricsListenAddress:              defaultMetricsListenAddress,
		UpstreamCacheMaxSize:              defaultUpstreamCacheMaxSize,
		UpstreamCacheTTL:                  defaultUpstreamCacheTTL,
		UpstreamTimeout:                   defaultUpstreamTimeout,
		OpenIDProviderRefreshInterval:     defaultOpenIDRefreshInterval,
		HTTPClientTimeout:                 defaultHTTPClientTimeout,
		HTTPClientTLSTimeout:              defaultHTTPClientTLSTimeout,
		RevocationCacheTTL:                defaultRevocationCacheTTL,
		RevocationProviderRefreshInterval: defaultRevokeProviderRefreshInterval,
		RevocationRefreshTolerance:        defaultRevocationRereshTolerance,
		HashingSalt:                       defaultHashingSalt,
		JwtProcessors:                     make(map[string]processor.JwtProcessor),
	}
}

// LoadFromEnvironment will try to load all the options from environment variables.
// It will return an error if the required options are not available. The required environment
// variables are:
//
//      UPSTREAM_TOKENINFO_URL
//      OPENID_PROVIDER_CONFIGURATION_URL
//	REVOCATION_PROVIDER_URL
//
// The remaining options have sane defaults and are not mandatory
func LoadFromEnvironment() error {
	settings := defaultSettings()

	if s := getString("UPSTREAM_TOKENINFO_URL", ""); s != "" {
		tokeninfoURL, err := getURL("UPSTREAM_TOKENINFO_URL")
		if err != nil {
			return fmt.Errorf("Error with UPSTREAM_TOKENINFO_URL: %v\n", err)
		}
		settings.UpstreamTokenInfoURL = tokeninfoURL
	}

	openIDConfiguration, err := getURL("OPENID_PROVIDER_CONFIGURATION_URL")
	if err != nil || openIDConfiguration == nil {
		return fmt.Errorf("Invalid OPENID_PROVIDER_CONFIGURATION_URL: %v\n", err)
	}
	settings.OpenIDProviderConfigurationURL = openIDConfiguration

	revocationURL, err := getURL("REVOCATION_PROVIDER_URL")
	if err != nil || revocationURL == nil {
		return fmt.Errorf("Invalid REVOCATION_PROVIDER_URL: %v\n", err)
	}
	settings.RevocationProviderUrl = revocationURL

	if s := getString("REVOCATION_HASHING_SALT", ""); s != "" {
		settings.HashingSalt = s
	}

	if s := getString("LISTEN_ADDRESS", ""); s != "" {
		settings.ListenAddress = s
	}

	if s := getString("METRICS_LISTEN_ADDRESS", ""); s != "" {
		settings.MetricsListenAddress = s
	}

	if i := getInt("UPSTREAM_CACHE_MAX_SIZE", -1); i > -1 {
		settings.UpstreamCacheMaxSize = int64(i)
	}

	if d := getDuration("UPSTREAM_CACHE_TTL", -1); d > -1 {
		settings.UpstreamCacheTTL = d
	}

	if d := getDuration("UPSTREAM_TIMEOUT", -1); d > -1 {
		settings.UpstreamTimeout = d
	}

	if d := getDuration("OPENID_PROVIDER_REFRESH_INTERVAL", 0); d > 0 {
		settings.OpenIDProviderRefreshInterval = d
	}

	if d := getDuration("HTTP_CLIENT_TIMEOUT", 0); d > 0 {
		settings.HTTPClientTimeout = d
	}

	if d := getDuration("HTTP_CLIENT_TLS_TIMEOUT", 0); d > 0 {
		settings.HTTPClientTLSTimeout = d
	}

	if d := getDuration("REVOCATION_CACHE_TTL", 0); d > 0 {
		settings.RevocationCacheTTL = d
	}

	if d := getDuration("REVOCATION_PROVIDER_REFRESH_INTERVAL", 0); d > 0 {
		settings.RevocationProviderRefreshInterval = d
	}

	if d := getDuration("REVOCATION_REFRESH_TOLERANCE", 0); d > 0 {
		settings.RevocationRefreshTolerance = d
	}

	AppSettings = settings
	return nil
}

func getString(v string, def string) string {
	s, ok := os.LookupEnv(v)
	if !ok {
		return def
	}
	return s
}

func getURL(v string) (*url.URL, error) {
	u, ok := os.LookupEnv(v)
	if !ok || u == "" {
		return nil, fmt.Errorf("Missing URL setting: %q", v)
	}
	return url.Parse(u)
}

func getInt(v string, def int) int {
	s, ok := os.LookupEnv(v)
	if !ok {
		return def
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return i
}

func getDuration(v string, def time.Duration) time.Duration {
	s, ok := os.LookupEnv(v)
	if !ok || s == "" {
		return def
	}

	if d, err := time.ParseDuration(s); err == nil {
		return d
	}

	if seconds, err := strconv.Atoi(s); err == nil {
		return time.Duration(seconds) * time.Second
	}

	return def
}
