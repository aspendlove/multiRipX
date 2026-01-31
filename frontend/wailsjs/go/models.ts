export namespace config {
	
	export class CD {
	    Musicbrainz: string;
	    Offset: number;
	
	    static createFrom(source: any = {}) {
	        return new CD(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Musicbrainz = source["Musicbrainz"];
	        this.Offset = source["Offset"];
	    }
	}
	export class Movie {
	    Name: string;
	    Title: number;
	
	    static createFrom(source: any = {}) {
	        return new Movie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Title = source["Title"];
	    }
	}
	export class Show {
	    Name: string;
	    Season: number;
	    Episode: number;
	    Title: number;
	
	    static createFrom(source: any = {}) {
	        return new Show(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Season = source["Season"];
	        this.Episode = source["Episode"];
	        this.Title = source["Title"];
	    }
	}
	export class JobDefinition {
	    Drive: string;
	    DiscType: string;
	    OutputDir: string;
	    Shows: Show[];
	    Movies: Movie[];
	    CD: CD[];
	
	    static createFrom(source: any = {}) {
	        return new JobDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Drive = source["Drive"];
	        this.DiscType = source["DiscType"];
	        this.OutputDir = source["OutputDir"];
	        this.Shows = this.convertValues(source["Shows"], Show);
	        this.Movies = this.convertValues(source["Movies"], Movie);
	        this.CD = this.convertValues(source["CD"], CD);
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
	    OutputDir: string;
	    Jobs: JobDefinition[];
	
	    static createFrom(source: any = {}) {
	        return new JobsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OutputDir = source["OutputDir"];
	        this.Jobs = this.convertValues(source["Jobs"], JobDefinition);
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

