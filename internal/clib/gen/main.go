package main

import (
	"reflect"

	barpathphysdata "code.barbellmath.net/barbell-math/providentia/internal/models/barPathPhysData"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcgoglue "code.barbellmath.net/barbell-math/smoothbrain-cgoGlue"
)

func main() {
	g := sbcgoglue.New(sbcgoglue.Opts{
		ExitOnErr: true,
		Rename: map[string]string{
			reflect.TypeFor[types.BarPathCalcConf]().Name():          "barPathCalcConf",
			reflect.TypeFor[types.Vec2[types.Meter]]().Name():        "posVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec]]().Name():  "velVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec2]]().Name(): "accVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec3]]().Name(): "jerkVec2",
			reflect.TypeFor[types.Vec2[types.Newton]]().Name():       "forceVec2",
			reflect.TypeFor[types.Vec2[types.NewtonSec]]().Name():    "impulseVec2",
			reflect.TypeFor[types.Split]().Name():                    "split",
			reflect.TypeFor[barpathphysdata.Data]().Name():           "barPathData",
		},
	})
	sbcgoglue.RegisterEnum(
		g,
		types.ApproximationErrorNames(),
		types.ApproximationErrorValues(),
	)
	sbcgoglue.RegisterEnum(
		g,
		barpathphysdata.BarPathCalcErrCodeNames(),
		barpathphysdata.BarPathCalcErrCodeValues(),
	)
	sbcgoglue.RegisterStruct[types.Vec2[types.Meter]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.MeterPerSec]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.MeterPerSec2]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.MeterPerSec3]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.Newton]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.NewtonSec]](g)
	sbcgoglue.RegisterStruct[types.Split](g)
	sbcgoglue.RegisterStruct[barpathphysdata.Data](g)
	sbcgoglue.RegisterStruct[types.BarPathCalcConf](g)
	g.WriteTo("./glue.h")
}
