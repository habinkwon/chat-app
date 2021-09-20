CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	name STRING,
	username STRING,
	email STRING NOT NULL,
	password STRING NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	UNIQUE (username),
	UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS chats (
	id UUID PRIMARY KEY,
	name STRING NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	created_by UUID REFERENCES users (id) ON DELETE SET NULL,
);

CREATE TABLE IF NOT EXISTS chat_members (
	chat_id UUID REFERENCES chats (id) ON DELETE CASCADE,
	member_id UUID REFERENCES users (id) ON DELETE CASCADE,
	added_at TIMESTAMPTZ NOT NULL,
	added_by UUID REFERENCES users (id) ON DELETE SET NULL,
	PRIMARY KEY (chat_id, member_id)
);

CREATE TABLE IF NOT EXISTS chat_messages (
	chat_id UUID REFERENCES chats (id) ON DELETE CASCADE,
	id UUID,
	type ENUM ('message', 'event'),
	content BYTES NOT NULL,
	event STRING,
	sender_id UUID REFERENCES users (id) ON DELETE SET NULL,
	reply_to UUID REFERENCES chat_messages (id) ON DELETE SET NULL,
	created_at TIMESTAMPTZ NOT NULL,
	edited_at TIMESTAMPTZ,
	PRIMARY KEY (chat_id, id),
	FOREIGN KEY (chat_id, reply_to) REFERENCES chat_messages (chat_id, id) ON DELETE SET NULL
);
