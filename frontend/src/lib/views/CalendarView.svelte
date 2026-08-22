<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { fade, fly } from 'svelte/transition'
  import { CalendarSync } from 'lucide-svelte'
  import { Button, WidgetNavigation } from '../components'
  import { transitionDuration } from '../transition'

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

  const calendarTitle = 'Calendar'
  const calendarActionIconSize = 36
  const currentDateAriaValue = 'date'
  const chooseMonthLabel = 'Choose month'
  const todayButtonLabel = 'Go to today'
  const monthPickerSelector = '.calendar__month-picker'
  const calendarContainerType = 'size'
  const calendarColumnCount = 7
  const calendarMonthCount = 12
  const previousMonthOffset = -1
  const nextMonthOffset = 1
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
  const monthNameFormatter = new Intl.DateTimeFormat(undefined, { month: 'long' })
  const weekdayFormatter = new Intl.DateTimeFormat(undefined, { weekday: 'short' })
  const monthNames = Array.from({ length: calendarMonthCount }, (_, monthIndex) =>
    monthNameFormatter.format(new Date(referenceMondayYear, monthIndex, firstDayOfMonth))
  )
  const weekdayNames = Array.from({ length: calendarColumnCount }, (_, dayIndex) =>
    weekdayFormatter.format(new Date(referenceMondayYear, referenceMondayMonth, referenceMondayDate + dayIndex))
  )

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

  onMount(() => {
    calendarDateRefreshTimer = setInterval(() => {
      currentDate = new Date()
    }, calendarDateRefreshIntervalMs)
  })

  onDestroy(() => clearInterval(calendarDateRefreshTimer))

  $: calendarView = {
    days: getCalendarDays(calendarCursor, currentDate),
    key: calendarCursor.getTime(),
    weeks: getCalendarWeeks(calendarCursor)
  } as CalendarView
</script>

<svelte:window on:pointerdown={handleWindowPointerDown} on:keydown={handleWindowKeydown} />

<div class="calendar-widget">
  <WidgetNavigation on:home>
    <svelte:fragment slot="actions">
      <Button ghost on:click={goToToday} aria-label={todayButtonLabel}>
        <CalendarSync size={calendarActionIconSize} />
      </Button>
    </svelte:fragment>
  </WidgetNavigation>
  <div class="calendar-widget__content">
    <section
      class="calendar"
      style:container-type={calendarContainerType}
      aria-label={calendarTitle}
      on:pointerdown={handleCalendarPointerDown}
      on:pointerup={handleCalendarPointerUp}
      on:pointercancel={cancelCalendarSwipe}
      out:fade={{ duration: transitionDuration }}
      in:fade={{ duration: transitionDuration, delay: transitionDuration }}
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
                style:container-type={calendarContainerType}
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
  </div>
</div>

<style>
  .calendar-widget {
    --calendar-widget-edge: 0;
    --calendar-widget-fill: 100%;
    --calendar-widget-minimum: 0;
    --calendar-widget-gap: 1rem;
    --calendar-widget-padding: 1rem;
    position: absolute;
    inset: var(--calendar-widget-edge);
    display: flex;
    flex-direction: column;
    width: var(--calendar-widget-fill);
    height: var(--calendar-widget-fill);
    min-height: var(--calendar-widget-minimum);
    box-sizing: border-box;
    padding: var(--calendar-widget-padding);
    gap: var(--calendar-widget-gap);
  }

  .calendar-widget__content {
    display: flex;
    flex: 1;
    align-items: center;
    justify-content: center;
    min-height: var(--calendar-widget-minimum);
  }

  .calendar {
    --calendar-width: 100%;
    --calendar-height: 100%;
    --calendar-gap: min(0.75rem, 2.5cqh);
    --calendar-header-spacing: min(1rem, 3cqh);
    --calendar-padding: min(1.25rem, 4cqh);
    --calendar-radius: 1rem;
    --calendar-background: transparent;
    --calendar-border: none;
    --calendar-day-radius: min(0.75rem, 6cqh);
    --calendar-day-background: transparent;
    --calendar-day-border: 1px solid rgba(74, 63, 102, 0.16);
    --calendar-weekend-background: rgba(217, 236, 249, 0.58);
    --calendar-title-size: min(clamp(1.25rem, 3vw, 2rem), 8cqh);
    --calendar-weekday-size: min(0.8rem, 4cqh);
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
    --calendar-week-number-size: min(0.7rem, 4cqh);
    --calendar-grid-columns: var(--calendar-week-width) repeat(7, minmax(0, 1fr));
    --calendar-week-column: 1;
    --calendar-first-day-column: 2;
    --calendar-event-gap: min(0.25rem, 5cqw);
    --calendar-event-size: clamp(1.5rem, 3vw, 2.5rem);
    --calendar-day-size: 1rem;
    --calendar-row-min-height: min(2rem, 10cqh);
    --calendar-day-height-size: 45cqh;
    --calendar-event-width-size: 40cqw;
    --calendar-event-height-size: 35cqh;
    --calendar-event-height-gap: 5cqh;
    --calendar-day-line-height: 1;
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
    min-height: var(--calendar-zero);
    box-sizing: border-box;
    padding-inline: var(--calendar-padding);
    padding-block-start: var(--calendar-zero);
    padding-block-end: var(--calendar-padding);
    border: var(--calendar-border);
    border-radius: var(--calendar-radius);
    background: var(--calendar-background);
    touch-action: pan-y;
    user-select: none;
  }

  .calendar__header {
    position: relative;
    z-index: var(--calendar-menu-layer);
    margin-block-end: var(--calendar-header-spacing);
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
    overflow-x: hidden;
    overflow-y: auto;
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
    grid-auto-rows: minmax(var(--calendar-row-min-height), 1fr);
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
    gap: var(--calendar-zero);
    border: var(--calendar-day-border);
    border-radius: var(--calendar-day-radius);
    background: var(--calendar-day-background);
    font-weight: var(--calendar-day-weight);
  }

  .calendar__day-number {
    font-size: min(var(--calendar-day-size), var(--calendar-day-height-size));
    line-height: var(--calendar-day-line-height);
  }

  .calendar__event {
    display: flex;
    align-items: center;
    justify-content: center;
    width: var(--calendar-width);
    box-sizing: border-box;
    gap: var(--calendar-event-gap);
    margin-block: min(var(--calendar-event-gap), var(--calendar-event-height-gap));
    font-size: min(var(--calendar-event-size), var(--calendar-event-width-size), var(--calendar-event-height-size));
    line-height: var(--calendar-day-line-height);
  }

  .calendar__day--today {
    color: var(--calendar-today-color);
    background: var(--calendar-today-background);
    font-weight: var(--calendar-today-weight);
  }

  .calendar__day--weekend:not(.calendar__day--today) {
    background: var(--calendar-weekend-background);
  }
</style>
