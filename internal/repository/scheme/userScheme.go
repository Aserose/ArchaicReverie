package scheme

import "fmt"

func CreateSchemaUser(numberCharLimit int) string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		username varchar(255) not null unique,
		password varchar(255) not null,
		numberOfCharacters smallint CHECK (numberOfCharacters < %d)
	);`, numberCharLimit)
}
