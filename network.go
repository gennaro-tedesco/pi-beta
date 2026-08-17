package main

import (
	"errors"
	"net"
	"net/http"
	"time"
)

const (
	networkAvailabilityAvailable   = "available"
	networkAvailabilityPartial     = "partial"
	networkAvailabilityUnavailable = "unavailable"
	networkConnectionNone          = "none"
	networkConnectionUnknown       = "unknown"
	networkConnectionWiFi          = "wifi"
	networkConnectionEthernet      = "ethernet"
	networkInternetOnline          = "online"
	networkInternetTimeout         = "timeout"
	networkInternetDNSFailure      = "dns-failure"
	networkInternetTransportError  = "transport-error"
	networkInternetCaptivePortal   = "captive-portal"
	networkInternetUnknown         = "unknown"
	networkSignalUnknown           = "unknown"
	networkSignalPoor              = "poor"
	networkSignalFair              = "fair"
	networkSignalGood              = "good"
	networkSignalExcellent         = "excellent"
	networkQualityMeasuring        = "measuring"
	networkQualityExcellent        = "excellent"
	networkQualityGood             = "good"
	networkQualityFair             = "fair"
	networkQualityPoor             = "poor"
	networkQualityOffline          = "offline"
	networkStabilityMeasuring      = "measuring"
	networkStabilityStable         = "stable"
	networkStabilityVariable       = "variable"
	networkStabilityUnstable       = "unstable"
	networkStabilityOffline        = "offline"
	networkSignalMinimum           = 0
	networkSignalPoorLimit         = 25
	networkSignalFairLimit         = 50
	networkSignalGoodLimit         = 75
	networkSignalMaximum           = 100
	networkRSSIMinimum             = -90
	networkRSSIMaximum             = -30
	networkProbeURL                = "https://cp.cloudflare.com/generate_204"
	networkProbeTimeout            = 5 * time.Second
	networkExpectedStatusCode      = http.StatusNoContent
)

var networkHTTPClient = &http.Client{Timeout: networkProbeTimeout}

type NetworkStatus struct {
	Availability      string   `json:"availability"`
	ConnectionType    string   `json:"connectionType"`
	InternetStatus    string   `json:"internetStatus"`
	InternetReachable *bool    `json:"internetReachable"`
	NetworkName       *string  `json:"networkName"`
	SignalPercent     *int     `json:"signalPercent"`
	SignalDbm         *int     `json:"signalDbm"`
	SignalQuality     string   `json:"signalQuality"`
	LinkBitrateMbps   *float64 `json:"linkBitrateMbps"`
	LatencyMs         *int64   `json:"latencyMs"`
	AverageLatencyMs  *float64 `json:"averageLatencyMs"`
	JitterMs          *float64 `json:"jitterMs"`
	ProbeLossPercent  float64  `json:"probeLossPercent"`
	QualityScore      int      `json:"qualityScore"`
	Quality           string   `json:"quality"`
	Stability         string   `json:"stability"`
	SampleCount       int      `json:"sampleCount"`
	CheckedAt         string   `json:"checkedAt"`
}

type networkLink struct {
	available       bool
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

func buildNetworkStatus(link networkLink, probe networkProbe, checkedAt time.Time) NetworkStatus {
	var reachable *bool
	if probe.status != networkInternetUnknown {
		reachable = &probe.reachable
	}
	connectionType := link.connectionType
	if connectionType == networkConnectionNone && probe.reachable {
		connectionType = networkConnectionEthernet
	}
	availability := networkAvailabilityUnavailable
	if link.available {
		availability = networkAvailabilityAvailable
	} else if probe.status != networkInternetUnknown {
		availability = networkAvailabilityPartial
	}

	return NetworkStatus{
		Availability:      availability,
		ConnectionType:    connectionType,
		InternetStatus:    probe.status,
		InternetReachable: reachable,
		NetworkName:       link.networkName,
		SignalPercent:     link.signalPercent,
		SignalDbm:         link.signalDbm,
		SignalQuality:     classifySignalQuality(link.signalPercent),
		LinkBitrateMbps:   link.linkBitrateMbps,
		LatencyMs:         probe.latencyMs,
		Quality:           networkQualityMeasuring,
		Stability:         networkStabilityMeasuring,
		CheckedAt:         checkedAt.UTC().Format(time.RFC3339),
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

func classifySignalQuality(signalPercent *int) string {
	if signalPercent == nil || *signalPercent < networkSignalMinimum || *signalPercent > networkSignalMaximum {
		return networkSignalUnknown
	}
	if *signalPercent < networkSignalPoorLimit {
		return networkSignalPoor
	}
	if *signalPercent < networkSignalFairLimit {
		return networkSignalFair
	}
	if *signalPercent < networkSignalGoodLimit {
		return networkSignalGood
	}
	return networkSignalExcellent
}
