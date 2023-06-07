export namespace app {
	
	export class DatabaseInfo {
	    collections: string[];
	    stats: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new DatabaseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.collections = source["collections"];
	        this.stats = source["stats"];
	    }
	}
	export class EnvironmentInfo {
	    arch: string;
	    buildType: string;
	    platform: string;
	    version: string;
	    hasMongoExport: boolean;
	    hasMongoDump: boolean;
	    homeDirectory: string;
	    dataDirectory: string;
	    logDirectory: string;
	    downloadDirectory: string;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.arch = source["arch"];
	        this.buildType = source["buildType"];
	        this.platform = source["platform"];
	        this.version = source["version"];
	        this.hasMongoExport = source["hasMongoExport"];
	        this.hasMongoDump = source["hasMongoDump"];
	        this.homeDirectory = source["homeDirectory"];
	        this.dataDirectory = source["dataDirectory"];
	        this.logDirectory = source["logDirectory"];
	        this.downloadDirectory = source["downloadDirectory"];
	    }
	}
	export class HostInfo {
	    databases: string[];
	    status: {[key: string]: any};
	    systemInfo: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new HostInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.databases = source["databases"];
	        this.status = source["status"];
	        this.systemInfo = source["systemInfo"];
	    }
	}
	export class QueryResult {
	    total: number;
	    results: string[];
	
	    static createFrom(source: any = {}) {
	        return new QueryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.results = source["results"];
	    }
	}
	export class Settings {
	    defaultLimit: number;
	    defaultSort: string;
	    autosubmitQuery: boolean;
	    defaultExportDirectory: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultLimit = source["defaultLimit"];
	        this.defaultSort = source["defaultSort"];
	        this.autosubmitQuery = source["autosubmitQuery"];
	        this.defaultExportDirectory = source["defaultExportDirectory"];
	    }
	}

}

