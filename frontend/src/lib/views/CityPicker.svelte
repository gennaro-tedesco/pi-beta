<script context="module" lang="ts">
  import type { FeatureCollection, MultiPolygon, Point, Polygon } from 'geojson'

  const countryDataUrl = 'https://cdn.jsdelivr.net/gh/nvkelso/natural-earth-vector@master/geojson/ne_50m_admin_0_countries.geojson'
  const cityDataUrl = 'https://cdn.jsdelivr.net/gh/nvkelso/natural-earth-vector@master/geojson/ne_10m_populated_places_simple.geojson'

  interface CityProperties {
    adm0name: string
    min_zoom: number
    name: string
    pop_max: number
  }

  type CountryData = FeatureCollection<Polygon | MultiPolygon>
  type CityData = FeatureCollection<Point, CityProperties>
  type CityFeature = CityData['features'][number]
  type SearchableCity = { feature: CityFeature; normalizedName: string }

  let countryDataPromise: Promise<CountryData> | null = null
  let cityDataPromise: Promise<CityData> | null = null
  let searchableCities: SearchableCity[] = []
  let citySearchIndex = new Map<string, SearchableCity[]>()
  let cityBuckets = new Map<string, CityFeature[]>()

  const coordinateBucketSize = 10
  const latitudeOffset = 90
  const longitudeOffset = 180
  const maximumSearchIndexLength = 3

  function coordinateBucket(value: number, offset: number): number {
    return Math.floor((value + offset) / coordinateBucketSize)
  }

  function bucketKey(latitudeBucket: number, longitudeBucket: number): string {
    return `${latitudeBucket}:${longitudeBucket}`
  }

  function prepareCityData(data: CityData): CityData {
    searchableCities = data.features.map((feature) => ({
      feature,
      normalizedName: feature.properties.name.toLocaleLowerCase()
    }))
    citySearchIndex = new Map<string, SearchableCity[]>()
    cityBuckets = new Map<string, CityFeature[]>()

    for (const searchableCity of searchableCities) {
      const { feature, normalizedName } = searchableCity
      const searchKeys = new Set<string>()
      for (let length = 1; length <= maximumSearchIndexLength; length += 1) {
        for (let start = 0; start <= normalizedName.length - length; start += 1) {
          searchKeys.add(normalizedName.slice(start, start + length))
        }
      }
      for (const searchKey of searchKeys) {
        const candidates = citySearchIndex.get(searchKey)
        if (candidates) candidates.push(searchableCity)
        else citySearchIndex.set(searchKey, [searchableCity])
      }

      const [longitude, latitude] = feature.geometry.coordinates
      const key = bucketKey(
        coordinateBucket(latitude, latitudeOffset),
        coordinateBucket(longitude, longitudeOffset)
      )
      const bucket = cityBuckets.get(key)
      if (bucket) bucket.push(feature)
      else cityBuckets.set(key, [feature])
    }

    return data
  }

  function fetchJson<T>(url: string): Promise<T> {
    return fetch(url).then((response) => {
      if (!response.ok) throw new Error(`Request to ${url} failed: ${response.status}`)
      return response.json() as Promise<T>
    })
  }

  function loadCountryData(): Promise<CountryData> {
    if (countryDataPromise === null) {
      countryDataPromise = fetchJson<CountryData>(countryDataUrl)
      countryDataPromise.catch(() => {
        countryDataPromise = null
      })
    }
    return countryDataPromise
  }

  function loadCityData(): Promise<CityData> {
    if (cityDataPromise === null) {
      cityDataPromise = fetchJson<CityData>(cityDataUrl).then(prepareCityData)
      cityDataPromise.catch(() => {
        cityDataPromise = null
      })
    }
    return cityDataPromise
  }
</script>

<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import { scale } from 'svelte/transition'
  import L from 'leaflet'
  import 'leaflet/dist/leaflet.css'
  import { X } from 'lucide-svelte'
  import { Button, Message } from '../components'
  import { transitionDuration } from '../transition'

  const dispatch = createEventDispatcher<{ select: { lat: number; lon: number }; close: void }>()

  const defaultZoom = 2
  const defaultCenter: L.LatLngExpression = [20, 0]
  const worldBounds = L.latLngBounds([-85, -180], [85, 180])
  const maxBoundsViscosity = 1
  const countryStyle: L.PathOptions = {
    color: 'rgba(74, 63, 102, 0.55)',
    fillColor: 'rgba(255, 255, 255, 0.4)',
    fillOpacity: 1,
    weight: 1
  }
  const cityStyle: L.CircleMarkerOptions = {
    color: '#4a3f66',
    fillColor: '#4a3f66',
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
  const searchDebounceMs = 150
  const cityDataErrorMessage = 'Unable to load city data. Close and reopen the city picker to retry.'

  let mapContainer: HTMLDivElement
  let map: L.Map | null = null
  let cityData: CityData | null = null
  let cityLayer: L.LayerGroup | null = null
  let cityQuery = ''
  let debouncedCityQuery = ''
  let cityDataError: string | null = null
  let searchDebounceTimer: ReturnType<typeof setTimeout>
  const visibleMarkers = new Map<string, L.CircleMarker>()

  function selectCity(feature: CityData['features'][number]): void {
    const [lon, lat] = feature.geometry.coordinates
    dispatch('select', { lat, lon })
  }

  function findCities(query: string, data: CityData | null): CityData['features'] {
    const normalizedQuery = query.trim().toLocaleLowerCase()
    if (normalizedQuery === '' || data === null) return []

    const matches: SearchableCity[] = []
    const searchKey = normalizedQuery.slice(0, maximumSearchIndexLength)
    const candidates = citySearchIndex.get(searchKey) ?? []
    for (const candidate of candidates) {
      if (!candidate.normalizedName.includes(normalizedQuery)) continue

      const insertionIndex = matches.findIndex((match) => {
        const prefixDifference = Number(candidate.normalizedName.startsWith(normalizedQuery))
          - Number(match.normalizedName.startsWith(normalizedQuery))
        return prefixDifference > 0
          || (prefixDifference === 0 && candidate.feature.properties.pop_max > match.feature.properties.pop_max)
      })
      matches.splice(insertionIndex < 0 ? matches.length : insertionIndex, 0, candidate)
      if (matches.length > maxSearchResults) matches.pop()
    }

    return matches.map((match) => match.feature)
  }

  function cityKey(feature: CityData['features'][number]): string {
    return feature.properties.name + feature.geometry.coordinates.join(',')
  }

  function renderCities(): void {
    if (map === null || cityData === null || cityLayer === null) return

    const bounds = map.getBounds()
    const zoom = map.getZoom()
    const nextVisible = new Set<string>()
    const southBucket = coordinateBucket(bounds.getSouth(), latitudeOffset)
    const northBucket = coordinateBucket(bounds.getNorth(), latitudeOffset)
    const westBucket = coordinateBucket(bounds.getWest(), longitudeOffset)
    const eastBucket = coordinateBucket(bounds.getEast(), longitudeOffset)

    for (let latitudeBucket = southBucket; latitudeBucket <= northBucket; latitudeBucket += 1) {
      for (let longitudeBucket = westBucket; longitudeBucket <= eastBucket; longitudeBucket += 1) {
        const bucket = cityBuckets.get(bucketKey(latitudeBucket, longitudeBucket)) ?? []
        for (const feature of bucket) {
          if (feature.properties.min_zoom > zoom) continue

          const [lon, lat] = feature.geometry.coordinates
          if (!bounds.contains([lat, lon])) continue

          const key = cityKey(feature)
          nextVisible.add(key)
          if (visibleMarkers.has(key)) continue

          const marker = L.circleMarker([lat, lon], cityStyle)
            .bindTooltip(feature.properties.name, cityTooltipOptions)
            .on('click', () => selectCity(feature))
            .addTo(cityLayer)
          visibleMarkers.set(key, marker)
        }
      }
    }

    for (const [key, marker] of visibleMarkers) {
      if (nextVisible.has(key)) continue
      marker.remove()
      visibleMarkers.delete(key)
    }
  }

  onMount(() => {
    map = L.map(mapContainer, { maxBounds: worldBounds, maxBoundsViscosity }).setView(defaultCenter, defaultZoom)
    cityLayer = L.layerGroup().addTo(map)
    map.on('moveend zoomend', renderCities)

    Promise.all([loadCountryData(), loadCityData()])
      .then(([countryData, loadedCityData]) => {
        if (map === null) return
        L.geoJSON(countryData, {
          interactive: false,
          style: countryStyle,
          filter: (feature) => feature.properties.CONTINENT !== 'Antarctica'
        }).addTo(map)
        cityData = loadedCityData
        renderCities()
      })
      .catch(() => {
        cityDataError = cityDataErrorMessage
      })

    setTimeout(() => map?.invalidateSize(), transitionDuration + 50)
  })

  onDestroy(() => {
    map?.remove()
    clearTimeout(searchDebounceTimer)
  })

  $: {
    clearTimeout(searchDebounceTimer)
    searchDebounceTimer = setTimeout(() => {
      debouncedCityQuery = cityQuery
    }, searchDebounceMs)
  }

  $: citySearchResults = findCities(debouncedCityQuery, cityData)
</script>

<div
  class="picker"
  out:scale={{ duration: transitionDuration }}
  in:scale={{ duration: transitionDuration, delay: transitionDuration }}
>
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
          {#each citySearchResults as city (cityKey(city))}
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
  {#if cityDataError}
    <Message variant="error" message={cityDataError} />
  {/if}
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
    --header-minimum-size: 0;
    --header-side-column: minmax(0, 1fr);
    --header-search-column: 2;
    --search-width: min(24rem, 50vw);
    position: relative;
    display: grid;
    grid-template-columns: var(--header-side-column) minmax(0, var(--search-width)) var(--header-side-column);
    align-items: center;
    padding: var(--header-padding);
  }

  .picker__title {
    min-width: var(--header-minimum-size);
    font-size: 1.1rem;
    font-weight: 600;
  }

  .picker__search {
    --search-layer: 1000;
    --search-gap: 0.25rem;
    --search-radius: 6px;
    position: relative;
    z-index: var(--search-layer);
    grid-column: var(--header-search-column);
    width: var(--search-width);
    max-width: var(--picker-width);
  }

  .picker__header > :global(.button) {
    justify-self: end;
  }

  .picker__search input {
    --search-background: rgba(255, 255, 255, 0.6);
    --search-border: rgba(74, 63, 102, 0.25);
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
    --results-background: rgba(255, 255, 255, 0.96);
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
    --city-label-color: #4a3f66;
    --city-label-shadow: 0 1px 3px rgba(255, 255, 255, 0.9);
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
