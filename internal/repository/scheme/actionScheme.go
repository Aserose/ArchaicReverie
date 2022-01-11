package scheme

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

var SchemaDamageAndResult = `CREATE TABLE IF NOT EXISTS action_result (
		name varchar(25) PRIMARY KEY
	);
	CREATE TABLE IF NOT EXISTS damage_type (
		name varchar(25) UNIQUE,
		damage_hp smallint,
		damage_mp smallint,
			FOREIGN KEY (name) REFERENCES action_result (name) ON DELETE CASCADE
	);
	WITH actionInfo(name, damage_hp, damage_mp) AS (
		VALUES ('fall', 10, 0)),
		ins1 AS (
		INSERT INTO action_result (name)
		SELECT name FROM actionInfo
			ON CONFLICT DO NOTHING
			RETURNING name)
		INSERT INTO damage_type (name, damage_hp, damage_mp)
		SELECT ins1.name, a.damage_hp, a.damage_mp
		FROM actionInfo a
		JOIN ins1 USING (name)
		ON CONFLICT (name) DO UPDATE SET
			damage_hp = EXCLUDED.damage_hp,
			damage_mp = EXCLUDED.damage_mp;`

var SchemaEnemy = `CREATE TABLE IF NOT EXISTS enemy (
		name char(25) PRIMARY KEY,
		class smallint
	);
		INSERT INTO enemy AS e (name, class) VALUES
			('hooligan',2), ('yakuza',1), ('drunkard',3)
		ON CONFLICT (name) DO UPDATE SET
			class = EXCLUDED.class`
