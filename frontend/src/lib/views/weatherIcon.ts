import clearDay from '@bybas/weather-icons/production/fill/all/clear-day.svg'
import clearNight from '@bybas/weather-icons/production/fill/all/clear-night.svg'
import partlyCloudyDay from '@bybas/weather-icons/production/fill/all/partly-cloudy-day.svg'
import partlyCloudyNight from '@bybas/weather-icons/production/fill/all/partly-cloudy-night.svg'
import overcastDay from '@bybas/weather-icons/production/fill/all/overcast-day.svg'
import overcastNight from '@bybas/weather-icons/production/fill/all/overcast-night.svg'
import fogDay from '@bybas/weather-icons/production/fill/all/fog-day.svg'
import fogNight from '@bybas/weather-icons/production/fill/all/fog-night.svg'
import drizzle from '@bybas/weather-icons/production/fill/all/drizzle.svg'
import sleet from '@bybas/weather-icons/production/fill/all/sleet.svg'
import rain from '@bybas/weather-icons/production/fill/all/rain.svg'
import partlyCloudyDayRain from '@bybas/weather-icons/production/fill/all/partly-cloudy-day-rain.svg'
import partlyCloudyNightRain from '@bybas/weather-icons/production/fill/all/partly-cloudy-night-rain.svg'
import snow from '@bybas/weather-icons/production/fill/all/snow.svg'
import partlyCloudyDaySnow from '@bybas/weather-icons/production/fill/all/partly-cloudy-day-snow.svg'
import partlyCloudyNightSnow from '@bybas/weather-icons/production/fill/all/partly-cloudy-night-snow.svg'
import thunderstormsDay from '@bybas/weather-icons/production/fill/all/thunderstorms-day.svg'
import thunderstormsNight from '@bybas/weather-icons/production/fill/all/thunderstorms-night.svg'
import thunderstormsDayRain from '@bybas/weather-icons/production/fill/all/thunderstorms-day-rain.svg'
import thunderstormsNightRain from '@bybas/weather-icons/production/fill/all/thunderstorms-night-rain.svg'

interface IconEntry {
  codes: number[]
  day: string
  night: string
}

const iconEntries: IconEntry[] = [
  { codes: [0], day: clearDay, night: clearNight },
  { codes: [1, 2], day: partlyCloudyDay, night: partlyCloudyNight },
  { codes: [3], day: overcastDay, night: overcastNight },
  { codes: [45, 48], day: fogDay, night: fogNight },
  { codes: [51, 53, 55, 56, 57], day: drizzle, night: drizzle },
  { codes: [66, 67], day: sleet, night: sleet },
  { codes: [61, 63, 65], day: rain, night: rain },
  { codes: [80, 81, 82], day: partlyCloudyDayRain, night: partlyCloudyNightRain },
  { codes: [71, 73, 75, 77], day: snow, night: snow },
  { codes: [85, 86], day: partlyCloudyDaySnow, night: partlyCloudyNightSnow },
  { codes: [95], day: thunderstormsDay, night: thunderstormsNight },
  { codes: [96, 99], day: thunderstormsDayRain, night: thunderstormsNightRain }
]

const fallbackEntry = iconEntries[2]

export function getWeatherIcon(code: number, isDay: boolean): string {
  const entry = iconEntries.find((candidate) => candidate.codes.includes(code)) ?? fallbackEntry
  return isDay ? entry.day : entry.night
}
