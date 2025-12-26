package migrations

import (
	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

var (
	ModelSetupData = []dal.CreateModelsWithIDOpts{
		{
			ModelID: types.BarPathCalc,
			Name:    types.BarPathCalc.String(),
			Desc:    "Controls for the algorithims that calculate physics data from bar path time series position data.",
		},
		{
			ModelID: types.BarPathTracker,
			Name:    types.BarPathTracker.String(),
			Desc:    "Controls for the algorithims that generate bar path time series position data from a video.",
		},
	}

	BarPathCalcHyperparamsSetupData = []types.BarPathCalcHyperparams{
		{
			Version:         0,
			MinNumSamples:   100,
			TimeDeltaEps:    0.02,
			ApproxErr:       types.FourthOrder,
			NearZeroFilter:  0.1,
			NoiseFilter:     3,
			SmootherWeight1: 0.5,
			SmootherWeight2: 0.5,
			SmootherWeight3: 1,
			SmootherWeight4: 0.5,
			SmootherWeight5: 0.5,
		},
	}

	BarPathTrackerHyperparamsSetupData = []types.BarPathTrackerHyperparams{
		{
			Version:     0,
			MinLength:   5,
			MinFileSize: 5e7, // 50MB
			MaxFileSize: 5e8, // 500MB
		},
	}

	ExerciseFocusSetupData = []dal.CreateExerciseFocusWithIDOpts{
		{
			ExerciseFocus: types.UnknownExerciseFocus,
			Desc:          types.UnknownExerciseFocus.String(),
		},
		{
			ExerciseFocus: types.Squat,
			Desc:          types.Squat.String(),
		},
		{
			ExerciseFocus: types.Bench,
			Desc:          types.Bench.String(),
		},
		{
			ExerciseFocus: types.Deadlift,
			Desc:          types.Deadlift.String(),
		},
	}

	ExerciseKindSetupData = []dal.CreateExerciseKindWithIDOpts{
		{
			ExerciseKind: types.MainCompound,
			Name:         types.MainCompound.String(),
			Desc:         "The squat, bench, and deadlift.",
		},
		{
			ExerciseKind: types.MainCompoundAccessory,
			Name:         types.MainCompoundAccessory.String(),
			Desc:         "Variations of the squat, bench, and deadlift.",
		},
		{
			ExerciseKind: types.CompoundAccessory,
			Name:         types.CompoundAccessory.String(),
			Desc:         "Multi-joint accessories that are not part of the main compound accessory group.",
		},
		{
			ExerciseKind: types.Accessory,
			Name:         types.Accessory.String(),
			Desc:         "Single joint lifts and core work.",
		},
	}

	ExerciseSetupData = []types.IdWrapper[types.Exercise]{
		{
			Id: 1, Val: types.Exercise{
				Name:    "45 Degree Hyperextensions",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 2, Val: types.Exercise{
				Name:    "Banded Deadlift",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 3, Val: types.Exercise{
				Name:    "Banded Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 4, Val: types.Exercise{
				Name:    "Barbell Rows",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 5, Val: types.Exercise{
				Name:    "Bench",
				KindId:  types.MainCompound,
				FocusId: types.Bench,
			},
		},
		{
			Id: 6, Val: types.Exercise{
				Name:    "Block Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 7, Val: types.Exercise{
				Name:    "Box Jump",
				KindId:  types.Accessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 8, Val: types.Exercise{
				Name:    "Bulgarian Split Squat",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 9, Val: types.Exercise{
				Name:    "Cable Crunches",
				KindId:  types.Accessory,
				FocusId: types.UnknownExerciseFocus,
			},
		},
		{
			Id: 10, Val: types.Exercise{
				Name:    "Mid Cable Fly",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 11, Val: types.Exercise{
				Name:    "Cable Row",
				KindId:  types.Accessory,
				FocusId: types.UnknownExerciseFocus,
			},
		},
		{
			Id: 12, Val: types.Exercise{
				Name:    "Cannonball Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 13, Val: types.Exercise{
				Name:    "Close Grip Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 14, Val: types.Exercise{
				Name:    "Deadbugs",
				KindId:  types.Accessory,
				FocusId: types.UnknownExerciseFocus,
			},
		},
		{
			Id: 15, Val: types.Exercise{
				Name:    "Deadlift",
				KindId:  types.MainCompound,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 16, Val: types.Exercise{
				Name:    "Deadlift Row",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 17, Val: types.Exercise{
				Name:    "Deficit Romanian Deadlift",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 18, Val: types.Exercise{
				Name:    "Dip",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 19, Val: types.Exercise{
				Name:    "Dumbbell Bulgarian Split Squat",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 20, Val: types.Exercise{
				Name:    "Dumbbell Lateral Side Raise",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 21, Val: types.Exercise{
				Name:    "Dumbbell Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 22, Val: types.Exercise{
				Name:    "Dumbbell RDL",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 23, Val: types.Exercise{
				Name:    "Face Pull",
				KindId:  types.Accessory,
				FocusId: types.UnknownExerciseFocus,
			},
		},
		{
			Id: 24, Val: types.Exercise{
				Name:    "Flat Dumbbell Press",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 25, Val: types.Exercise{
				Name:    "Goblet Squat",
				KindId:  types.Accessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 26, Val: types.Exercise{
				Name:    "Halfway Paused Squats",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 27, Val: types.Exercise{
				Name:    "Hamstring Curls",
				KindId:  types.Accessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 28, Val: types.Exercise{
				Name:    "Heeled Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 29, Val: types.Exercise{
				Name:    "Incline Bench",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 30, Val: types.Exercise{
				Name:    "Incline Dumbbell Press",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 31, Val: types.Exercise{
				Name:    "Larson Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 32, Val: types.Exercise{
				Name:    "Lat Pulldown",
				KindId:  types.Accessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 33, Val: types.Exercise{
				Name:    "Narrow Grip Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 34, Val: types.Exercise{
				Name:    "Overhead Press",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 35, Val: types.Exercise{
				Name:    "Paused Deadlift",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 36, Val: types.Exercise{
				Name:    "Paused Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 37, Val: types.Exercise{
				Name:    "Pin Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 38, Val: types.Exercise{
				Name:    "Plank",
				KindId:  types.Accessory,
				FocusId: types.UnknownExerciseFocus,
			},
		},
		{
			Id: 39, Val: types.Exercise{
				Name:    "Pullup",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 40, Val: types.Exercise{
				Name:    "Pushup",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 41, Val: types.Exercise{
				Name:    "Romanian Deadlift",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 42, Val: types.Exercise{
				Name:    "Saftey Bar Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 43, Val: types.Exercise{
				Name:    "Seated Dumbbell Overhead Press",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 44, Val: types.Exercise{
				Name:    "Skull Crusher",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 45, Val: types.Exercise{
				Name:    "Squat",
				KindId:  types.MainCompound,
				FocusId: types.Squat,
			},
		},
		{
			Id: 46, Val: types.Exercise{
				Name:    "Squat Static Hold",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 47, Val: types.Exercise{
				Name:    "Tempo Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 48, Val: types.Exercise{
				Name:    "Tricep Pushdowns",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 49, Val: types.Exercise{
				Name:    "V Bar Rows",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 50, Val: types.Exercise{
				Name:    "YTW",
				KindId:  types.Accessory,
				FocusId: types.UnknownExerciseFocus,
			},
		},
		{
			Id: 51, Val: types.Exercise{
				Name:    "Low Cable Fly",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 52, Val: types.Exercise{
				Name:    "High Cable Fly",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 53, Val: types.Exercise{
				Name:    "Overhead Cable Tricep Extension",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 54, Val: types.Exercise{
				Name:    "Chest Supported Row",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 55, Val: types.Exercise{
				Name:    "Hack Squat",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 56, Val: types.Exercise{
				Name:    "Single Arm Overhead Cable Tricep Extension",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 57, Val: types.Exercise{
				Name:    "Cable Lateral Side Raise",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 58, Val: types.Exercise{
				Name:    "Heeled Saftey Bar Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 59, Val: types.Exercise{
				Name:    "Supinated Grip Lat Pulldown",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 60, Val: types.Exercise{
				Name:    "5-1-0 Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 61, Val: types.Exercise{
				Name:    "Spoto Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 62, Val: types.Exercise{
				Name:    "Single Arm Dumbbell Row",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 63, Val: types.Exercise{
				Name:    "Paused Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 64, Val: types.Exercise{
				Name:    "Duffalo Bar Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 65, Val: types.Exercise{
				Name:    "Quad Extension",
				KindId:  types.Accessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 66, Val: types.Exercise{
				Name:    "Seated Barbell Overhead Press",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 67, Val: types.Exercise{
				Name:    "Banded Pushup",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 68, Val: types.Exercise{
				Name:    "Single Leg Quad Extension",
				KindId:  types.Accessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 69, Val: types.Exercise{
				Name:    "Hyperextension",
				KindId:  types.Accessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 70, Val: types.Exercise{
				Name:    "Close Grip Larson Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 71, Val: types.Exercise{
				Name:    "Dumbbell Pec Fly",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 72, Val: types.Exercise{
				Name:    "Wall Sit",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 73, Val: types.Exercise{
				Name:    "Staggered Pushup",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 74, Val: types.Exercise{
				Name:    "Single Leg Hamstring Curl",
				KindId:  types.Accessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 75, Val: types.Exercise{
				Name:    "3-0-1 Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 76, Val: types.Exercise{
				Name:    "2-0-1 Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 77, Val: types.Exercise{
				Name:    "Single Leg Squat",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 78, Val: types.Exercise{
				Name:    "Reverse Hyperextension",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 79, Val: types.Exercise{
				Name:    "Hand Release Pushup",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 80, Val: types.Exercise{
				Name:    "Bison Bench Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 81, Val: types.Exercise{
				Name:    "T Bar Row",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 82, Val: types.Exercise{
				Name:    "Tricep Rollbacks",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 83, Val: types.Exercise{
				Name:    "Tibialis Raise",
				KindId:  types.Accessory,
				FocusId: types.UnknownExerciseFocus,
			},
		},
		{
			Id: 84, Val: types.Exercise{
				Name:    "Cossack Squat",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 85, Val: types.Exercise{
				Name:    "Reverse Nordic Curl",
				KindId:  types.Accessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 86, Val: types.Exercise{
				Name:    "Heeled Banded Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 87, Val: types.Exercise{
				Name:    "Pylo Pushups",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 88, Val: types.Exercise{
				Name:    "Cross Cable Tricep Extensions",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 89, Val: types.Exercise{
				Name:    "Lying Hamstring Curl",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 90, Val: types.Exercise{
				Name:    "0-1-0 Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 91, Val: types.Exercise{
				Name:    "Banded Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 92, Val: types.Exercise{
				Name:    "3-1-0 Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 93, Val: types.Exercise{
				Name:    "Walking Dumbbell Lunges",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 94, Val: types.Exercise{
				Name:    "3-3-0 Bench",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 95, Val: types.Exercise{
				Name:    "Sternal Pec Fly",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 96, Val: types.Exercise{
				Name:    "Behind The Back Row",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 97, Val: types.Exercise{
				Name:    "3-1-0 Larson Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 98, Val: types.Exercise{
				Name:    "Wide Grip Lat Pulldown",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 99, Val: types.Exercise{
				Name:    "0-1-0 Larson Press",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 100, Val: types.Exercise{
				Name:    "Single Leg Landmine Romanian Deadlift",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 101, Val: types.Exercise{
				Name:    "Single Arm Landmine Row",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 102, Val: types.Exercise{
				Name:    "Dumbbell Row",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 103, Val: types.Exercise{
				Name:    "Chest Supported Dumbbell Row",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 104, Val: types.Exercise{
				Name:    "Good Morning",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 105, Val: types.Exercise{
				Name:    "Kickstand Romanian Deadlift",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 106, Val: types.Exercise{
				Name:    "JM Press",
				KindId:  types.CompoundAccessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 107, Val: types.Exercise{
				Name:    "Belt Squat",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 108, Val: types.Exercise{
				Name:    "Close Grip Pullups",
				KindId:  types.CompoundAccessory,
				FocusId: types.Deadlift,
			},
		},
		{
			Id: 109, Val: types.Exercise{
				Name:    "Single Arm Tricep Pushdown",
				KindId:  types.Accessory,
				FocusId: types.Bench,
			},
		},
		{
			Id: 110, Val: types.Exercise{
				Name:    "Constant Tension SSB Squat",
				KindId:  types.MainCompoundAccessory,
				FocusId: types.Squat,
			},
		},
		{
			Id: 111, Val: types.Exercise{
				Name:    "Constant Tension Belt Squat",
				KindId:  types.CompoundAccessory,
				FocusId: types.Squat,
			},
		},
	}
)
