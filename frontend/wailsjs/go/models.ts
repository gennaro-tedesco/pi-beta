export namespace main {
	
	export class Weather {
	    city: string;
	    temperature: number;
	    weatherCode: number;
	    isDay: boolean;
	    windSpeed: number;
	    windDirection: number;
	    pressureMsl: number;
	    rainProbability: number;
	    uvIndex: number;
	    dailyTime: string[];
	    dailyWeatherCode: number[];
	    dailyTemperatureMax: number[];
	    dailyTemperatureMin: number[];
	
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
	        this.windDirection = source["windDirection"];
	        this.pressureMsl = source["pressureMsl"];
	        this.rainProbability = source["rainProbability"];
	        this.uvIndex = source["uvIndex"];
	        this.dailyTime = source["dailyTime"];
	        this.dailyWeatherCode = source["dailyWeatherCode"];
	        this.dailyTemperatureMax = source["dailyTemperatureMax"];
	        this.dailyTemperatureMin = source["dailyTemperatureMin"];
	    }
	}

}

