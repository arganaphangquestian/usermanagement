CREATE TABLE IF NOT EXISTS users
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
	username 		VARCHAR(100) NOT NULL,
	email 			VARCHAR(100) NOT NULL UNIQUE,
	role 			VARCHAR(100) NOT NULL,
	password 		TEXT NOT NULL,
	referral 		TEXT
);