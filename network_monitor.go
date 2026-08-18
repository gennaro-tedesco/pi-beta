package main

import (
	"math"
	"sync"
	"time"
)

const (
	networkSampleWindow               = 2 * time.Minute
	networkMinimumQualitySamples      = 3
	networkReliabilityWeight          = 45.0
	networkLatencyWeight              = 25.0
	networkJitterWeight               = 20.0
	networkSignalWeight               = 10.0
	networkLossScoreMaximumPercent    = 10.0
	networkLatencyBestMs              = 20.0
	networkLatencyWorstMs             = 300.0
	networkJitterBestMs               = 5.0
	networkJitterWorstMs              = 100.0
	networkQualityExcellentMinimum    = 85
	networkQualityGoodMinimum         = 70
	networkQualityFairMinimum         = 50
	networkUnstableLossMinimumPercent = 5.0
	networkUnstableJitterMinimumMs    = 50.0
	networkVariableJitterMinimumMs    = 20.0
	networkPercentageScale            = 100.0
	networkNoSamples                  = 0
	networkMinimumJitterSamples       = 2
	networkFirstVariationSample       = 1
	networkScoreMinimum               = 0.0
	networkScoreMaximum               = 1.0
)

type networkMonitor struct {
	mutex   sync.Mutex
	samples []networkSample
}

type networkSample struct {
	checkedAt time.Time
	reachable bool
	latencyMs *int64
}

type networkWindowMetrics struct {
	averageLatencyMs *float64
	jitterMs         *float64
	probeLossPercent float64
	sampleCount      int
}

func (monitor *networkMonitor) record(status NetworkStatus, checkedAt time.Time) NetworkStatus {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	reachable := status.InternetReachable != nil && *status.InternetReachable
	monitor.samples = append(monitor.samples, networkSample{
		checkedAt: checkedAt,
		reachable: reachable,
		latencyMs: status.LatencyMs,
	})
	monitor.prune(checkedAt.Add(-networkSampleWindow))
	metrics := calculateNetworkWindowMetrics(monitor.samples)

	status.AverageLatencyMs = metrics.averageLatencyMs
	status.Stability = classifyNetworkStability(metrics, reachable)
	status.QualityScore = calculateNetworkQualityScore(metrics, status.SignalPercent, reachable)
	status.Quality = classifyNetworkQuality(status.QualityScore, metrics.sampleCount, reachable)
	return status
}

func (monitor *networkMonitor) prune(cutoff time.Time) {
	firstRetained := networkNoSamples
	for firstRetained < len(monitor.samples) && monitor.samples[firstRetained].checkedAt.Before(cutoff) {
		firstRetained++
	}
	monitor.samples = monitor.samples[firstRetained:]
}

func calculateNetworkWindowMetrics(samples []networkSample) networkWindowMetrics {
	metrics := networkWindowMetrics{sampleCount: len(samples)}
	if len(samples) == networkNoSamples {
		return metrics
	}

	failedSamples := networkNoSamples
	latencies := make([]float64, 0, len(samples))
	for _, sample := range samples {
		if !sample.reachable {
			failedSamples++
		}
		if sample.latencyMs != nil {
			latencies = append(latencies, float64(*sample.latencyMs))
		}
	}
	metrics.probeLossPercent = float64(failedSamples) / float64(len(samples)) * networkPercentageScale
	if len(latencies) == networkNoSamples {
		return metrics
	}

	totalLatency := networkScoreMinimum
	for _, latency := range latencies {
		totalLatency += latency
	}
	averageLatency := totalLatency / float64(len(latencies))
	metrics.averageLatencyMs = &averageLatency
	if len(latencies) < networkMinimumJitterSamples {
		return metrics
	}

	totalVariation := networkScoreMinimum
	for index := networkFirstVariationSample; index < len(latencies); index++ {
		totalVariation += math.Abs(latencies[index] - latencies[index-networkFirstVariationSample])
	}
	jitter := totalVariation / float64(len(latencies)-networkFirstVariationSample)
	metrics.jitterMs = &jitter
	return metrics
}

func calculateNetworkQualityScore(metrics networkWindowMetrics, signalPercent *int, reachable bool) int {
	if !reachable || metrics.sampleCount < networkMinimumQualitySamples {
		return networkSignalMinimum
	}

	reliabilityScore := networkReliabilityWeight * clampScore(networkScoreMaximum-metrics.probeLossPercent/networkLossScoreMaximumPercent)
	latencyScore := networkLatencyWeight
	if metrics.averageLatencyMs != nil {
		latencyScore *= inverseRangeScore(*metrics.averageLatencyMs, networkLatencyBestMs, networkLatencyWorstMs)
	}
	jitterScore := networkJitterWeight
	if metrics.jitterMs != nil {
		jitterScore *= inverseRangeScore(*metrics.jitterMs, networkJitterBestMs, networkJitterWorstMs)
	}
	signalScore := networkSignalWeight
	if signalPercent != nil {
		signalScore *= clampScore(float64(*signalPercent) / networkPercentageScale)
	}

	return int(math.Round(reliabilityScore + latencyScore + jitterScore + signalScore))
}

func classifyNetworkQuality(score, sampleCount int, reachable bool) string {
	if !reachable {
		return networkQualityOffline
	}
	if sampleCount < networkMinimumQualitySamples {
		return networkQualityMeasuring
	}
	if score >= networkQualityExcellentMinimum {
		return networkQualityExcellent
	}
	if score >= networkQualityGoodMinimum {
		return networkQualityGood
	}
	if score >= networkQualityFairMinimum {
		return networkQualityFair
	}
	return networkQualityPoor
}

func classifyNetworkStability(metrics networkWindowMetrics, reachable bool) string {
	if !reachable {
		return networkStabilityOffline
	}
	if metrics.sampleCount < networkMinimumQualitySamples || metrics.jitterMs == nil {
		return networkStabilityMeasuring
	}
	if metrics.probeLossPercent > networkUnstableLossMinimumPercent || *metrics.jitterMs >= networkUnstableJitterMinimumMs {
		return networkStabilityUnstable
	}
	if metrics.probeLossPercent > networkSignalMinimum || *metrics.jitterMs >= networkVariableJitterMinimumMs {
		return networkStabilityVariable
	}
	return networkStabilityStable
}

func inverseRangeScore(value, best, worst float64) float64 {
	return clampScore(networkScoreMaximum - (value-best)/(worst-best))
}

func clampScore(value float64) float64 {
	if value < networkScoreMinimum {
		return networkScoreMinimum
	}
	if value > networkScoreMaximum {
		return networkScoreMaximum
	}
	return value
}
