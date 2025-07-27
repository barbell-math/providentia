CREATE SCHEMA IF NOT EXISTS providentia;

CREATE TABLE IF NOT EXISTS providentia.exercise_focus (
	id SERIAL PRIMARY KEY NOT NULL,
	focus TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.exercise_kind (
	id SERIAL NOT NULL PRIMARY KEY,
	kind TEXT NOT NULL,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.model (
	id SERIAL NOT NULL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.exercise (
	id SERIAL NOT NULL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	kind_id INT NOT NULL REFERENCES providentia.exercise_kind(id),
	focus_id INT NOT NULL REFERENCES providentia.exercise_focus(id)
);

CREATE TABLE IF NOT EXISTS providentia.client (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS providentia.video_data (
	id BIGINT NOT NULL PRIMARY KEY CHECK (id>=0),
	path TEXT NOT NULL UNIQUE,
	position FLOAT[][] NOT NULL,
	velocity FLOAT[][] NOT NULL,
	acceleration FLOAT[][] NOT NULL,
	force FLOAT[][] NOT NULL,
	impulse FLOAT[][] NOT NULL
);

CREATE TABLE IF NOT EXISTS providentia.training_log (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	exercise_id INT NOT NULL REFERENCES providentia.exercise(id),
	exercise_kind_id INT NOT NULL REFERENCES providentia.exercise_kind(id),
	exercise_focus_id INT NOT NULL REFERENCES providentia.exercise_focus(id),
	client_id BIGINT NOT NULL REFERENCES providentia.client(id),
	video_id BIGINT NOT NULL REFERENCES providentia.video_data(id),

	date_performed DATE NOT NULL,
	weight FLOAT NOT NULL CHECK (weight>=0),
	sets FLOAT NOT NULL CHECK (sets>=0),
	reps INT NOT NULL CHECK (reps>=0),
	effort FLOAT NOT NULL CHECK (effort>=0 AND effort<=10),

	inter_session_cntr INT NOT NULL CHECK (inter_session_cntr>0),
	inter_workout_cntr INT NOT NULL CHECK (inter_workout_cntr>0),

	volume FLOAT NOT NULL CHECK (volume>=0) GENERATED ALWAYS AS (weight*sets*reps) STORED,
	exertion FLOAT NOT NULL CHECK (exertion>=0) GENERATED ALWAYS AS (effort*sets*reps) STORED,
	total_reps FLOAT NOT NULL CHECK (total_reps>=0) GENERATED ALWAYS AS (sets*reps) STORED
);

CREATE TABLE IF NOT EXISTS providentia.model_state (
	id BIGSERIAL NOT NULL PRIMARY KEY,
	client_id BIGINT NOT NULL REFERENCES providentia.client(id),
	training_log_id BIGINT NOT NULL REFERENCES providentia.training_log(id),
	model_id INT NOT NULL REFERENCES providentia.model(id),

	v1 FLOAT NOT NULL CHECK (v1>=0),
	v2 FLOAT NOT NULL CHECK (v2>=0),
	v3 FLOAT NOT NULL CHECK (v3>=0),
	v4 FLOAT NOT NULL CHECK (v4>=0),
	v5 FLOAT NOT NULL CHECK (v5>=0),
	v6 FLOAT NOT NULL CHECK (v6>=0),
	v7 FLOAT NOT NULL CHECK (v7>=0),
	v8 FLOAT NOT NULL CHECK (v8>=0),
	v9 FLOAT NOT NULL CHECK (v9>=0),
	v10 FLOAT NOT NULL CHECK (v10>=0),

	time_frame INT NOT NULL CHECK (time_frame>=0),
	mse FLOAT NOT NULL CHECK (mse>=0),
	pred_weight FLOAT NOT NULL CHECK (pred_weight>=0),

	UNIQUE (training_log_id, model_id)
);
