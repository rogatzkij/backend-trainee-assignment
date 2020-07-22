package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"time"
)
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type ConfigMySQL struct {
	Login    string `default:"root"`
	Password string `default:"password"`
	Host     string `default:"127.0.0.1"`
	Port     string `default:"3306"`
	Database string `default:"chat"`
}

func initConfigMySQL() (*ConfigMySQL, error) {
	config := &ConfigMySQL{}
	err := envconfig.Process("MySQL", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

type ConnectorMySQL struct {
	config *ConfigMySQL
	db     *sql.DB
}

func (cp *ConnectorMySQL) connect() error {
	sourceAddr := fmt.Sprintf("%s:%s@/%s", cp.config.Login, cp.config.Password, cp.config.Database)
	db, err := sql.Open("mysql", sourceAddr)
	if err != nil {
		return err
	}

	cp.db = db
	return nil
}

func (cp *ConnectorMySQL) createUser(username string) (User, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return User{}, err
		}
	}
	_, err := cp.db.Query("INSERT INTO E1_Users (username, created_at) VALUE(?,?)", username, time.Now().String())
	if err != nil {
		return User{}, err
	}

	user := User{}
	rows, err := cp.db.Query("SELECT * FROM E1_Users WHERE username = ?", username)
	if err != nil {
		return User{}, err
	}

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.CreatedAt)
		if err != nil {
			return User{}, err
		}
		break
	}
	rows.Close()

	return user, nil
}

func (cp *ConnectorMySQL) checkUsername(username string) (bool, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return false, err
		}
	}

	rows, err := cp.db.Query("SELECT * FROM E1_Users WHERE username = ?", username)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func (cp *ConnectorMySQL) checkUserID(user uint64) (bool, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return false, err
		}
	}

	rows, err := cp.db.Query("SELECT * FROM E1_Users WHERE id = ?", user)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func (cp *ConnectorMySQL) createChart(name string, users []uint64) (Chat, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return Chat{}, err
		}
	}
	_, err := cp.db.Query("INSERT INTO E2_Chat (name, created_at) VALUE(?,?)", name, time.Now().String())
	if err != nil {
		return Chat{}, err
	}

	chat := Chat{}
	rows, err := cp.db.Query("SELECT * FROM E2_Chat WHERE name = ?", name)
	if err != nil {
		return Chat{}, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
		break
	}

	for _, userID := range users {
		_, err := cp.db.Query("INSERT INTO E3_Chatroom (id_user, id_chat) VALUE (?,?)", userID, chat.ID)
		if err != nil {
			return Chat{}, err
		}
	}

	return chat, nil
}

func (cp *ConnectorMySQL) checkChartName(name string) (bool, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return false, err
		}
	}

	rows, err := cp.db.Query("SELECT * FROM E2_Chat WHERE name = ?", name)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func (cp *ConnectorMySQL) checkChartID(chat uint64) (bool, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return false, err
		}
	}

	rows, err := cp.db.Query("SELECT * FROM E2_Chat WHERE id = ?", chat)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func (cp *ConnectorMySQL) getCharts(user uint64) ([]Chat, error) {
	querry := `SELECT
E2_Chat.id,
E2_Chat.name
FROM E2_Chat
JOIN E3_Chatroom E3C on E2_Chat.id = E3C.id_chat
JOIN
    (SELECT id_chat, created_at FROM E4_Messages ORDER BY created_at DESC LIMIT 1) as E4M on E2_Chat.id = E4M.id_chat
WHERE E3C.id_user = ?`

	var result []Chat

	rows, err := cp.db.Query(querry, user)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		chat := Chat{}
		err = rows.Scan(&chat.ID, &chat.Name)
		if err != nil {
			return nil, err
		}

		rowsUserID, err := cp.db.Query("SELECT id_user FROM E3_Chatroom WHERE id_chat = ?", chat.ID)
		if err != nil {
			return nil, err
		}

		for rowsUserID.Next() {
			var userID uint64
			err = rowsUserID.Scan(&userID)
			if err != nil {
				return nil, err
			}
			chat.Users = append(chat.Users, userID)
		}

		result = append(result, chat)
	}

	return result, nil
}

func (cp *ConnectorMySQL) sendMessage(chatID uint64, authorID uint64, text string) (Message, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return Message{}, err
		}
	}

	createTime := time.Now().String()
	_, err := cp.db.Query("INSERT INTO E4_Messages(id_chat, id_user, text , created_at) VALUE(?,?,?,?)",
		chatID, authorID, text, createTime)
	if err != nil {
		return Message{}, err
	}

	rows, err := cp.db.Query("SELECT * FROM E4_Messages WHERE id_chat = ? AND id_user = ? AND created_at=?",
		chatID, authorID, createTime)

	message := Message{}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&message.ID, &message.Chat, &message.Author, &message.CreatedAt)
		break
	}

	return message, nil
}

func (cp *ConnectorMySQL) getMessages(chatID uint64) ([]Message, error) {
	if cp.db == nil {
		if err := cp.connect(); err != nil {
			return nil, err
		}
	}

	var result []Message
	rows, err := cp.db.Query("SELECT * FROM E4_Messages WHERE id_chat = ? ORDER BY created_at ASC",
		chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		chat := Message{}
		err = rows.Scan(&chat.ID, &chat.Chat, &chat.Author, &chat.Text, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, chat)
	}

	return result, nil
}
