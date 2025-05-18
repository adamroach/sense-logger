package sense

type DeviceDetails struct {
	Alerts   Alerts   `json:"alerts"`
	Device   Device   `json:"device"`
	Info     string   `json:"info"`
	Notes    *string  `json:"notes"`
	Timeline Timeline `json:"timeline"`
	Usage    Usage    `json:"usage"`
}

type Alerts struct {
	Allowed bool `json:"allowed"`
	Enabled bool `json:"enabled"`
}

type Timeline struct {
	Allowed bool `json:"allowed"`
	Visible bool `json:"visible"`
}

type Usage struct {
	AvgDuration      float64 `json:"avg_duration"`
	AvgMonthlyKWH    float64 `json:"avg_monthly_KWH"`
	AvgMonthlyCost   int     `json:"avg_monthly_cost"`
	AvgMonthlyPct    float64 `json:"avg_monthly_pct"`
	AvgMonthlyRuns   int     `json:"avg_monthly_runs"`
	AvgWatts         float64 `json:"avg_watts"`
	CurrentMonthKWH  float64 `json:"current_month_KWH"`
	CurrentMonthCost int     `json:"current_month_cost"`
	CurrentMonthRuns int     `json:"current_month_runs"`
	YearlyKWH        float64 `json:"yearly_KWH"`
	YearlyCost       int     `json:"yearly_cost"`
	YearlyText       string  `json:"yearly_text"`
}
