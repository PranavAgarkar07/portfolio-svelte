package analytics

type SessionPayload struct {
	ID       string `json:"id"`
	IPHash   string `json:"ip_hash"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Referrer string `json:"referrer"`
	Source   string `json:"source"`
	Device   string `json:"device"`
	OS       string `json:"os"`
	Browser  string `json:"browser"`
	Theme    string `json:"theme"`
}

type EventPayload struct {
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
	Target    string `json:"target"`
	Value     string `json:"value"`
	Ts        int64  `json:"ts"`
}

type ReferrerCount struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

type TargetCount struct {
	Target string `json:"target"`
	Count  int    `json:"count"`
}

type CountryCount struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

type SourceCount struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

type DashboardStats struct {
	TotalViews      int             `json:"total_views"`
	UniqueVisitors  int             `json:"unique_visitors"`
	TopReferrers    []ReferrerCount `json:"top_referrers"`
	TopTargets      []TargetCount   `json:"top_targets"`
	CountryBreakdown []CountryCount `json:"country_breakdown"`
	SourceBreakdown []SourceCount   `json:"source_breakdown"`
	AvgTimeOnSite   float64         `json:"avg_time_on_site"`
	ResumeDownloads int             `json:"resume_downloads"`
	FormSubmissions int             `json:"form_submissions"`
}
