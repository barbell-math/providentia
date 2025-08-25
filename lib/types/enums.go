package types

//go:generate go-enum --marshal --names --values --nocase --noprefix

type (
	// ENUM(UnknownExerciseFocus, Squat, Bench, Deadlift)
	ExerciseFocus int32

	// ENUM(
	//	MainCompound = 1
	//	MainCompoundAccessory
	//	CompoundAccessory
	//	Accessory
	// )
	ExerciseKind int32

	// ENUM(
	//	UnknownModel,
	//	SimplifiedNegativeSpace,
	// )
	ModelID int32

	// ENUM(
	//	SecondOrder = 2
	//	FourthOrder = 4
	// )
	ApproximationError int32
)
