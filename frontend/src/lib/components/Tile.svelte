<script lang="ts">
  import { createEventDispatcher } from 'svelte'

  export let title: string
  export let expanded: boolean = false

  const dispatch = createEventDispatcher<{ toggle: boolean }>()

  function handleActivate(): void {
    expanded = !expanded
    dispatch('toggle', expanded)
  }

  function handleKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      handleActivate()
    }
  }
</script>

<div
  class="tile"
  class:tile--expanded={expanded}
  role="button"
  tabindex="0"
  aria-expanded={expanded}
  on:click={handleActivate}
  on:keydown={handleKeydown}
>
  <div class="tile__header">
    <slot name="icon" />
    <span class="tile__title">{title}</span>
  </div>

  {#if expanded}
    <div class="tile__content">
      <slot />
    </div>
  {/if}
</div>

<style>
  .tile {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    background-color: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 10px;
    padding: 1rem;
    cursor: pointer;
    user-select: none;
  }

  .tile--expanded {
    border-color: rgba(255, 255, 255, 0.3);
  }

  .tile__header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .tile__title {
    font-weight: 600;
  }
</style>
