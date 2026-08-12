<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import {
    BarController,
    BarElement,
    CategoryScale,
    Chart as ChartJS,
    Legend,
    LinearScale,
    LineController,
    LineElement,
    PointElement,
    Tooltip
  } from 'chart.js'
  import type { ChartConfiguration } from 'chart.js'

  ChartJS.register(
    BarController,
    BarElement,
    CategoryScale,
    LinearScale,
    LineController,
    LineElement,
    PointElement,
    Legend,
    Tooltip
  )

  export let config: ChartConfiguration

  let canvas: HTMLCanvasElement
  let chart: ChartJS | null = null

  onMount(() => {
    chart = new ChartJS(canvas, config)
  })

  onDestroy(() => {
    chart?.destroy()
  })

  $: if (chart) {
    chart.data = config.data
    chart.options = config.options ?? {}
    chart.update()
  }
</script>

<canvas bind:this={canvas} />
