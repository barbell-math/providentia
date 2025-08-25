package main

import (
	"reflect"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	structgen "code.barbellmath.net/barbell-math/smoothbrain-cgostructgen"
)

func main() {
	sg := structgen.New(structgen.Opts{
		ExitOnErr: true,
		StructRename: map[string]string{
			reflect.TypeFor[types.PhysicsDataConf]().Name():          "physDataConf",
			reflect.TypeFor[types.BarPathCalcConf]().Name():          "barPathCalcConf",
			reflect.TypeFor[types.Vec2[types.Meter]]().Name():        "posVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec]]().Name():  "velVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec2]]().Name(): "accVec2",
			reflect.TypeFor[types.Vec2[types.MeterPerSec3]]().Name(): "jerkVec2",
			reflect.TypeFor[types.Vec2[types.Newton]]().Name():       "forceVec2",
			reflect.TypeFor[types.Vec2[types.NewtonSec]]().Name():    "impulseVec2",
			reflect.TypeFor[types.Vec2[types.Joule]]().Name():        "workVec2",
		},
	})
	structgen.GenerateFor[types.BarPathCalcConf](sg)
	structgen.GenerateFor[types.PhysicsDataConf](sg)
	structgen.GenerateFor[types.Vec2[types.Meter]](sg)
	structgen.GenerateFor[types.Vec2[types.MeterPerSec]](sg)
	structgen.GenerateFor[types.Vec2[types.MeterPerSec2]](sg)
	structgen.GenerateFor[types.Vec2[types.MeterPerSec3]](sg)
	structgen.GenerateFor[types.Vec2[types.Newton]](sg)
	structgen.GenerateFor[types.Vec2[types.NewtonSec]](sg)
	structgen.GenerateFor[types.Vec2[types.Joule]](sg)
	sg.WriteTo("./cgoStructs.h", "CGO_STRUCTS")
}
