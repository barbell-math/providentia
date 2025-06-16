package migrations

import (
	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	pubenums "github.com/barbell-math/providentia/lib/pubEnums"
)

var (
	VideoDataSetupData = []dal.BulkCreateVideoDataWithIDParams{
		{
			ID:           0,
			Path:         "",
			Position:     [][]float64{},
			Velocity:     [][]float64{},
			Acceleration: [][]float64{},
			Force:        [][]float64{},
			Impulse:      [][]float64{},
		},
	}

	ModelSetupData = []dal.BulkCreateModelsParams{
		{
			ID:          int32(pubenums.ModelIDSimplifiedNegativeSpace),
			Name:        "SimplifiedNegativeSpaceModel",
			Description: "A simplified version of the fill negative space model that can use linear regression.",
		},
	}

	ExerciseFocusSetupData = []dal.BulkCreateExerciseFocusWithIDParams{
		{
			ID:    int32(pubenums.ExerciseFocusUnknown),
			Focus: pubenums.ExerciseFocusUnknown.String(),
		},
		{
			ID:    int32(pubenums.ExerciseFocusSquat),
			Focus: pubenums.ExerciseFocusSquat.String(),
		},
		{
			ID:    int32(pubenums.ExerciseFocusBench),
			Focus: pubenums.ExerciseFocusBench.String(),
		},
		{
			ID:    int32(pubenums.ExerciseFocusDeadlift),
			Focus: pubenums.ExerciseFocusDeadlift.String(),
		},
	}

	ExerciseKindSetupData = []dal.BulkCreateExerciseKindWithIDParams{
		{
			ID:          int32(pubenums.ExerciseKindMainCompound),
			Kind:        pubenums.ExerciseKindMainCompound.String(),
			Description: "The squat, bench, and deadlift.",
		},
		{
			ID:          int32(pubenums.ExerciseKindMainCompoundAccessory),
			Kind:        pubenums.ExerciseKindMainCompoundAccessory.String(),
			Description: "Variations of the squat, bench, and deadlift.",
		},
		{
			ID:          int32(pubenums.ExerciseKindCompoundAccessory),
			Kind:        pubenums.ExerciseKindCompoundAccessory.String(),
			Description: "Multi-joint accessories that are not part of the main compound accessory group.",
		},
		{
			ID:          int32(pubenums.ExerciseKindAccessory),
			Kind:        pubenums.ExerciseKindAccessory.String(),
			Description: "Single joint lifts and core work.",
		},
	}

	ExerciseSetupData = []dal.BulkCreateExerciseWithIDParams{
		{
			ID: 1, Name: "45 Degree Hyperextensions",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 2, Name: "Banded Deadlift",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 3, Name: "Banded Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 4, Name: "Barbell Rows",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 5, Name: "Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompound),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 6, Name: "Block Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 7, Name: "Box Jump",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 8, Name: "Bulgarian Split Squat",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 9, Name: "Cable Crunches",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusUnknown),
		},
		{
			ID: 10, Name: "Mid Cable Fly",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 11, Name: "Cable Row",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusUnknown),
		},
		{
			ID: 12, Name: "Cannonball Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 13, Name: "Close Grip Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 14, Name: "Deadbugs",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusUnknown),
		},
		{
			ID: 15, Name: "Deadlift",
			KindID:  int32(pubenums.ExerciseKindMainCompound),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 16, Name: "Deadlift Row",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 17, Name: "Deficit Romanian Deadlift",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 18, Name: "Dip",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 19, Name: "Dumbbell Bulgarian Split Squat",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 20, Name: "Dumbbell Lateral Side Raise",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 21, Name: "Dumbbell Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 22, Name: "Dumbbell RDL",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 23, Name: "Face Pull",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusUnknown),
		},
		{
			ID: 24, Name: "Flat Dumbbell Press",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 25, Name: "Goblet Squat",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 26, Name: "Halfway Paused Squats",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 27, Name: "Hamstring Curls",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 28, Name: "Heeled Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 29, Name: "Incline Bench",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 30, Name: "Incline Dumbbell Press",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 31, Name: "Larson Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 32, Name: "Lat Pulldown",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 33, Name: "Narrow Grip Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 34, Name: "Overhead Press",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 35, Name: "Paused Deadlift",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 36, Name: "Paused Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 37, Name: "Pin Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 38, Name: "Plank",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusUnknown),
		},
		{
			ID: 39, Name: "Pullup",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 40, Name: "Pushup",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 41, Name: "Romanian Deadlift",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 42, Name: "Saftey Bar Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 43, Name: "Seated Dumbbell Overhead Press",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 44, Name: "Skull Crusher",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 45, Name: "Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompound),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 46, Name: "Squat Static Hold",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 47, Name: "Tempo Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 48, Name: "Tricep Pushdowns",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 49, Name: "V Bar Rows",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 50, Name: "YTW",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusUnknown),
		},
		{
			ID: 51, Name: "Low Cable Fly",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 52, Name: "High Cable Fly",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 53, Name: "Overhead Cable Tricep Extension",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 54, Name: "Chest Supported Row",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 55, Name: "Hack Squat",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 56, Name: "Single Arm Overhead Cable Tricep Extension",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 57, Name: "Cable Lateral Side Raise",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 58, Name: "Heeled Saftey Bar Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 59, Name: "Supinated Grip Lat Pulldown",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 60, Name: "5-1-0 Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 61, Name: "Spoto Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 62, Name: "Single Arm Dumbbell Row",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 63, Name: "Paused Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 64, Name: "Duffalo Bar Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 65, Name: "Quad Extension",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 66, Name: "Seated Barbell Overhead Press",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 67, Name: "Banded Pushup",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 68, Name: "Single Leg Quad Extension",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 69, Name: "Hyperextension",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 70, Name: "Close Grip Larson Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 71, Name: "Dumbbell Pec Fly",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 72, Name: "Wall Sit",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 73, Name: "Staggered Pushup",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 74, Name: "Single Leg Hamstring Curl",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 75, Name: "3-0-1 Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 76, Name: "2-0-1 Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 77, Name: "Single Leg Squat",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 78, Name: "Reverse Hyperextension",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 79, Name: "Hand Release Pushup",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 80, Name: "Bison Bench Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 81, Name: "T Bar Row",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 82, Name: "Tricep Rollbacks",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 83, Name: "Tibialis Raise",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusUnknown),
		},
		{
			ID: 84, Name: "Cossack Squat",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 85, Name: "Reverse Nordic Curl",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 86, Name: "Heeled Banded Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 87, Name: "Pylo Pushups",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 88, Name: "Cross Cable Tricep Extensions",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 89, Name: "Lying Hamstring Curl",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 90, Name: "0-1-0 Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 91, Name: "Banded Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 92, Name: "3-1-0 Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 93, Name: "Walking Dumbbell Lunges",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 94, Name: "3-3-0 Bench",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 95, Name: "Sternal Pec Fly",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 96, Name: "Behind The Back Row",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 97, Name: "3-1-0 Larson Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 98, Name: "Wide Grip Lat Pulldown",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 99, Name: "0-1-0 Larson Press",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 100, Name: "Single Leg Landmine Romanian Deadlift",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 101, Name: "Single Arm Landmine Row",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 102, Name: "Dumbbell Row",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 103, Name: "Chest Supported Dumbbell Row",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 104, Name: "Good Morning",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 105, Name: "Kickstand Romanian Deadlift",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 106, Name: "JM Press",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 107, Name: "Belt Squat",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 108, Name: "Close Grip Pullups",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusDeadlift),
		},
		{
			ID: 109, Name: "Single Arm Tricep Pushdown",
			KindID:  int32(pubenums.ExerciseKindAccessory),
			FocusID: int32(pubenums.ExerciseFocusBench),
		},
		{
			ID: 110, Name: "Constant Tension SSB Squat",
			KindID:  int32(pubenums.ExerciseKindMainCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
		{
			ID: 111, Name: "Constant Tension Belt Squat",
			KindID:  int32(pubenums.ExerciseKindCompoundAccessory),
			FocusID: int32(pubenums.ExerciseFocusSquat),
		},
	}
)
