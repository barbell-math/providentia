package types

type (
	// Represents a client from the database
	Client struct {
		FirstName string `db:"first_name"` // The first name of the client
		LastName  string `db:"last_name"`  // The last name of the client
		Email     string `db:"email"`      // The clients email
	}
)
