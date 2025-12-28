export namespace config {
	
	export class Movie {
	    name: string;
	    title: number;
	
	    static createFrom(source: any = {}) {
	        return new Movie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.title = source["title"];
	    }
	}
	export class Show {
	    name: string;
	    season: number;
	    episode: number;
	    title: number;
	
	    static createFrom(source: any = {}) {
	        return new Show(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.season = source["season"];
	        this.episode = source["episode"];
	        this.title = source["title"];
	    }
	}
	export class JobDefinition {
	    drive: string;
	    discType: string;
	    outputDir: string;
	    shows: Show[];
	    movies: Movie[];
	
	    static createFrom(source: any = {}) {
	        return new JobDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.drive = source["drive"];
	        this.discType = source["discType"];
	        this.outputDir = source["outputDir"];
	        this.shows = this.convertValues(source["shows"], Show);
	        this.movies = this.convertValues(source["movies"], Movie);
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
	export class JobsConfig {
	    outputDir: string;
	    jobs: JobDefinition[];
	
	    static createFrom(source: any = {}) {
	        return new JobsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outputDir = source["outputDir"];
	        this.jobs = this.convertValues(source["jobs"], JobDefinition);
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
	

}

