package main

import (
	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func main() {
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterUpdateDepsTarget()
	sbbs.RegisterGoMarkDocTargets()
	sbbs.RegisterSqlcTargets("./internal/db")
	sbbs.RegisterGoEnumTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.GoTargets{
		GenericTestTarget:     true,
		GenericBenchTarget:    true,
		GenericFmtTarget:      true,
		GenericGenerateTarget: true,
	})
	sbbs.RegisterMergegateTarget(sbbs.MergegateTargets{
		PreStages: []sbbs.StageFunc{
			sbbs.TargetAsStage("goenumInstall"),
			sbbs.TargetAsStage("sqlcInstall"),
		},
		CheckDepsUpdated:     true,
		CheckReadmeGomarkdoc: true,
		CheckFmt:             true,
		CheckUnitTests:       true,
		CheckGeneratedCode:   true,
	})

	sbbs.Main("bs")
}
