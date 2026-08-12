<script lang="ts">
  export let disabled: boolean = false
  export let ghost: boolean = false
</script>

<button class="button" class:button--ghost={ghost} {disabled} on:click {...$$restProps}>
  <slot />
</button>

<style>
  .button {
    --border-color: rgba(255, 255, 255, 0.2);
    --background-color: rgba(255, 255, 255, 0.05);
    --hover-background-color: rgba(255, 255, 255, 0.15);
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    border: 1px solid var(--border-color);
    background-color: var(--background-color);
    color: inherit;
    font: inherit;
    cursor: pointer;
  }

  .button--ghost {
    border-color: transparent;
    background-color: transparent;
  }

  .button:hover {
    background-color: var(--hover-background-color);
  }

  .button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
