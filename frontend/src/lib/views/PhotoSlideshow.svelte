<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { fade } from 'svelte/transition'
  import { FolderOpen, Rabbit, Turtle } from 'lucide-svelte'
  import { Button, Message, WidgetNavigation } from '../components'
  import { transitionDuration } from '../transition'
  import { ChoosePhotosFolder, ListPhotos } from '../../../wailsjs/go/main/App'

  const minSlideDuration = 2000
  const maxSlideDuration = 12000
  const swipeThreshold = 60
  const speedIconSize = 12
  const pickerIconSize = 18
  const chooseFolderLabel = 'Choose photos folder'
  const changePhotosFolderLabel = 'Change photos folder'
  const photosRoutePrefix = '/photos/'

  let photos: string[] = []
  let slideDuration = 6000
  let activeIndex = 0
  let timer: ReturnType<typeof setInterval>
  let dragStartX: number | null = null
  let error: string | null = null

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

  async function loadPhotos(): Promise<void> {
    error = null
    try {
      const names = await ListPhotos()
      photos = names.map((name) => `${photosRoutePrefix}${name}`)
      activeIndex = 0
    } catch (err) {
      error = String(err)
    }
  }

  export async function chooseFolder(): Promise<void> {
    error = null
    try {
      const selected = await ChoosePhotosFolder()
      if (!selected) return
      await loadPhotos()
    } catch (err) {
      error = String(err)
    }
  }

  onMount(loadPhotos)

  $: photos, slideDuration, restartTimer()

  onDestroy(() => clearInterval(timer))
</script>

<div
  class="slideshow"
  out:fade={{ duration: transitionDuration }}
  in:fade={{ duration: transitionDuration, delay: transitionDuration }}
  on:mousedown={handleDragStart}
  on:mouseup={handleDragEnd}
  on:touchstart={handleDragStart}
  on:touchend={handleDragEnd}
>
  <WidgetNavigation on:home>
    <svelte:fragment slot="actions">
      <Button ghost on:click={chooseFolder} aria-label={changePhotosFolderLabel}>
        <FolderOpen size={pickerIconSize} />
      </Button>
    </svelte:fragment>
  </WidgetNavigation>

  {#if error}
    <div class="slideshow__message">
      <Message variant="error" message={error} />
    </div>
  {:else if photos.length === 0}
    <div class="slideshow__picker">
      <Button on:click={chooseFolder}>
        <FolderOpen size={pickerIconSize} />
        <span>{chooseFolderLabel}</span>
      </Button>
    </div>
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
    --slideshow-padding: 1rem;
    position: fixed;
    inset: 0;
    z-index: 0;
    box-sizing: border-box;
    padding: var(--slideshow-padding);
    overflow: hidden;
    touch-action: none;
  }

  .slideshow__picker {
    --slideshow-picker-gap: 0.5rem;
    --slideshow-picker-font-size: clamp(0.9rem, 2vw, 1.1rem);
    --slideshow-picker-padding-x: 2rem;
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 var(--slideshow-picker-padding-x);
    color: #4a3f66;
  }

  .slideshow__picker :global(.button) {
    gap: var(--slideshow-picker-gap);
    font-size: var(--slideshow-picker-font-size);
  }

  .slideshow__message {
    --slideshow-message-edge: 0;
    --slideshow-message-layer: 3;
    --slideshow-message-padding: 2rem;
    position: absolute;
    inset: var(--slideshow-message-edge);
    z-index: var(--slideshow-message-layer);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--slideshow-message-padding);
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
