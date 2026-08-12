<script lang="ts">
  import { tick } from 'svelte'
  import { House, LocateFixed, MapPin, Moon, RotateCw } from 'lucide-svelte'
  import { Button, Message, Tile } from './lib/components'
  import WeatherView from './lib/views/WeatherView.svelte'
  import CityPicker from './lib/views/CityPicker.svelte'
  import type { WeatherScene } from './lib/views/weatherScene'
  import weatherTileAnimation from '@bybas/weather-icons/production/fill/all/partly-cloudy-day-rain.svg'

  type WidgetId = 'photos' | 'weather'
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

  const windEffects = ['wind-one', 'wind-two', 'wind-three']
  const weatherActionIconSize = 18
  const currentLocationLabel = 'Use current IP location'

  let activeWidget: WidgetId | null = null
  let weatherScene: WeatherScene | null = null
  let weatherViewRef: WeatherView
  let selectedCoords: Coords | null = null
  let showCityPicker = false

  function maximize(id: WidgetId): void {
    activeWidget = id
  }

  function goHome(): void {
    activeWidget = null
  }

  function handleWeatherScene(event: CustomEvent<WeatherScene | null>): void {
    weatherScene = event.detail
  }

  function handleRefresh(): void {
    weatherViewRef?.refresh()
  }

  async function handleCurrentLocation(): Promise<void> {
    const wasCityPickerOpen = showCityPicker
    selectedCoords = null
    showCityPicker = false
    await tick()
    if (!wasCityPickerOpen) weatherViewRef?.refresh()
  }

  function openCityPicker(): void {
    showCityPicker = true
  }

  function closeCityPicker(): void {
    showCityPicker = false
  }

  function handleCitySelect(event: CustomEvent<Coords>): void {
    selectedCoords = event.detail
    showCityPicker = false
  }

  $: if (activeWidget !== 'weather') weatherScene = null
</script>

<main
  class="screen"
  class:screen--night={weatherScene?.night}
  style={weatherScene ? `background-image: ${weatherScene.gradient}` : ''}
>
  {#if weatherScene && !showCityPicker}
    <div class="sky" aria-hidden="true">
      {#if weatherScene.night}
        <Moon class="sky__moon" size={48} />
        {#each stars as star}
          <span class="sky__star" style="top: {star.top}%; left: {star.left}%; animation-delay: {star.delay}s;" />
        {/each}
      {/if}
      {#each windEffects.slice(0, weatherScene.windEffectCount) as effect}
        <img class="sky__weather sky__weather--{effect}" src={weatherScene.windAnimation} alt="" />
      {/each}
      {#if weatherScene.rainAnimation}
        <img class="sky__weather sky__weather--rain" src={weatherScene.rainAnimation} alt="" />
      {/if}
    </div>
  {/if}

  {#if activeWidget !== null}
    <div class="nav">
      <Button ghost on:click={goHome} aria-label="Home">
        <House />
      </Button>
      {#if activeWidget === 'weather'}
        <div class="nav__actions">
          <Button ghost on:click={handleRefresh} aria-label="Refresh weather">
            <RotateCw size={weatherActionIconSize} />
          </Button>
          <Button ghost on:click={handleCurrentLocation} aria-label={currentLocationLabel}>
            <LocateFixed size={weatherActionIconSize} />
          </Button>
          <Button ghost on:click={openCityPicker} aria-label="Change city">
            <MapPin size={weatherActionIconSize} />
          </Button>
        </div>
      {/if}
    </div>
  {/if}

  <div class="content">
    {#if activeWidget === null}
      <div class="tiles">
        <Tile title="Photos" on:toggle={() => maximize('photos')} />
        <Tile title="Weather" on:toggle={() => maximize('weather')}>
          <img slot="visual" class="tile__weather-scene" src={weatherTileAnimation} alt="" />
        </Tile>
      </div>
    {:else if activeWidget === 'weather' && showCityPicker}
      <CityPicker on:select={handleCitySelect} on:close={closeCityPicker} />
    {:else if activeWidget === 'weather'}
      <WeatherView bind:this={weatherViewRef} coords={selectedCoords} on:scene={handleWeatherScene} />
    {:else}
      <Message variant="warning" message="Not implemented yet" />
    {/if}
  </div>
</main>

<style>
  .screen {
    position: relative;
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    box-sizing: border-box;
    padding: 1rem;
    gap: 1rem;
    background-image: var(--app-background);
    background-size: 200% 200%;
    animation: driftGradient 30s ease-in-out infinite;
  }

  .screen--night {
    filter: brightness(0.82);
  }

  .nav {
    position: relative;
    z-index: 1;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .nav__actions {
    display: flex;
    gap: 0.5rem;
  }

  .sky {
    --wind-animation-size: clamp(8rem, 18vw, 16rem);
    --rain-animation-size: clamp(8rem, 15vw, 13rem);
    --weather-animation-opacity: 0.58;
    position: absolute;
    inset: 0;
    overflow: hidden;
    pointer-events: none;
  }

  .sky__weather {
    position: absolute;
    width: var(--wind-animation-size);
    height: var(--wind-animation-size);
    opacity: var(--weather-animation-opacity);
    object-fit: contain;
  }

  .sky__weather--wind-one {
    top: 12%;
    left: 14%;
    transform: rotate(-8deg);
  }

  .sky__weather--wind-two {
    top: 9%;
    right: 12%;
    transform: rotate(12deg) scale(0.82);
  }

  .sky__weather--wind-three {
    bottom: 7%;
    left: 9%;
    transform: rotate(7deg) scale(1.12);
  }

  .sky__weather--rain {
    right: 10%;
    bottom: 6%;
    width: var(--rain-animation-size);
    height: var(--rain-animation-size);
  }

  .sky :global(.sky__moon) {
    position: absolute;
    top: 4rem;
    left: 50%;
    transform: translateX(-50%);
    color: #f4f1de;
  }

  .sky__star {
    position: absolute;
    width: 3px;
    height: 3px;
    border-radius: 50%;
    background-color: #ffffff;
    animation: twinkle 3s ease-in-out infinite;
  }

  @keyframes twinkle {
    0%,
    100% {
      opacity: 0.2;
    }
    50% {
      opacity: 1;
    }
  }

  .content {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .tiles {
    --tile-height: 50vh;
    --tile-width: 30vw;
    display: flex;
    gap: 1rem;
  }

  .tiles :global(.tile) {
    flex: 0 0 auto;
    height: var(--tile-height);
    width: var(--tile-width);
  }

  .tiles :global(.tile__weather-scene) {
    width: 80%;
    height: 80%;
    object-fit: contain;
    filter: drop-shadow(0 0 28px rgba(255, 255, 255, 0.18));
    animation: floatIcon 6s ease-in-out infinite;
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

  @keyframes driftGradient {
    0%,
    100% {
      background-position: 0% 0%;
    }
    50% {
      background-position: 100% 100%;
    }
  }
</style>
