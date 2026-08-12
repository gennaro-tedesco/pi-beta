<script lang="ts">
  import { tick } from 'svelte'
  import { scale } from 'svelte/transition'
  import { House, LocateFixed, MapPin, Moon, RotateCw } from 'lucide-svelte'
  import { Button, Tile } from './lib/components'
  import WeatherView from './lib/views/WeatherView.svelte'
  import CityPicker from './lib/views/CityPicker.svelte'
  import PhotoSlideshow from './lib/views/PhotoSlideshow.svelte'
  import type { WeatherScene } from './lib/views/weatherScene'
  import cameraIllustration from './assets/images/camera-color.svg'
  import filmIllustration from './assets/images/film-frames-color.svg'
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

  const transitionDuration = 250

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
      <div
        class="tiles"
        out:scale={{ duration: transitionDuration }}
        in:scale={{ duration: transitionDuration, delay: transitionDuration }}
      >
        <Tile title="Photos" on:toggle={() => maximize('photos')}>
          <div slot="visual" class="tile__slideshow" aria-hidden="true">
            <img class="tile__ambient tile__ambient--camera" src={cameraIllustration} alt="" />
            <img class="tile__ambient tile__ambient--film" src={filmIllustration} alt="" />
            <div class="tile__landscapes">
              <div class="tile__landscape tile__landscape--alpine">
                <span class="tile__orb tile__orb--sun" />
                <span class="tile__ridge tile__ridge--far" />
                <span class="tile__ridge tile__ridge--near" />
              </div>
              <div class="tile__landscape tile__landscape--coast">
                <span class="tile__orb tile__orb--moon" />
                <span class="tile__wave tile__wave--far" />
                <span class="tile__wave tile__wave--near" />
              </div>
              <div class="tile__landscape tile__landscape--meadow">
                <span class="tile__cloud tile__cloud--one" />
                <span class="tile__cloud tile__cloud--two" />
                <span class="tile__hill tile__hill--far" />
                <span class="tile__hill tile__hill--near" />
              </div>
            </div>
          </div>
        </Tile>
        <Tile title="Weather" on:toggle={() => maximize('weather')}>
          <img slot="visual" class="tile__weather-scene" src={weatherTileAnimation} alt="" />
        </Tile>
      </div>
    {:else if activeWidget === 'weather' && showCityPicker}
      <CityPicker on:select={handleCitySelect} on:close={closeCityPicker} />
    {:else if activeWidget === 'weather'}
      <WeatherView bind:this={weatherViewRef} coords={selectedCoords} on:scene={handleWeatherScene} />
    {:else if activeWidget === 'photos'}
      <PhotoSlideshow />
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
    pointer-events: none;
  }

  .nav :global(button) {
    pointer-events: auto;
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

  .tiles :global(.tile__slideshow) {
    --ambient-edge-offset: 1%;
    --ambient-opacity: 0.72;
    --ambient-size: clamp(4.5rem, 10vw, 7rem);
    --ambient-travel: 0.5rem;
    --ambient-phase-delay: -3s;
    --ambient-duration: 6s;
    --landscape-duration: 15s;
    --landscape-radius: 1rem;
    --landscape-shadow: 0 1rem 2.5rem rgba(74, 63, 102, 0.2);
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
  }

  .tiles :global(.tile__ambient) {
    position: absolute;
    width: var(--ambient-size);
    height: var(--ambient-size);
    object-fit: contain;
    opacity: var(--ambient-opacity);
    animation: ambientFloat var(--ambient-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__ambient--camera) {
    top: var(--ambient-edge-offset);
    left: var(--ambient-edge-offset);
  }

  .tiles :global(.tile__ambient--film) {
    right: var(--ambient-edge-offset);
    bottom: var(--ambient-edge-offset);
    animation-delay: var(--ambient-phase-delay);
  }

  .tiles :global(.tile__landscapes) {
    position: relative;
    z-index: 1;
    width: min(92%, 18rem);
    aspect-ratio: 4 / 3;
    overflow: hidden;
    border-radius: var(--landscape-radius);
    box-shadow: var(--landscape-shadow);
  }

  .tiles :global(.tile__landscape) {
    position: absolute;
    inset: 0;
    overflow: hidden;
    opacity: 0;
    animation: landscapeSequence var(--landscape-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__landscape--alpine) {
    background: linear-gradient(180deg, #83bfd1 0%, #f3cfaa 62%, #486b61 63%, #294a43 100%);
  }

  .tiles :global(.tile__landscape--coast) {
    background: linear-gradient(180deg, #533c70 0%, #d17a86 54%, #345f76 55%, #173d50 100%);
    animation-delay: calc(var(--landscape-duration) / -3);
  }

  .tiles :global(.tile__landscape--meadow) {
    background: linear-gradient(180deg, #8fc9d5 0%, #d8ead9 62%, #75a767 63%, #487b50 100%);
    animation-delay: calc(var(--landscape-duration) / -1.5);
  }

  .tiles :global(.tile__orb) {
    position: absolute;
    width: 16%;
    aspect-ratio: 1;
    border-radius: 50%;
  }

  .tiles :global(.tile__orb--sun) {
    top: 16%;
    right: 18%;
    background: #fff0b5;
    box-shadow: 0 0 2rem rgba(255, 240, 181, 0.8);
    animation: sunDrift var(--landscape-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__orb--moon) {
    top: 14%;
    left: 18%;
    background: #f8e8cf;
    box-shadow: 0 0 1.5rem rgba(248, 232, 207, 0.55);
    animation: moonDrift var(--landscape-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__ridge) {
    position: absolute;
    bottom: -24%;
    aspect-ratio: 1;
    transform: rotate(45deg);
    animation: ridgeDrift var(--landscape-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__ridge--far) {
    left: -2%;
    width: 72%;
    background: #67847b;
  }

  .tiles :global(.tile__ridge--near) {
    right: -12%;
    width: 82%;
    background: #365b52;
    animation-direction: reverse;
  }

  .tiles :global(.tile__wave) {
    position: absolute;
    right: -10%;
    left: -10%;
    height: 32%;
    border-radius: 50%;
    animation: waveDrift var(--landscape-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__wave--far) {
    bottom: 2%;
    background: #3f7285;
  }

  .tiles :global(.tile__wave--near) {
    bottom: -14%;
    background: #214d62;
    animation-direction: reverse;
  }

  .tiles :global(.tile__hill) {
    position: absolute;
    width: 90%;
    aspect-ratio: 2 / 1;
    border-radius: 50%;
    animation: hillDrift var(--landscape-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__hill--far) {
    bottom: -8%;
    left: -20%;
    background: #66945c;
  }

  .tiles :global(.tile__hill--near) {
    right: -24%;
    bottom: -18%;
    background: #477c50;
    animation-direction: reverse;
  }

  .tiles :global(.tile__cloud) {
    position: absolute;
    width: 22%;
    height: 7%;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.72);
    animation: cloudDrift var(--landscape-duration) linear infinite;
  }

  .tiles :global(.tile__cloud--one) {
    top: 22%;
    left: 16%;
  }

  .tiles :global(.tile__cloud--two) {
    top: 36%;
    right: 14%;
    width: 14%;
    animation-direction: reverse;
  }

  @keyframes landscapeSequence {
    0%,
    28% {
      opacity: 1;
      transform: scale(1);
    }
    34%,
    94% {
      opacity: 0;
      transform: scale(1.04);
    }
    100% {
      opacity: 1;
      transform: scale(1);
    }
  }

  @keyframes sunDrift {
    0%,
    100% {
      transform: translate(0, 0);
    }
    50% {
      transform: translate(-20%, 18%);
    }
  }

  @keyframes moonDrift {
    0%,
    100% {
      transform: translate(0, 0);
    }
    50% {
      transform: translate(24%, -12%);
    }
  }

  @keyframes ridgeDrift {
    0%,
    100% {
      translate: 0 0;
    }
    50% {
      translate: 4% -2%;
    }
  }

  @keyframes waveDrift {
    0%,
    100% {
      translate: -3% 0;
    }
    50% {
      translate: 3% -8%;
    }
  }

  @keyframes hillDrift {
    0%,
    100% {
      translate: -2% 0;
    }
    50% {
      translate: 2% -4%;
    }
  }

  @keyframes cloudDrift {
    0%,
    100% {
      translate: -18% 0;
    }
    50% {
      translate: 18% 0;
    }
  }

  @keyframes ambientFloat {
    0%,
    100% {
      transform: translateY(0) rotate(-2deg);
    }
    50% {
      transform: translateY(calc(var(--ambient-travel) * -1)) rotate(2deg);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .tiles :global(.tile__ambient),
    .tiles :global(.tile__landscape),
    .tiles :global(.tile__landscape span) {
      animation-play-state: paused;
    }
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
