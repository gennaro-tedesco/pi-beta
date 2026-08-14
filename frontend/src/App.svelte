<script lang="ts">
  import { onDestroy, onMount, tick } from 'svelte'
  import { fly, scale } from 'svelte/transition'
  import { CalendarDays, LocateFixed, MapPin, Moon, RotateCw } from 'lucide-svelte'
  import { Button, Tile } from './lib/components'
  import WeatherView from './lib/views/WeatherView.svelte'
  import CityPicker from './lib/views/CityPicker.svelte'
  import PhotoSlideshow from './lib/views/PhotoSlideshow.svelte'
  import type { WeatherScene } from './lib/views/weatherScene'
  import { transitionDuration } from './lib/transition'
  import cameraIllustration from './assets/images/camera-color.svg'
  import filmIllustration from './assets/images/film-frames-color.svg'
  import weatherTileAnimation from '@bybas/weather-icons/production/fill/all/partly-cloudy-day-rain.svg'

  type WidgetId = 'calendar' | 'photos' | 'weather'
  type Coords = { lat: number; lon: number }
  type CalendarDay = {
    day: number
    event: CalendarEvent | null
    gridColumnStart: number
    gridRow: number
    isToday: boolean
    isWeekend: boolean
    key: string
  }
  type CalendarWeek = {
    key: string
    number: number
    row: number
  }
  type CalendarView = {
    days: CalendarDay[]
    key: number
    weeks: CalendarWeek[]
  }
  type CalendarEvent = {
    day: number
    icons: string[]
    label: string
    month: number
  }

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
  const homeLogoSrc = '/images/logo.png'
  const calendarTitle = 'Calendar'
  const calendarWidgetId = 'calendar'
  const currentDateAriaValue = 'date'
  const chooseMonthLabel = 'Choose month'
  const todayButtonLabel = 'Go to today'
  const monthPickerSelector = '.calendar__month-picker'
  const calendarColumnCount = 7
  const calendarMonthCount = 12
  const previousMonthOffset = -1
  const nextMonthOffset = 1
  const nextDayOffset = 1
  const calendarSwipeThresholdPx = 48
  const calendarMonthTransitionDistancePx = 56
  const calendarMonthTransitionDurationMs = 300
  const calendarDropdownFadeOutDurationMs = 1000
  const calendarJumpTransitionDirection = 0
  const calendarDateRefreshIntervalMs = 60000
  const escapeKey = 'Escape'
  const firstDayOfMonth = 1
  const lastDayOfPreviousMonth = 0
  const firstGridRow = 1
  const calendarWeekColumnOffset = 2
  const datePartLength = 2
  const datePartFill = '0'
  const datePartSeparator = '-'
  const referenceMondayYear = 2024
  const referenceMondayMonth = 0
  const referenceMondayDate = 8
  const sundayToMondayOffset = 6
  const weekendStartIndex = 5
  const isoSunday = 7
  const isoThursday = 4
  const isoYearStartMonth = 0
  const isoYearStartDay = 1
  const millisecondsPerDay = 86400000
  const recurringCalendarEvents: CalendarEvent[] = [
    { month: 2, day: 19, icons: ['🐧', '🎂'], label: "Penguin's birthday" },
    { month: 3, day: 13, icons: ['🐟', '🎂'], label: "Fish's birthday" },
    { month: 12, day: 25, icons: ['🎄'], label: 'Christmas' },
    { month: 6, day: 22, icons: ['💍'], label: 'Wedding anniversary' }
  ]
  const monthFormatter = new Intl.DateTimeFormat(undefined, { month: 'long', year: 'numeric' })
  const monthNameFormatter = new Intl.DateTimeFormat(undefined, { month: 'long' })
  const weekdayFormatter = new Intl.DateTimeFormat(undefined, { weekday: 'short' })
  const monthNames = Array.from({ length: calendarMonthCount }, (_, monthIndex) =>
    monthNameFormatter.format(new Date(referenceMondayYear, monthIndex, firstDayOfMonth))
  )
  const weekdayNames = Array.from({ length: calendarColumnCount }, (_, dayIndex) =>
    weekdayFormatter.format(new Date(referenceMondayYear, referenceMondayMonth, referenceMondayDate + dayIndex))
  )

  let activeWidget: WidgetId | null = null
  let weatherScene: WeatherScene | null = null
  let weatherViewRef: WeatherView
  let selectedCoords: Coords | null = null
  let showCityPicker = false
  let currentDate = new Date()
  let calendarCursor = new Date(currentDate.getFullYear(), currentDate.getMonth(), firstDayOfMonth)
  let calendarSwipeStartX: number | null = null
  let calendarTransitionDirection = calendarJumpTransitionDirection
  let calendarTransitionOutDurationMs = calendarMonthTransitionDurationMs
  let showMonthPicker = false
  let calendarDateRefreshTimer: ReturnType<typeof setInterval>
  let monthPickerTrigger: HTMLButtonElement

  function isSameDay(left: Date, right: Date): boolean {
    return left.getFullYear() === right.getFullYear()
      && left.getMonth() === right.getMonth()
      && left.getDate() === right.getDate()
  }

  function formatDateKey(date: Date): string {
    const month = String(date.getMonth() + firstDayOfMonth).padStart(datePartLength, datePartFill)
    const day = String(date.getDate()).padStart(datePartLength, datePartFill)
    return [date.getFullYear(), month, day].join(datePartSeparator)
  }

  function getMondayBasedDayIndex(date: Date): number {
    return (date.getDay() + sundayToMondayOffset) % calendarColumnCount
  }

  function getDaysInMonth(cursor: Date): number {
    return new Date(
      cursor.getFullYear(),
      cursor.getMonth() + firstDayOfMonth,
      lastDayOfPreviousMonth
    ).getDate()
  }

  function getCalendarEvent(date: Date): CalendarEvent | null {
    const month = date.getMonth() + firstDayOfMonth
    return recurringCalendarEvents.find((event) => event.month === month && event.day === date.getDate()) ?? null
  }

  function getIsoWeekNumber(date: Date): number {
    const normalizedDate = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
    const isoWeekday = normalizedDate.getUTCDay() || isoSunday
    normalizedDate.setUTCDate(normalizedDate.getUTCDate() + isoThursday - isoWeekday)
    const yearStart = new Date(Date.UTC(normalizedDate.getUTCFullYear(), isoYearStartMonth, isoYearStartDay))
    const elapsedDays = (normalizedDate.getTime() - yearStart.getTime()) / millisecondsPerDay
    return Math.ceil((elapsedDays + firstDayOfMonth) / calendarColumnCount)
  }

  function getCalendarWeeks(cursor: Date): CalendarWeek[] {
    const firstDayOffset = getMondayBasedDayIndex(cursor)
    const firstVisibleMonday = new Date(
      cursor.getFullYear(),
      cursor.getMonth(),
      firstDayOfMonth - firstDayOffset
    )
    const weekCount = Math.ceil((firstDayOffset + getDaysInMonth(cursor)) / calendarColumnCount)

    return Array.from({ length: weekCount }, (_, weekIndex) => {
      const monday = new Date(
        firstVisibleMonday.getFullYear(),
        firstVisibleMonday.getMonth(),
        firstVisibleMonday.getDate() + weekIndex * calendarColumnCount
      )
      return {
        key: formatDateKey(monday),
        number: getIsoWeekNumber(monday),
        row: weekIndex + firstGridRow
      }
    })
  }

  function getCalendarDays(cursor: Date, referenceDate: Date): CalendarDay[] {
    const firstDayOffset = getMondayBasedDayIndex(cursor)
    const daysInMonth = getDaysInMonth(cursor)

    return Array.from({ length: daysInMonth }, (_, dayIndex) => {
      const day = dayIndex + firstDayOfMonth
      const date = new Date(cursor.getFullYear(), cursor.getMonth(), day)
      return {
        day,
        event: getCalendarEvent(date),
        gridColumnStart: getMondayBasedDayIndex(date) + calendarWeekColumnOffset,
        gridRow: Math.floor((firstDayOffset + dayIndex) / calendarColumnCount) + firstGridRow,
        isToday: isSameDay(date, referenceDate),
        isWeekend: getMondayBasedDayIndex(date) >= weekendStartIndex,
        key: formatDateKey(date)
      }
    })
  }

  function shiftCalendarMonth(offset: number): void {
    calendarTransitionDirection = offset
    calendarTransitionOutDurationMs = calendarMonthTransitionDurationMs
    calendarCursor = new Date(calendarCursor.getFullYear(), calendarCursor.getMonth() + offset, firstDayOfMonth)
  }

  function handleCalendarPointerDown(event: PointerEvent): void {
    if (event.target instanceof Element && event.target.closest(monthPickerSelector)) return

    calendarSwipeStartX = event.clientX
    const calendar = event.currentTarget as HTMLElement
    calendar.setPointerCapture(event.pointerId)
  }

  function handleCalendarPointerUp(event: PointerEvent): void {
    if (calendarSwipeStartX === null) return

    const horizontalDistance = event.clientX - calendarSwipeStartX
    calendarSwipeStartX = null

    if (horizontalDistance >= calendarSwipeThresholdPx) shiftCalendarMonth(previousMonthOffset)
    if (horizontalDistance <= -calendarSwipeThresholdPx) shiftCalendarMonth(nextMonthOffset)
  }

  function cancelCalendarSwipe(): void {
    calendarSwipeStartX = null
  }

  function toggleMonthPicker(): void {
    showMonthPicker = !showMonthPicker
  }

  function selectCalendarMonth(monthIndex: number): void {
    calendarTransitionDirection = calendarJumpTransitionDirection
    calendarTransitionOutDurationMs = calendarDropdownFadeOutDurationMs
    calendarCursor = new Date(calendarCursor.getFullYear(), monthIndex, firstDayOfMonth)
    showMonthPicker = false
  }

  function goToToday(): void {
    calendarTransitionDirection = calendarJumpTransitionDirection
    calendarTransitionOutDurationMs = calendarMonthTransitionDurationMs
    calendarCursor = new Date(currentDate.getFullYear(), currentDate.getMonth(), firstDayOfMonth)
    showMonthPicker = false
  }

  function handleWindowPointerDown(event: PointerEvent): void {
    if (event.target instanceof Element && event.target.closest(monthPickerSelector)) return
    showMonthPicker = false
  }

  function handleWindowKeydown(event: KeyboardEvent): void {
    if (event.key !== escapeKey || !showMonthPicker) return
    showMonthPicker = false
    monthPickerTrigger.focus()
  }

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

  onMount(() => {
    calendarDateRefreshTimer = setInterval(() => {
      currentDate = new Date()
    }, calendarDateRefreshIntervalMs)
  })

  onDestroy(() => clearInterval(calendarDateRefreshTimer))

  $: if (activeWidget !== 'weather') weatherScene = null
  $: tomorrow = new Date(currentDate.getFullYear(), currentDate.getMonth(), currentDate.getDate() + nextDayOffset)
  $: calendarView = {
    days: getCalendarDays(calendarCursor, currentDate),
    key: calendarCursor.getTime(),
    weeks: getCalendarWeeks(calendarCursor)
  } as CalendarView
</script>

<svelte:window on:pointerdown={handleWindowPointerDown} on:keydown={handleWindowKeydown} />

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
        <img class="nav__logo" src={homeLogoSrc} alt="" />
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
      {:else if activeWidget === calendarWidgetId}
        <div class="nav__actions">
          <Button ghost on:click={goToToday} aria-label={todayButtonLabel}>
            <CalendarDays size={weatherActionIconSize} />
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
      </div>
    {:else if activeWidget === 'weather' && showCityPicker}
      <CityPicker on:select={handleCitySelect} on:close={closeCityPicker} />
    {:else if activeWidget === 'weather'}
      <WeatherView bind:this={weatherViewRef} coords={selectedCoords} on:scene={handleWeatherScene} />
    {:else if activeWidget === 'photos'}
      <PhotoSlideshow />
    {:else if activeWidget === calendarWidgetId}
      <section
        class="calendar"
        aria-label={calendarTitle}
        on:pointerdown={handleCalendarPointerDown}
        on:pointerup={handleCalendarPointerUp}
        on:pointercancel={cancelCalendarSwipe}
        out:scale={{ duration: transitionDuration }}
        in:scale={{ duration: transitionDuration, delay: transitionDuration }}
      >
        <header class="calendar__header">
          <h1>
            <span class="calendar__month-picker">
              <button
                class="calendar__month-trigger"
                type="button"
                bind:this={monthPickerTrigger}
                aria-label={chooseMonthLabel}
                aria-haspopup="true"
                aria-expanded={showMonthPicker}
                on:click={toggleMonthPicker}
              >{monthNameFormatter.format(calendarCursor)}</button>
              {#if showMonthPicker}
                <span class="calendar__month-menu" role="group" aria-label={chooseMonthLabel}>
                  {#each monthNames as monthName, monthIndex}
                    <button
                      class="calendar__month-option"
                      type="button"
                      aria-pressed={calendarCursor.getMonth() === monthIndex}
                      on:click={() => selectCalendarMonth(monthIndex)}
                    >{monthName}</button>
                  {/each}
                </span>
              {/if}
            </span>
            <span>{calendarCursor.getFullYear()}</span>
          </h1>
        </header>
        <div class="calendar__viewport">
          {#each [calendarView] as visibleMonth (visibleMonth.key)}
            <div
              class="calendar__month"
              in:fly={{
                x: calendarTransitionDirection * calendarMonthTransitionDistancePx,
                duration: calendarMonthTransitionDurationMs
              }}
              out:fly={{
                x: calendarTransitionDirection * -calendarMonthTransitionDistancePx,
                duration: calendarTransitionOutDurationMs
              }}
            >
              <div class="calendar__weekdays" aria-hidden="true">
                {#each weekdayNames as weekday}
                  <span>{weekday}</span>
                {/each}
              </div>
              <div class="calendar__grid">
                {#each visibleMonth.weeks as calendarWeek (calendarWeek.key)}
                  <span class="calendar__week-number" style:grid-row={calendarWeek.row}>{calendarWeek.number}</span>
                {/each}
                {#each visibleMonth.days as calendarDay (calendarDay.key)}
                  <time
                    class="calendar__day"
                    class:calendar__day--today={calendarDay.isToday}
                    class:calendar__day--weekend={calendarDay.isWeekend}
                    datetime={calendarDay.key}
                    aria-current={calendarDay.isToday ? currentDateAriaValue : undefined}
                    title={calendarDay.event?.label}
                    style:grid-column-start={calendarDay.gridColumnStart}
                    style:grid-row={calendarDay.gridRow}
                  >
                    <span class="calendar__day-number">{calendarDay.day}</span>
                    {#if calendarDay.event}
                      <span class="calendar__event" role="img" aria-label={calendarDay.event.label}>
                        {#each calendarDay.event.icons as icon}
                          <span>{icon}</span>
                        {/each}
                      </span>
                    {/if}
                  </time>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      </section>
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

  .nav > :global(button) {
    --nav-home-button-size: 3rem;
    width: var(--nav-home-button-size);
    height: var(--nav-home-button-size);
    padding: 0;
  }

  .nav__logo {
    --nav-logo-fill: 90%;
    width: var(--nav-logo-fill);
    height: var(--nav-logo-fill);
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
    min-height: 0;
  }

  .tiles {
    --tile-height: 50vh;
    --tile-width: clamp(13rem, 27vw, 24rem);
    --tile-column-count: 3;
    --tiles-width: 100%;
    display: grid;
    grid-template-columns: repeat(var(--tile-column-count), var(--tile-width));
    justify-content: center;
    width: var(--tiles-width);
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
    .tiles :global(.tile__landscape span),
    .tiles :global(.tile__calendar-page--turning) {
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

  .tiles :global(.tile__calendar) {
    --calendar-preview-width: min(82%, 15rem);
    --calendar-preview-radius: 1rem;
    --calendar-preview-padding: 1rem;
    --calendar-preview-gap: 0.45rem;
    --calendar-preview-background: rgba(255, 255, 255, 0.78);
    --calendar-preview-shadow: 0 1rem 2.5rem rgba(74, 63, 102, 0.18);
    --calendar-preview-month-size: 0.8rem;
    --calendar-preview-day-size: clamp(3.5rem, 8vw, 5.5rem);
    --calendar-preview-weekday-size: 0.9rem;
    --calendar-preview-aspect-ratio: 1;
    --calendar-preview-scale: 0.8;
    --calendar-page-fold-duration: 4s;
    --calendar-page-fold-angle: 180deg;
    --calendar-page-fold-scale: 1;
    --calendar-preview-perspective: 120rem;
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
    --calendar-binding-offset: -0.65rem;
    --calendar-binding-horizontal-position: 50%;
    --calendar-binding-horizontal-offset: -50%;
    --calendar-binding-size: 0.75rem;
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

  .calendar {
    --calendar-width: 100%;
    --calendar-height: 100%;
    --calendar-gap: 0.75rem;
    --calendar-padding: 1.25rem;
    --calendar-radius: 1rem;
    --calendar-background: transparent;
    --calendar-border: none;
    --calendar-day-radius: 0.75rem;
    --calendar-day-background: transparent;
    --calendar-day-border: 1px solid rgba(74, 63, 102, 0.16);
    --calendar-weekend-background: rgba(217, 236, 249, 0.58);
    --calendar-title-size: clamp(1.25rem, 3vw, 2rem);
    --calendar-weekday-size: 0.8rem;
    --calendar-menu-width: 10rem;
    --calendar-menu-background: rgba(255, 255, 255, 0.92);
    --calendar-menu-shadow: 0 0.75rem 2rem rgba(74, 63, 102, 0.16);
    --calendar-menu-padding: 0.5rem;
    --calendar-menu-radius: 0.75rem;
    --calendar-menu-font-size: 0.9rem;
    --calendar-menu-horizontal-position: 50%;
    --calendar-menu-horizontal-offset: -50%;
    --calendar-menu-layer: 3;
    --calendar-today-background: #ffe1ec;
    --calendar-today-color: #4a3f66;
    --calendar-week-width: 1.25rem;
    --calendar-week-number-size: 0.7rem;
    --calendar-grid-columns: var(--calendar-week-width) repeat(7, minmax(0, 1fr));
    --calendar-week-column: 1;
    --calendar-first-day-column: 2;
    --calendar-event-gap: 0.25rem;
    --calendar-event-size: clamp(1.5rem, 3vw, 2.5rem);
    --calendar-zero: 0;
    --calendar-semibold-weight: 600;
    --calendar-day-weight: 500;
    --calendar-today-weight: 700;
    --calendar-muted-opacity: 0.72;
    display: flex;
    flex-direction: column;
    gap: var(--calendar-gap);
    width: var(--calendar-width);
    height: var(--calendar-height);
    box-sizing: border-box;
    padding: var(--calendar-padding);
    border: var(--calendar-border);
    border-radius: var(--calendar-radius);
    background: var(--calendar-background);
    touch-action: pan-y;
    user-select: none;
  }

  .calendar__header {
    position: relative;
    z-index: var(--calendar-menu-layer);
    text-align: center;
  }

  .calendar__header h1 {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--calendar-gap);
    margin: var(--calendar-zero);
    font-size: var(--calendar-title-size);
  }

  .calendar__month-picker {
    position: relative;
    display: inline-flex;
  }

  .calendar__month-trigger {
    padding: var(--calendar-zero);
    border: none;
    background: transparent;
    color: inherit;
    font: inherit;
    cursor: pointer;
  }

  .calendar__month-menu {
    position: absolute;
    top: calc(100% + var(--calendar-gap));
    left: var(--calendar-menu-horizontal-position);
    display: grid;
    width: var(--calendar-menu-width);
    box-sizing: border-box;
    gap: var(--calendar-menu-padding);
    padding: var(--calendar-menu-padding);
    border: var(--calendar-day-border);
    border-radius: var(--calendar-menu-radius);
    background: var(--calendar-menu-background);
    box-shadow: var(--calendar-menu-shadow);
    transform: translateX(var(--calendar-menu-horizontal-offset));
  }

  .calendar__month-option {
    padding: var(--calendar-menu-padding);
    border: none;
    border-radius: var(--calendar-menu-radius);
    background: transparent;
    color: inherit;
    font: inherit;
    font-size: var(--calendar-menu-font-size);
    cursor: pointer;
  }

  .calendar__month-option[aria-pressed='true'] {
    background: var(--calendar-today-background);
  }

  .calendar__viewport {
    position: relative;
    flex: 1;
    min-height: var(--calendar-zero);
    overflow: hidden;
  }

  .calendar__month {
    position: absolute;
    inset: var(--calendar-zero);
    display: flex;
    flex-direction: column;
    gap: var(--calendar-gap);
  }

  .calendar__weekdays,
  .calendar__grid {
    display: grid;
    grid-template-columns: var(--calendar-grid-columns);
    gap: var(--calendar-gap);
  }

  .calendar__grid {
    flex: 1;
    min-height: var(--calendar-zero);
    grid-auto-rows: minmax(var(--calendar-zero), 1fr);
  }

  .calendar__weekdays {
    font-size: var(--calendar-weekday-size);
    font-weight: var(--calendar-semibold-weight);
    opacity: var(--calendar-muted-opacity);
  }

  .calendar__weekdays > span:first-child {
    grid-column: var(--calendar-first-day-column);
  }

  .calendar__week-number {
    display: flex;
    align-items: center;
    justify-content: center;
    grid-column: var(--calendar-week-column);
    font-size: var(--calendar-week-number-size);
    font-weight: var(--calendar-semibold-weight);
    opacity: var(--calendar-muted-opacity);
  }

  .calendar__day {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: var(--calendar-event-gap);
    border: var(--calendar-day-border);
    border-radius: var(--calendar-day-radius);
    background: var(--calendar-day-background);
    font-weight: var(--calendar-day-weight);
  }

  .calendar__event {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--calendar-event-gap);
    font-size: var(--calendar-event-size);
    line-height: 1;
  }

  .calendar__day--today {
    color: var(--calendar-today-color);
    background: var(--calendar-today-background);
    font-weight: var(--calendar-today-weight);
  }

  .calendar__day--weekend:not(.calendar__day--today) {
    background: var(--calendar-weekend-background);
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
