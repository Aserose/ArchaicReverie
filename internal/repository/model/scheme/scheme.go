package scheme

var SchemaUser = `CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		username varchar(255) not null unique,
		password varchar(255) not null,
		numberOfCharacters smallint CHECK (numberOfCharacters < 4)
	);`

var SchemaCharacter = `CREATE TABLE IF NOT EXISTS characters (
		charId serial not null unique,
		ownerId integer not null,
		name varchar(255) not null,
		growth smallint CHECK (growth>144) CHECK (growth<201),
		weight smallint CHECK (weight>38) CHECK (weight<123),
			FOREIGN KEY (ownerId) REFERENCES users (id) ON DELETE CASCADE
	);`

var SchemaLocation = ` 
CREATE TABLE IF NOT EXISTS times (
		name varchar(25) not null unique,
		clarity smallint
	); INSERT INTO times (name, clarity) VALUES
		('night', -1),('day', 1),
		('sunset', 1),('sunrise', 1) 
	ON CONFLICT (name) DO UPDATE SET
		clarity = EXCLUDED.clarity;

CREATE TABLE IF NOT EXISTS weathers (
		name varchar(25) not null unique,
		clarity smallint,
		difficulty_movement smallint
	); INSERT INTO weathers (name, clarity, difficulty_movement) VALUES
		('fog', -1, 1),('rain', -1,-1),
		('clear', 1,1),('snowfall', -1,-1) 
	ON CONFLICT (name) DO UPDATE SET 
		clarity = EXCLUDED.clarity,
		difficulty_movement = EXCLUDED.difficulty_movement;

CREATE TABLE IF NOT EXISTS places (
		name varchar(25) not null unique,
		difficulty_movement smallint
	); INSERT INTO places (name, difficulty_movement) VALUES
		('road', 1),('rough surface', -1) 
	ON CONFLICT (name) DO UPDATE SET
		difficulty_movement = EXCLUDED.difficulty_movement;

CREATE TABLE IF NOT EXISTS obstacles (
		name varchar(25) not null unique,
		height smallint,
		length smallint
	); INSERT INTO obstacles AS ob (name, height, length) VALUES
		('small pit',1,1),('beam',-1,-1) 
	ON CONFLICT (name) DO UPDATE SET
		height = EXCLUDED.height,
		length = EXCLUDED.length;
`
