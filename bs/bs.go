package main

import (
	sbbs "code.barbellmath.net/barbell-math/smoothbrain-bs"
)

func main() {
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterDepTargets()
	sbbs.RegisterGoMarkDocTargets(sbbs.NewReadmeOpts().
		SetRepoUrl("dummy").
		SetDirToRead("./logic").
		SetPostStages(sbbs.SedStage("dummy/blob/main", "", "README.md")),
	)
	sbbs.RegisterSqlcTargets()
	sbbs.RegisterGoEnumTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.AllGoTargets())
	sbbs.RegisterMergegateTarget(sbbs.NewMergegateTargets().
		SetPreStages(
			sbbs.TargetAsStage("install.goenum"),
			sbbs.TargetAsStage("install.sqlc"),
		),
	)
	sbbs.Main()
}
