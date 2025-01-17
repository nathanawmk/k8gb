/*
Copyright 2021 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/
package depresolver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AbsaOSS/gopkg/env"
	"github.com/rs/zerolog"
)

// Environment variables keys
const (
	ReconcileRequeueSecondsKey = "RECONCILE_REQUEUE_SECONDS"
	ClusterGeoTagKey           = "CLUSTER_GEO_TAG"
	ExtClustersGeoTagsKey      = "EXT_GSLB_CLUSTERS_GEO_TAGS"
	Route53EnabledKey          = "ROUTE53_ENABLED"
	NS1EnabledKey              = "NS1_ENABLED"
	EdgeDNSServerKey           = "EDGE_DNS_SERVER"
	EdgeDNSServerPortKey       = "EDGE_DNS_SERVER_PORT"
	EdgeDNSZoneKey             = "EDGE_DNS_ZONE"
	DNSZoneKey                 = "DNS_ZONE"
	InfobloxGridHostKey        = "INFOBLOX_GRID_HOST"
	InfobloxVersionKey         = "INFOBLOX_WAPI_VERSION"
	InfobloxPortKey            = "INFOBLOX_WAPI_PORT"
	InfobloxUsernameKey        = "EXTERNAL_DNS_INFOBLOX_WAPI_USERNAME"
	// #nosec G101; ignore false positive gosec; see: https://securego.io/docs/rules/g101.html
	InfobloxPasswordKey            = "EXTERNAL_DNS_INFOBLOX_WAPI_PASSWORD"
	InfobloxHTTPRequestTimeoutKey  = "INFOBLOX_HTTP_REQUEST_TIMEOUT"
	InfobloxHTTPPoolConnectionsKey = "INFOBLOX_HTTP_POOL_CONNECTIONS"
	OverrideFakeInfobloxKey        = "FAKE_INFOBLOX"
	K8gbNamespaceKey               = "POD_NAMESPACE"
	CoreDNSExposedKey              = "COREDNS_EXPOSED"
	LogLevelKey                    = "LOG_LEVEL"
	LogFormatKey                   = "LOG_FORMAT"
	LogNoColorKey                  = "NO_COLOR"
	SplitBrainCheckKey             = "SPLIT_BRAIN_CHECK"
	MetricsAddressKey              = "METRICS_ADDRESS"
)

// ResolveOperatorConfig executes once. It reads operator's configuration
// from environment variables into &Config and validates
func (dr *DependencyResolver) ResolveOperatorConfig() (*Config, error) {
	var recognizedDNSTypes []EdgeDNSType
	dr.onceConfig.Do(func() {
		dr.config = &Config{}
		dr.config.ReconcileRequeueSeconds, _ = env.GetEnvAsIntOrFallback(ReconcileRequeueSecondsKey, 30)
		dr.config.ClusterGeoTag = env.GetEnvAsStringOrFallback(ClusterGeoTagKey, "")
		dr.config.ExtClustersGeoTags = env.GetEnvAsArrayOfStringsOrFallback(ExtClustersGeoTagsKey, []string{})
		dr.config.route53Enabled = env.GetEnvAsBoolOrFallback(Route53EnabledKey, false)
		dr.config.ns1Enabled = env.GetEnvAsBoolOrFallback(NS1EnabledKey, false)
		dr.config.CoreDNSExposed = env.GetEnvAsBoolOrFallback(CoreDNSExposedKey, false)
		dr.config.EdgeDNSServer = env.GetEnvAsStringOrFallback(EdgeDNSServerKey, "")
		dr.config.EdgeDNSServerPort, _ = env.GetEnvAsIntOrFallback(EdgeDNSServerPortKey, 53)
		dr.config.EdgeDNSZone = env.GetEnvAsStringOrFallback(EdgeDNSZoneKey, "")
		dr.config.DNSZone = env.GetEnvAsStringOrFallback(DNSZoneKey, "")
		dr.config.K8gbNamespace = env.GetEnvAsStringOrFallback(K8gbNamespaceKey, "")
		dr.config.Infoblox.Host = env.GetEnvAsStringOrFallback(InfobloxGridHostKey, "")
		dr.config.Infoblox.Version = env.GetEnvAsStringOrFallback(InfobloxVersionKey, "")
		dr.config.Infoblox.Port, _ = env.GetEnvAsIntOrFallback(InfobloxPortKey, 0)
		dr.config.Infoblox.Username = env.GetEnvAsStringOrFallback(InfobloxUsernameKey, "")
		dr.config.Infoblox.Password = env.GetEnvAsStringOrFallback(InfobloxPasswordKey, "")
		dr.config.Infoblox.HTTPPoolConnections, _ = env.GetEnvAsIntOrFallback(InfobloxHTTPPoolConnectionsKey, 10)
		dr.config.Infoblox.HTTPRequestTimeout, _ = env.GetEnvAsIntOrFallback(InfobloxHTTPRequestTimeoutKey, 20)
		dr.config.Override.FakeInfobloxEnabled = env.GetEnvAsBoolOrFallback(OverrideFakeInfobloxKey, false)
		dr.config.Log.Level, _ = zerolog.ParseLevel(strings.ToLower(env.GetEnvAsStringOrFallback(LogLevelKey, zerolog.InfoLevel.String())))
		dr.config.Log.Format = parseLogOutputFormat(strings.ToLower(env.GetEnvAsStringOrFallback(LogFormatKey, SimpleFormat.String())))
		dr.config.Log.NoColor = env.GetEnvAsBoolOrFallback(LogNoColorKey, false)
		dr.config.MetricsAddress = env.GetEnvAsStringOrFallback(MetricsAddressKey, "0.0.0.0:8080")
		dr.config.SplitBrainCheck = env.GetEnvAsBoolOrFallback(SplitBrainCheckKey, false)
		dr.config.EdgeDNSType, recognizedDNSTypes = getEdgeDNSType(dr.config)
		dr.errorConfig = dr.validateConfig(dr.config, recognizedDNSTypes)
	})
	return dr.config, dr.errorConfig
}

func (dr *DependencyResolver) validateConfig(config *Config, recognizedDNSTypes []EdgeDNSType) (err error) {
	const dnsNameMax = 253
	const dnsLabelMax = 63
	if config.Log.Level == zerolog.NoLevel {
		return fmt.Errorf("invalid '%s', allowed values ['','%s','%s','%s','%s','%s','%s','%s']", LogLevelKey,
			zerolog.TraceLevel, zerolog.DebugLevel, zerolog.InfoLevel, zerolog.WarnLevel, zerolog.FatalLevel,
			zerolog.DebugLevel, zerolog.PanicLevel)
	}
	if config.Log.Format == NoFormat {
		return fmt.Errorf("invalid '%s', allowed values ['','%s','%s']", LogFormatKey, JSONFormat, SimpleFormat)
	}
	if config.EdgeDNSType == DNSTypeMultipleProviders {
		return fmt.Errorf("several EdgeDNS recognized %s", recognizedDNSTypes)
	}
	err = field(K8gbNamespaceKey, config.K8gbNamespace).isNotEmpty().matchRegexp(k8sNamespaceRegex).err
	if err != nil {
		return err
	}
	err = field(ReconcileRequeueSecondsKey, config.ReconcileRequeueSeconds).isHigherThanZero().err
	if err != nil {
		return err
	}
	err = field(ClusterGeoTagKey, config.ClusterGeoTag).isNotEmpty().matchRegexp(geoTagRegex).err
	if err != nil {
		return err
	}
	err = field(ExtClustersGeoTagsKey, config.ExtClustersGeoTags).hasItems().hasUniqueItems().err
	if err != nil {
		return err
	}
	for i, geoTag := range config.ExtClustersGeoTags {
		err = field(fmt.Sprintf("%s[%v]", ExtClustersGeoTagsKey, i), geoTag).
			isNotEmpty().matchRegexp(geoTagRegex).isNotEqualTo(config.ClusterGeoTag).err
		if err != nil {
			return err
		}
	}
	err = field(EdgeDNSServerKey, config.EdgeDNSServer).isNotEmpty().matchRegexps(hostNameRegex, ipAddressRegex).err
	if err != nil {
		return err
	}
	err = field(EdgeDNSServerPortKey, config.EdgeDNSServerPort).isHigherThanZero().err
	if err != nil {
		return err
	}
	err = field(EdgeDNSZoneKey, config.EdgeDNSZone).isNotEmpty().matchRegexp(hostNameRegex).err
	if err != nil {
		return err
	}
	err = field(DNSZoneKey, config.DNSZone).isNotEmpty().matchRegexp(hostNameRegex).err
	if err != nil {
		return err
	}
	// do full Infoblox validation only in case that Host exists
	if isNotEmpty(config.Infoblox.Host) {
		err = field(InfobloxGridHostKey, config.Infoblox.Host).matchRegexps(hostNameRegex, ipAddressRegex).err
		if err != nil {
			return err
		}
		err = field(InfobloxVersionKey, config.Infoblox.Version).isNotEmpty().matchRegexp(versionNumberRegex).err
		if err != nil {
			return err
		}
		err = field(InfobloxPortKey, config.Infoblox.Port).isHigherThanZero().isLessOrEqualTo(65535).err
		if err != nil {
			return err
		}
		err = field(InfobloxUsernameKey, config.Infoblox.Username).isNotEmpty().err
		if err != nil {
			return err
		}
		err = field(InfobloxPasswordKey, config.Infoblox.Password).isNotEmpty().err
		if err != nil {
			return err
		}
		err = field(InfobloxHTTPPoolConnectionsKey, config.Infoblox.HTTPPoolConnections).isHigherOrEqualToZero().err
		if err != nil {
			return err
		}
		err = field(InfobloxHTTPRequestTimeoutKey, config.Infoblox.HTTPRequestTimeout).isHigherThanZero().err
		if err != nil {
			return err
		}
	}
	validateLabels := func(label string) error {
		labels := strings.Split(label, ".")
		for _, l := range labels {
			if len(l) > dnsLabelMax {
				return fmt.Errorf("%s exceeds %v characters limit", l, dnsLabelMax)
			}
		}
		return nil
	}

	serverNames := config.GetExternalClusterNSNames()
	serverNames[config.ClusterGeoTag] = config.GetClusterNSName()
	for geoTag, nsName := range serverNames {
		if len(nsName) > dnsNameMax {
			return fmt.Errorf("ns name '%s' exceeds %v charactes limit for [GeoTag: '%s', %s: '%s', %s: '%s']",
				nsName, dnsLabelMax, geoTag, EdgeDNSZoneKey, config.EdgeDNSZone, DNSZoneKey, config.DNSZone)
		}
		if err := validateLabels(nsName); err != nil {
			return fmt.Errorf("error for geo tag: %s. %s in ns name %s", geoTag, err, nsName)
		}
	}

	mHost, mPort, err := parseMetricsAddr(config.MetricsAddress)
	if err != nil {
		return fmt.Errorf("invalid %s: expecting MetricsAddress in form {host}:port (%s)", MetricsAddressKey, err)
	}
	err = field(MetricsAddressKey, mHost).matchRegexps(hostNameRegex, ipAddressRegex).err
	if err != nil {
		return err
	}
	err = field(MetricsAddressKey, mPort).isLessOrEqualTo(65535).isHigherThan(1024).err
	if err != nil {
		return err
	}
	return nil
}

func parseMetricsAddr(metricsAddr string) (host string, port int, err error) {
	ma := strings.Split(metricsAddr, ":")
	if len(ma) != 2 {
		err = fmt.Errorf("invalid format {host}:port (%s)", metricsAddr)
		return
	}
	host = ma[0]
	port, err = strconv.Atoi(ma[1])
	return
}

// getEdgeDNSType contains logic retrieving EdgeDNSType.
func getEdgeDNSType(config *Config) (EdgeDNSType, []EdgeDNSType) {
	recognized := make([]EdgeDNSType, 0)
	if config.ns1Enabled {
		recognized = append(recognized, DNSTypeNS1)
	}
	if config.route53Enabled {
		recognized = append(recognized, DNSTypeRoute53)
	}
	if isNotEmpty(config.Infoblox.Host) {
		recognized = append(recognized, DNSTypeInfoblox)
	}
	switch len(recognized) {
	case 0:
		return DNSTypeNoEdgeDNS, recognized
	case 1:
		return recognized[0], recognized
	}
	return DNSTypeMultipleProviders, recognized
}

func parseLogOutputFormat(value string) LogFormat {
	switch value {
	case json:
		return JSONFormat
	case simple:
		return SimpleFormat
	}
	return NoFormat
}

func (c *Config) GetExternalClusterNSNames() (m map[string]string) {
	m = make(map[string]string, len(c.ExtClustersGeoTags))
	for _, tag := range c.ExtClustersGeoTags {
		m[tag] = getNsName(tag, c.DNSZone, c.EdgeDNSZone, c.EdgeDNSServer)
	}
	return
}

func (c *Config) GetClusterNSName() string {
	return getNsName(c.ClusterGeoTag, c.DNSZone, c.EdgeDNSZone, c.EdgeDNSServer)
}

func (c *Config) GetClusterOldNSName() string {
	dnsZoneIntoNS := strings.ReplaceAll(c.DNSZone, ".", "-")
	return fmt.Sprintf("gslb-ns-%s-%s.%s", dnsZoneIntoNS, c.ClusterGeoTag, c.EdgeDNSZone)
}

func (c *Config) GetExternalClusterHeartbeatFQDNs(gslbName string) (m map[string]string) {
	m = make(map[string]string, len(c.ExtClustersGeoTags))
	for _, tag := range c.ExtClustersGeoTags {
		m[tag] = getHeartbeatFQDN(gslbName, tag, c.EdgeDNSZone)
	}
	return
}

func (c *Config) GetClusterHeartbeatFQDN(gslbName string) string {
	return getHeartbeatFQDN(gslbName, c.ClusterGeoTag, c.EdgeDNSZone)
}

// getNsName returns NS for geo tag.
// The values is combination of DNSZone, EdgeDNSZone and (Ext)ClusterGeoTag, see:
// DNS_ZONE k8gb-test.gslb.cloud.example.com
// EDGE_DNS_ZONE: cloud.example.com
// CLUSTER_GEOTAG: us
// will generate "gslb-ns-us-k8gb-test-gslb.cloud.example.com"
// If edgeDNSServer == localhost or 127.0.0.1 than edgeDNSServer is returned.
// The function is private and expects only valid inputs.
func getNsName(tag, dnsZone, edgeDNSZone, edgeDNSServer string) string {
	if edgeDNSServer == "127.0.0.1" || edgeDNSServer == "localhost" {
		return edgeDNSServer
	}
	const prefix = "gslb-ns"
	d := strings.TrimSuffix(dnsZone, "."+edgeDNSZone)
	domainX := strings.ReplaceAll(d, ".", "-")
	return fmt.Sprintf("%s-%s-%s.%s", prefix, tag, domainX, edgeDNSZone)
}

// getHeartbeatFQDN returns heartbeat for geo tag.
// The values is combination of EdgeDNSZone and (Ext)ClusterGeoTag, and GSLB name see:
// EDGE_DNS_ZONE: cloud.example.com
// CLUSTER_GEOTAG: us
// gslb.Name: test-gslb-1
// will generate "test-gslb-1-heartbeat-us.cloud.example.com"
// The function is private and expects only valid inputs.
func getHeartbeatFQDN(name, geoTag, edgeDNSZone string) string {
	return fmt.Sprintf("%s-heartbeat-%s.%s", name, geoTag, edgeDNSZone)
}
