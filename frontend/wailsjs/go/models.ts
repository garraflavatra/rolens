export namespace app {
	
	export class Settings {
	    defaultLimit: number;
	    defaultSort: string;
	    autosubmitQuery: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultLimit = source["defaultLimit"];
	        this.defaultSort = source["defaultSort"];
	        this.autosubmitQuery = source["autosubmitQuery"];
	    }
	}
	export class findResult {
	    total: number;
	    results: string[];
	
	    static createFrom(source: any = {}) {
	        return new findResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.results = source["results"];
	    }
	}

}

