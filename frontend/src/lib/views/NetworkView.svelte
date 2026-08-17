<script lang="ts">
  import { scale } from 'svelte/transition'
  import { Activity, Cable, RadioTower, Server, Wifi, WifiOff } from 'lucide-svelte'
  import type { main } from '../../../wailsjs/go/models'
  import { Message } from '../components'
  import { transitionDuration } from '../transition'

  export let status: main.MachineStatus | null = null
  export let loading = true
  export let error: string | null = null

  const machineContainerType = 'size'
  const machineAriaLabel = 'Machine status'
  const loadingLabel = 'Reading machine status…'
  const unavailableLabel = 'Unavailable'
  const unknownMachineLabel = 'This computer'
  const overviewLabel = 'System overview'
  const networkLabel = 'Network'
  const processesLabel = 'Largest processes'
  const processCountLabel = 'running'
  const signalLabel = 'Signal'
  const latencyLabel = 'Latency'
  const linkSpeedLabel = 'Link speed'
  const connectedLabel = 'Connected'
  const offlineLabel = 'Offline'
  const wifiLabel = 'Wi-Fi'
  const ethernetLabel = 'Ethernet'
  const noConnectionLabel = 'No connection'
  const processIdLabel = 'PID'
  const processNameLabel = 'Process'
  const processMemoryLabel = 'Memory use'
  const processRAMLabel = 'RAM share'
  const totalRAMLabel = 'Total RAM'
  const notApplicableLabel = '-'
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

  $: networkQualityScore = clampPercentage(status?.network.qualityScore)
  $: networkGaugeOffset = percentageMaximum - networkQualityScore
</script>

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
        <span class="machine__eyebrow">{overviewLabel}</span>
        <h1>{status.hostname || unknownMachineLabel}</h1>
        <p>{formatPlatform(status)}</p>
      </header>

      <div class="machine__details">
        <article class="panel network-panel">
          <div class="panel__title">
            {#if status.network.connectionType === connectionTypeWiFi}
              <Wifi aria-hidden="true" />
            {:else if status.network.connectionType === connectionTypeEthernet}
              <Cable aria-hidden="true" />
            {:else}
              <WifiOff aria-hidden="true" />
            {/if}
            <div>
              <span>{networkLabel}</span>
              <h2>{status.network.networkName ?? formatConnection(status.network.connectionType)}</h2>
            </div>
            <strong class:network-panel__online={status.network.internetReachable}>
              {status.network.internetReachable ? connectedLabel : offlineLabel}
            </strong>
          </div>
          <div class="network-panel__body">
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
                <strong>{formatPercentage(networkQualityScore)}</strong>
                <span>{formatLabel(status.network.quality)}</span>
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
            </div>
          </div>
        </article>

        <article class="panel process-panel">
          <div class="panel__title">
            <Server aria-hidden="true" />
            <div>
              <span>{processesLabel}</span>
              <h2>{status.processCount} {processCountLabel}</h2>
            </div>
          </div>
          <div class="process-panel__table-wrap">
            <table>
              <thead>
                <tr>
                  <th scope="col">{processNameLabel}</th>
                  <th scope="col">{processIdLabel}</th>
                  <th scope="col">{processMemoryLabel}</th>
                  <th scope="col">{processRAMLabel}</th>
                </tr>
              </thead>
              <tbody>
                {#each status.largestProcesses as process (process.pid)}
                  <tr>
                    <th scope="row">{process.name}</th>
                    <td>{process.pid}</td>
                    <td>{formatBytes(process.memoryBytes)}</td>
                    <td>{formatPercentage(process.memoryPercent)}</td>
                  </tr>
                {/each}
              </tbody>
              <tfoot>
                <tr>
                  <th scope="row">{totalRAMLabel}</th>
                  <td>{notApplicableLabel}</td>
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

<style>
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
    --machine-double-span: 2;
    --machine-title-size: min(6cqh, 5cqw);
    --machine-heading-size: min(3.5cqh, 2cqw);
    --machine-label-size: min(2.5cqh, 1.4cqw);
    --machine-icon-size: min(4cqh, 2cqw);
    --machine-gauge-stroke: min(4cqh, 1.4cqw);
    --machine-gauge-linecap: round;
    --machine-gauge-dash-size: 100;
    --machine-gauge-reading-top: 46%;
    --machine-gauge-value-size: min(5cqh, 3cqw);
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
    --machine-process-name-width: 24cqw;
    --machine-gauge-transition-duration: 700ms;
    --machine-gauge-transition-easing: cubic-bezier(0.22, 1, 0.36, 1);
    --machine-line-height: 1.1;
    --machine-eyebrow-spacing: 0.14em;
    display: flex;
    align-items: center;
    justify-content: center;
    width: var(--machine-fill);
    height: var(--machine-fill);
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
  .machine__eyebrow,
  .panel__title span,
  .network-panel__metrics span,
  .network-gauge__reading span,
  .process-panel thead {
    color: var(--machine-muted);
    font-size: var(--machine-label-size);
  }

  .machine__eyebrow {
    text-transform: uppercase;
    letter-spacing: var(--machine-eyebrow-spacing);
  }

  .machine__details {
    display: grid;
    gap: var(--machine-gap);
  }

  .machine__details {
    flex: 1 1 auto;
    grid-template-columns: var(--machine-detail-columns);
    grid-auto-rows: var(--machine-detail-rows);
    min-height: var(--machine-zero);
    overflow: hidden;
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

  .panel__title {
    display: grid;
    grid-template-columns: auto minmax(0, 1fr) auto;
    align-items: center;
    gap: calc(var(--machine-gap) / 2);
    text-align: left;
  }

  .panel h2 {
    overflow: hidden;
    font-size: var(--machine-heading-size);
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .panel__title > strong {
    color: var(--machine-muted);
  }

  .panel__title > .network-panel__online {
    color: var(--machine-online);
  }

  .network-panel__body {
    display: grid;
    grid-template-columns: var(--machine-network-body-columns);
    align-items: center;
    gap: var(--machine-gap);
    min-height: var(--machine-zero);
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
    font-size: var(--machine-gauge-value-size);
    line-height: var(--machine-line-height);
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

  .process-panel table {
    width: var(--machine-fill);
    height: var(--machine-fill);
    border-collapse: collapse;
    font-size: var(--machine-table-font-size);
    line-height: var(--machine-line-height);
    table-layout: fixed;
    text-align: left;
  }

  .process-panel th,
  .process-panel td {
    padding: var(--machine-table-cell-padding);
    border-bottom: var(--machine-table-border);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .process-panel tbody th {
    max-width: var(--machine-process-name-width);
    overflow: hidden;
    text-overflow: ellipsis;
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
  }
</style>
