package plugins

import "github.com/ekimeel/sabal-db/pkg/services"

type Plugin interface {
	Run(env *Environment)
	Install()
}

type Environment struct {
	MetricService services.MetricService
	EquipService  services.EquipService
	PointService  services.PointService
}

func weightedAverage(value1, weight1, value2, weight2 float64) float64 {
	return (value1*weight1 + value2*weight2) / (weight1 + weight2)
}
