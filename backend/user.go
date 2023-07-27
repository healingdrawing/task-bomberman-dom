package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type signupData struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Dob         string `json:"dob"`
	Avatar      string `json:"avatar"`
	avatarBytes []byte `sqlite3:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"about_me"`
	Public      bool   `json:"public"`
	Privacy     string `sqlite3:"privacy"`
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UUIDData struct {
	UUID string `json:"UUID"`
}

func createSession(email string) (UUID string, err error) {
	random, _ := uuid.NewV4()
	UUID = random.String()
	log.Println("HERE WAS ADD SESSION QUERY")
	return UUID, nil
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer recovery(w)

	var data signupData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusBadRequest, "")
		return
	}

	// if dob is empty, return error
	if data.Dob == "" {
		log.Println("dob is empty")
		jsonResponse(w, http.StatusUnprocessableEntity, "Invalid date of birth")
		return
	}

	if len([]rune(strings.TrimSpace(data.Nickname))) > 15 {
		data.Nickname =
			string([]rune(strings.TrimSpace(data.Nickname))[:15])
	}

	log.Println("HERE WAS ADD USER QUERY")

	// TODO: check/simplify , we need only uuid for map as key
	UUID, err := createSession(data.Email)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, " create session failed")
		return
	}

	w.WriteHeader(200)
	jsonResponseObj, _ := json.Marshal(map[string]string{
		"UUID":  UUID,
		"email": data.Email,
	})
	_, err = w.Write(jsonResponseObj)
	if err != nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusInternalServerError, "w.Write(jsonResponseObj)<-UUID,email failed")
		return
	}

}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer recovery(w)

	cookie, err := r.Cookie("user_uuid")
	if err != nil || cookie.Value == "" || cookie == nil {
		log.Println(err.Error())
		jsonResponse(w, http.StatusOK, "You are not logged in")
		return
	}

	uuid := cookie.Value

	log.Println(uuid)
	log.Println("HERE WAS REMOVE SESSION DATA BASE QUERY")

	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	jsonResponse(w, http.StatusOK, "Session deleted")
}
