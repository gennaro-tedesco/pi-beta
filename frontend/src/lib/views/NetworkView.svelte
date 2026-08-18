<script lang="ts">
  import { onMount } from 'svelte'
  import { scale } from 'svelte/transition'
  import { Activity, MemoryStick, Percent, RadioTower, RotateCw, Waves, Wifi } from 'lucide-svelte'
  import { GetMachineStatus } from '../../../wailsjs/go/main/App'
  import type { main } from '../../../wailsjs/go/models'
  import { Button, Message, WidgetNavigation } from '../components'
  import { transitionDuration } from '../transition'

  type MatrixLine = {
    delay: string
    duration: string
    id: string
    opacity: string
    text: string
    top: string
  }

  const machineContainerType = 'size'
  const networkRefreshLabel = 'Refresh machine status'
  const networkPollIntervalMs = 5000
  const networkActionIconSize = 18
  const machineAriaLabel = 'Machine status'
  const loadingLabel = 'Reading machine status…'
  const unavailableLabel = 'Unavailable'
  const unknownMachineLabel = 'This computer'
  const cpuUtilisationLabel = 'CPU utilisation'
  const signalLabel = 'Signal'
  const stabilityLabel = 'Stability'
  const latencyLabel = 'Latency'
  const linkSpeedLabel = 'Link speed'
  const connectedLabel = 'Connected'
  const offlineLabel = 'Offline'
  const wifiLabel = 'Wi-Fi'
  const ethernetLabel = 'Ethernet'
  const noConnectionLabel = 'No connection'
  const processNameLabel = 'Process'
  const processMemoryLabel = 'Memory use'
  const processRAMLabel = 'RAM share'
  const totalRAMLabel = 'Total RAM'
  const platformSeparator = ' · '
  const byteUnits = ['B', 'KB', 'MB', 'GB', 'TB']
  const percentUnit = '%'
  const millisecondsUnit = 'ms'
  const megabitsUnit = 'Mbps'
  const decibelsUnit = 'dBm'
  const bytesPerUnit = 1024
  const percentageMinimum = 0
  const percentageMaximum = 100
  const firstByteUnitIndex = 0
  const nextByteUnitOffset = 1
  const byteFractionDigits = 1
  const percentFractionDigits = 1
  const metricFractionDigits = 0
  const connectionTypeWiFi = 'wifi'
  const connectionTypeEthernet = 'ethernet'
  const gaugeViewBox = '0 0 200 120'
  const gaugeAspectRatio = 'xMidYMid meet'
  const gaugePath = 'M 20 100 A 80 80 0 0 1 180 100'
  const gaugePathLength = percentageMaximum
  const labelFirstCharacterIndex = 0
  const labelRemainderIndex = 1
  const matrixLines: MatrixLine[] = [
    { id: 'matrix-01', text: '$ ping -c 4 1.1.1.1', top: '2%', delay: '-1.8s', duration: '4.2s', opacity: '0.72' },
    { id: 'matrix-02', text: '$ curl -sI https://api', top: '9%', delay: '-3.1s', duration: '3.7s', opacity: '0.9' },
    { id: 'matrix-03', text: '$ top -o cpu', top: '16%', delay: '-0.6s', duration: '4.8s', opacity: '0.64' },
    { id: 'matrix-04', text: '$ netstat -an | grep LISTEN', top: '23%', delay: '-2.5s', duration: '3.9s', opacity: '0.82' },
    { id: 'matrix-05', text: '$ whoami', top: '30%', delay: '-4.2s', duration: '5.1s', opacity: '0.58' },
    { id: 'matrix-06', text: '$ uptime', top: '37%', delay: '-1.3s', duration: '4.4s', opacity: '0.88' },
    { id: 'matrix-07', text: '$ ifconfig wlan0', top: '44%', delay: '-3.8s', duration: '4.9s', opacity: '0.7' },
    { id: 'matrix-08', text: '$ ps aux | grep node', top: '51%', delay: '-0.9s', duration: '3.6s', opacity: '0.94' },
    { id: 'matrix-09', text: '$ df -h /', top: '58%', delay: '-2.9s', duration: '4.6s', opacity: '0.62' },
    { id: 'matrix-10', text: '$ free -m', top: '65%', delay: '-4.5s', duration: '5.2s', opacity: '0.84' },
    { id: 'matrix-11', text: '$ traceroute 8.8.8.8', top: '72%', delay: '-1.6s', duration: '4.1s', opacity: '0.76' },
    { id: 'matrix-12', text: '$ systemctl status network', top: '79%', delay: '-3.4s', duration: '4.7s', opacity: '0.92' },
    { id: 'matrix-13', text: '$ journalctl -f', top: '86%', delay: '-2.1s', duration: '3.8s', opacity: '0.66' },
    { id: 'matrix-14', text: '$ ssh admin@10.0.0.4', top: '93%', delay: '-4.1s', duration: '5s', opacity: '0.86' }
  ]

  let status: main.MachineStatus | null = null
  let loading = true
  let error: string | null = null
  let refreshing = false

  async function refreshNetworkStatus(): Promise<void> {
    if (refreshing) return
    refreshing = true
    error = null
    try {
      status = await GetMachineStatus()
    } catch (err) {
      error = String(err)
    } finally {
      loading = false
      refreshing = false
    }
  }

  onMount(() => {
    refreshNetworkStatus()
    const refreshTimer = setInterval(refreshNetworkStatus, networkPollIntervalMs)
    return () => clearInterval(refreshTimer)
  })

  function clampPercentage(value: number | undefined): number {
    return Math.min(percentageMaximum, Math.max(percentageMinimum, value ?? percentageMinimum))
  }

  function formatPercentage(value: number | undefined): string {
    return value == null ? unavailableLabel : `${value.toFixed(percentFractionDigits)}${percentUnit}`
  }

  function formatBytes(value: number): string {
    if (value <= percentageMinimum) return unavailableLabel
    let scaledValue = value
    let unitIndex = firstByteUnitIndex
    while (scaledValue >= bytesPerUnit && unitIndex < byteUnits.length - nextByteUnitOffset) {
      scaledValue /= bytesPerUnit
      unitIndex += nextByteUnitOffset
    }
    return `${scaledValue.toFixed(byteFractionDigits)} ${byteUnits[unitIndex]}`
  }

  function formatPlatform(value: main.MachineStatus): string {
    return [value.platform, value.platformVersion, value.architecture].filter(Boolean).join(platformSeparator)
  }

  function formatSystemSummary(value: main.MachineStatus): string {
    return [formatPlatform(value), `${cpuUtilisationLabel} ${formatPercentage(value.cpuPercent)}`]
      .filter(Boolean)
      .join(platformSeparator)
  }

  function formatConnection(value: string): string {
    if (value === connectionTypeWiFi) return wifiLabel
    if (value === connectionTypeEthernet) return ethernetLabel
    return noConnectionLabel
  }

  function formatMetric(value: number | undefined, unit: string): string {
    return value == null ? unavailableLabel : `${value.toFixed(metricFractionDigits)} ${unit}`
  }

  function formatLabel(value: string): string {
    return value.charAt(labelFirstCharacterIndex).toUpperCase() + value.slice(labelRemainderIndex)
  }

  function formatNetworkQuality(score: number, quality: string): string {
    return [formatPercentage(score), formatLabel(quality)].join(platformSeparator)
  }

  $: networkQualityScore = clampPercentage(status?.network.qualityScore)
  $: networkGaugeOffset = percentageMaximum - networkQualityScore
</script>

<div class="network-widget">
  <WidgetNavigation on:home>
    <svelte:fragment slot="actions">
      <Button ghost disabled={refreshing} on:click={refreshNetworkStatus} aria-label={networkRefreshLabel}>
        <RotateCw size={networkActionIconSize} />
      </Button>
    </svelte:fragment>
  </WidgetNavigation>
  <section
    class="machine"
    style:container-type={machineContainerType}
    aria-label={machineAriaLabel}
    out:scale={{ duration: transitionDuration }}
    in:scale={{ duration: transitionDuration, delay: transitionDuration }}
  >
  {#if error}
    <Message variant="error" message={error} />
  {:else if loading || status === null}
    <p class="machine__loading">{loadingLabel}</p>
  {:else}
    <div class="machine__content">
      <header class="machine__header">
        <h1>{status.hostname || unknownMachineLabel}</h1>
        <p>{formatSystemSummary(status)}</p>
      </header>

      <div class="machine__details">
        <article class="panel network-panel">
          <div class="network-panel__body">
            <div class="network-gauge-group">
              <h2 class="network-gauge__name">{status.network.networkName ?? formatConnection(status.network.connectionType)}</h2>
              <div class="network-gauge" class:network-gauge--offline={!status.network.internetReachable}>
                <svg viewBox={gaugeViewBox} preserveAspectRatio={gaugeAspectRatio} aria-hidden="true">
                  <path class="network-gauge__track" d={gaugePath} pathLength={gaugePathLength} />
                  <path
                    class="network-gauge__value network-gauge__value--{status.network.quality}"
                    d={gaugePath}
                    pathLength={gaugePathLength}
                    style:stroke-dashoffset={networkGaugeOffset}
                  />
                </svg>
                <div class="network-gauge__reading">
                  <strong class:network-panel__online={status.network.internetReachable}>
                    {status.network.internetReachable ? connectedLabel : offlineLabel}
                  </strong>
                  <span>{formatNetworkQuality(networkQualityScore, status.network.quality)}</span>
                </div>
              </div>
            </div>
            <div class="network-panel__metrics">
              <div>
                <RadioTower aria-hidden="true" />
                <span>{linkSpeedLabel}</span>
                <strong>{formatMetric(status.network.linkBitrateMbps, megabitsUnit)}</strong>
              </div>
              <div>
                <Activity aria-hidden="true" />
                <span>{latencyLabel}</span>
                <strong>{formatMetric(status.network.averageLatencyMs, millisecondsUnit)}</strong>
              </div>
              <div>
                <Wifi aria-hidden="true" />
                <span>{signalLabel}</span>
                <strong>{status.network.signalDbm == null ? unavailableLabel : formatMetric(status.network.signalDbm, decibelsUnit)}</strong>
              </div>
              <div>
                <Waves aria-hidden="true" />
                <span>{stabilityLabel}</span>
                <strong>{formatLabel(status.network.stability)}</strong>
              </div>
            </div>
          </div>
        </article>

        <span class="machine__matrix" aria-hidden="true">
          {#each matrixLines as line (line.id)}
            <span
              class="machine__matrix-line"
              style:top={line.top}
              style:opacity={line.opacity}
              style:animation-delay={line.delay}
              style:animation-duration={line.duration}
              style:animation-timing-function={`steps(${line.text.length}, end)`}
              style:--machine-matrix-line-width={`${line.text.length}ch`}
            >{line.text}</span>
          {/each}
        </span>

        <article class="panel process-panel">
          <div class="process-panel__table-wrap">
            <table>
              <colgroup>
                <col class="process-panel__name-column" />
                <col class="process-panel__memory-column" />
                <col class="process-panel__percentage-column" />
              </colgroup>
              <thead>
                <tr>
                  <th scope="col">{processNameLabel}</th>
                  <th scope="col" aria-label={processMemoryLabel}><MemoryStick aria-hidden="true" /></th>
                  <th scope="col" aria-label={processRAMLabel}><Percent aria-hidden="true" /></th>
                </tr>
              </thead>
              <tbody>
                {#each status.largestProcesses as process (process.pid)}
                  <tr>
                    <th scope="row">{process.name}</th>
                    <td>{formatBytes(process.memoryBytes)}</td>
                    <td>{formatPercentage(process.memoryPercent)}</td>
                  </tr>
                {/each}
              </tbody>
              <tfoot>
                <tr>
                  <th scope="row">{totalRAMLabel}</th>
                  <td>{formatBytes(status.memoryUsedBytes)}</td>
                  <td>{formatPercentage(status.memoryUsedPercent)}</td>
                </tr>
              </tfoot>
            </table>
          </div>
        </article>
      </div>
    </div>
  {/if}
  </section>
</div>

<style>
  .network-widget {
    --network-widget-edge: 0;
    --network-widget-fill: 100%;
    --network-widget-grow: 1;
    --network-widget-minimum: 0;
    --network-widget-padding: 1rem;
    --network-widget-gap: 1rem;
    position: absolute;
    inset: var(--network-widget-edge);
    display: flex;
    flex-direction: column;
    width: var(--network-widget-fill);
    height: var(--network-widget-fill);
    min-height: var(--network-widget-minimum);
    box-sizing: border-box;
    padding: var(--network-widget-padding);
    gap: var(--network-widget-gap);
  }

  .machine {
    --machine-fill: 100%;
    --machine-zero: 0;
    --machine-gap: 2.5cqh;
    --machine-padding: 3cqh;
    --machine-small-padding: 2cqh;
    --machine-primary: #4a3f66;
    --machine-muted: rgba(74, 63, 102, 0.65);
    --machine-online: #3f8a65;
    --machine-offline: #b9525f;
    --machine-poor: #d46a6a;
    --machine-fair: #d49d51;
    --machine-good: #69a679;
    --machine-excellent: #3c8d72;
    --machine-track: rgba(74, 63, 102, 0.12);
    --machine-detail-columns: minmax(0, 3fr) minmax(0, 2fr);
    --machine-detail-rows: minmax(0, 1fr);
    --machine-detail-layout-rows: auto auto;
    --machine-detail-areas: ". matrix" "network processes";
    --machine-double-span: 2;
    --machine-title-size: min(6cqh, 5cqw);
    --machine-heading-size: min(3.5cqh, 2cqw);
    --machine-label-size: min(2.5cqh, 1.4cqw);
    --machine-icon-size: min(4cqh, 2cqw);
    --machine-gauge-stroke: min(4cqh, 1.4cqw);
    --machine-gauge-linecap: round;
    --machine-gauge-dash-size: 100;
    --machine-gauge-reading-top: 46%;
    --machine-visual-height: min(100%, 22.2cqw);
    --machine-network-body-columns: minmax(0, 3fr) minmax(0, 2fr);
    --machine-panel-rows: auto minmax(0, 1fr);
    --machine-table-border-width: min(0.1cqh, 0.075cqw);
    --machine-table-footer-border-width: min(0.2cqh, 0.15cqw);
    --machine-table-border: var(--machine-table-border-width) solid rgba(74, 63, 102, 0.1);
    --machine-table-footer-border: var(--machine-table-footer-border-width) solid rgba(74, 63, 102, 0.24);
    --machine-table-footer-weight: 700;
    --machine-table-cell-padding: min(1.2cqh, 0.7cqw);
    --machine-table-font-size: min(2.5cqh, 1.4cqw);
    --machine-process-name-column-width: 60%;
    --machine-process-memory-column-width: 24%;
    --machine-process-percentage-column-width: 16%;
    --machine-gauge-transition-duration: 700ms;
    --machine-gauge-transition-easing: cubic-bezier(0.22, 1, 0.36, 1);
    --machine-line-height: 1.1;
    --machine-matrix-width: min(36cqw, 30cqh);
    --machine-matrix-aspect-ratio: 16 / 9;
    --machine-matrix-background: radial-gradient(ellipse at center, rgba(184, 224, 232, 0.64), transparent 68%), radial-gradient(ellipse at center, rgba(224, 204, 240, 0.58), transparent 82%);
    --machine-matrix-background-blend-mode: soft-light;
    --machine-matrix-edge-mask: radial-gradient(ellipse at center, #000000 42%, rgba(0, 0, 0, 0.82) 64%, transparent 100%);
    --machine-matrix-blend-mode: multiply;
    --machine-matrix-color: #315f73;
    --machine-matrix-font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
    --machine-matrix-font-size: min(1.65cqh, 0.93cqw);
    --machine-matrix-font-weight: 700;
    --machine-matrix-line-height: 1;
    --machine-matrix-shadow-blur: min(1.2cqh, 0.675cqw);
    --machine-matrix-text-shadow: 0 0 var(--machine-matrix-shadow-blur) currentColor;
    --machine-matrix-horizontal-offset: min(1.5cqw, 1.5cqh);
    display: flex;
    flex: var(--network-widget-grow);
    align-items: center;
    justify-content: center;
    width: var(--machine-fill);
    height: var(--machine-fill);
    min-height: var(--machine-zero);
    box-sizing: border-box;
    padding: var(--machine-padding);
    color: var(--machine-primary);
    overflow: hidden;
  }

  .machine__content {
    display: flex;
    flex-direction: column;
    width: var(--machine-fill);
    height: var(--machine-fill);
    min-height: var(--machine-zero);
    gap: var(--machine-gap);
  }

  .machine__header {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: calc(var(--machine-gap) / 4);
    text-align: center;
  }

  .machine__header h1,
  .machine__header p,
  .machine__loading,
  .panel h2 {
    margin: var(--machine-zero);
  }

  .machine__header h1 {
    font-size: var(--machine-title-size);
    line-height: var(--machine-line-height);
  }

  .machine__header p,
  .network-panel__metrics span,
  .network-gauge__reading span,
  .process-panel thead {
    color: var(--machine-muted);
    font-size: var(--machine-label-size);
  }

  .machine__matrix {
    position: relative;
    display: block;
    width: var(--machine-matrix-width);
    aspect-ratio: var(--machine-matrix-aspect-ratio);
    box-sizing: border-box;
    overflow: hidden;
    background: var(--machine-matrix-background);
    background-blend-mode: var(--machine-matrix-background-blend-mode);
    color: var(--machine-matrix-color);
    mix-blend-mode: var(--machine-matrix-blend-mode);
    -webkit-mask-image: var(--machine-matrix-edge-mask);
    mask-image: var(--machine-matrix-edge-mask);
  }

  .machine__matrix-line {
    position: absolute;
    left: var(--machine-zero);
    width: var(--machine-zero);
    overflow: hidden;
    font-family: var(--machine-matrix-font-family);
    font-size: var(--machine-matrix-font-size);
    font-weight: var(--machine-matrix-font-weight);
    line-height: var(--machine-matrix-line-height);
    text-shadow: var(--machine-matrix-text-shadow);
    white-space: nowrap;
    animation-name: matrixType;
    animation-iteration-count: infinite;
    will-change: width;
  }

  .machine__details {
    display: grid;
    gap: var(--machine-gap);
  }

  .machine__details {
    grid-template-areas: var(--machine-detail-areas);
    grid-template-columns: var(--machine-detail-columns);
    grid-template-rows: var(--machine-detail-layout-rows);
    min-height: var(--machine-zero);
    overflow: hidden;
  }

  .machine__details > .machine__matrix {
    right: var(--machine-matrix-horizontal-offset);
    grid-area: matrix;
    justify-self: center;
  }

  .panel {
    box-sizing: border-box;
    padding: var(--machine-padding);
    height: var(--machine-fill);
    min-height: var(--machine-zero);
    overflow: hidden;
  }

  .panel {
    display: grid;
    grid-template-rows: var(--machine-panel-rows);
    gap: var(--machine-gap);
    min-width: var(--machine-zero);
  }

  .panel h2 {
    overflow: hidden;
    font-size: var(--machine-heading-size);
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .network-panel__online {
    color: var(--machine-online);
  }

  .network-panel__body {
    display: grid;
    grid-template-columns: var(--machine-network-body-columns);
    align-items: center;
    gap: var(--machine-gap);
    min-height: var(--machine-zero);
  }

  .network-panel {
    grid-area: network;
    grid-template-rows: var(--machine-detail-rows);
  }

  .network-gauge-group {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: var(--machine-fill);
    height: var(--machine-fill);
    min-height: var(--machine-zero);
    gap: var(--machine-small-padding);
  }

  .network-gauge {
    position: relative;
    align-self: center;
    width: var(--machine-fill);
    height: var(--machine-visual-height);
    min-height: var(--machine-zero);
  }

  .network-gauge svg {
    display: block;
    width: var(--machine-fill);
    height: var(--machine-fill);
    overflow: visible;
  }

  .network-gauge__name {
    width: var(--machine-fill);
    text-align: center;
  }

  .network-gauge path {
    fill: none;
    stroke-width: var(--machine-gauge-stroke);
    stroke-linecap: var(--machine-gauge-linecap);
  }

  .network-gauge__track {
    stroke: var(--machine-track);
  }

  .network-gauge__value {
    stroke: var(--machine-muted);
    stroke-dasharray: var(--machine-gauge-dash-size);
    transition: stroke-dashoffset var(--machine-gauge-transition-duration) var(--machine-gauge-transition-easing);
  }

  .network-gauge__value--poor { stroke: var(--machine-poor); }
  .network-gauge__value--fair { stroke: var(--machine-fair); }
  .network-gauge__value--good { stroke: var(--machine-good); }
  .network-gauge__value--excellent { stroke: var(--machine-excellent); }
  .network-gauge--offline .network-gauge__value { stroke: var(--machine-offline); }

  .network-gauge__reading {
    position: absolute;
    top: var(--machine-gauge-reading-top);
    right: var(--machine-zero);
    left: var(--machine-zero);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: calc(var(--machine-gap) / 3);
    text-align: center;
  }

  .network-gauge__reading strong {
    font-size: var(--machine-heading-size);
    line-height: var(--machine-line-height);
  }

  .network-gauge--offline .network-gauge__reading strong {
    color: var(--machine-offline);
  }

  .network-panel__metrics {
    display: flex;
    flex-direction: column;
    gap: var(--machine-small-padding);
  }

  .network-panel__metrics > div {
    display: grid;
    grid-template-columns: auto minmax(0, 1fr);
    align-items: center;
    gap: calc(var(--machine-gap) / 3);
    text-align: left;
  }

  .network-panel__metrics :global(.lucide) {
    grid-row: span var(--machine-double-span);
  }

  .process-panel__table-wrap {
    align-self: center;
    width: var(--machine-fill);
    height: var(--machine-visual-height);
    min-height: var(--machine-zero);
    overflow: hidden;
  }

  .process-panel {
    grid-area: processes;
    grid-template-rows: var(--machine-detail-rows);
  }

  .process-panel table {
    width: var(--machine-fill);
    height: var(--machine-fill);
    border-collapse: collapse;
    font-size: var(--machine-table-font-size);
    line-height: var(--machine-line-height);
    table-layout: fixed;
    text-align: left;
  }

  .process-panel__name-column {
    width: var(--machine-process-name-column-width);
  }

  .process-panel__memory-column {
    width: var(--machine-process-memory-column-width);
  }

  .process-panel__percentage-column {
    width: var(--machine-process-percentage-column-width);
  }

  .process-panel th,
  .process-panel td {
    padding: var(--machine-table-cell-padding);
    border-bottom: var(--machine-table-border);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .process-panel tfoot {
    border-top: var(--machine-table-footer-border);
    font-weight: var(--machine-table-footer-weight);
  }

  .process-panel th:not(:first-child),
  .process-panel td {
    text-align: right;
  }

  .machine :global(.lucide) {
    width: var(--machine-icon-size);
    height: var(--machine-icon-size);
    color: var(--machine-muted);
  }

  .machine__loading {
    text-align: center;
  }

  @media (prefers-reduced-motion: reduce) {
    .network-gauge__value { transition: none; }
    .machine__matrix-line {
      animation: none;
      width: var(--machine-matrix-line-width);
    }
  }

  @keyframes matrixType {
    from {
      width: var(--machine-zero);
    }
    to {
      width: var(--machine-matrix-line-width);
    }
  }
</style>
