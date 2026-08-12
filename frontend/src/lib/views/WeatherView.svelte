<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte'
  import { Droplets, Sun } from 'lucide-svelte'
  import { GetWeather, GetWeatherAt } from '../../../wailsjs/go/main/App'
  import type { main } from '../../../wailsjs/go/models'
  import { Message } from '../components'
  import { getWeatherIcon } from './weatherIcon'
  import { getWeatherScene, type WeatherScene } from './weatherScene'

  export let coords: { lat: number; lon: number } | null = null

  const dispatch = createEventDispatcher<{ scene: WeatherScene | null }>()

  let weather: main.Weather | null = null
  let error: string | null = null
  let loading = true

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

  function formatDayLabel(dateString: string, index: number): string {
    if (index === 0) {
      return 'Today'
    }
    return new Date(dateString).toLocaleDateString(undefined, { weekday: 'short' })
  }

  onMount(loadWeather)

  $: scene = weather
      ? getWeatherScene(weather.weatherCode, weather.isDay, {
        windSpeed: weather.windSpeed,
        rainProbability: weather.rainProbability
      })
    : null
  $: iconUrl = weather ? getWeatherIcon(weather.weatherCode, weather.isDay) : null
  $: dispatch('scene', scene)
  $: if (coords) loadWeather()
</script>

<div class="weather">
  {#if error}
    <Message variant="error" message={error} />
  {:else if loading || weather === null || scene === null || iconUrl === null}
    <p class="loading">Loading weather…</p>
  {:else}
    <div class="weather__main">
      <div class="card">
        <div class="card__header">
          <span class="card__city">{weather.city}</span>
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

    <div class="forecast">
      {#each weather.dailyTime as date, index (date)}
        <div class="forecast__row">
          <span class="forecast__day">{formatDayLabel(date, index)}</span>
          <img class="forecast__icon" src={getWeatherIcon(weather.dailyWeatherCode[index], true)} alt="" />
          <span class="forecast__high">{Math.round(weather.dailyTemperatureMax[index])}°</span>
          <span class="forecast__low">{Math.round(weather.dailyTemperatureMin[index])}°</span>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .weather {
    --panel-width: 300px;
    --details-height: 1.5rem;
    --panel-gap: 1rem;
    --panel-height: calc(var(--panel-width) * 4 / 3 + var(--panel-gap) + var(--details-height));
    --text-color: #4a3f66;
    --icon-glow-color: rgba(255, 255, 255, 0.35);
    --row-background: rgba(255, 255, 255, 0.5);
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-evenly;
    width: 100%;
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

  .card {
    position: relative;
    display: flex;
    flex-direction: column;
    width: 100%;
    aspect-ratio: 3 / 4;
    padding: 1.25rem;
    color: var(--text-color);
    overflow: hidden;
  }

  .card__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .card__city {
    font-size: 2.25rem;
    font-weight: 600;
  }

  .card__scene {
    position: relative;
    flex: 1;
    min-height: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0.75rem 0;
  }

  .card__icon {
    width: 92%;
    height: 92%;
    object-fit: contain;
    filter: drop-shadow(0 0 28px var(--icon-glow-color));
    animation: floatIcon 6s ease-in-out infinite;
  }

  .card__temperature {
    position: absolute;
    top: 0.5rem;
    left: 0;
    font-size: 3rem;
    font-weight: 700;
    line-height: 1;
  }

  .card__footer {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .card__range {
    font-size: 0.9rem;
    opacity: 0.9;
  }

  .card__condition {
    font-size: 1.1rem;
    font-weight: 600;
  }

  .details {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1.5rem;
    width: 100%;
    height: var(--details-height);
  }

  .detail {
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .forecast {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: 0.5rem;
    width: 100%;
    max-width: var(--panel-width);
    height: var(--panel-height);
  }

  .forecast__row {
    display: grid;
    grid-template-columns: 3rem 2rem 1fr 1fr;
    align-items: center;
    padding: 0.5rem 0.75rem;
    border-radius: 8px;
    background-color: var(--row-background);
  }

  .forecast__day {
    font-weight: 600;
  }

  .forecast__icon {
    width: 24px;
    height: 24px;
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
      transform: translateY(-6px);
    }
  }
</style>
