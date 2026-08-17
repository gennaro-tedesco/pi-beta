export namespace main {
	
	export class MachineProcess {
	    pid: number;
	    name: string;
	    memoryBytes: number;
	    memoryPercent: number;
	
	    static createFrom(source: any = {}) {
	        return new MachineProcess(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.memoryBytes = source["memoryBytes"];
	        this.memoryPercent = source["memoryPercent"];
	    }
	}
	export class NetworkStatus {
	    availability: string;
	    connectionType: string;
	    internetStatus: string;
	    internetReachable?: boolean;
	    networkName?: string;
	    signalPercent?: number;
	    signalDbm?: number;
	    signalQuality: string;
	    linkBitrateMbps?: number;
	    latencyMs?: number;
	    averageLatencyMs?: number;
	    jitterMs?: number;
	    probeLossPercent: number;
	    qualityScore: number;
	    quality: string;
	    stability: string;
	    sampleCount: number;
	    checkedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.availability = source["availability"];
	        this.connectionType = source["connectionType"];
	        this.internetStatus = source["internetStatus"];
	        this.internetReachable = source["internetReachable"];
	        this.networkName = source["networkName"];
	        this.signalPercent = source["signalPercent"];
	        this.signalDbm = source["signalDbm"];
	        this.signalQuality = source["signalQuality"];
	        this.linkBitrateMbps = source["linkBitrateMbps"];
	        this.latencyMs = source["latencyMs"];
	        this.averageLatencyMs = source["averageLatencyMs"];
	        this.jitterMs = source["jitterMs"];
	        this.probeLossPercent = source["probeLossPercent"];
	        this.qualityScore = source["qualityScore"];
	        this.quality = source["quality"];
	        this.stability = source["stability"];
	        this.sampleCount = source["sampleCount"];
	        this.checkedAt = source["checkedAt"];
	    }
	}
	export class MachineStatus {
	    hostname: string;
	    platform: string;
	    platformVersion: string;
	    architecture: string;
	    uptimeSeconds: number;
	    cpuCount: number;
	    cpuPercent?: number;
	    memoryTotalBytes: number;
	    memoryUsedBytes: number;
	    memoryUsedPercent?: number;
	    diskTotalBytes: number;
	    diskUsedBytes: number;
	    diskUsedPercent?: number;
	    processCount: number;
	    largestProcesses: MachineProcess[];
	    network: NetworkStatus;
	
	    static createFrom(source: any = {}) {
	        return new MachineStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.platform = source["platform"];
	        this.platformVersion = source["platformVersion"];
	        this.architecture = source["architecture"];
	        this.uptimeSeconds = source["uptimeSeconds"];
	        this.cpuCount = source["cpuCount"];
	        this.cpuPercent = source["cpuPercent"];
	        this.memoryTotalBytes = source["memoryTotalBytes"];
	        this.memoryUsedBytes = source["memoryUsedBytes"];
	        this.memoryUsedPercent = source["memoryUsedPercent"];
	        this.diskTotalBytes = source["diskTotalBytes"];
	        this.diskUsedBytes = source["diskUsedBytes"];
	        this.diskUsedPercent = source["diskUsedPercent"];
	        this.processCount = source["processCount"];
	        this.largestProcesses = this.convertValues(source["largestProcesses"], MachineProcess);
	        this.network = this.convertValues(source["network"], NetworkStatus);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
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

