<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import L from 'leaflet'
  import 'leaflet/dist/leaflet.css'
  import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png'
  import markerIcon from 'leaflet/dist/images/marker-icon.png'
  import markerShadow from 'leaflet/dist/images/marker-shadow.png'
  import { X } from 'lucide-svelte'
  import { Button } from '../components'

  const dispatch = createEventDispatcher<{ select: { lat: number; lon: number }; close: void }>()

  const defaultZoom = 2
  const defaultCenter: L.LatLngExpression = [20, 0]
  const tileUrl = 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png'
  const tileAttribution = '&copy; OpenStreetMap contributors'
  const maxZoom = 18

  let mapContainer: HTMLDivElement
  let map: L.Map | null = null

  onMount(() => {
    L.Icon.Default.mergeOptions({
      iconRetinaUrl: markerIcon2x,
      iconUrl: markerIcon,
      shadowUrl: markerShadow
    })

    map = L.map(mapContainer).setView(defaultCenter, defaultZoom)
    L.tileLayer(tileUrl, { attribution: tileAttribution, maxZoom }).addTo(map)

    map.on('click', (event: L.LeafletMouseEvent) => {
      dispatch('select', { lat: event.latlng.lat, lon: event.latlng.lng })
    })
  })

  onDestroy(() => {
    map?.remove()
  })
</script>

<div class="picker">
  <div class="picker__header">
    <span class="picker__title">Choose a location</span>
    <Button ghost on:click={() => dispatch('close')} aria-label="Close">
      <X size={20} />
    </Button>
  </div>
  <div class="picker__map" bind:this={mapContainer}></div>
</div>

<style>
  .picker {
    position: fixed;
    inset: 0;
    z-index: 100;
    display: flex;
    flex-direction: column;
    background-color: var(--app-background);
  }

  .picker__header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem;
  }

  .picker__title {
    font-size: 1.1rem;
    font-weight: 600;
  }

  .picker__map {
    flex: 1;
  }
</style>
