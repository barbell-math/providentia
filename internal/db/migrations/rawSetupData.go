package migrations

import (
	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	"github.com/barbell-math/providentia/lib/types"
)

var (
	ModelSetupData = []dal.BulkCreateModelsWithIDParams{
		{
			ID:          int32(types.SimplifiedNegativeSpace),
			Name:        "SimplifiedNegativeSpaceModel",
			Description: "A simplified version of the fill negative space model that can use linear regression.",
		},
	}

	ExerciseFocusSetupData = []dal.BulkCreateExerciseFocusWithIDParams{
		{
			ID:    int32(types.UnknownExerciseFocus),
			Focus: types.UnknownExerciseFocus.String(),
		},
		{
			ID:    int32(types.Squat),
			Focus: types.Squat.String(),
		},
		{
			ID:    int32(types.Bench),
			Focus: types.Bench.String(),
		},
		{
			ID:    int32(types.Deadlift),
			Focus: types.Deadlift.String(),
		},
	}

	ExerciseKindSetupData = []dal.BulkCreateExerciseKindWithIDParams{
		{
			ID:          int32(types.MainCompound),
			Kind:        types.MainCompound.String(),
			Description: "The squat, bench, and deadlift.",
		},
		{
			ID:          int32(types.MainCompoundAccessory),
			Kind:        types.MainCompoundAccessory.String(),
			Description: "Variations of the squat, bench, and deadlift.",
		},
		{
			ID:          int32(types.CompoundAccessory),
			Kind:        types.CompoundAccessory.String(),
			Description: "Multi-joint accessories that are not part of the main compound accessory group.",
		},
		{
			ID:          int32(types.Accessory),
			Kind:        types.Accessory.String(),
			Description: "Single joint lifts and core work.",
		},
	}

	ExerciseSetupData = []dal.BulkCreateExerciseWithIDParams{
		{
			ID: 1, Name: "45 Degree Hyperextensions",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 2, Name: "Banded Deadlift",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 3, Name: "Banded Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 4, Name: "Barbell Rows",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 5, Name: "Bench",
			KindID:  types.MainCompound,
			FocusID: types.Bench,
		},
		{
			ID: 6, Name: "Block Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 7, Name: "Box Jump",
			KindID:  types.Accessory,
			FocusID: types.Squat,
		},
		{
			ID: 8, Name: "Bulgarian Split Squat",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 9, Name: "Cable Crunches",
			KindID:  types.Accessory,
			FocusID: types.UnknownExerciseFocus,
		},
		{
			ID: 10, Name: "Mid Cable Fly",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 11, Name: "Cable Row",
			KindID:  types.Accessory,
			FocusID: types.UnknownExerciseFocus,
		},
		{
			ID: 12, Name: "Cannonball Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 13, Name: "Close Grip Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 14, Name: "Deadbugs",
			KindID:  types.Accessory,
			FocusID: types.UnknownExerciseFocus,
		},
		{
			ID: 15, Name: "Deadlift",
			KindID:  types.MainCompound,
			FocusID: types.Deadlift,
		},
		{
			ID: 16, Name: "Deadlift Row",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 17, Name: "Deficit Romanian Deadlift",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 18, Name: "Dip",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 19, Name: "Dumbbell Bulgarian Split Squat",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 20, Name: "Dumbbell Lateral Side Raise",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 21, Name: "Dumbbell Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 22, Name: "Dumbbell RDL",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 23, Name: "Face Pull",
			KindID:  types.Accessory,
			FocusID: types.UnknownExerciseFocus,
		},
		{
			ID: 24, Name: "Flat Dumbbell Press",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 25, Name: "Goblet Squat",
			KindID:  types.Accessory,
			FocusID: types.Squat,
		},
		{
			ID: 26, Name: "Halfway Paused Squats",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 27, Name: "Hamstring Curls",
			KindID:  types.Accessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 28, Name: "Heeled Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 29, Name: "Incline Bench",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 30, Name: "Incline Dumbbell Press",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 31, Name: "Larson Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 32, Name: "Lat Pulldown",
			KindID:  types.Accessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 33, Name: "Narrow Grip Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 34, Name: "Overhead Press",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 35, Name: "Paused Deadlift",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 36, Name: "Paused Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 37, Name: "Pin Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 38, Name: "Plank",
			KindID:  types.Accessory,
			FocusID: types.UnknownExerciseFocus,
		},
		{
			ID: 39, Name: "Pullup",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 40, Name: "Pushup",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 41, Name: "Romanian Deadlift",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 42, Name: "Saftey Bar Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 43, Name: "Seated Dumbbell Overhead Press",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 44, Name: "Skull Crusher",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 45, Name: "Squat",
			KindID:  types.MainCompound,
			FocusID: types.Squat,
		},
		{
			ID: 46, Name: "Squat Static Hold",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 47, Name: "Tempo Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 48, Name: "Tricep Pushdowns",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 49, Name: "V Bar Rows",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 50, Name: "YTW",
			KindID:  types.Accessory,
			FocusID: types.UnknownExerciseFocus,
		},
		{
			ID: 51, Name: "Low Cable Fly",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 52, Name: "High Cable Fly",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 53, Name: "Overhead Cable Tricep Extension",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 54, Name: "Chest Supported Row",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 55, Name: "Hack Squat",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 56, Name: "Single Arm Overhead Cable Tricep Extension",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 57, Name: "Cable Lateral Side Raise",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 58, Name: "Heeled Saftey Bar Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 59, Name: "Supinated Grip Lat Pulldown",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 60, Name: "5-1-0 Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 61, Name: "Spoto Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 62, Name: "Single Arm Dumbbell Row",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 63, Name: "Paused Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 64, Name: "Duffalo Bar Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 65, Name: "Quad Extension",
			KindID:  types.Accessory,
			FocusID: types.Squat,
		},
		{
			ID: 66, Name: "Seated Barbell Overhead Press",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 67, Name: "Banded Pushup",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 68, Name: "Single Leg Quad Extension",
			KindID:  types.Accessory,
			FocusID: types.Squat,
		},
		{
			ID: 69, Name: "Hyperextension",
			KindID:  types.Accessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 70, Name: "Close Grip Larson Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 71, Name: "Dumbbell Pec Fly",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 72, Name: "Wall Sit",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 73, Name: "Staggered Pushup",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 74, Name: "Single Leg Hamstring Curl",
			KindID:  types.Accessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 75, Name: "3-0-1 Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 76, Name: "2-0-1 Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 77, Name: "Single Leg Squat",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 78, Name: "Reverse Hyperextension",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 79, Name: "Hand Release Pushup",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 80, Name: "Bison Bench Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 81, Name: "T Bar Row",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 82, Name: "Tricep Rollbacks",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 83, Name: "Tibialis Raise",
			KindID:  types.Accessory,
			FocusID: types.UnknownExerciseFocus,
		},
		{
			ID: 84, Name: "Cossack Squat",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 85, Name: "Reverse Nordic Curl",
			KindID:  types.Accessory,
			FocusID: types.Squat,
		},
		{
			ID: 86, Name: "Heeled Banded Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 87, Name: "Pylo Pushups",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 88, Name: "Cross Cable Tricep Extensions",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 89, Name: "Lying Hamstring Curl",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 90, Name: "0-1-0 Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 91, Name: "Banded Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 92, Name: "3-1-0 Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 93, Name: "Walking Dumbbell Lunges",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 94, Name: "3-3-0 Bench",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 95, Name: "Sternal Pec Fly",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 96, Name: "Behind The Back Row",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 97, Name: "3-1-0 Larson Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 98, Name: "Wide Grip Lat Pulldown",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 99, Name: "0-1-0 Larson Press",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 100, Name: "Single Leg Landmine Romanian Deadlift",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 101, Name: "Single Arm Landmine Row",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 102, Name: "Dumbbell Row",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 103, Name: "Chest Supported Dumbbell Row",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 104, Name: "Good Morning",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 105, Name: "Kickstand Romanian Deadlift",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 106, Name: "JM Press",
			KindID:  types.CompoundAccessory,
			FocusID: types.Bench,
		},
		{
			ID: 107, Name: "Belt Squat",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 108, Name: "Close Grip Pullups",
			KindID:  types.CompoundAccessory,
			FocusID: types.Deadlift,
		},
		{
			ID: 109, Name: "Single Arm Tricep Pushdown",
			KindID:  types.Accessory,
			FocusID: types.Bench,
		},
		{
			ID: 110, Name: "Constant Tension SSB Squat",
			KindID:  types.MainCompoundAccessory,
			FocusID: types.Squat,
		},
		{
			ID: 111, Name: "Constant Tension Belt Squat",
			KindID:  types.CompoundAccessory,
			FocusID: types.Squat,
		},
	}
)
