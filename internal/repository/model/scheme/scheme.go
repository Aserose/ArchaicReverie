package scheme

var SchemaUser = `CREATE TABLE IF NOT EXISTS users (
		id serial not null unique,
		username varchar(255) not null unique,
		password varchar(255) not null
	);`

var SchemaCharacter = `CREATE TABLE IF NOT EXISTS characters (
		charId serial not null unique,
		ownerId integer not null,
		name varchar(255) not null,
		growth smallint CHECK (growth>0) CHECK (growth<200),
		weight smallint CHECK (weight>0) CHECK (weight<200)
	);`
