import compass from '@bybas/weather-icons/production/fill/all/compass.svg'
import pressureHigh from '@bybas/weather-icons/production/fill/all/pressure-high.svg'
import pressureLow from '@bybas/weather-icons/production/fill/all/pressure-low.svg'
import umbrella from '@bybas/weather-icons/production/fill/all/umbrella.svg'
import wind from '@bybas/weather-icons/production/fill/all/wind.svg'

export interface AtmosphericAnimation {
  label: string
  source: string
  rotation: number
}

export interface WeatherScene {
  label: string
  gradient: string
  night: boolean
  animations: AtmosphericAnimation[]
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
    gradient: 'linear-gradient(160deg, #1c4d73 0%, #4fa8e0 55%, #8ed6c9 100%)'
  },
  {
    codes: [1, 2, 3],
    label: 'Cloudy',
    gradient: 'linear-gradient(160deg, #3b4a5a 0%, #7c93ab 55%, #b8c6d6 100%)'
  },
  {
    codes: [45, 48],
    label: 'Fog',
    gradient: 'linear-gradient(160deg, #3f4750 0%, #97a3ab 55%, #c7cfd4 100%)'
  },
  {
    codes: [51, 53, 55, 56, 57, 61, 63, 65, 66, 67, 80, 81, 82],
    label: 'Rain',
    gradient: 'linear-gradient(160deg, #1f2c3a 0%, #4a6178 55%, #6f8ea3 100%)'
  },
  {
    codes: [71, 73, 75, 77, 85, 86],
    label: 'Snow',
    gradient: 'linear-gradient(160deg, #4a5560 0%, #93a8b8 55%, #d7e3ea 100%)'
  },
  {
    codes: [95, 96, 99],
    label: 'Thunderstorm',
    gradient: 'linear-gradient(160deg, #14141c 0%, #2c3347 55%, #4a5068 100%)'
  }
]

const fallbackGroup = conditionGroups[1]
const calmWindSpeed = 1
const likelyRainProbability = 30
const standardPressureMsl = 1013.25
const noRotation = 0
const highPressureLabel = 'High pressure'
const lowPressureLabel = 'Low pressure'
const windDirectionLabel = 'Wind direction'
const windLabel = 'Wind'
const rainExpectedLabel = 'Rain expected'

interface AtmosphericConditions {
  windSpeed: number
  windDirection: number
  pressureMsl: number
  rainProbability: number
}

function getAtmosphericAnimations(conditions: AtmosphericConditions): AtmosphericAnimation[] {
  const animations: AtmosphericAnimation[] = [
    {
      label: conditions.pressureMsl >= standardPressureMsl ? highPressureLabel : lowPressureLabel,
      source: conditions.pressureMsl >= standardPressureMsl ? pressureHigh : pressureLow,
      rotation: noRotation
    },
    { label: windDirectionLabel, source: compass, rotation: conditions.windDirection }
  ]

  if (conditions.windSpeed > calmWindSpeed) {
    animations.unshift({ label: windLabel, source: wind, rotation: noRotation })
  }
  if (conditions.rainProbability >= likelyRainProbability) {
    animations.push({ label: rainExpectedLabel, source: umbrella, rotation: noRotation })
  }

  return animations
}

export function getWeatherScene(
  code: number,
  isDay: boolean,
  conditions: AtmosphericConditions
): WeatherScene {
  const group = conditionGroups.find((candidate) => candidate.codes.includes(code)) ?? fallbackGroup
  return {
    label: group.label,
    gradient: group.gradient,
    night: !isDay,
    animations: getAtmosphericAnimations(conditions)
  }
}
