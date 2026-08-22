<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import Button from './Button.svelte'
  import homeLogoSrc from '../../assets/images/logo.png'

  const dispatch = createEventDispatcher<{ home: void }>()
  const homeLabel = 'Home'
</script>

<nav class="navigation">
  <Button ghost on:click={() => dispatch('home')} aria-label={homeLabel}>
    <img class="navigation__logo" src={homeLogoSrc} alt="" />
  </Button>
  <div class="navigation__actions">
    <slot name="actions" />
  </div>
</nav>

<style>
  .navigation {
    --navigation-layer: 1;
    --navigation-zero: 0;
    --navigation-actions-gap: 0.5rem;
    --navigation-home-size: 3rem;
    --navigation-logo-fill: 90%;
    position: relative;
    z-index: var(--navigation-layer);
    display: flex;
    justify-content: space-between;
    align-items: center;
    pointer-events: none;
  }

  .navigation :global(button) {
    pointer-events: auto;
  }

  .navigation__actions {
    display: flex;
    gap: var(--navigation-actions-gap);
  }

  .navigation > :global(button) {
    width: var(--navigation-home-size);
    height: var(--navigation-home-size);
    padding: var(--navigation-zero);
  }

  .navigation__logo {
    width: var(--navigation-logo-fill);
    height: var(--navigation-logo-fill);
  }
</style>
