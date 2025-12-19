package types

type (
	// Represents a client from the database
	Client struct {
		FirstName string `db:"first_name"` // The first name of the client
		LastName  string `db:"last_name"`  // The last name of the client
		Email     string `db:"email"`      // The clients email
	}

	// Represents an exercise from the database
	Exercise struct {
		Name    string        `db:"name"`     // The exercise name
		KindId  ExerciseKind  `db:"kind_id"`  // The kind of exercise
		FocusId ExerciseFocus `db:"focus_id"` // The focus of the exercise
	}
)
