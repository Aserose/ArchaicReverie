package scheme

import (
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/config"
)

func CreateSchemaUser(numberCharLimit int) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		username varchar(255) not null unique,
		password varchar(255) not null,
		numberOfCharacters smallint CHECK (numberOfCharacters < %d)
	);`, numberCharLimit)
}

func CreateSchemaCharacter(charConfig config.CharacterConfig) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS characters (
		charId serial not null unique,
		ownerId integer not null,
		name varchar(255) not null,
		growth smallint CHECK (growth>%d) CHECK (growth<%d),
		weight smallint CHECK (weight>%d) CHECK (weight<%d),
			FOREIGN KEY (ownerId) REFERENCES users (id) ON DELETE CASCADE
	);`, charConfig.MinCharGrowth, charConfig.MaxCharGrowth,
		charConfig.MinCharWeight, charConfig.MaxCharWeight)
}

var SchemaLocation = `CREATE TABLE IF NOT EXISTS times (
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

var SchemaFood = `CREATE TABLE IF NOT EXISTS foods (
		name varchar(25) not null unique,
		price smallint,
		restore_hp smallint
	); INSERT INTO foods AS f (name, price, restore_hp) VALUES
		('apple',3,10),('beef',15,25)
	ON CONFLICT (name) DO UPDATE SET
		price = EXCLUDED.price,
		restore_hp = EXCLUDED.restore_hp;
`
