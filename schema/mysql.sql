CREATE TABLE IF NOT EXISTS users (
	id BIGINT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	created_at DATETIME NOT NULL,
	UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS chats (
	id BIGINT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	created_at DATETIME NOT NULL,
	created_by BIGINT -- REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS chat_members (
	chat_id BIGINT NOT NULL, -- REFERENCES chats (id) ON DELETE CASCADE
	user_id BIGINT NOT NULL, -- REFERENCES users (id) ON DELETE CASCADE
	added_at DATETIME NOT NULL,
	added_by BIGINT, -- REFERENCES users (id) ON DELETE SET NULL
	PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
	chat_id BIGINT, -- REFERENCES chats (id) ON DELETE CASCADE
	id BIGINT,
	type ENUM ('message', 'event') NOT NULL,
	content MEDIUMBLOB NOT NULL,
	event TEXT,
	sender_id BIGINT, -- REFERENCES users (id) ON DELETE SET NULL
	reply_to BIGINT, -- REFERENCES chat_messages (id) ON DELETE SET NULL
	created_at DATETIME NOT NULL,
	edited_at DATETIME,
	PRIMARY KEY (chat_id, id)
	-- FOREIGN KEY (chat_id, reply_to) REFERENCES chat_messages (chat_id, id) ON DELETE SET NULL
);
