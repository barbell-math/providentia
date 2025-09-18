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
			reflect.TypeFor[types.BarPathCalcConf]().Name():                               "barPathCalcConf",
			reflect.TypeFor[types.Vec2[types.Meter, types.Meter]]().Name():                "posVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec, types.MeterPerSec]]().Name():    "velVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec2, types.MeterPerSec2]]().Name():  "accVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec3, types.MeterPerSec3]]().Name():  "jerkVec2",
			reflect.TypeFor[types.Vec2[types.Newton, types.Newton]]().Name():              "forceVec2",
			reflect.TypeFor[types.Vec2[types.NewtonSec, types.NewtonSec]]().Name():        "impulseVec2",
			reflect.TypeFor[types.PointInTime[types.Second, types.MeterPerSec]]().Name():  "velPointInTime",
			reflect.TypeFor[types.PointInTime[types.Second, types.MeterPerSec2]]().Name(): "accPointInTime",
			reflect.TypeFor[types.PointInTime[types.Second, types.MeterPerSec3]]().Name(): "jerkPointInTime",
			reflect.TypeFor[types.PointInTime[types.Second, types.Newton]]().Name():       "newtonPointInTime",
			reflect.TypeFor[types.PointInTime[types.Second, types.NewtonSec]]().Name():    "newtonSecPointInTime",
			reflect.TypeFor[types.PointInTime[types.Second, types.Joule]]().Name():        "joulePointInTime",
			reflect.TypeFor[types.PointInTime[types.Second, types.Watt]]().Name():         "wattPointInTime",
			reflect.TypeFor[types.Split]().Name():                                         "split",
			reflect.TypeFor[barpathphysdata.Data]().Name():                                "barPathData",
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
	sbcgoglue.RegisterStruct[types.Vec2[types.Meter, types.Meter]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.MeterPerSec, types.MeterPerSec]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.MeterPerSec2, types.MeterPerSec2]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.MeterPerSec3, types.MeterPerSec3]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.Newton, types.Newton]](g)
	sbcgoglue.RegisterStruct[types.Vec2[types.NewtonSec, types.NewtonSec]](g)

	sbcgoglue.RegisterStruct[types.PointInTime[types.Second, types.MeterPerSec]](g)
	sbcgoglue.RegisterStruct[types.PointInTime[types.Second, types.MeterPerSec2]](g)
	sbcgoglue.RegisterStruct[types.PointInTime[types.Second, types.MeterPerSec3]](g)
	sbcgoglue.RegisterStruct[types.PointInTime[types.Second, types.Newton]](g)
	sbcgoglue.RegisterStruct[types.PointInTime[types.Second, types.NewtonSec]](g)
	sbcgoglue.RegisterStruct[types.PointInTime[types.Second, types.Joule]](g)
	sbcgoglue.RegisterStruct[types.PointInTime[types.Second, types.Watt]](g)

	sbcgoglue.RegisterStruct[types.Split](g)
	sbcgoglue.RegisterStruct[barpathphysdata.Data](g)
	sbcgoglue.RegisterStruct[types.BarPathCalcConf](g)
	g.WriteTo("./glue.h")
}
