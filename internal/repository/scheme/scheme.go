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
	);`, charConfig.Restriction.MinCharGrowth, charConfig.Restriction.MaxCharGrowth,
		charConfig.Restriction.MinCharWeight, charConfig.Restriction.MaxCharWeight)
}
