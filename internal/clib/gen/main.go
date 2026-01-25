package main

import (
	"fmt"
	"reflect"

	barpathphysdata "code.barbellmath.net/barbell-math/providentia/internal/models/barPathPhysData"
	barpathtracker "code.barbellmath.net/barbell-math/providentia/internal/models/barPathTracker"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	testgen "code.barbellmath.net/barbell-math/smoothbrain-cgoGlue/testGen"
	typegen "code.barbellmath.net/barbell-math/smoothbrain-cgoGlue/typeGen"
)

func main() {
	testgen.Generate(&testgen.Opts{
		ExitOnErr:       true,
		SearchPath:      []string{"."},
		OutputPath:      ".",
		HeaderGuardName: "CLIB",
		CXXFlags:        []string{"-Wall", "-march=native", "-std=c++23"},
		LDFlags:         []string{"-lstdc++"},
		AddAssertHeader: true,
		GoPackage:       "clib",
	})

	g := typegen.New(typegen.Opts{
		ExitOnErr: true,
		Rename: map[string]string{
			reflect.TypeFor[types.BarPathCalcHyperparams]().Name():                        "barPathCalcHyperparams",
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
			reflect.TypeFor[barpathphysdata.CData]().Name():                               "barPathData",
		},
	})
	typegen.RegisterEnum(
		g,
		types.ApproximationErrorNames(),
		types.ApproximationErrorValues(),
	)
	typegen.RegisterEnum(
		g,
		barpathphysdata.BarPathCalcErrCodeNames(),
		barpathphysdata.BarPathCalcErrCodeValues(),
	)
	typegen.RegisterEnum(
		g,
		barpathtracker.BarPathTrackerErrCodeNames(),
		barpathtracker.BarPathTrackerErrCodeValues(),
	)
	typegen.RegisterStruct[types.Vec2[types.Meter, types.Meter]](g)
	typegen.RegisterStruct[types.Vec2[types.MeterPerSec, types.MeterPerSec]](g)
	typegen.RegisterStruct[types.Vec2[types.MeterPerSec2, types.MeterPerSec2]](g)
	typegen.RegisterStruct[types.Vec2[types.MeterPerSec3, types.MeterPerSec3]](g)
	typegen.RegisterStruct[types.Vec2[types.Newton, types.Newton]](g)
	typegen.RegisterStruct[types.Vec2[types.NewtonSec, types.NewtonSec]](g)

	typegen.RegisterStruct[types.PointInTime[types.Second, types.MeterPerSec]](g)
	typegen.RegisterStruct[types.PointInTime[types.Second, types.MeterPerSec2]](g)
	typegen.RegisterStruct[types.PointInTime[types.Second, types.MeterPerSec3]](g)
	typegen.RegisterStruct[types.PointInTime[types.Second, types.Newton]](g)
	typegen.RegisterStruct[types.PointInTime[types.Second, types.NewtonSec]](g)
	typegen.RegisterStruct[types.PointInTime[types.Second, types.Joule]](g)
	typegen.RegisterStruct[types.PointInTime[types.Second, types.Watt]](g)

	typegen.RegisterStruct[types.Split](g)
	typegen.RegisterStruct[barpathphysdata.CData](g)
	typegen.RegisterStruct[types.BarPathCalcHyperparams](g)
	g.WriteTo("./glue.h")

	fmt.Println("Generated cgo-glue")
}
