<script lang="ts">
  export let disabled: boolean = false
</script>

<button class="button" {disabled} on:click {...$$restProps}>
  <slot />
</button>

<style>
  .button {
    padding: 0.5rem 1rem;
    border-radius: 6px;
    border: 1px solid rgba(255, 255, 255, 0.2);
    background-color: rgba(255, 255, 255, 0.05);
    color: inherit;
    font: inherit;
    cursor: pointer;
  }

  .button:hover {
    background-color: rgba(255, 255, 255, 0.15);
  }

  .button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
