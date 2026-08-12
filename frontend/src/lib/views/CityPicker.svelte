<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import L from 'leaflet'
  import type { FeatureCollection, MultiPolygon, Point, Polygon } from 'geojson'
  import 'leaflet/dist/leaflet.css'
  import { X } from 'lucide-svelte'
  import { Button } from '../components'

  const dispatch = createEventDispatcher<{ select: { lat: number; lon: number }; close: void }>()

  const defaultZoom = 2
  const defaultCenter: L.LatLngExpression = [20, 0]
  const countryDataUrl = 'https://cdn.jsdelivr.net/gh/nvkelso/natural-earth-vector@master/geojson/ne_110m_admin_0_countries.geojson'
  const cityDataUrl = 'https://cdn.jsdelivr.net/gh/nvkelso/natural-earth-vector@master/geojson/ne_10m_populated_places_simple.geojson'
  const worldBounds = L.latLngBounds([-85, -180], [85, 180])
  const maxBoundsViscosity = 1
  const countryStyle: L.PathOptions = {
    color: 'rgba(255, 255, 255, 0.72)',
    fillColor: 'rgba(255, 255, 255, 0.08)',
    fillOpacity: 1,
    weight: 1
  }
  const cityStyle: L.CircleMarkerOptions = {
    color: '#ffffff',
    fillColor: '#ffffff',
    fillOpacity: 0.9,
    radius: 6,
    weight: 6
  }
  const cityTooltipOptions: L.TooltipOptions = {
    direction: 'right',
    offset: [6, 0],
    opacity: 1,
    permanent: true,
    interactive: true
  }

  interface CityProperties {
    min_zoom: number
    name: string
  }

  type CountryData = FeatureCollection<Polygon | MultiPolygon>
  type CityData = FeatureCollection<Point, CityProperties>

  let mapContainer: HTMLDivElement
  let map: L.Map | null = null
  let cityData: CityData | null = null
  let cityLayer: L.LayerGroup | null = null

  function renderCities(): void {
    if (map === null || cityData === null || cityLayer === null) return

    const bounds = map.getBounds()
    const zoom = map.getZoom()
    cityLayer.clearLayers()

    for (const feature of cityData.features) {
      if (feature.properties.min_zoom > zoom) continue

      const [lon, lat] = feature.geometry.coordinates
      if (!bounds.contains([lat, lon])) continue

      L.circleMarker([lat, lon], cityStyle)
        .bindTooltip(feature.properties.name, cityTooltipOptions)
        .on('click', () => dispatch('select', { lat, lon }))
        .addTo(cityLayer)
    }
  }

  onMount(() => {
    map = L.map(mapContainer, { maxBounds: worldBounds, maxBoundsViscosity }).setView(defaultCenter, defaultZoom)
    cityLayer = L.layerGroup().addTo(map)
    map.on('moveend zoomend', renderCities)

    Promise.all([
      fetch(countryDataUrl).then((response) => response.json() as Promise<CountryData>),
      fetch(cityDataUrl).then((response) => response.json() as Promise<CityData>)
    ]).then(([countryData, loadedCityData]) => {
      if (map === null) return
      L.geoJSON(countryData, { interactive: false, style: countryStyle }).addTo(map)
      cityData = loadedCityData
      renderCities()
    })
  })

  onDestroy(() => {
    map?.remove()
  })
</script>

<div class="picker">
  <div class="picker__header">
    <span class="picker__title">Choose a city</span>
    <Button ghost on:click={() => dispatch('close')} aria-label="Close">
      <X size={20} />
    </Button>
  </div>
  <div class="picker__map" bind:this={mapContainer}></div>
</div>

<style>
  .picker {
    --picker-alignment: stretch;
    --picker-flex: 1;
    --picker-width: 100%;
    display: flex;
    flex-direction: column;
    align-self: var(--picker-alignment);
    flex: var(--picker-flex);
    width: var(--picker-width);
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
    --map-background: transparent;
    flex: 1;
    background: var(--map-background);
  }

  .picker__map:global(.leaflet-container) {
    background: var(--map-background);
  }

  .picker__map :global(.leaflet-control-attribution) {
    background: var(--map-background);
    color: inherit;
  }

  .picker__map :global(.leaflet-tooltip) {
    --city-label-background: transparent;
    --city-label-border: transparent;
    --city-label-color: #ffffff;
    --city-label-shadow: 0 1px 3px rgba(0, 0, 0, 0.8);
    background: var(--city-label-background);
    border-color: var(--city-label-border);
    box-shadow: none;
    color: var(--city-label-color);
    cursor: pointer;
    text-shadow: var(--city-label-shadow);
  }

  .picker__map :global(.leaflet-interactive) {
    cursor: pointer;
    stroke-opacity: 0.35;
  }

  .picker__map :global(.leaflet-tooltip-right::before) {
    border-right-color: transparent;
  }
</style>
