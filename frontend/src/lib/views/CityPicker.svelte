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
  const countryDataUrl = 'https://cdn.jsdelivr.net/gh/nvkelso/natural-earth-vector@master/geojson/ne_50m_admin_0_countries.geojson'
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
  const maxSearchResults = 8
  const citySearchLabel = 'Search cities'
  const citySearchPlaceholder = 'Type a city name'

  interface CityProperties {
    adm0name: string
    min_zoom: number
    name: string
    pop_max: number
  }

  type CountryData = FeatureCollection<Polygon | MultiPolygon>
  type CityData = FeatureCollection<Point, CityProperties>

  let mapContainer: HTMLDivElement
  let map: L.Map | null = null
  let cityData: CityData | null = null
  let cityLayer: L.LayerGroup | null = null
  let cityQuery = ''

  function selectCity(feature: CityData['features'][number]): void {
    const [lon, lat] = feature.geometry.coordinates
    dispatch('select', { lat, lon })
  }

  function findCities(query: string, data: CityData | null): CityData['features'] {
    const normalizedQuery = query.trim().toLocaleLowerCase()
    if (normalizedQuery === '' || data === null) return []

    return data.features
      .filter((feature) => feature.properties.name.toLocaleLowerCase().includes(normalizedQuery))
      .sort((first, second) => {
        const firstName = first.properties.name.toLocaleLowerCase()
        const secondName = second.properties.name.toLocaleLowerCase()
        const prefixDifference = Number(secondName.startsWith(normalizedQuery)) - Number(firstName.startsWith(normalizedQuery))
        return prefixDifference || second.properties.pop_max - first.properties.pop_max
      })
      .slice(0, maxSearchResults)
  }

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
        .on('click', () => selectCity(feature))
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

  $: citySearchResults = findCities(cityQuery, cityData)
</script>

<div class="picker">
  <div class="picker__header">
    <span class="picker__title">Choose a city</span>
    <div class="picker__search">
      <form autocomplete="off" on:submit|preventDefault={() => citySearchResults[0] && selectCity(citySearchResults[0])}>
        <input
          type="search"
          bind:value={cityQuery}
          aria-label={citySearchLabel}
          placeholder={citySearchPlaceholder}
          autocomplete="off"
          autocapitalize="off"
          spellcheck="false"
        />
      </form>
      {#if citySearchResults.length > 0}
        <div class="picker__results" role="listbox" aria-label={citySearchLabel}>
          {#each citySearchResults as city (city.properties.name + city.geometry.coordinates.join(','))}
            <Button ghost on:click={() => selectCity(city)} role="option">
              <span>{city.properties.name}</span>
              <span class="picker__country">{city.properties.adm0name}</span>
            </Button>
          {/each}
        </div>
      {/if}
    </div>
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
    --header-padding: 1rem;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: var(--header-padding);
  }

  .picker__title {
    font-size: 1.1rem;
    font-weight: 600;
  }

  .picker__search {
    --search-width: min(24rem, 50vw);
    --search-layer: 1000;
    --search-gap: 0.25rem;
    --search-radius: 6px;
    position: absolute;
    z-index: var(--search-layer);
    left: 50%;
    width: var(--search-width);
    transform: translateX(-50%);
  }

  .picker__search input {
    --search-background: rgba(255, 255, 255, 0.12);
    --search-border: rgba(255, 255, 255, 0.28);
    --search-spacing: 0.6rem;
    box-sizing: border-box;
    width: 100%;
    padding: var(--search-spacing);
    border: 1px solid var(--search-border);
    border-radius: var(--search-radius);
    background: var(--search-background);
    color: inherit;
    font: inherit;
  }

  .picker__search input::placeholder {
    color: inherit;
    opacity: 0.7;
  }

  .picker__results {
    --results-background: rgba(27, 38, 54, 0.96);
    --results-gap: 0.25rem;
    --results-padding: 0.25rem;
    position: absolute;
    top: calc(100% + var(--search-gap));
    right: 0;
    left: 0;
    display: flex;
    flex-direction: column;
    gap: var(--results-gap);
    padding: var(--results-padding);
    border-radius: var(--search-radius);
    background: var(--results-background);
  }

  .picker__results :global(.button) {
    justify-content: space-between;
    width: 100%;
  }

  .picker__country {
    opacity: 0.65;
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
