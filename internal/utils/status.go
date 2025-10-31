package utils

func GetStatusBarData(percent int) []Bar {
	var statuses []Bar

	for i := range 50 {
		if i < percent/2 {
			statuses = append(statuses, Bar{Status: "ok"})
			continue
		}
		statuses = append(statuses, Bar{Status: "fail"})
	}

	return statuses
}

type Bar struct {
	Status string // "ok", "fail"
}
