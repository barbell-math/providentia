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
	sbbs.RegisterCommonGoCmdTargets(sbbs.NewGoTargets().
		DefaultFmtTarget().
		DefaultGenerateTarget().
		// TODO - eventually replace with default target once old is deleted
		SetTestTarget(sbbs.DefaultTestTargetName, "-v", "./lib/logic/..."),
	)
	sbbs.RegisterMergegateTarget(sbbs.MergegateTargets{
		PreStages: []sbbs.StageFunc{
			sbbs.TargetAsStage("goenumInstall"),
			sbbs.TargetAsStage("sqlcInstall"),
		},
		CheckDepsUpdated:     true,
		CheckReadmeGomarkdoc: true,
		FmtTarget:            sbbs.DefaultFmtTargetName,
		TestTarget:           sbbs.DefaultTestTargetName,
		GenerateTarget:       sbbs.DefaultGenerateTargetName,
	})

	sbbs.Main("bs")
}
