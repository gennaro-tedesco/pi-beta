import rain from '@bybas/weather-icons/production/fill/all/rain.svg'
import wind from '@bybas/weather-icons/production/fill/all/wind.svg'

export interface WeatherScene {
  label: string
  gradient: string
  night: boolean
  windAnimation: string
  windEffectCount: number
  rainAnimation: string | null
}

interface ConditionGroup {
  codes: number[]
  label: string
  gradient: string
}

const conditionGroups: ConditionGroup[] = [
  {
    codes: [0],
    label: 'Clear',
    gradient: 'linear-gradient(160deg, #ffe8b8 0%, #ffd3e0 55%, #cdeaf9 100%)'
  },
  {
    codes: [1, 2, 3],
    label: 'Cloudy',
    gradient: 'linear-gradient(160deg, #e7e6f5 0%, #d8e1ee 55%, #eef1f7 100%)'
  },
  {
    codes: [45, 48],
    label: 'Fog',
    gradient: 'linear-gradient(160deg, #eef2f0 0%, #e3ecec 55%, #f2f6f5 100%)'
  },
  {
    codes: [51, 53, 55, 56, 57, 61, 63, 65, 66, 67, 80, 81, 82],
    label: 'Rain',
    gradient: 'linear-gradient(160deg, #cdeaf6 0%, #b9dcf0 55%, #dceaf7 100%)'
  },
  {
    codes: [71, 73, 75, 77, 85, 86],
    label: 'Snow',
    gradient: 'linear-gradient(160deg, #f3f9fd 0%, #e2eef8 55%, #eef6fb 100%)'
  },
  {
    codes: [95, 96, 99],
    label: 'Thunderstorm',
    gradient: 'linear-gradient(160deg, #e3d9f5 0%, #d3c3ea 55%, #ece3f7 100%)'
  }
]

const fallbackGroup = conditionGroups[1]
const likelyRainProbability = 30
const rainCodes = conditionGroups[3].codes
const mediumWindSpeed = 20
const strongWindSpeed = 39
const lightWindEffectCount = 1
const mediumWindEffectCount = 2
const strongWindEffectCount = 3

interface WeatherConditions {
  windSpeed: number
  rainProbability: number
}

function getWindEffectCount(windSpeed: number): number {
  if (windSpeed >= strongWindSpeed) {
    return strongWindEffectCount
  }
  if (windSpeed >= mediumWindSpeed) {
    return mediumWindEffectCount
  }
  return lightWindEffectCount
}

export function getWeatherScene(
  code: number,
  isDay: boolean,
  conditions: WeatherConditions
): WeatherScene {
  const group = conditionGroups.find((candidate) => candidate.codes.includes(code)) ?? fallbackGroup
  const rainExpected = rainCodes.includes(code) || conditions.rainProbability >= likelyRainProbability

  return {
    label: group.label,
    gradient: group.gradient,
    night: !isDay,
    windAnimation: wind,
    windEffectCount: getWindEffectCount(conditions.windSpeed),
    rainAnimation: rainExpected ? rain : null
  }
}
