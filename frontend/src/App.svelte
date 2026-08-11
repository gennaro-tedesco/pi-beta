<script lang="ts">
  import { Message, Tile } from './lib/components'

  const widgets = ['Photos', 'Weather', 'Stocks']
</script>

<main class="tiles">
  {#each widgets as title (title)}
    <Tile {title}>
      <Message variant="warning" message="Not implemented yet" />
    </Tile>
  {/each}
</main>

<style>
  .tiles {
    display: flex;
    gap: 1rem;
    padding: 1rem;
    overflow-x: auto;
  }

  .tiles :global(.tile) {
    flex: 0 0 auto;
    min-width: 220px;
  }
</style>
