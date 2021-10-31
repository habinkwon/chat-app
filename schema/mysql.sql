CREATE TABLE IF NOT EXISTS user (
	id BIGINT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	nickname VARCHAR(255) NOT NULL,
	email VARCHAR(255),
	picture VARCHAR(255),
	livingPlace VARCHAR(255),
	preference1 VARCHAR(255),
	preference2 VARCHAR(255),
	preference3 VARCHAR(255),
	SelfIntroduction VARCHAR(255),
	role ENUM('USER', 'GUEST') NOT NULL DEFAULT 'USER',
	UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS chats (
	id BIGINT PRIMARY KEY,
	created_by BIGINT NOT NULL, -- REFERENCES users (id) ON DELETE SET NULL
	created_at DATETIME NOT NULL,
	last_posted_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_members (
	chat_id BIGINT, -- REFERENCES chats (id) ON DELETE CASCADE
	user_id BIGINT, -- REFERENCES users (id) ON DELETE CASCADE
	added_by BIGINT NOT NULL, -- REFERENCES users (id) ON DELETE SET NULL
	added_at DATETIME NOT NULL,
	PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
	id BIGINT PRIMARY KEY,
	chat_id BIGINT NOT NULL, -- REFERENCES chats (id) ON DELETE CASCADE
	type ENUM ('message', 'event') NOT NULL,
	content TEXT,
	event TEXT,
	sender_id BIGINT NOT NULL, -- REFERENCES users (id) ON DELETE SET NULL
	reply_to BIGINT, -- REFERENCES chat_messages (id) ON DELETE SET NULL
	created_at DATETIME NOT NULL,
	edited_at DATETIME
	-- FOREIGN KEY (chat_id, reply_to) REFERENCES chat_messages (chat_id, id) ON DELETE SET NULL
);

ALTER TABLE chat_messages ADD INDEX (chat_id);
