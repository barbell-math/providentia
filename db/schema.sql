CREATE SCHEMA IF NOT EXISTS providentia;

CREATE TABLE IF NOT EXISTS providentia.ExerciseFocus (
	ID SERIAL PRIMARY KEY NOT NULL,
	Focus TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.ExerciseKind (
	ID SERIAL NOT NULL PRIMARY KEY,
	Kind TEXT NOT NULL,
	Description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.Exercise (
	ID SERIAL NOT NULL PRIMARY KEY,
	Name TEXT NOT NULL,
	KindID INT NOT NULL REFERENCES providentia.ExerciseKind(ID),
	FocusID INT NOT NULL REFERENCES providentia.ExerciseFocus(ID)
);

CREATE TABLE IF NOT EXISTS providentia.Client (
	ID BIGSERIAL NOT NULL PRIMARY KEY,
	FirstName TEXT NOT NULL,
	LastName TEXT NOT NULL,
	Email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.VideoData (
	ID BIGINT NOT NULL PRIMARY KEY CHECK (ID>=0),
	Position FLOAT[][] NOT NULL,
	Velocity FLOAT[][] NOT NULL,
	Acceleration FLOAT[][] NOT NULL,
	Force FLOAT[][] NOT NULL,
	Impulse FLOAT[][] NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.TrainingLog (
	ID BIGSERIAL NOT NULL PRIMARY KEY,
	ExerciseKindID INT NOT NULL REFERENCES providentia.ExerciseKind(ID),
	ExerciseFocusID INT NOT NULL REFERENCES providentia.ExerciseFocus(ID),
	ClientID INT NOT NULL REFERENCES providentia.Client(ID),
	VideoID BIGINT NOT NULL REFERENCES providentia.VideoData(ID),

	DatePerformed DATE NOT NULL,
	Weight FLOAT NOT NULL CHECK (Weight>=0),
	Sets FLOAT NOT NULL CHECK (Weight>0),
	Reps INT NOT NULL CHECK (Weight>=0),
	Effort FLOAT NOT NULL CHECK (Effort >=0 AND Effort<=10),
	InterExerciseCntr INT NOT NULL CHECK (InterExerciseCntr>=0),
	InterWorkoutCntr INT NOT NULL CHECK (InterWorkoutCntr>=0),

	Volume FLOAT NOT NULL CHECK (Volume>=0) GENERATED ALWAYS AS (Weight*Sets*Reps) STORED,
	Efforts FLOAT NOT NULL CHECK (Efforts>=0) GENERATED ALWAYS AS (Effort*Sets*Reps) STORED
);
