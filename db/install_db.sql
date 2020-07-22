CREATE DATABASE chat DEFAULT charset utf8;
USE chat;

-- Пользователь приложения
CREATE TABLE E1_Users
(
    id         INTEGER AUTO_INCREMENT, -- уникальный идентификатор пользователя
    username   VARCHAR(32),            -- уникальное имя пользователя
    created_at DATETIME,                   -- время создания пользователя

    PRIMARY KEY (id),
    UNIQUE (username)
);

CREATE TABLE E2_Chat
(
    id         INTEGER AUTO_INCREMENT, -- уникальный идентификатор чата
    name       VARCHAR(32),            -- уникальное имя чата
    created_at DATETIME,                   -- время создания

    PRIMARY KEY (id),
    UNIQUE (name)
);

-- список пользователей в чате, отношение многие-ко-многим
CREATE TABLE E3_Chatroom
(
	id_user  INTEGER NOT NULL,
	id_chat  INTEGER NOT NULL,

	PRIMARY KEY (id_user,id_chat),
	FOREIGN KEY (id_user) REFERENCES E1_Users(id),
	FOREIGN KEY (id_chat) REFERENCES E2_Chat(id)
);

-- Сообщение в чате. Имеет следующие свойства:
CREATE TABLE E4_Messages
(
    id         INTEGER AUTO_INCREMENT, -- уникальный идентификатор сообщения
    id_chat    INTEGER NOT NULL,       -- ссылка на идентификатор чата, в который было отправлено сообщение
    id_user    INTEGER NOT NULL,       -- ссылка на идентификатор отправителя сообщения, отношение многие-к-одному
    text       VARCHAR(32),            -- текст отправленного сообщения
    created_at DATETIME,               -- время создания

    PRIMARY KEY (id),
    FOREIGN KEY (id_user) REFERENCES E1_Users(id),
	FOREIGN KEY (id_chat) REFERENCES E2_Chat(id)
);
