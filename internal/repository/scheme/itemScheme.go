package scheme

var SchemaWeaponAndInventory = `CREATE TABLE IF NOT EXISTS weapon (
		weapon_id serial PRIMARY KEY,
		name varchar(25) unique,
		price smallint,
		weight smallint,
		weapon_class smallint,
		sharp smallint,
		weight smallint
	);
	ALTER TABLE weapon
	ADD COLUMN IF NOT EXISTS price smallint;

	CREATE TABLE IF NOT EXISTS inventory (
		char_id smallint REFERENCES characters (charId),
		weapon_id smallint[] CHECK (cardinality(weapon_id) < 3),
		coin_amount smallint
	);

	INSERT INTO weapon AS w (name, price, weapon_class, sharp, weight) VALUES 
		('knife', 50, 2, 1, 0),('baseball_bat', 30, 3, 0, 1),('brass_knuckles', 50, 2, 0, 0),
		('katana', 100 ,1, 3, 1),('sai', 100, 1, 2, 0),('glassing', 10, 3, 1, 0)
	ON CONFLICT (name) DO UPDATE SET
		weapon_class = EXCLUDED.weapon_class,
		price = EXCLUDED.price,
		sharp = EXCLUDED.sharp,
		weight = EXCLUDED.weight`

var SchemaFood = `CREATE TABLE IF NOT EXISTS foods (
		name varchar(25) not null unique,
		price smallint,
		restore_hp smallint
	); INSERT INTO foods AS f (name, price, restore_hp) VALUES
		('apple', 3, 10),('beef', 15, 25)
	ON CONFLICT (name) DO UPDATE SET
		price = EXCLUDED.price,
		restore_hp = EXCLUDED.restore_hp;
`
