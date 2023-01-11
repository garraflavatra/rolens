export namespace main {
	
	export class findResult {
	    total: number;
	    results: any;
	
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

