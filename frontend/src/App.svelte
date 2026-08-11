<script lang="ts">
  import { House } from 'lucide-svelte'
  import { Button, Message, Tile } from './lib/components'

  type WidgetId = 'photos' | 'weather'

  const widgets: { id: WidgetId; title: string }[] = [
    { id: 'photos', title: 'Photos' },
    { id: 'weather', title: 'Weather' }
  ]

  let activeWidget: WidgetId | null = null

  function maximize(id: WidgetId): void {
    activeWidget = id
  }

  function goHome(): void {
    activeWidget = null
  }
</script>

<main class="screen">
  {#if activeWidget !== null}
    <div class="nav">
      <Button on:click={goHome} aria-label="Home">
        <House />
      </Button>
    </div>
  {/if}

  <div class="content">
    {#if activeWidget === null}
      <div class="tiles">
        {#each widgets as widget (widget.id)}
          <Tile title={widget.title} on:toggle={() => maximize(widget.id)} />
        {/each}
      </div>
    {:else}
      <Message variant="warning" message="Not implemented yet" />
    {/if}
  </div>
</main>

<style>
  .screen {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    padding: 1rem;
    gap: 1rem;
  }

  .nav {
    display: flex;
    justify-content: center;
  }

  .content {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .tiles {
    --tile-height: 50vh;
    --tile-width: 30vw;
    display: flex;
    gap: 1rem;
  }

  .tiles :global(.tile) {
    flex: 0 0 auto;
    height: var(--tile-height);
    width: var(--tile-width);
  }
</style>
