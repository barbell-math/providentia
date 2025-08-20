CREATE SCHEMA IF NOT EXISTS providentia;

CREATE TABLE IF NOT EXISTS providentia.exercise_focus (
	id SERIAL4 PRIMARY KEY NOT NULL,
	focus TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.exercise_kind (
	id SERIAL4 NOT NULL PRIMARY KEY,
	kind TEXT NOT NULL,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.model (
	id SERIAL4 NOT NULL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.exercise (
	id SERIAL4 NOT NULL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	kind_id INT4 NOT NULL REFERENCES providentia.exercise_kind(id),
	focus_id INT4 NOT NULL REFERENCES providentia.exercise_focus(id)
);

CREATE TABLE IF NOT EXISTS providentia.client (
	id SERIAL8 NOT NULL PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS providentia.physics_data (
	id SERIAL8 NOT NULL PRIMARY KEY,
	path TEXT[] UNIQUE,

	time FLOAT8[][] NOT NULL,
	position FLOAT8[][] NOT NULL,
	velocity FLOAT8[][] NOT NULL,
	acceleration FLOAT8[][] NOT NULL,
	jerk FLOAT8[][] NOT NULL,

	force FLOAT8[][] NOT NULL,
	impulse FLOAT8[][] NOT NULL,
	work FLOAT8[][] NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.training_log (
	id SERIAL8 NOT NULL PRIMARY KEY,
	exercise_id INT4 NOT NULL REFERENCES providentia.exercise(id),
	client_id INT8 NOT NULL REFERENCES providentia.client(id),
	physics_id INT8 REFERENCES providentia.physics_data(id),

	date_performed DATE NOT NULL,
	inter_session_cntr INT2 NOT NULL CHECK (inter_session_cntr>0),
	inter_workout_cntr INT2 NOT NULL CHECK (inter_workout_cntr>0),

	weight FLOAT8 NOT NULL CHECK (weight>=0),
	sets FLOAT8 NOT NULL CHECK (sets>=0),
	reps INT4 NOT NULL CHECK (reps>=0),
	effort FLOAT8 NOT NULL CHECK (effort>=0 AND effort<=10),

	volume FLOAT8 NOT NULL CHECK (volume>=0) GENERATED ALWAYS AS (weight*sets*reps) STORED,
	exertion FLOAT8 NOT NULL CHECK (exertion>=0) GENERATED ALWAYS AS (effort*sets*reps) STORED,
	total_reps FLOAT8 NOT NULL CHECK (total_reps>=0) GENERATED ALWAYS AS (sets*reps) STORED,

	UNIQUE (client_id, date_performed, inter_session_cntr, inter_workout_cntr)
);

CREATE TABLE IF NOT EXISTS providentia.model_state (
	id SERIAL8 NOT NULL PRIMARY KEY,
	client_id INT8 NOT NULL REFERENCES providentia.client(id),
	training_log_id INT8 NOT NULL REFERENCES providentia.training_log(id),
	model_id INT4 NOT NULL REFERENCES providentia.model(id),

	v1 FLOAT8 NOT NULL CHECK (v1>=0),
	v2 FLOAT8 NOT NULL CHECK (v2>=0),
	v3 FLOAT8 NOT NULL CHECK (v3>=0),
	v4 FLOAT8 NOT NULL CHECK (v4>=0),
	v5 FLOAT8 NOT NULL CHECK (v5>=0),
	v6 FLOAT8 NOT NULL CHECK (v6>=0),
	v7 FLOAT8 NOT NULL CHECK (v7>=0),
	v8 FLOAT8 NOT NULL CHECK (v8>=0),
	v9 FLOAT8 NOT NULL CHECK (v9>=0),
	v10 FLOAT8 NOT NULL CHECK (v10>=0),

	time_frame INT8 NOT NULL CHECK (time_frame>=0),
	mse FLOAT8 NOT NULL CHECK (mse>=0),
	pred_weight FLOAT8 NOT NULL CHECK (pred_weight>=0),

	UNIQUE (training_log_id, model_id)
);
