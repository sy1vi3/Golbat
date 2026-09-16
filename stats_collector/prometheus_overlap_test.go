package stats_collector

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"golbat/geo"
)

func TestOverlappingAreasCountOnceGlobally(t *testing.T) {
	col := &promCollector{}
	for _, areas := range [][]geo.AreaName{
		{{Parent: "Michigan", Name: "Michigan"}, {Parent: "Detroit", Name: "Detroit"}, {Parent: "Detroit", Name: "Detroit"}},
		{{Parent: "world", Name: "world"}, {Parent: "Detroit", Name: "Detroit"}},
		nil,
	} {
		raid := counterValue(t, raidCount.WithLabelValues("world", "5"))
		incident := counterValue(t, incidentCount.WithLabelValues("world"))
		fort := counterValue(t, fortCount.WithLabelValues("world", "pokestop", "addition"))
		city := counterValue(t, raidCount.WithLabelValues("Detroit", "5"))
		col.UpdateRaidCount(areas, 5)
		col.UpdateIncidentCount(areas)
		col.UpdateFortCount(areas, "pokestop", "addition")
		for name, delta := range map[string]float64{
			"raid":     counterValue(t, raidCount.WithLabelValues("world", "5")) - raid,
			"incident": counterValue(t, incidentCount.WithLabelValues("world")) - incident,
			"fort":     counterValue(t, fortCount.WithLabelValues("world", "pokestop", "addition")) - fort,
		} {
			if delta != 1 {
				t.Errorf("%s global increment = %v; want 1", name, delta)
			}
		}
		if len(areas) > 0 && counterValue(t, raidCount.WithLabelValues("Detroit", "5"))-city != 1 {
			t.Error("regional series should increment once despite duplicate area names")
		}
	}
}

func counterValue(t *testing.T, c prometheus.Counter) float64 {
	t.Helper()
	var metric dto.Metric
	if err := c.Write(&metric); err != nil {
		t.Fatal(err)
	}
	return metric.GetCounter().GetValue()
}
