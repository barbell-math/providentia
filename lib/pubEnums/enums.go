package pubenums

//go:generate go-enum --marshal --names --values --nocase

type (
	// ENUM(Unknown, Squat, Bench, Deadlift)
	ExerciseFocus int

	// ENUM(
	//	Unknown,
	//	MainCompound,
	//	MainCompoundAccessory,
	//	CompoundAccessory,
	//	Accessory,
	// )
	ExerciseKind int

	// ENUM(
	//	Unknown,
	//	SimplifiedNegativeSpace,
	// )
	ModelID int32
)
