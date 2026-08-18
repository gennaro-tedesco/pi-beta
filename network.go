package main

import (
	"errors"
	"net"
	"net/http"
	"time"
)

const (
	networkConnectionNone         = "none"
	networkConnectionUnknown      = "unknown"
	networkConnectionWiFi         = "wifi"
	networkConnectionEthernet     = "ethernet"
	networkInternetOnline         = "online"
	networkInternetTimeout        = "timeout"
	networkInternetDNSFailure     = "dns-failure"
	networkInternetTransportError = "transport-error"
	networkInternetCaptivePortal  = "captive-portal"
	networkInternetUnknown        = "unknown"
	networkQualityMeasuring       = "measuring"
	networkQualityExcellent       = "excellent"
	networkQualityGood            = "good"
	networkQualityFair            = "fair"
	networkQualityPoor            = "poor"
	networkQualityOffline         = "offline"
	networkStabilityMeasuring     = "measuring"
	networkStabilityStable        = "stable"
	networkStabilityVariable      = "variable"
	networkStabilityUnstable      = "unstable"
	networkStabilityOffline       = "offline"
	networkSignalMinimum          = 0
	networkSignalMaximum          = 100
	networkRSSIMinimum            = -90
	networkRSSIMaximum            = -30
	networkProbeURL               = "https://cp.cloudflare.com/generate_204"
	networkProbeTimeout           = 5 * time.Second
	networkExpectedStatusCode     = http.StatusNoContent
)

var networkHTTPClient = &http.Client{Timeout: networkProbeTimeout}

type NetworkStatus struct {
	ConnectionType    string   `json:"connectionType"`
	InternetReachable *bool    `json:"internetReachable"`
	NetworkName       *string  `json:"networkName"`
	SignalDbm         *int     `json:"signalDbm"`
	LinkBitrateMbps   *float64 `json:"linkBitrateMbps"`
	AverageLatencyMs  *float64 `json:"averageLatencyMs"`
	QualityScore      int      `json:"qualityScore"`
	Quality           string   `json:"quality"`
	Stability         string   `json:"stability"`
	signalPercent     *int
	latencyMs         *int64
}

type networkLink struct {
	connectionType  string
	networkName     *string
	signalPercent   *int
	signalDbm       *int
	linkBitrateMbps *float64
}

type networkProbe struct {
	status    string
	reachable bool
	latencyMs *int64
}

func buildNetworkStatus(link networkLink, probe networkProbe) NetworkStatus {
	var reachable *bool
	if probe.status != networkInternetUnknown {
		reachable = &probe.reachable
	}
	connectionType := link.connectionType
	if connectionType == networkConnectionNone && probe.reachable {
		connectionType = networkConnectionEthernet
	}
	return NetworkStatus{
		ConnectionType:    connectionType,
		InternetReachable: reachable,
		NetworkName:       link.networkName,
		SignalDbm:         link.signalDbm,
		LinkBitrateMbps:   link.linkBitrateMbps,
		Quality:           networkQualityMeasuring,
		Stability:         networkStabilityMeasuring,
		signalPercent:     link.signalPercent,
		latencyMs:         probe.latencyMs,
	}
}

func probeInternet(client *http.Client, endpoint string) networkProbe {
	startedAt := time.Now()
	response, err := client.Get(endpoint)
	if err != nil {
		return networkProbe{status: classifyProbeError(err)}
	}
	response.Body.Close()

	if response.StatusCode != networkExpectedStatusCode {
		return networkProbe{status: networkInternetCaptivePortal}
	}

	latency := time.Since(startedAt).Milliseconds()
	return networkProbe{status: networkInternetOnline, reachable: true, latencyMs: &latency}
}

func classifyProbeError(err error) string {
	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		return networkInternetDNSFailure
	}

	var networkError net.Error
	if errors.As(err, &networkError) && networkError.Timeout() {
		return networkInternetTimeout
	}

	return networkInternetTransportError
}

func normalizeRSSI(rssi int) int {
	if rssi <= networkRSSIMinimum {
		return networkSignalMinimum
	}
	if rssi >= networkRSSIMaximum {
		return networkSignalMaximum
	}
	return (rssi - networkRSSIMinimum) * networkSignalMaximum / (networkRSSIMaximum - networkRSSIMinimum)
}
