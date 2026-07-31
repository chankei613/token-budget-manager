export namespace api {
	
	export class BudgetStatus {
	    budget: db.Budget;
	    tokens_used: number;
	    cost_used_usd: number;
	    percent_used: number;
	    alert_level: string;
	    // Go type: time
	    period_start: any;
	
	    static createFrom(source: any = {}) {
	        return new BudgetStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.budget = this.convertValues(source["budget"], db.Budget);
	        this.tokens_used = source["tokens_used"];
	        this.cost_used_usd = source["cost_used_usd"];
	        this.percent_used = source["percent_used"];
	        this.alert_level = source["alert_level"];
	        this.period_start = this.convertValues(source["period_start"], null);
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
	export class DailyUsagePoint {
	    date: string;
	    cost_usd: number;
	    tokens: number;
	
	    static createFrom(source: any = {}) {
	        return new DailyUsagePoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.cost_usd = source["cost_usd"];
	        this.tokens = source["tokens"];
	    }
	}
	export class IssueKeyResult {
	    id: string;
	    name: string;
	    api_key: string;
	
	    static createFrom(source: any = {}) {
	        return new IssueKeyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.api_key = source["api_key"];
	    }
	}
	export class ListUsageResult {
	    items: db.UsageEvent[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ListUsageResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], db.UsageEvent);
	        this.total = source["total"];
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
	export class UsageSummary {
	    total_input_tokens: number;
	    total_output_tokens: number;
	    total_cost_usd: number;
	    by_agent_cost_usd: Record<string, number>;
	    by_model_cost_usd: Record<string, number>;
	    daily_series: DailyUsagePoint[];
	
	    static createFrom(source: any = {}) {
	        return new UsageSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_input_tokens = source["total_input_tokens"];
	        this.total_output_tokens = source["total_output_tokens"];
	        this.total_cost_usd = source["total_cost_usd"];
	        this.by_agent_cost_usd = source["by_agent_cost_usd"];
	        this.by_model_cost_usd = source["by_model_cost_usd"];
	        this.daily_series = this.convertValues(source["daily_series"], DailyUsagePoint);
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

export namespace db {
	
	export class AgentKey {
	    id: string;
	    name: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    revoked_at?: any;
	
	    static createFrom(source: any = {}) {
	        return new AgentKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.revoked_at = this.convertValues(source["revoked_at"], null);
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
	export class Budget {
	    id: string;
	    name: string;
	    source: string;
	    scope_key: string;
	    period: string;
	    limit_tokens: number;
	    limit_usd: number;
	    last_alert_level: string;
	    // Go type: time
	    alert_period_start: any;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Budget(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.source = source["source"];
	        this.scope_key = source["scope_key"];
	        this.period = source["period"];
	        this.limit_tokens = source["limit_tokens"];
	        this.limit_usd = source["limit_usd"];
	        this.last_alert_level = source["last_alert_level"];
	        this.alert_period_start = this.convertValues(source["alert_period_start"], null);
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class ModelPricing {
	    model_id: string;
	    input_price_per_1m: number;
	    output_price_per_1m: number;
	
	    static createFrom(source: any = {}) {
	        return new ModelPricing(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model_id = source["model_id"];
	        this.input_price_per_1m = source["input_price_per_1m"];
	        this.output_price_per_1m = source["output_price_per_1m"];
	    }
	}
	export class UsageEvent {
	    id: string;
	    // Go type: time
	    received_at: any;
	    source: string;
	    agent_id: string;
	    scope_key: string;
	    provider: string;
	    model_id: string;
	    input_tokens: number;
	    output_tokens: number;
	    cost_usd: number;
	
	    static createFrom(source: any = {}) {
	        return new UsageEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.received_at = this.convertValues(source["received_at"], null);
	        this.source = source["source"];
	        this.agent_id = source["agent_id"];
	        this.scope_key = source["scope_key"];
	        this.provider = source["provider"];
	        this.model_id = source["model_id"];
	        this.input_tokens = source["input_tokens"];
	        this.output_tokens = source["output_tokens"];
	        this.cost_usd = source["cost_usd"];
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

