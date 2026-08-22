<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { fade } from 'svelte/transition'
  import { Wifi } from 'lucide-svelte'
  import { Tile } from './lib/components'
  import WeatherView from './lib/views/WeatherView.svelte'
  import PhotoSlideshow from './lib/views/PhotoSlideshow.svelte'
  import NetworkView from './lib/views/NetworkView.svelte'
  import CalendarView from './lib/views/CalendarView.svelte'
  import { transitionDuration } from './lib/transition'
  import cameraIllustration from './assets/images/camera-color.svg'
  import filmIllustration from './assets/images/film-frames-color.svg'
  import weatherTileAnimation from '@bybas/weather-icons/production/fill/all/partly-cloudy-day-rain.svg'

  type WidgetId = 'calendar' | 'network' | 'photos' | 'weather'

  const calendarTitle = 'Calendar'
  const calendarWidgetId = 'calendar'
  const networkWidgetId = 'network'
  const networkTitle = 'Machine'
  const nextDayOffset = 1
  const tileSwipeThresholdPx = 48
  const previousTileOffset = -1
  const nextTileOffset = 1
  const tileSwipeClickResetDelayMs = 0
  const tileScrollBehavior: ScrollBehavior = 'smooth'
  const calendarDateRefreshIntervalMs = 60000
  const monthFormatter = new Intl.DateTimeFormat(undefined, { month: 'long', year: 'numeric' })
  const weekdayFormatter = new Intl.DateTimeFormat(undefined, { weekday: 'short' })

  let activeWidget: WidgetId | null = null
  let currentDate = new Date()
  let tileSwipeStartX: number | null = null
  let tileSwipeConsumed = false
  let calendarDateRefreshTimer: ReturnType<typeof setInterval>

  function handleTilesPointerDown(event: PointerEvent): void {
    tileSwipeStartX = event.clientX
    tileSwipeConsumed = false
  }

  function scrollTiles(tiles: HTMLElement, offset: number): void {
    const tileElements = Array.from(tiles.children) as HTMLElement[]
    const firstTile = tileElements[0]
    const secondTile = tileElements[1]
    if (!firstTile || !secondTile) return

    const tileStride = secondTile.offsetLeft - firstTile.offsetLeft
    tiles.scrollBy({ left: offset * tileStride, behavior: tileScrollBehavior })
  }

  function handleTilesPointerUp(event: PointerEvent): void {
    if (tileSwipeStartX === null) return

    const horizontalDistance = event.clientX - tileSwipeStartX
    tileSwipeStartX = null

    if (horizontalDistance >= tileSwipeThresholdPx) {
      tileSwipeConsumed = true
      scrollTiles(event.currentTarget as HTMLElement, previousTileOffset)
    } else if (horizontalDistance <= -tileSwipeThresholdPx) {
      tileSwipeConsumed = true
      scrollTiles(event.currentTarget as HTMLElement, nextTileOffset)
    }

    if (tileSwipeConsumed) window.setTimeout(() => tileSwipeConsumed = false, tileSwipeClickResetDelayMs)
  }

  function cancelTileSwipe(): void {
    tileSwipeStartX = null
    tileSwipeConsumed = false
  }

  function maximize(id: WidgetId): void {
    if (tileSwipeConsumed) {
      tileSwipeConsumed = false
      return
    }

    activeWidget = id
  }

  function goHome(): void {
    activeWidget = null
  }

  onMount(() => {
    calendarDateRefreshTimer = setInterval(() => {
      currentDate = new Date()
    }, calendarDateRefreshIntervalMs)
  })

  onDestroy(() => {
    clearInterval(calendarDateRefreshTimer)
  })

  $: tomorrow = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate() + nextDayOffset)
</script>

<main class="screen">
  <div class="content">
    {#if activeWidget === null}
      <div
        class="tiles"
        on:pointerdown={handleTilesPointerDown}
        on:pointerup={handleTilesPointerUp}
        on:pointercancel={cancelTileSwipe}
        out:fade={{ duration: transitionDuration }}
        in:fade={{ duration: transitionDuration, delay: transitionDuration }}
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
        <Tile title={calendarTitle} on:toggle={() => maximize(calendarWidgetId)}>
          <div slot="visual" class="tile__calendar" aria-hidden="true">
            <div class="tile__calendar-binding"><span /><span /></div>
            <div class="tile__calendar-page">
              <span class="tile__calendar-month">{monthFormatter.format(tomorrow)}</span>
              <strong class="tile__calendar-day">{tomorrow.getDate()}</strong>
              <span class="tile__calendar-weekday">{weekdayFormatter.format(tomorrow)}</span>
            </div>
            <div class="tile__calendar-page tile__calendar-page--turning">
              <span class="tile__calendar-month">{monthFormatter.format(currentDate)}</span>
              <strong class="tile__calendar-day">{currentDate.getDate()}</strong>
              <span class="tile__calendar-weekday">{weekdayFormatter.format(currentDate)}</span>
            </div>
          </div>
        </Tile>
        <Tile title={networkTitle} on:toggle={() => maximize(networkWidgetId)}>
          <div slot="visual" class="tile__network" aria-hidden="true">
            <div class="tile__network-icon">
              <Wifi />
            </div>
          </div>
        </Tile>
      </div>
    {:else if activeWidget === 'weather'}
      <WeatherView on:home={goHome} />
    {:else if activeWidget === 'photos'}
      <PhotoSlideshow on:home={goHome} />
    {:else if activeWidget === networkWidgetId}
      <NetworkView on:home={goHome} />
    {:else if activeWidget === calendarWidgetId}
      <CalendarView on:home={goHome} />
    {/if}
  </div>
</main>

<style>
  .screen {
    position: relative;
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    height: 100vh;
    box-sizing: border-box;
    padding: 1rem;
    gap: 1rem;
    background-image: var(--app-background);
  }

  .content {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 0;
  }

  .tiles {
    --tile-height: 50vh;
    --tile-width: clamp(13rem, 27vw, 24rem);
    --tile-column-count: 4;
    --tiles-width: 100%;
    --tiles-horizontal-overflow: auto;
    --tiles-justification: safe center;
    --tiles-overscroll-behavior: contain;
    --tiles-scroll-snap-type: x mandatory;
    --tiles-scrollbar-visibility: none;
    --tiles-touch-action: pan-y;
    display: grid;
    grid-template-columns: repeat(var(--tile-column-count), var(--tile-width));
    justify-content: var(--tiles-justification);
    width: var(--tiles-width);
    gap: 1rem;
    overflow-x: var(--tiles-horizontal-overflow);
    overscroll-behavior-x: var(--tiles-overscroll-behavior);
    scrollbar-width: var(--tiles-scrollbar-visibility);
    -ms-overflow-style: var(--tiles-scrollbar-visibility);
    scroll-snap-type: var(--tiles-scroll-snap-type);
    touch-action: var(--tiles-touch-action);
  }

  .tiles::-webkit-scrollbar {
    display: var(--tiles-scrollbar-visibility);
  }

  .tiles :global(.tile) {
    --tile-scroll-snap-alignment: start;
    --tile-scroll-snap-stop: always;
    flex: 0 0 auto;
    height: var(--tile-height);
    width: var(--tile-width);
    scroll-snap-align: var(--tile-scroll-snap-alignment);
    scroll-snap-stop: var(--tile-scroll-snap-stop);
  }

  .tiles :global(.tile__slideshow) {
    --ambient-edge-offset: 1%;
    --ambient-opacity: 0.72;
    --ambient-size: min(28cqmin, 7rem);
    --ambient-travel: min(2cqmin, 0.5rem);
    --ambient-phase-delay: -3s;
    --ambient-duration: 6s;
    --landscape-duration: 15s;
    --landscape-width: min(92%, 18rem, 100cqh);
    --landscape-radius: min(4cqmin, 1rem);
    --landscape-shadow-offset: min(4cqmin, 1rem);
    --landscape-shadow-blur: min(10cqmin, 2.5rem);
    --landscape-shadow: 0 var(--landscape-shadow-offset) var(--landscape-shadow-blur) rgba(74, 63, 102, 0.2);
    --sun-shadow-blur: min(8cqmin, 2rem);
    --moon-shadow-blur: min(6cqmin, 1.5rem);
    --round-shape-radius: 50%;
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
    width: var(--landscape-width);
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
    border-radius: var(--round-shape-radius);
  }

  .tiles :global(.tile__orb--sun) {
    top: 16%;
    right: 18%;
    background: #fff0b5;
    box-shadow: 0 0 var(--sun-shadow-blur) rgba(255, 240, 181, 0.8);
    animation: sunDrift var(--landscape-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__orb--moon) {
    top: 14%;
    left: 18%;
    background: #f8e8cf;
    box-shadow: 0 0 var(--moon-shadow-blur) rgba(248, 232, 207, 0.55);
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
    border-radius: var(--round-shape-radius);
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
    border-radius: 50%;
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
    .tiles :global(.tile__landscape span),
    .tiles :global(.tile__calendar-page--turning),
    .tiles :global(.tile__network-icon) {
      animation-play-state: paused;
    }
  }

  .tiles :global(.tile__weather-scene) {
    --weather-shadow-blur: min(7cqmin, 28px);
    --weather-float-distance: min(1.5cqmin, 6px);
    width: 80%;
    height: 80%;
    object-fit: contain;
    filter: drop-shadow(0 0 var(--weather-shadow-blur) rgba(255, 255, 255, 0.18));
    animation: floatIcon 6s ease-in-out infinite;
  }

  .tiles :global(.tile__network) {
    --network-tile-fill: 100%;
    --network-tile-icon-size: min(38cqmin, 9rem);
    --network-tile-icon-padding: min(8cqmin, 2rem);
    --network-tile-radius: 50%;
    --network-tile-background: rgba(215, 238, 229, 0.72);
    --network-tile-color: #3c8d72;
    --network-tile-broadcast-duration: 2.4s;
    --network-tile-broadcast-start-scale: 0.82;
    --network-tile-broadcast-end-scale: 1;
    --network-tile-broadcast-start-opacity: 0.32;
    --network-tile-broadcast-end-opacity: 1;
    --network-tile-broadcast-origin: center bottom;
    display: flex;
    align-items: center;
    justify-content: center;
    width: var(--network-tile-fill);
    height: var(--network-tile-fill);
  }

  .tiles :global(.tile__network-icon) {
    display: flex;
    align-items: center;
    justify-content: center;
    width: var(--network-tile-icon-size);
    height: var(--network-tile-icon-size);
    box-sizing: border-box;
    padding: var(--network-tile-icon-padding);
    border-radius: var(--network-tile-radius);
    background: var(--network-tile-background);
    color: var(--network-tile-color);
    animation: wifiBroadcast var(--network-tile-broadcast-duration) ease-out infinite;
    transform-origin: var(--network-tile-broadcast-origin);
  }

  .tiles :global(.tile__network-icon .lucide) {
    width: var(--network-tile-fill);
    height: var(--network-tile-fill);
  }

  @keyframes wifiBroadcast {
    from {
      opacity: var(--network-tile-broadcast-start-opacity);
      transform: scale(var(--network-tile-broadcast-start-scale));
    }
    to {
      opacity: var(--network-tile-broadcast-end-opacity);
      transform: scale(var(--network-tile-broadcast-end-scale));
    }
  }

  .tiles :global(.tile__calendar) {
    --calendar-preview-width: min(82%, 82cqh, 15rem);
    --calendar-preview-radius: min(4cqmin, 1rem);
    --calendar-preview-padding: min(4cqmin, 1rem);
    --calendar-preview-gap: min(2cqmin, 0.45rem);
    --calendar-preview-background: rgba(255, 255, 255, 0.78);
    --calendar-preview-shadow-offset: min(4cqmin, 1rem);
    --calendar-preview-shadow-blur: min(10cqmin, 2.5rem);
    --calendar-preview-shadow: 0 var(--calendar-preview-shadow-offset) var(--calendar-preview-shadow-blur) rgba(74, 63, 102, 0.18);
    --calendar-preview-month-size: min(4cqmin, 0.8rem);
    --calendar-preview-day-size: min(30cqmin, 5.5rem);
    --calendar-preview-weekday-size: min(4.5cqmin, 0.9rem);
    --calendar-preview-aspect-ratio: 1;
    --calendar-preview-scale: 0.8;
    --calendar-page-fold-duration: 4s;
    --calendar-page-fold-angle: 180deg;
    --calendar-page-fold-scale: 1;
    --calendar-preview-perspective: min(480cqmin, 120rem);
    --calendar-zero: 0;
    --calendar-preview-font-weight: 600;
    --calendar-preview-line-height: 1;
    --calendar-preview-muted-opacity: 0.72;
    position: relative;
    width: var(--calendar-preview-width);
    aspect-ratio: var(--calendar-preview-aspect-ratio);
    box-sizing: border-box;
    padding: var(--calendar-preview-padding);
    border-radius: var(--calendar-preview-radius);
    background: var(--calendar-preview-background);
    box-shadow: var(--calendar-preview-shadow);
    perspective: var(--calendar-preview-perspective);
    transform: scale(var(--calendar-preview-scale));
  }

  .tiles :global(.tile__calendar-binding) {
    --calendar-binding-width: 62%;
    --calendar-binding-offset: max(-3cqmin, -0.65rem);
    --calendar-binding-horizontal-position: 50%;
    --calendar-binding-horizontal-offset: -50%;
    --calendar-binding-size: min(4cqmin, 0.75rem);
    --calendar-binding-color: #d17a86;
    position: absolute;
    top: var(--calendar-binding-offset);
    left: var(--calendar-binding-horizontal-position);
    display: flex;
    justify-content: space-between;
    width: var(--calendar-binding-width);
    z-index: 2;
    transform: translateX(var(--calendar-binding-horizontal-offset));
  }

  .tiles :global(.tile__calendar-binding span) {
    width: var(--calendar-binding-size);
    height: var(--calendar-binding-size);
    border-radius: var(--calendar-binding-size);
    background: var(--calendar-binding-color);
  }

  .tiles :global(.tile__calendar-page) {
    position: absolute;
    inset: var(--calendar-zero);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--calendar-preview-gap);
    padding: var(--calendar-preview-padding);
    border-radius: var(--calendar-preview-radius);
    box-sizing: border-box;
    background: var(--calendar-preview-background);
    backface-visibility: visible;
  }

  .tiles :global(.tile__calendar-page--turning) {
    transform-origin: top center;
    animation: calendarPageFold var(--calendar-page-fold-duration) ease-in-out infinite;
  }

  .tiles :global(.tile__calendar-month) {
    font-size: var(--calendar-preview-month-size);
    font-weight: var(--calendar-preview-font-weight);
  }

  .tiles :global(.tile__calendar-day) {
    font-size: var(--calendar-preview-day-size);
    line-height: var(--calendar-preview-line-height);
  }

  .tiles :global(.tile__calendar-weekday) {
    font-size: var(--calendar-preview-weekday-size);
    opacity: var(--calendar-preview-muted-opacity);
  }

  @keyframes calendarPageFold {
    0%,
    20% {
      transform: rotateX(0) scaleY(1);
      opacity: 1;
    }
    65% {
      transform: rotateX(var(--calendar-page-fold-angle)) scaleY(var(--calendar-page-fold-scale));
      opacity: 1;
    }
    80% {
      transform: rotateX(var(--calendar-page-fold-angle)) scaleY(var(--calendar-page-fold-scale));
      opacity: 0;
    }
    81% {
      transform: rotateX(0) scaleY(1);
      opacity: 0;
    }
    100% {
      transform: rotateX(0) scaleY(1);
      opacity: 1;
    }
  }

  @keyframes floatIcon {
    0%,
    100% {
      transform: translateY(0);
    }
    50% {
      transform: translateY(calc(var(--weather-float-distance) * -1));
    }
  }

</style>
