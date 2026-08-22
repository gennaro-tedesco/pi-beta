<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { scale } from 'svelte/transition'
  import { CircleGauge, Droplets, Leaf, LocateFixed, MapPin, Moon, MoonStar, RotateCw, Sun, SunMedium, Wind } from 'lucide-svelte'
  import { GetWeather, GetWeatherAt } from '../../../wailsjs/go/main/App'
  import type { main } from '../../../wailsjs/go/models'
  import { Button, Message, WidgetNavigation } from '../components'
  import CityPicker from './CityPicker.svelte'
  import { getWeatherIcon } from './weatherIcon'
  import { getWeatherScene } from './weatherScene'
  import { transitionDuration } from '../transition'

  type Coords = { lat: number; lon: number }

  const stars = [
    { top: 10, left: 15, delay: 0 },
    { top: 18, left: 32, delay: 0.6 },
    { top: 8, left: 52, delay: 1.2 },
    { top: 22, left: 68, delay: 0.3 },
    { top: 14, left: 85, delay: 1.8 },
    { top: 30, left: 8, delay: 2.1 },
    { top: 34, left: 45, delay: 0.9 },
    { top: 28, left: 92, delay: 1.5 }
  ]
  const weatherActionIconSize = 18
  const weatherMoonIconSize = 48
  const refreshWeatherLabel = 'Refresh weather'
  const currentLocationLabel = 'Use current IP location'
  const changeCityLabel = 'Change city'
  const cityClockUpdateIntervalMs = 15000
  const millisecondsPerSecond = 1000
  const minimumSunProgress = 0
  const maximumSunProgress = 1
  const sunArcHeightPercent = 94
  const sunIconSize = 18
  const detailIconSize = 18
  const sunPath = 'M 1 63 Q 50 -58 99 63'
  const sunViewBox = '0 0 100 64'
  const sunAspectRatio = 'none'
  const sunCycleLabel = "Today's sunrise and sunset"
  const weatherContainerType = 'size'
  const forecastCardContainerType = 'inline-size'
  const timestampHourStart = 11
  const timestampMinuteEnd = 16
  const timestampUTCDesignator = ':00Z'
  const percentageScale = 100
  const chartViewBoxWidth = 100
  const chartViewBoxHeight = 40
  const chartVerticalPadding = 0.18
  const chartLabelIntervalHours = 2
  const chartGradientId = 'weather-chart-gradient'
  const windEffects = ['wind-one', 'wind-two', 'wind-three']
  const initialWeatherViewKey = 0
  const weatherViewKeyIncrement = 1
  const reloadTransitionDuration = transitionDuration / 2

  interface ChartPoint {
    x: number
    y: number
  }

  interface ChartLabel {
    x: number
    label: string
    align: 'start' | 'center'
  }

  let weather: main.Weather | null = null
  let error: string | null = null
  let loading = true
  let cityClockNow = Date.now()
  let cityClockTimer: ReturnType<typeof setInterval>
  let selectedCoords: Coords | null = null
  let showCityPicker = false
  let weatherViewKey = initialWeatherViewKey
  let weatherViewTransitionDuration = transitionDuration

  export async function refresh(): Promise<void> {
    if (selectedCoords === null) {
      await loadCurrentLocationWeather(reloadTransitionDuration)
      return
    }

    await loadWeather()
  }

  async function loadCurrentLocationWeather(duration: number): Promise<void> {
    error = null
    try {
      const currentLocationWeather = await GetWeather()
      selectedCoords = null
      weather = currentLocationWeather
      weatherViewTransitionDuration = duration
      weatherViewKey += weatherViewKeyIncrement
    } catch (err) {
      error = String(err)
    }
  }

  async function loadWeather(): Promise<void> {
    loading = true
    error = null
    try {
      weather = selectedCoords ? await GetWeatherAt(selectedCoords.lat, selectedCoords.lon) : await GetWeather()
    } catch (err) {
      error = String(err)
    } finally {
      loading = false
    }
  }

  function formatDayLabel(dateString: string): string {
    return new Date(dateString).toLocaleDateString(undefined, { weekday: 'short' })
  }

  function formatCityTime(nowMs: number, utcOffsetSeconds: number): string {
    const shifted = new Date(nowMs + utcOffsetSeconds * millisecondsPerSecond)
    return shifted.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit', timeZone: 'UTC', hourCycle: 'h23' })
  }

  function formatCityDate(nowMs: number, utcOffsetSeconds: number): string {
    const shifted = new Date(nowMs + utcOffsetSeconds * millisecondsPerSecond)
    return shifted.toLocaleDateString(undefined, { weekday: 'long', month: 'long', day: '2-digit', timeZone: 'UTC' })
  }

  function formatSunTime(timestamp: string): string {
    return timestamp.slice(timestampHourStart, timestampMinuteEnd)
  }

  function formatChartHourLabel(dateString: string): string {
    return new Date(dateString).getHours().toString()
  }

  function buildChartPoints(temperatures: number[]): ChartPoint[] {
    if (temperatures.length === 0) return []
    const minTemperature = Math.min(...temperatures)
    const maxTemperature = Math.max(...temperatures)
    const range = maxTemperature - minTemperature
    return temperatures.map((temperature, index) => {
      const x = (index / (temperatures.length - 1 || 1)) * chartViewBoxWidth
      const normalized = range === 0 ? 0.5 : (temperature - minTemperature) / range
      const eased = chartVerticalPadding + normalized * (1 - 2 * chartVerticalPadding)
      const y = chartViewBoxHeight - eased * chartViewBoxHeight
      return { x, y }
    })
  }

  function buildChartLinePath(points: ChartPoint[]): string {
    if (points.length === 0) return ''
    if (points.length === 1) return `M ${points[0].x} ${points[0].y}`
    let path = `M ${points[0].x} ${points[0].y}`
    for (let index = 1; index < points.length; index++) {
      const previous = points[index - 1]
      const current = points[index]
      const controlX = (previous.x + current.x) / 2
      const controlY = (previous.y + current.y) / 2
      path += ` Q ${previous.x} ${previous.y} ${controlX} ${controlY}`
    }
    const last = points[points.length - 1]
    path += ` Q ${last.x} ${last.y} ${last.x} ${last.y}`
    return path
  }

  function buildChartAreaPath(linePath: string, points: ChartPoint[]): string {
    if (points.length === 0 || linePath === '') return ''
    const first = points[0]
    const last = points[points.length - 1]
    return `${linePath} L ${last.x} ${chartViewBoxHeight} L ${first.x} ${chartViewBoxHeight} Z`
  }

  function buildChartLabels(times: string[], points: ChartPoint[]): ChartLabel[] {
    return points
      .map((point, index) => ({ point, index }))
      .filter(({ index }) => index % chartLabelIntervalHours === 0)
      .map(({ point, index }) => ({
        x: point.x,
        label: formatChartHourLabel(times[index]),
        align: index === 0 ? 'start' as const : 'center' as const
      }))
  }

  function getSunPosition(nowMs: number, utcOffsetSeconds: number, sunrise: string, sunset: string) {
    const cityNowMs = nowMs + utcOffsetSeconds * millisecondsPerSecond
    const sunriseMs = Date.parse(`${sunrise}${timestampUTCDesignator}`)
    const sunsetMs = Date.parse(`${sunset}${timestampUTCDesignator}`)
    const daylightMs = sunsetMs - sunriseMs
    const progress = Math.min(
      maximumSunProgress,
      Math.max(minimumSunProgress, (cityNowMs - sunriseMs) / daylightMs)
    )

    return {
      left: progress * percentageScale,
      bottom: Math.sin(progress * Math.PI) * sunArcHeightPercent,
      daylight: cityNowMs >= sunriseMs && cityNowMs <= sunsetMs
    }
  }

  async function handleCurrentLocation(): Promise<void> {
    if (selectedCoords === null) {
      await loadCurrentLocationWeather(reloadTransitionDuration)
      return
    }

    showCityPicker = false
    await loadCurrentLocationWeather(transitionDuration)
  }

  function openCityPicker(): void {
    weatherViewTransitionDuration = transitionDuration
    showCityPicker = true
  }

  function closeCityPicker(): void {
    showCityPicker = false
  }

  function handleCitySelect(event: CustomEvent<Coords>): void {
    selectedCoords = event.detail
    showCityPicker = false
    loadWeather()
  }

  onMount(() => {
    loadWeather()
    cityClockTimer = setInterval(() => {
      cityClockNow = Date.now()
    }, cityClockUpdateIntervalMs)
  })

  onDestroy(() => clearInterval(cityClockTimer))

  $: scene = weather
      ? getWeatherScene(weather.weatherCode, weather.isDay, {
        windSpeed: weather.windSpeed,
        rainProbability: weather.rainProbability
      })
    : null
  $: iconUrl = weather ? getWeatherIcon(weather.weatherCode, weather.isDay) : null
  $: cityTime = weather ? formatCityTime(cityClockNow, weather.utcOffsetSeconds) : null
  $: cityDate = weather ? formatCityDate(cityClockNow, weather.utcOffsetSeconds) : null
  $: sunPosition = weather
    ? getSunPosition(cityClockNow, weather.utcOffsetSeconds, weather.sunrise, weather.sunset)
    : null
  $: chartPoints = weather ? buildChartPoints(weather.hourlyTemperature) : []
  $: chartLinePath = buildChartLinePath(chartPoints)
  $: chartAreaPath = buildChartAreaPath(chartLinePath, chartPoints)
  $: chartLabels = weather ? buildChartLabels(weather.hourlyTime, chartPoints) : []
  $: chartMinTemperature = weather && weather.hourlyTemperature.length ? Math.min(...weather.hourlyTemperature) : null
  $: chartMaxTemperature = weather && weather.hourlyTemperature.length ? Math.max(...weather.hourlyTemperature) : null
</script>

<div class="weather-widget" style={scene ? `background-image: ${scene.gradient}` : ''}>
  {#if scene && !showCityPicker}
    <div class="sky" aria-hidden="true">
      {#if scene.night}
        <Moon class="sky__moon" size={weatherMoonIconSize} />
        {#each stars as star}
          <span class="sky__star" style="top: {star.top}%; left: {star.left}%; animation-delay: {star.delay}s;" />
        {/each}
      {/if}
      {#if scene.rainAnimation}
        <img class="sky__weather sky__weather--rain" src={scene.rainAnimation} alt="" />
      {/if}
    </div>
  {/if}

  <WidgetNavigation on:home>
    <svelte:fragment slot="actions">
      <Button ghost on:click={refresh} aria-label={refreshWeatherLabel}>
        <RotateCw size={weatherActionIconSize} />
      </Button>
      <Button ghost on:click={handleCurrentLocation} aria-label={currentLocationLabel}>
        <LocateFixed size={weatherActionIconSize} />
      </Button>
      <Button ghost on:click={openCityPicker} aria-label={changeCityLabel}>
        <MapPin size={weatherActionIconSize} />
      </Button>
    </svelte:fragment>
  </WidgetNavigation>

  {#if showCityPicker}
    <CityPicker on:select={handleCitySelect} on:close={closeCityPicker} />
  {:else}
    {#key weatherViewKey}
    <div
      class="weather"
      style:container-type={weatherContainerType}
      out:scale={{ duration: weatherViewTransitionDuration }}
      in:scale={{ duration: weatherViewTransitionDuration, delay: weatherViewTransitionDuration }}
    >
  {#if error}
    <Message variant="error" message={error} />
  {:else if loading || weather === null || scene === null || iconUrl === null}
    <p class="loading">Loading weather…</p>
  {:else}
    <div class="weather__content">
      <div class="sun-cycle" aria-label={sunCycleLabel}>
        <div class="sun-cycle__visual" aria-hidden="true">
          <svg class="sun-cycle__arc" viewBox={sunViewBox} preserveAspectRatio={sunAspectRatio}>
            <path d={sunPath} />
          </svg>
          {#if sunPosition}
            <span
              class="sun-cycle__sun"
              class:sun-cycle__sun--daylight={sunPosition.daylight}
              style="left: {sunPosition.left}%; bottom: {sunPosition.bottom}%;"
            >
              <Sun size={sunIconSize} />
            </span>
          {/if}
        </div>
        <div class="sun-cycle__times">
          <span class="sun-cycle__event">
            <span class="sun-cycle__event-icon sun-cycle__event-icon--wake"><SunMedium size={sunIconSize} /></span>
            <span>{formatSunTime(weather.sunrise)}</span>
          </span>
          <span class="sun-cycle__event">
            <span class="sun-cycle__event-icon sun-cycle__event-icon--sleep"><MoonStar size={sunIconSize} /></span>
            <span>{formatSunTime(weather.sunset)}</span>
          </span>
        </div>
      </div>

      <div class="weather__header">
        <span class="weather__city">{weather.city}</span>
        <span class="weather__date">
          {cityDate}{#if cityTime}<span class="weather__time"> · {cityTime}</span>{/if}
        </span>
      </div>

      <div class="weather__main">
        <div class="card">
          <div class="card__scene">
            <img class="card__icon" src={iconUrl} alt={scene.label} />
          </div>

          <div class="card__readout">
            <span class="card__temperature">{Math.round(weather.temperature)}°</span>
            <span class="card__condition">{scene.label}</span>
            <span class="card__range">
              {Math.round(weather.dailyTemperatureMin[0])}°/{Math.round(weather.dailyTemperatureMax[0])}°
            </span>
          </div>

          <div class="details">
            {#each windEffects.slice(0, scene.windEffectCount) as effect}
              <img class="details__weather details__weather--{effect}" src={scene.windAnimation} alt="" />
            {/each}
            <div class="detail">
              <Wind size={detailIconSize} />
              <span class="detail__label">Wind</span>
              <span class="detail__value">{weather.windSpeed} km/h</span>
            </div>
            <div class="detail">
              <Droplets size={detailIconSize} />
              <span class="detail__label">Rain</span>
              <span class="detail__value">{weather.rainProbability}%</span>
            </div>
            <div class="detail">
              <Sun size={detailIconSize} />
              <span class="detail__label">UV Index</span>
              <span class="detail__value">{weather.uvIndex}</span>
            </div>
            <div class="detail">
              <CircleGauge size={detailIconSize} />
              <span class="detail__label">Pressure</span>
              <span class="detail__value">{Math.round(weather.pressure)} hPa</span>
            </div>
            <div class="detail">
              <Leaf size={detailIconSize} />
              <span class="detail__label">Air Quality</span>
              <span class="detail__value">{Math.round(weather.airQuality)} AQI</span>
            </div>
          </div>
        </div>
      </div>

      <div class="chart" aria-label="Hourly temperature trend">
        <div class="chart__scale">
          {#if chartMaxTemperature !== null}<span class="chart__scale-value">{Math.round(chartMaxTemperature)}°</span>{/if}
          {#if chartMinTemperature !== null}<span class="chart__scale-value">{Math.round(chartMinTemperature)}°</span>{/if}
        </div>
        <div class="chart__plot">
          <svg class="chart__svg" viewBox="0 0 {chartViewBoxWidth} {chartViewBoxHeight}" preserveAspectRatio="none">
            <defs>
              <linearGradient id={chartGradientId} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stop-color="var(--chart-fill-color-top)" />
                <stop offset="100%" stop-color="var(--chart-fill-color-bottom)" />
              </linearGradient>
            </defs>
            <path class="chart__area" d={chartAreaPath} fill="url(#{chartGradientId})" />
            <path class="chart__line" d={chartLinePath} />
          </svg>
          <div class="chart__labels">
            {#each chartLabels as label (label.x)}
              <span
                class="chart__label"
                style="left: {(label.x / chartViewBoxWidth) * 100}%; transform: translateX({label.align === 'start' ? '0' : '-50%'});"
              >{label.label}</span>
            {/each}
          </div>
        </div>
      </div>

      <div class="forecast">
        <div class="forecast__row">
          {#each weather.dailyTime as date, index (date)}
            {#if index > 0}
              <div class="forecast__card" style:container-type={forecastCardContainerType}>
                <span class="forecast__day">{formatDayLabel(date)}</span>
                <img class="forecast__icon" src={getWeatherIcon(weather.dailyWeatherCode[index], true)} alt="" />
                <span class="forecast__temps">
                  <span class="forecast__high">{Math.round(weather.dailyTemperatureMax[index])}°</span>/<span
                    class="forecast__low">{Math.round(weather.dailyTemperatureMin[index])}°</span>
                </span>
              </div>
            {/if}
          {/each}
        </div>
      </div>
    </div>
  {/if}
    </div>
    {/key}
  {/if}
  {#if scene?.night}
    <div class="weather-widget__night-overlay" aria-hidden="true" />
  {/if}
</div>

<style>
  .weather-widget {
    --weather-widget-edge: 0;
    --weather-widget-fill: 100%;
    --weather-widget-minimum: 0;
    --weather-widget-padding: 1rem;
    --weather-widget-gap: 1rem;
    --weather-widget-night-overlay-color: #000000;
    --weather-widget-night-overlay-opacity: 0.18;
    --weather-widget-night-overlay-layer: 1001;
    position: absolute;
    inset: var(--weather-widget-edge);
    display: flex;
    flex-direction: column;
    width: var(--weather-widget-fill);
    height: var(--weather-widget-fill);
    min-height: var(--weather-widget-minimum);
    box-sizing: border-box;
    padding: var(--weather-widget-padding);
    gap: var(--weather-widget-gap);
  }

  .weather-widget__night-overlay {
    position: absolute;
    inset: var(--weather-widget-edge);
    z-index: var(--weather-widget-night-overlay-layer);
    background-color: var(--weather-widget-night-overlay-color);
    opacity: var(--weather-widget-night-overlay-opacity);
    pointer-events: none;
  }

  .sky {
    --rain-animation-size: clamp(8rem, 15vw, 13rem);
    --weather-animation-opacity: 0.58;
    --moon-block-position: clamp(3.5rem, 8vh, 6rem);
    --moon-inline-position: clamp(2rem, 8vw, 8rem);
    --rain-inline-position: 10%;
    --rain-block-position: 6%;
    --moon-color: #f4f1de;
    --star-size: 3px;
    --star-radius: 50%;
    --star-color: #ffffff;
    --star-animation-duration: 3s;
    --star-dimmed-opacity: 0.2;
    --star-full-opacity: 1;
    position: absolute;
    inset: var(--weather-widget-edge);
    overflow: hidden;
    pointer-events: none;
  }

  .sky__weather {
    position: absolute;
    opacity: var(--weather-animation-opacity);
    object-fit: contain;
  }

  .sky__weather--rain {
    right: var(--rain-inline-position);
    bottom: var(--rain-block-position);
    width: var(--rain-animation-size);
    height: var(--rain-animation-size);
  }

  .sky :global(.sky__moon) {
    position: absolute;
    top: var(--moon-block-position);
    left: var(--moon-inline-position);
    color: var(--moon-color);
  }

  .sky__star {
    position: absolute;
    width: var(--star-size);
    height: var(--star-size);
    border-radius: var(--star-radius);
    background-color: var(--star-color);
    animation: twinkle var(--star-animation-duration) ease-in-out infinite;
  }

  @keyframes twinkle {
    0%,
    100% {
      opacity: var(--star-dimmed-opacity);
    }
    50% {
      opacity: var(--star-full-opacity);
    }
  }

  .weather {
    --content-edge-margin: min(1rem, 2cqw);
    --content-width: calc(100cqw - (2 * var(--content-edge-margin)));
    --content-row-width: 100%;
    --content-row-grow: 0;
    --content-row-shrink: 0;
    --dashboard-shrink: 1;
    --dashboard-main-grow: 0;
    --dashboard-main-basis: 100%;
    --minimum-size: 0;
    --grid-flexible-track: 1fr;
    --forecast-day-count: 7;
    --weather-overflow: hidden;
    --chart-height: min(3rem, 7cqh);
    --chart-gap: min(0.75rem, 1.5cqw);
    --chart-scale-width: min(3rem, 4cqw);
    --chart-label-font-size: clamp(0.65rem, 1.8cqh, 1rem);
    --chart-label-row-height: min(1.1rem, 3cqh);
    --chart-line-color: var(--sun-color);
    --chart-line-width: 0.6;
    --chart-fill-color-top: rgba(246, 183, 60, 0.45);
    --chart-fill-color-bottom: rgba(246, 183, 60, 0);
    --header-gap: min(0.25rem, 0.5cqh);
    --sun-cycle-width: min(42rem, 52cqw);
    --sun-cycle-gap: min(0.6rem, 1cqh);
    --panel-gap: min(0.75rem, 1.2cqh, 1.5cqw);
    --content-bottom-padding: min(0.75rem, 1.5cqh);
    --text-color: #4a3f66;
    --icon-glow-color: rgba(255, 255, 255, 0.35);
    --row-background: rgba(255, 255, 255, 0.5);
    --sun-arc-color: rgba(74, 63, 102, 0.3);
    --sun-color: #f6b73c;
    --sun-glow-color: rgba(246, 183, 60, 0.55);
    --sun-glow-size: min(2rem, 4cqh);
    --sun-motion-duration: 15s;
    --sun-entry-duration: 1.2s;
    --sun-visual-height: min(3rem, 7cqh);
    --sun-arc-stroke-width: 1;
    --sun-arc-length: 160;
    --minimum-sun-arc-offset: 0;
    --sun-marker-horizontal-offset: -50%;
    --sun-marker-vertical-offset: 50%;
    --sun-glow-duration: 2.4s;
    --sun-time-font-size: clamp(0.65rem, 1.6cqh, 0.9rem);
    --sun-time-gap: min(0.35rem, 0.5cqh);
    --sun-event-animation-duration: 3.6s;
    --sun-event-rise-distance: max(-0.18rem, -0.72cqmin);
    --sun-event-set-distance: min(0.18rem, 0.72cqmin);
    --sun-event-dimmed-opacity: 0.72;
    --sun-dimmed-opacity: 0.8;
    --sun-full-opacity: 1;
    --weather-font-size: clamp(0.7rem, 1.8cqh, 1rem);
    --weather-icon-size: clamp(0.9rem, 3cqh, 1.5rem);
    --card-padding: min(1rem, 2cqh, 2cqw);
    --card-gap: min(3rem, 4cqw);
    --card-section-min-width: min(10rem, 100%);
    --city-font-size: clamp(1.35rem, 6cqh, 2.5rem);
    --secondary-font-size: clamp(0.65rem, 1.8cqh, 1rem);
    --scene-size: min(16rem, 32cqh, 20cqw);
    --scene-art-size: 135%;
    --scene-art-max-width: none;
    --scene-shadow-blur: min(60px, 8cqh);
    --temperature-font-size: clamp(2.75rem, 14cqh, 6rem);
    --compact-gap: min(0.6rem, 1cqh);
    --condition-font-size: clamp(1rem, 3.5cqh, 1.75rem);
    --detail-row-gap: min(1rem, 2.5cqh);
    --detail-gap: min(1rem, 1.5cqw);
    --detail-icon-column-width: min(2rem, 4cqw);
    --detail-label-column-width: min(7rem, 10cqw);
    --forecast-gap: min(1rem, 1.5cqw);
    --forecast-card-gap: min(0.35rem, 0.5cqh);
    --forecast-card-padding-block: min(0.5rem, 1cqh);
    --forecast-card-padding-inline: min(0.75rem, 1cqw);
    --forecast-card-border-width: 1px;
    --forecast-card-border-color: rgba(74, 63, 102, 0.18);
    --forecast-row-radius: min(16px, 2cqmin);
    --forecast-icon-size: min(2.5rem, 5cqh, 54cqw);
    --wind-decoration-opacity: 0.28;
    --wind-decoration-size: min(12rem, 28cqw, 32cqh);
    --wind-decoration-layer: 0;
    --detail-layer: 1;
    --wind-one-transform: rotate(-8deg) scale(1.05);
    --wind-two-transform: rotate(6deg) scale(0.82);
    --wind-three-transform: rotate(-3deg) scale(0.68);
    --forecast-card-font-size: clamp(0.65rem, 1.8cqh, 0.85rem);
    --float-distance: min(6px, 1.5cqmin);
    display: flex;
    flex: 1;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--panel-gap);
    width: 100%;
    height: var(--weather-widget-fill);
    min-height: var(--minimum-size);
    overflow: var(--weather-overflow);
    font-size: var(--weather-font-size);
  }

  .loading {
    opacity: 0.7;
  }

  .weather__main {
    display: flex;
    flex-direction: column;
    flex: var(--dashboard-main-grow) var(--dashboard-shrink) var(--dashboard-main-basis);
    min-width: var(--minimum-size);
  }

  .weather__content {
    position: relative;
    display: flex;
    flex-flow: row wrap;
    align-items: center;
    align-content: space-between;
    justify-content: space-between;
    gap: var(--panel-gap);
    width: var(--content-width);
    max-width: 100%;
    height: 100%;
    min-height: var(--minimum-size);
    box-sizing: border-box;
    padding-bottom: var(--content-bottom-padding);
  }

  .weather__header {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--header-gap);
    flex: var(--content-row-grow) var(--content-row-shrink) var(--content-row-width);
    color: var(--text-color);
    text-align: center;
  }

  .weather__city {
    font-size: var(--city-font-size);
    font-weight: 600;
  }

  .weather__date {
    font-size: var(--secondary-font-size);
    opacity: 0.9;
  }

  .card {
    position: relative;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(var(--card-section-min-width), var(--grid-flexible-track)));
    align-items: center;
    width: 100%;
    box-sizing: border-box;
    padding: var(--card-padding);
    gap: var(--card-gap);
    color: var(--text-color);
  }

  .card__scene {
    position: relative;
    flex: 0 0 auto;
    width: var(--scene-size);
    height: var(--scene-size);
    display: flex;
    align-items: center;
    justify-content: center;
    justify-self: center;
  }

  .card__icon {
    flex-shrink: var(--content-row-shrink);
    width: var(--scene-art-size);
    height: var(--scene-art-size);
    max-width: var(--scene-art-max-width);
    object-fit: contain;
    filter: drop-shadow(0 0 var(--scene-shadow-blur) var(--icon-glow-color));
    animation: floatIcon 6s ease-in-out infinite;
  }

  .card__readout {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: var(--compact-gap);
    flex: 0 1 auto;
    min-width: 0;
  }

  .card__temperature {
    font-size: var(--temperature-font-size);
    font-weight: 700;
    line-height: 1;
  }

  .card__range {
    font-size: var(--secondary-font-size);
    opacity: 0.9;
  }

  .card__condition {
    font-size: var(--condition-font-size);
    font-weight: 600;
  }

  .details {
    position: relative;
    isolation: isolate;
    display: flex;
    flex-direction: column;
    gap: var(--detail-row-gap);
    flex: 0 0 auto;
  }

  .details__weather {
    position: absolute;
    z-index: var(--wind-decoration-layer);
    top: 50%;
    left: 50%;
    width: var(--wind-decoration-size);
    height: var(--wind-decoration-size);
    opacity: var(--wind-decoration-opacity);
    object-fit: contain;
    pointer-events: none;
  }

  .details__weather--wind-one {
    transform: translate(-50%, -50%) var(--wind-one-transform);
  }

  .details__weather--wind-two {
    transform: translate(-50%, -50%) var(--wind-two-transform);
  }

  .details__weather--wind-three {
    transform: translate(-50%, -50%) var(--wind-three-transform);
  }

  .sun-cycle {
    position: relative;
    align-self: center;
    display: flex;
    flex-direction: column;
    flex: var(--content-row-grow) var(--content-row-shrink) var(--content-row-width);
    align-items: center;
    gap: var(--sun-cycle-gap);
    width: var(--sun-cycle-width);
    max-width: 100%;
    color: var(--text-color);
    pointer-events: none;
  }

  .sun-cycle__visual {
    position: relative;
    width: var(--sun-cycle-width);
    max-width: 100%;
    height: var(--sun-visual-height);
  }

  .sun-cycle__arc {
    width: 100%;
    height: 100%;
    overflow: visible;
  }

  .sun-cycle__arc path {
    fill: none;
    stroke: var(--sun-arc-color);
    stroke-width: var(--sun-arc-stroke-width);
    stroke-dasharray: var(--sun-arc-length);
    animation: drawSunArc var(--sun-entry-duration) ease-out both;
  }

  .sun-cycle__sun {
    position: absolute;
    display: flex;
    color: var(--sun-color);
    transform: translate(var(--sun-marker-horizontal-offset), var(--sun-marker-vertical-offset));
    transition: left var(--sun-motion-duration) linear, bottom var(--sun-motion-duration) linear;
  }

  .sun-cycle__sun--daylight {
    filter: drop-shadow(0 0 var(--sun-glow-size) var(--sun-glow-color));
    animation: sunGlow var(--sun-glow-duration) ease-in-out infinite;
  }

  .sun-cycle__times {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: var(--sun-cycle-width);
    max-width: 100%;
    font-size: var(--sun-time-font-size);
  }

  .sun-cycle__event {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--sun-time-gap);
  }

  .sun-cycle__event-icon {
    display: flex;
  }

  .sun-cycle__event-icon--wake {
    animation: wakeIcon var(--sun-event-animation-duration) ease-in-out infinite;
  }

  .sun-cycle__event-icon--sleep {
    animation: sleepIcon var(--sun-event-animation-duration) ease-in-out infinite;
  }

  .detail {
    position: relative;
    z-index: var(--detail-layer);
    display: grid;
    grid-template-columns: var(--detail-icon-column-width) var(--detail-label-column-width) minmax(var(--minimum-size), var(--grid-flexible-track));
    align-items: center;
    gap: var(--detail-gap);
  }

  .detail :global(.lucide) {
    justify-self: end;
  }

  .detail__label {
    opacity: 0.8;
    text-align: left;
    white-space: nowrap;
  }

  .detail__value {
    font-weight: 600;
    text-align: left;
    white-space: nowrap;
  }

  .chart {
    display: flex;
    align-items: stretch;
    gap: var(--chart-gap);
    flex: var(--content-row-grow) var(--content-row-shrink) var(--content-row-width);
    width: var(--content-row-width);
    color: var(--text-color);
  }

  .chart__scale {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: flex-end;
    flex-shrink: 0;
    width: var(--chart-scale-width);
    padding-bottom: var(--chart-label-row-height);
  }

  .chart__scale-value {
    font-size: var(--chart-label-font-size);
    opacity: 0.8;
  }

  .chart__plot {
    position: relative;
    flex: 1 1 auto;
    min-width: 0;
  }

  .chart__svg {
    display: block;
    width: 100%;
    height: var(--chart-height);
    overflow: visible;
  }

  .chart__area {
    stroke: none;
  }

  .chart__line {
    fill: none;
    stroke: var(--chart-line-color);
    stroke-width: var(--chart-line-width);
    stroke-linecap: round;
    vector-effect: non-scaling-stroke;
  }

  .chart__labels {
    position: relative;
    width: 100%;
    height: var(--chart-label-row-height);
  }

  .chart__label {
    position: absolute;
    top: 0;
    font-size: var(--chart-label-font-size);
    opacity: 0.7;
    white-space: nowrap;
  }

  .forecast {
    display: flex;
    flex-direction: column;
    gap: var(--forecast-gap);
    flex: var(--content-row-grow) var(--content-row-shrink) var(--content-row-width);
    width: var(--content-row-width);
  }

  .forecast__row {
    display: grid;
    grid-template-columns: repeat(var(--forecast-day-count), minmax(var(--minimum-size), var(--grid-flexible-track)));
    gap: var(--forecast-gap);
    width: 100%;
  }

  .forecast__card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--forecast-card-gap);
    min-width: 0;
    box-sizing: border-box;
    padding: var(--forecast-card-padding-block) var(--forecast-card-padding-inline);
    border: var(--forecast-card-border-width) solid var(--forecast-card-border-color);
    border-radius: var(--forecast-row-radius);
    background-color: var(--row-background);
    color: var(--text-color);
  }

  .forecast__temps {
    display: flex;
    gap: var(--compact-gap);
    font-size: var(--forecast-card-font-size);
  }

  .forecast__day {
    font-size: var(--forecast-card-font-size);
    font-weight: 600;
  }

  .forecast__icon {
    width: var(--forecast-icon-size);
    height: var(--forecast-icon-size);
  }

  .weather :global(.lucide) {
    width: var(--weather-icon-size);
    height: var(--weather-icon-size);
    flex-shrink: 0;
  }

  .forecast__high {
    font-weight: 600;
  }

  .forecast__low {
    opacity: 0.7;
  }

  @keyframes floatIcon {
    0%,
    100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(calc(var(--float-distance) * -1));
    }
  }

  @keyframes drawSunArc {
    from {
      stroke-dashoffset: var(--sun-arc-length);
    }
    to {
      stroke-dashoffset: var(--minimum-sun-arc-offset);
    }
  }

  @keyframes sunGlow {
    0%,
    100% {
      opacity: var(--sun-dimmed-opacity);
    }
    50% {
      opacity: var(--sun-full-opacity);
    }
  }

  @keyframes wakeIcon {
    0%,
    100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(var(--sun-event-rise-distance));
    }
  }

  @keyframes sleepIcon {
    0%,
    100% {
      transform: translateY(0);
      opacity: var(--sun-full-opacity);
    }
    50% {
      transform: translateY(var(--sun-event-set-distance));
      opacity: var(--sun-event-dimmed-opacity);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .sky__star,
    .card__icon,
    .sun-cycle__arc path,
    .sun-cycle__sun--daylight,
    .sun-cycle__event-icon {
      animation: none;
    }

    .sun-cycle__sun {
      transition: none;
    }
  }
</style>
