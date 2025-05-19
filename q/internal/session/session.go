package session

import (
	"net/http"
)

const (
	UserIDKey = "user_id"
)

func GetUserID(r *http.Request) int {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		return 0
	}

	userID, ok := session.Values[UserIDKey].(int)
	if !ok {
		return 0
	}

	return userID
}

func SetUserID(w http.ResponseWriter, r *http.Request, userID int) error {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		return err
	}

	session.Values[UserIDKey] = userID
	return session.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		return err
	}

	session.Values = make(map[interface{}]interface{})
	return session.Save(r, w)
}
