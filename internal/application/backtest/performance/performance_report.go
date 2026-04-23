package performance

type Report struct {
	TotalReturn float64
	MaxDrawdown float64
}

func GenerateReport(points []EquityPoint) Report {
	if len(points) == 0 {
		return Report{}
	}

	start := points[0].Equity
	end := points[len(points)-1].Equity

	totalReturn := (end - start) / start

	maxDD := calculateMaxDrawdown(points)

	return Report{
		TotalReturn: totalReturn,
		MaxDrawdown: maxDD,
	}
}

func calculateMaxDrawdown(points []EquityPoint) float64 {
	var peak float64
	var maxDD float64

	for _, p := range points {
		if p.Equity > peak {
			peak = p.Equity
		}

		dd := (peak - p.Equity) / peak

		if dd > maxDD {
			maxDD = dd
		}
	}
	return maxDD
}
