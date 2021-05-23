package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Chat ...
type Chat struct {
	ChatID    int64  `json:"chat_id"`
	Name      string `json:"name"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Token     string `json:"token"`
	Date      string `json:"date"`
	MsgCount  int    `json:"mcount"`
	Admin     bool   `json:"admin"`
}

var (
	chats []Chat
)

func dbOpen() error {

	chats = make([]Chat, 0)

	fpath := filepath.Join(ArgWorkdir, "register_db.json")
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		if err == os.ErrNotExist {
			return dbFlush()
		}
		log.Println(err)
		return nil
	}
	err = json.Unmarshal(data, &chats)
	if err != nil {
		return err
	}

	return nil
}

func dbClose() {
}

func dbListUsers() string {

	var strs []string
	strs = append(strs, "Nutzer Liste")

	for _, chat := range chats {
		strs = append(strs, "")
		strs = append(strs, fmt.Sprintf("%s (%v, %v)", chat.Name, chat.ChatID, chat.Admin))
		strs = append(strs, fmt.Sprintf("  - Name:        %s, %s", chat.FirstName, chat.LastName))
		strs = append(strs, fmt.Sprintf("  - Registriert: %s", chat.Date))
		strs = append(strs, fmt.Sprintf("  - Nachrichten: %d", chat.MsgCount))
	}

	return strings.Join(strs, "\n")
}

func dbGetChats() []int64 {
	var ids []int64
	for _, chat := range chats {
		ids = append(ids, chat.ChatID)
	}
	return ids
}

func dbRegisterChat(id int64, user string, fname string, lname string, token string) error {
	now := time.Now()

	if dbIsRegisteredChat(id) == nil {
		return errors.New("Chat allready registered")
	}
	chats = append(chats, Chat{
		ChatID:    id,
		Name:      user,
		FirstName: fname,
		LastName:  lname,
		Token:     token,
		Date:      now.Format("2006-Jan-2T15:04:05-0700"),
		MsgCount:  0,
	})

	return dbFlush()
}

func dbFlush() error {
	fpath := filepath.Join(ArgWorkdir, "register_db.json")
	file, err := os.OpenFile(fpath+".new", os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	jstring, err := json.MarshalIndent(chats, "", "  ")
	if err != nil {
		return err
	}
	file.Write(jstring)
	file.Close()

	return os.Rename(fpath+".new", fpath)
}

func dbIsRegisteredChat(id int64) error {
	for _, chat := range chats {
		if chat.ChatID == id {
			return nil
		}
	}
	return errors.New("no registered user")
}

func dbIsAdmin(id int64) bool {
	first := true
	for _, chat := range chats {
		if chat.ChatID == id {
			return chat.Admin || first
		}
		first = false
	}
	return false
}

func dbToggleAdmin(id int64) error {
	for _, chat := range chats {
		if chat.ChatID == id {
			chat.Admin = !chat.Admin
			return dbFlush()
		}
	}
	return errors.New("no registered user")
}

func dbLeave(id int64) string {

	for i, chat := range chats {
		if chat.ChatID == id {
			copy(chats[i:], chats[i+1:])
			chats = chats[:len(chats)-1]
			if err := dbFlush(); err != nil {
				return "Error: " + err.Error()
			}
			return "Successfully unregistered!"
		}
	}
	return "Unknown user"
}