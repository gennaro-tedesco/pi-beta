<script lang="ts">
  import { onDestroy } from 'svelte'
  import { fade, scale } from 'svelte/transition'
  import { Rabbit, Turtle } from 'lucide-svelte'

  const transitionDuration = 1200
  const minSlideDuration = 2000
  const maxSlideDuration = 12000
  const swipeThreshold = 60

  const photos = Object.values(
    import.meta.glob('../../assets/images/gallery/*.jpg', { eager: true, import: 'default' })
  ) as string[]

  let slideDuration = 6000
  let activeIndex = 0
  let timer: ReturnType<typeof setInterval>
  let dragStartX: number | null = null

  function restartTimer(): void {
    clearInterval(timer)
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
    if (dragStartX === null) return
    const deltaX = getClientX(event) - dragStartX
    dragStartX = null
    if (deltaX <= -swipeThreshold) {
      activeIndex = (activeIndex + 1) % photos.length
    } else if (deltaX >= swipeThreshold) {
      activeIndex = (activeIndex - 1 + photos.length) % photos.length
    }
  }

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

  <div
    class="slideshow__controls"
    on:mousedown|stopPropagation
    on:mouseup|stopPropagation
    on:touchstart|stopPropagation
    on:touchend|stopPropagation
  >
    <Rabbit size={18} />
    <input
      type="range"
      min={minSlideDuration}
      max={maxSlideDuration}
      step="500"
      bind:value={slideDuration}
      aria-label="Slide transition speed"
    />
    <Turtle size={18} />
  </div>
</div>

<style>
  .slideshow {
    position: fixed;
    inset: 0;
    z-index: 0;
    overflow: hidden;
    touch-action: none;
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
    position: absolute;
    z-index: 2;
    top: 1.5rem;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border-radius: 999px;
    border: 1px solid rgba(74, 63, 102, 0.2);
    background-color: rgba(255, 255, 255, 0.5);
    color: #4a3f66;
    touch-action: auto;
  }
</style>
