<script lang="ts">
  import { createEventDispatcher } from 'svelte'

  export let title: string

  const tileContainerType = 'size'
  const dispatch = createEventDispatcher<{ toggle: void }>()

  function handleActivate(): void {
    dispatch('toggle')
  }

  function handleKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      handleActivate()
    }
  }
</script>

<div class="tile" role="button" tabindex="0" aria-label={title} on:click={handleActivate} on:keydown={handleKeydown} style:container-type={tileContainerType}>
  <div class="tile__content">
    <div class="tile__visual" style:container-type={tileContainerType}>
      <slot name="visual" />
    </div>
  </div>
</div>

<style>
  .tile {
    display: flex;
    flex-direction: column;
    cursor: pointer;
    user-select: none;
  }

  .tile__content {
    --tile-spacing: 1rem;
    --tile-relative-spacing: 4cqmin;
    --tile-radius: 10px;
    --tile-relative-radius: 2.5cqmin;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
    min-height: 0;
    box-sizing: border-box;
    gap: min(var(--tile-spacing), var(--tile-relative-spacing));
    padding: min(var(--tile-spacing), var(--tile-relative-spacing));
    background-color: rgba(255, 255, 255, 0.55);
    border: 1px solid rgba(74, 63, 102, 0.15);
    border-radius: min(var(--tile-radius), var(--tile-relative-radius));
    overflow: hidden;
  }

  .tile__visual {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 0;
  }
</style>
