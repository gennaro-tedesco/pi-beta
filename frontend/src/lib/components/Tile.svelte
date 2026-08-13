<script lang="ts">
  import { createEventDispatcher } from 'svelte'

  export let title: string

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

<div class="tile" role="button" tabindex="0" aria-label={title} on:click={handleActivate} on:keydown={handleKeydown}>
  <div class="tile__visual">
    <slot name="visual" />
  </div>
</div>

<style>
  .tile {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    background-color: rgba(255, 255, 255, 0.55);
    border: 1px solid rgba(74, 63, 102, 0.15);
    border-radius: 10px;
    padding: 1rem;
    cursor: pointer;
    user-select: none;
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
