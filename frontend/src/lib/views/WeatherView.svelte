<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import { scale } from 'svelte/transition'
  import { CalendarDays, Clock, Droplets, MoonStar, Sun, SunMedium } from 'lucide-svelte'
  import { GetWeather, GetWeatherAt } from '../../../wailsjs/go/main/App'
  import type { main } from '../../../wailsjs/go/models'
  import { Button, Message } from '../components'
  import { getWeatherIcon } from './weatherIcon'
  import { getWeatherScene, type WeatherScene } from './weatherScene'
  import { transitionDuration } from '../transition'

  export let coords: { lat: number; lon: number } | null = null

  const dispatch = createEventDispatcher<{ scene: WeatherScene | null }>()
  const cityClockUpdateIntervalMs = 15000
  const millisecondsPerSecond = 1000
  const minimumSunProgress = 0
  const maximumSunProgress = 1
  const sunArcHeightPercent = 94
  const sunIconSize = 18
  const sunPath = 'M 1 63 Q 50 -58 99 63'
  const sunViewBox = '0 0 100 64'
  const sunAspectRatio = 'none'
  const sunCycleLabel = "Today's sunrise and sunset"
  const weatherContainerType = 'size'
  const timestampHourStart = 11
  const timestampMinuteEnd = 16
  const timestampUTCDesignator = ':00Z'
  const percentageScale = 100

  let weather: main.Weather | null = null
  let error: string | null = null
  let loading = true
  let forecastView: 'day' | 'week' = 'day'
  let hourlyList: HTMLDivElement
  let dragStartY: number | null = null
  let dragStartScrollTop = 0
  let cityClockNow = Date.now()
  let cityClockTimer: ReturnType<typeof setInterval>

  export async function refresh(): Promise<void> {
    await loadWeather()
  }

  async function loadWeather(): Promise<void> {
    loading = true
    error = null
    try {
      weather = coords ? await GetWeatherAt(coords.lat, coords.lon) : await GetWeather()
    } catch (err) {
      error = String(err)
    } finally {
      loading = false
    }
  }

  function formatDayLabel(dateString: string): string {
    return new Date(dateString).toLocaleDateString(undefined, { weekday: 'short' })
  }

  function formatHourLabel(dateString: string): string {
    return new Date(dateString).toLocaleTimeString(undefined, { hour: '2-digit', hourCycle: 'h23' })
  }

  function handleListPointerDown(event: PointerEvent): void {
    if (forecastView !== 'day') return
    dragStartY = event.clientY
    dragStartScrollTop = hourlyList.scrollTop
    hourlyList.setPointerCapture(event.pointerId)
  }

  function handleListPointerMove(event: PointerEvent): void {
    if (dragStartY === null) return
    hourlyList.scrollTop = dragStartScrollTop - (event.clientY - dragStartY)
  }

  function handleListPointerUp(): void {
    dragStartY = null
  }

  function formatCityTime(nowMs: number, utcOffsetSeconds: number): string {
    const shifted = new Date(nowMs + utcOffsetSeconds * millisecondsPerSecond)
    return shifted.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit', timeZone: 'UTC', hourCycle: 'h23' })
  }

  function formatSunTime(timestamp: string): string {
    return timestamp.slice(timestampHourStart, timestampMinuteEnd)
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
  $: sunPosition = weather
    ? getSunPosition(cityClockNow, weather.utcOffsetSeconds, weather.sunrise, weather.sunset)
    : null
  $: dispatch('scene', scene)
</script>

<div
  class="weather"
  style:container-type={weatherContainerType}
  out:scale={{ duration: transitionDuration }}
  in:scale={{ duration: transitionDuration, delay: transitionDuration }}
>
  {#if error}
    <Message variant="error" message={error} />
  {:else if loading || weather === null || scene === null || iconUrl === null}
    <p class="loading">Loading weather…</p>
  {:else}
    <div class="weather__content">
      <div class="weather__main">
        <div class="card">
          <div class="card__header">
            <span class="card__city">{weather.city}</span>
            {#if cityTime}<span class="card__time">{cityTime}</span>{/if}
          </div>

          <div class="card__scene">
            <img class="card__icon" src={iconUrl} alt={scene.label} />
            <span class="card__temperature">{Math.round(weather.temperature)}°</span>
          </div>

          <div class="card__footer">
            <span class="card__range">
              {Math.round(weather.dailyTemperatureMin[0])}°/{Math.round(weather.dailyTemperatureMax[0])}°
            </span>
            <span class="card__condition">{scene.label}</span>
          </div>
        </div>

        <div class="details">
          <div class="detail">
            <span>{weather.windSpeed} km/h</span>
          </div>
          <div class="detail">
            <Droplets size={18} />
            <span>{weather.rainProbability}%</span>
          </div>
          <div class="detail">
            <Sun size={18} />
            <span>UV {weather.uvIndex}</span>
          </div>
        </div>
      </div>

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

      <div class="forecast">
        <div class="forecast__toggle">
          <Button ghost={forecastView !== 'day'} on:click={() => (forecastView = 'day')} aria-label="Hourly forecast">
            <Clock size={18} />
          </Button>
          <Button ghost={forecastView !== 'week'} on:click={() => (forecastView = 'week')} aria-label="Weekly forecast">
            <CalendarDays size={18} />
          </Button>
        </div>

        <div
          class="forecast__list"
          class:forecast__list--scroll={forecastView === 'day'}
          bind:this={hourlyList}
          on:pointerdown={handleListPointerDown}
          on:pointermove={handleListPointerMove}
          on:pointerup={handleListPointerUp}
          on:pointercancel={handleListPointerUp}
        >
          {#if forecastView === 'week'}
            {#each weather.dailyTime as date, index (date)}
              {#if index > 0}
                <div class="forecast__row">
                  <span class="forecast__day">{formatDayLabel(date)}</span>
                  <img class="forecast__icon" src={getWeatherIcon(weather.dailyWeatherCode[index], true)} alt="" />
                  <span class="forecast__temps">
                    <span class="forecast__high">{Math.round(weather.dailyTemperatureMax[index])}°</span>/<span
                      class="forecast__low">{Math.round(weather.dailyTemperatureMin[index])}°</span>
                  </span>
                </div>
              {/if}
            {/each}
          {:else}
            {#each weather.hourlyTime as time, index (time)}
              <div class="forecast__row forecast__row--hourly">
                <span class="forecast__day">{formatHourLabel(time)}</span>
                <img class="forecast__icon" src={getWeatherIcon(weather.hourlyWeatherCode[index], true)} alt="" />
                <span class="forecast__high">{Math.round(weather.hourlyTemperature[index])}°</span>
              </div>
            {/each}
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .weather {
    --panel-width: min(300px, 38cqw, 62cqh);
    --details-height: min(1.5rem, 6cqmin);
    --sun-cycle-width: min(9rem, 22cqw, 36cqmin);
    --sun-cycle-height: min(4.25rem, 17cqmin);
    --sun-cycle-horizontal-position: 50%;
    --sun-cycle-horizontal-offset: -50%;
    --sun-cycle-vertical-position: 50%;
    --sun-cycle-vertical-offset: -50%;
    --panel-gap: min(1rem, 4cqmin);
    --panel-height: calc(var(--panel-width) * 4 / 3 + var(--panel-gap) + var(--details-height));
    --text-color: #4a3f66;
    --icon-glow-color: rgba(255, 255, 255, 0.35);
    --row-background: rgba(255, 255, 255, 0.5);
    --hourly-row-gap: min(1.95rem, 8cqmin);
    --sun-arc-color: rgba(74, 63, 102, 0.3);
    --sun-color: #f6b73c;
    --sun-glow-color: rgba(246, 183, 60, 0.55);
    --sun-glow-size: min(0.8rem, 3.2cqmin);
    --sun-motion-duration: 15s;
    --sun-entry-duration: 1.2s;
    --sun-visual-height: min(2.5rem, 10cqmin);
    --sun-arc-stroke-width: 1;
    --sun-arc-length: 160;
    --minimum-sun-arc-offset: 0;
    --sun-marker-horizontal-offset: -50%;
    --sun-marker-vertical-offset: 50%;
    --sun-glow-duration: 2.4s;
    --sun-time-font-size: min(0.8rem, 3.2cqmin);
    --sun-time-gap: min(0.35rem, 1.4cqmin);
    --sun-event-animation-duration: 3.6s;
    --sun-event-rise-distance: max(-0.18rem, -0.72cqmin);
    --sun-event-set-distance: min(0.18rem, 0.72cqmin);
    --sun-event-dimmed-opacity: 0.72;
    --sun-dimmed-opacity: 0.8;
    --sun-full-opacity: 1;
    --weather-font-size: min(1rem, 4cqmin);
    --weather-icon-size: min(18px, 4.5cqmin);
    --card-padding: min(1.25rem, 5cqmin);
    --city-font-size: min(2.25rem, 9cqmin);
    --secondary-font-size: min(0.9rem, 3.6cqmin);
    --scene-margin: min(0.75rem, 3cqmin);
    --scene-shadow-blur: min(28px, 7cqmin);
    --temperature-offset: min(0.5rem, 2cqmin);
    --temperature-font-size: min(3rem, 12cqmin);
    --compact-gap: min(0.25rem, 1cqmin);
    --condition-font-size: min(1.1rem, 4.4cqmin);
    --details-gap: min(1.5rem, 6cqmin);
    --detail-gap: min(0.4rem, 1.6cqmin);
    --forecast-gap: min(0.5rem, 2cqmin);
    --forecast-day-column: min(3rem, 12cqmin);
    --forecast-hour-column: min(4rem, 16cqmin);
    --forecast-icon-column: min(2rem, 8cqmin);
    --forecast-row-padding-block: min(0.5rem, 2cqmin);
    --forecast-row-padding-inline: min(0.75rem, 3cqmin);
    --forecast-row-radius: min(8px, 2cqmin);
    --forecast-icon-size: min(24px, 6cqmin);
    --float-distance: min(6px, 1.5cqmin);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--panel-gap);
    width: 100%;
    height: 100%;
    font-size: var(--weather-font-size);
  }

  .loading {
    opacity: 0.7;
  }

  .weather__main {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--panel-gap);
    width: var(--panel-width);
    height: var(--panel-height);
    flex-shrink: 0;
  }

  .weather__content {
    position: relative;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-evenly;
    width: 100%;
    height: 100%;
  }

  .card {
    position: relative;
    display: flex;
    flex-direction: column;
    width: 100%;
    aspect-ratio: 3 / 4;
    box-sizing: border-box;
    padding: var(--card-padding);
    color: var(--text-color);
    overflow: hidden;
  }

  .card__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .card__city {
    font-size: var(--city-font-size);
    font-weight: 600;
  }

  .card__time {
    font-size: var(--secondary-font-size);
    opacity: 0.9;
  }

  .card__scene {
    position: relative;
    flex: 1;
    min-height: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: var(--scene-margin) 0;
  }

  .card__icon {
    width: 92%;
    height: 92%;
    object-fit: contain;
    filter: drop-shadow(0 0 var(--scene-shadow-blur) var(--icon-glow-color));
    animation: floatIcon 6s ease-in-out infinite;
  }

  .card__temperature {
    position: absolute;
    top: var(--temperature-offset);
    left: 0;
    font-size: var(--temperature-font-size);
    font-weight: 700;
    line-height: 1;
  }

  .card__footer {
    display: flex;
    flex-direction: column;
    gap: var(--compact-gap);
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
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--details-gap);
    width: 100%;
    height: var(--details-height);
  }

  .sun-cycle {
    position: absolute;
    top: var(--sun-cycle-vertical-position);
    left: var(--sun-cycle-horizontal-position);
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    width: var(--sun-cycle-width);
    height: var(--sun-cycle-height);
    color: var(--text-color);
    transform: translate(var(--sun-cycle-horizontal-offset), var(--sun-cycle-vertical-offset));
    pointer-events: none;
  }

  .sun-cycle__visual {
    position: relative;
    width: 100%;
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
    display: flex;
    align-items: center;
    gap: var(--detail-gap);
  }

  .forecast {
    display: flex;
    flex-direction: column;
    gap: var(--forecast-gap);
    width: 100%;
    max-width: var(--panel-width);
    height: var(--panel-height);
  }

  .forecast__toggle {
    display: flex;
    justify-content: center;
    gap: var(--forecast-gap);
  }

  .forecast__list {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: var(--forecast-gap);
    flex: 1;
    min-height: 0;
  }

  .forecast__list--scroll {
    justify-content: flex-start;
    gap: var(--hourly-row-gap);
    overflow-y: auto;
    scrollbar-width: none;
    -ms-overflow-style: none;
    touch-action: pan-y;
    cursor: grab;
  }

  .forecast__list--scroll:active {
    cursor: grabbing;
  }

  .forecast__list--scroll::-webkit-scrollbar {
    display: none;
  }

  .forecast__row {
    display: grid;
    grid-template-columns: var(--forecast-day-column) var(--forecast-icon-column) 1fr;
    align-items: center;
    padding: var(--forecast-row-padding-block) var(--forecast-row-padding-inline);
    border-radius: var(--forecast-row-radius);
    background-color: var(--row-background);
  }

  .forecast__temps {
    text-align: right;
  }

  .forecast__row--hourly {
    grid-template-columns: var(--forecast-hour-column) var(--forecast-icon-column) 1fr;
  }

  .forecast__row--hourly .forecast__day {
    white-space: nowrap;
  }

  .forecast__day {
    font-weight: 600;
  }

  .forecast__icon {
    width: var(--forecast-icon-size);
    height: var(--forecast-icon-size);
  }

  .weather :global(.lucide) {
    width: var(--weather-icon-size);
    height: var(--weather-icon-size);
  }

  .forecast__toggle :global(.button) {
    padding: var(--forecast-row-padding-block) var(--forecast-row-padding-inline);
  }

  .forecast__high {
    text-align: right;
    font-weight: 600;
  }

  .forecast__low {
    text-align: right;
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
