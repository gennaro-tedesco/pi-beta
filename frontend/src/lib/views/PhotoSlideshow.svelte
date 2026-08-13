<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { fade, scale } from 'svelte/transition'
  import { Rabbit, Turtle } from 'lucide-svelte'
  import { transitionDuration } from '../transition'
  import { ListPhotos } from '../../../wailsjs/go/main/App'

  const minSlideDuration = 2000
  const maxSlideDuration = 12000
  const swipeThreshold = 60
  const speedIconSize = 12

  let photos: string[] = []
  let slideDuration = 6000
  let activeIndex = 0
  let timer: ReturnType<typeof setInterval>
  let dragStartX: number | null = null

  function restartTimer(): void {
    clearInterval(timer)
    if (photos.length === 0) return
    timer = setInterval(() => {
      activeIndex = (activeIndex + 1) % photos.length
    }, slideDuration)
  }

  function getClientX(event: MouseEvent | TouchEvent): number {
    if ('touches' in event) {
      const touch = event.touches[0] ?? event.changedTouches[0]
      return touch.clientX
    }
    return event.clientX
  }

  function handleDragStart(event: MouseEvent | TouchEvent): void {
    dragStartX = getClientX(event)
  }

  function handleDragEnd(event: MouseEvent | TouchEvent): void {
    if (dragStartX === null || photos.length === 0) return
    const deltaX = getClientX(event) - dragStartX
    dragStartX = null
    if (deltaX <= -swipeThreshold) {
      activeIndex = (activeIndex + 1) % photos.length
    } else if (deltaX >= swipeThreshold) {
      activeIndex = (activeIndex - 1 + photos.length) % photos.length
    }
  }

  onMount(async () => {
    const names = await ListPhotos()
    photos = names.map((name) => `/images/${name}`)
  })

  $: slideDuration, restartTimer()

  onDestroy(() => clearInterval(timer))
</script>

<div
  class="slideshow"
  out:scale={{ duration: transitionDuration }}
  in:scale={{ duration: transitionDuration, delay: transitionDuration }}
  on:mousedown={handleDragStart}
  on:mouseup={handleDragEnd}
  on:touchstart={handleDragStart}
  on:touchend={handleDragEnd}
>
  {#if photos.length === 0}
    <p class="slideshow__empty">Add photos to the images folder to start the slideshow</p>
  {/if}

  {#each photos as photo, index (photo)}
    {#if index === activeIndex}
      <img
        class="slideshow__photo"
        src={photo}
        alt=""
        draggable="false"
        transition:fade={{ duration: slideDuration / 3 }}
      />
    {/if}
  {/each}

  {#if photos.length > 1}
    <div
      class="slideshow__controls"
      on:mousedown|stopPropagation
      on:mouseup|stopPropagation
      on:touchstart|stopPropagation
      on:touchend|stopPropagation
    >
      <Rabbit size={speedIconSize} />
      <input
        class="slideshow__speed"
        type="range"
        min={minSlideDuration}
        max={maxSlideDuration}
        step="500"
        bind:value={slideDuration}
        aria-label="Slide transition speed"
      />
      <Turtle size={speedIconSize} />
    </div>
  {/if}
</div>

<style>
  .slideshow {
    position: fixed;
    inset: 0;
    z-index: 0;
    overflow: hidden;
    touch-action: none;
  }

  .slideshow__empty {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0;
    padding: 0 2rem;
    text-align: center;
    color: #4a3f66;
  }

  .slideshow__photo {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    -webkit-user-drag: none;
    user-select: none;
  }

  .slideshow__controls {
    --slideshow-controls-gap: 0.375rem;
    --slideshow-controls-padding-y: 0.25rem;
    --slideshow-controls-padding-x: 0.625rem;
    --slideshow-controls-fade-color: rgba(255, 255, 255, 0.55);
    --slideshow-controls-fade-extent: 70%;
    position: absolute;
    z-index: 2;
    top: 1.5rem;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    gap: var(--slideshow-controls-gap);
    padding: var(--slideshow-controls-padding-y) var(--slideshow-controls-padding-x);
    border: none;
    background: radial-gradient(ellipse, var(--slideshow-controls-fade-color) 0%, transparent var(--slideshow-controls-fade-extent));
    color: #4a3f66;
    touch-action: auto;
  }

  .slideshow__speed {
    --slideshow-speed-width: 4rem;
    --slideshow-speed-track-height: 3px;
    --slideshow-speed-thumb-size: 10px;
    width: var(--slideshow-speed-width);
    height: var(--slideshow-speed-thumb-size);
    margin: 0;
    -webkit-appearance: none;
    appearance: none;
    background: transparent;
  }

  .slideshow__speed::-webkit-slider-runnable-track {
    height: var(--slideshow-speed-track-height);
    border-radius: 999px;
    background-color: rgba(74, 63, 102, 0.35);
  }

  .slideshow__speed::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: var(--slideshow-speed-thumb-size);
    height: var(--slideshow-speed-thumb-size);
    border-radius: 50%;
    background-color: #4a3f66;
    margin-top: calc((var(--slideshow-speed-track-height) - var(--slideshow-speed-thumb-size)) / 2);
    cursor: pointer;
  }
</style>
