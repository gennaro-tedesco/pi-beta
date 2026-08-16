export namespace main {
	
	export class Weather {
	    city: string;
	    temperature: number;
	    weatherCode: number;
	    isDay: boolean;
	    windSpeed: number;
	    rainProbability: number;
	    uvIndex: number;
	    pressure: number;
	    airQuality: number;
	    dailyTime: string[];
	    dailyWeatherCode: number[];
	    dailyTemperatureMax: number[];
	    dailyTemperatureMin: number[];
	    sunrise: string;
	    sunset: string;
	    hourlyTime: string[];
	    hourlyTemperature: number[];
	    hourlyWeatherCode: number[];
	    utcOffsetSeconds: number;
	
	    static createFrom(source: any = {}) {
	        return new Weather(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.city = source["city"];
	        this.temperature = source["temperature"];
	        this.weatherCode = source["weatherCode"];
	        this.isDay = source["isDay"];
	        this.windSpeed = source["windSpeed"];
	        this.rainProbability = source["rainProbability"];
	        this.uvIndex = source["uvIndex"];
	        this.pressure = source["pressure"];
	        this.airQuality = source["airQuality"];
	        this.dailyTime = source["dailyTime"];
	        this.dailyWeatherCode = source["dailyWeatherCode"];
	        this.dailyTemperatureMax = source["dailyTemperatureMax"];
	        this.dailyTemperatureMin = source["dailyTemperatureMin"];
	        this.sunrise = source["sunrise"];
	        this.sunset = source["sunset"];
	        this.hourlyTime = source["hourlyTime"];
	        this.hourlyTemperature = source["hourlyTemperature"];
	        this.hourlyWeatherCode = source["hourlyWeatherCode"];
	        this.utcOffsetSeconds = source["utcOffsetSeconds"];
	    }
	}

}

