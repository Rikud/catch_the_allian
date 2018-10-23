package services

import (
	"IT-Berries_Go_server/auth/models"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	CookieName = "session_id"
)

type SessionData struct {
	Uid int
	Expires int64
}

var sessions = make(map[string]*SessionData)

func NewSession(w http.ResponseWriter, userId int) error {
	sid, err := sessionId()
	if err != nil {
		log.Println("Error while trying to generate session id:", err)
		 return err
	}
	loc, _ := time.LoadLocation("UTC")
	expirationDate := time.Now().In(loc).Add(time.Hour)
	cookie := http.Cookie {
		Name: CookieName,
		Value: sid,
		HttpOnly: true,
		Expires: expirationDate,
	}
	http.SetCookie(w, &cookie)
	sessions[sid] = &SessionData{userId, expirationDate.Unix()}
	return nil
}

func GetUserBySessionId(sessionID string) *models.User {
	userSession := sessions[sessionID]
	if userSession == nil {
		return nil
	}
	loc, _ := time.LoadLocation("UTC")
	sessionTime := userSession.Expires
	nowTime :=time.Now().In(loc).Unix()
	if sessionTime < nowTime {
		return  nil
	}
	user := findUserById(userSession.Uid)
	if user == nil {
		return nil
	}
	return user
}

func sessionId() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func DeleteSession(session *http.Cookie, w http.ResponseWriter) {
	delete(sessions, session.Value)
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}