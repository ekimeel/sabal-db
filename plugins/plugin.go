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
