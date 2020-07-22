package main

// User - Пользователь приложения. Имеет следующие свойства:
type User struct {
	ID        uint64 `json:"id"`         // уникальный идентификатор пользователя
	Username  string `json:"username"`   // уникальное имя пользователя
	CreatedAt string `json:"created_at"` // время создания пользователя
}

// Chat - Отдельный чат. Имеет следующие свойства:
type Chat struct {
	ID        uint64   `json:"id"`         //уникальный идентификатор чата
	Name      string   `json:"name"`       //уникальное имя чата
	Users     []uint64 `json:"users"`      //список пользователей в чате, отношение многие-ко-многим
	CreatedAt string   `json:"created_at"` //время создания
}

//Message - Сообщение в чате. Имеет следующие свойства:
type Message struct {
	ID        uint64 `json:"id"`         //уникальный идентификатор сообщения
	Chat      uint64 `json:"chat"`       //ссылка на идентификатор чата, в который было отправлено сообщение
	Author    string `json:"author"`     //ссылка на идентификатор отправителя сообщения, отношение многие-к-одному
	Text      string `json:"text"`       //текст отправленного сообщения
	CreatedAt string `json:"created_at"` //время создания
}
