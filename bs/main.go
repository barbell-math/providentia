package main

import (
	"code.barbellmath.net/barbell-math/providentia/bs/prov_exportedbs"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbbs"
)

func main() {
	targets := sbbs.NewTargets[prov_exportedbs.Conf]()

	sbbs.NewDefaultStages[prov_exportedbs.Conf]().
		SetFmtStages(prov_exportedbs.FmtStages...).
		SetGenStages(prov_exportedbs.GenStages...).
		SetTestStages(prov_exportedbs.TestStages...).
		SetUpdatePkgsTarget(prov_exportedbs.UpdateDepsStages...).
		Register(targets)
	targets.RegisterTarget(
		sbbs.DefaultMergegateTargetName,
		prov_exportedbs.MergegateStages...,
	)
	targets.RegisterTarget("dev", prov_exportedbs.DevSetup...)

	sbbs.Main(targets, func(args ...string) (prov_exportedbs.Conf, error) {
		return prov_exportedbs.NewConf(), nil
	})
}
