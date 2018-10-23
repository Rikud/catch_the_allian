package controllers

import (
	"IT-Berries_Go_server/auth/encoder"
	"IT-Berries_Go_server/auth/models"
	"IT-Berries_Go_server/auth/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const avatarPath = "/app/avatars/"

type JSONError struct {
	Err string `json:"error, string"`
	Code int
}

var Handlers = map[string]Handler{}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type LoginHandle struct{}
type RegistrationHandle struct{}
type MeHandle struct{}
type ScoreboardHandle struct{}
type LogOutHandle struct {}
type AvatarHandle struct{}

func (AvatarHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, avatarPath + string(r.URL.Query()["avatar"][0]))
}

type EntryScore struct {
	Scorelist []*models.ScoreRecord `json:"scorelist"`
	Length int `json:"length"`
}

func (LogOutHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Log out failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	log.Println("Trying to log out user.")
	cookie, err := r.Cookie(services.CookieName)
	if err != nil || cookie.Value == "" {
		log.Println("The user isn't authorized!")
		errorResponce(w, "The user isn't authorized!", http.StatusUnauthorized)
		return
	}
	services.DeleteSession(cookie, w)
	w.WriteHeader(http.StatusOK)
}

func (LoginHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Login failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	if !checkMethod(w, r) {return}
	log.Println("Trying to login user")
	loginForm := models.LoginForm{r.FormValue("email"), r.FormValue("password")}
	user := services.FindUserByEmail(loginForm.GetEmail())
	if user == nil {
		errorResponce(w, "Wrong email or password!", http.StatusBadRequest)
		log.Println("Couldn't find a user with this email")
		return
	}
	if encoder.ComparePasswords(user.GetPassword(), loginForm.GetPassword()) {
		log.Println("Login success")
		err := services.NewSession(w, user.GetId())
		if err != nil {
			panic("Error while trying to generate session id!")
		}
		user.SetScore(services.GetBestScoreForUserById(user.GetId()))
		result, _ := json.Marshal(user)
		fmt.Fprint(w, string(result))
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Wrong password")
		errorResponce(w, "Wrong email or password!", http.StatusBadRequest)
	}
}

func (RegistrationHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Registration failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	if !checkMethod(w, r) {return}
	log.Println("Trying to register new user")
	username := r.FormValue("username")
	if username == "" {
		log.Println("Empty username")
		errorResponce(w, "Specify a correct login!", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	check, err := regexp.MatchString("(.*)@(.*)", email)
	if err != nil {
		log.Println("Email matching error")
		errorResponce(w, "Email matching error!", http.StatusBadRequest)
		return
	}
	if email == "" || !check {
		log.Println("Empty email")
		errorResponce(w, "Specify a valid e-mail!", http.StatusBadRequest)
		return
	}
	if services.FindUserByEmail(email) != nil {
		log.Println("User with this email already exists!")
		errorResponce(w, "User with this email already exists!", http.StatusConflict)
		return
	}
	password := r.FormValue("password")
	if password == "" || len(password) < 4 {
		log.Println("Wrong password length")
		errorResponce(w, "The password field must contain more than 4 characters!", http.StatusBadRequest)
		return
	}
	repPassword := r.FormValue("password_repeat")
	if repPassword == "" || repPassword != password {
		log.Println("Passwords do not match")
		errorResponce(w, "Repeat password correctly!", http.StatusBadRequest)
		return
	}
	r.ParseMultipartForm(32 << 20)
	avatarFile, avatarHeader, err := r.FormFile("avatar")
	if err != nil{
		log.Println("Error reading avatar file!", err)
	} else {
		defer avatarFile.Close()
	}

	avatarName := username + "_avatar"
	if avatarFile != nil && avatarHeader.Filename != "" {
		avatarSave, err := os.OpenFile(avatarPath + avatarName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Save avatar error!")
			return
		}
		defer avatarSave.Close()
		io.Copy(avatarSave, avatarFile)
	}
	user := new(models.User)
	user.SetEmail(email)
	user.SetPassword(encoder.HashPassword(password))
	user.SetUsername(username)
	user.SetAvatar(avatarName)
	services.AddUser(*user)
	log.Println("Register success!")
	result, _ := json.Marshal(user)
	fmt.Fprint(w, string(result))
	w.WriteHeader(http.StatusCreated)
}

func (handle MeHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handle.authentication(w, r)
	} else {
		handle.profileChange(w, r)
	}
}

func (MeHandle) authentication(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Registration failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	log.Println("Trying to authenticate user.")
	cookie := checkAuth(r, w)
	if cookie == nil {
		return
	}
	user := services.GetUserBySessionId(cookie.Value)
	if user == nil {
		log.Println("Can't find such user")
		errorResponce(w, "The user isn't authorized!", http.StatusUnauthorized)
		return
	}
	user.SetScore(services.GetBestScoreForUserById(user.GetId()))
	result, _ := json.Marshal(user)
	fmt.Fprint(w, string(result))
	w.WriteHeader(http.StatusOK)
}

func (MeHandle) profileChange(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to change user profile data.")
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("Registration failed.", rec)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	if !checkMethod(w, r) {return}
	cookie := checkAuth(r, w)
	if cookie == nil {
		return
	}
	currentUser := services.GetUserBySessionId(cookie.Value)
	password := r.FormValue("current_password")
	if !encoder.ComparePasswords(currentUser.GetPassword(), password) {
		log.Println("Wrong password")
		errorResponce(w, "Wrong email or password!", http.StatusBadRequest)
	}
	username := r.FormValue("username")
	if username != "" && username != currentUser.GetUsername() {
		currentUser.SetUsername(username)
	}
	email := r.FormValue("email")
	if email != "" && email != currentUser.GetEmail() {
		check, err := regexp.MatchString("(.*)@(.*)", email)
		if err != nil {
			log.Println("Email matching error")
			errorResponce(w, "Specify a valid e-mail!", http.StatusBadRequest)
			return
		}
		if !check {
			log.Println("Not correct email!")
			errorResponce(w, "Specify a valid e-mail!", http.StatusBadRequest)
			return
		}
		if services.FindUserByEmail(email) != nil {
			log.Println("User with this email already exists!")
			errorResponce(w, "User with this email already exists!", http.StatusConflict)
			return
		}
		currentUser.SetEmail(email)
	}

	newPassword := r.FormValue("new_password")
	if newPassword != "" {
		if len(newPassword) < 4 {
			log.Println("Wrong password length")
			errorResponce(w, "The password field must contain more than 4 characters!", http.StatusBadRequest)
			return
		}
		repPassword := r.FormValue("new_password_repeat")
		if repPassword == "" || repPassword != newPassword {
			log.Println("Passwords do not match")
			errorResponce(w, "Repeat password correctly!", http.StatusBadRequest)
			return
		}
		currentUser.SetPassword(encoder.HashPassword(newPassword))
	}

	r.ParseMultipartForm(32 << 20)
	avatarFile, avatarHeader, err := r.FormFile("avatar")
	if err != nil{
		log.Println("Error reading avatar file!", err)
	} else {
		defer avatarFile.Close()
		if avatarFile != nil && avatarHeader.Filename != "" {
			avatarName := username + "_avatar"
			if err != nil {
				log.Println("Create avatar path", err)
				return
			}
			avatarSave, err := os.OpenFile(avatarPath + avatarName, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Println("Save avatar error!", err)
				return
			}
			defer avatarSave.Close()
			io.Copy(avatarSave, avatarFile)
			currentUser.SetAvatar(avatarName)
		}
	}

	services.SaveUser(*currentUser)
	log.Println("Successful user data change!")
	result, _ := json.Marshal(currentUser)
	fmt.Fprint(w, string(result))
	w.WriteHeader(http.StatusOK)
}

func (ScoreboardHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.FormValue("listNumber"))
	if err != nil {
		log.Println("Page value not int")
		errorResponce(w, "Page value not int!", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(r.FormValue("listSize"))
	if err != nil {
		log.Println("Page value not int")
		errorResponce(w, "Page value not int!", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Println("Page value not int")
		errorResponce(w, "Page value not int!", http.StatusBadRequest)
		return
	}
	if (page < 1) {
		log.Println("Page value < 1")
		errorResponce(w, "This sheet may not be formed!", http.StatusBadRequest)
		return
	}
	startPosition := (page - 1) * size
	usersForScoreBoard := services.FindAllUsersForScoreBoard()
	var numericResult []*models.ScoreRecord
	if len(usersForScoreBoard) < startPosition {
		numericResult = usersForScoreBoard[startPosition:]
	} else {
		numericResult = usersForScoreBoard[startPosition:startPosition+len(usersForScoreBoard)]
	}
	ScoreBoardData := EntryScore{numericResult, len(usersForScoreBoard)}
	result, _ := json.Marshal(ScoreBoardData)
	fmt.Fprint(w, string(result))
	w.WriteHeader(http.StatusOK)
}

func init() {
	fmt.Println("init in RestApiController.go")
	Handlers["/api/login"] = LoginHandle{}
	Handlers["/api/registration"] = RegistrationHandle{}
	Handlers["/api/me"] = MeHandle{}
	Handlers["/api/users/scoreboard"] = ScoreboardHandle{}
	Handlers["/api/logout"] = LogOutHandle{}
	Handlers["/avatar"] = AvatarHandle{}
}

func checkMethod(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		errorResponce(w, "method error", http.StatusBadRequest)
		return false
	}
	return true
}

func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, accept, authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func checkAuth(r *http.Request, w http.ResponseWriter) *http.Cookie {
	cookie, err := r.Cookie(services.CookieName)
	if err != nil || cookie.Value == "" {
		log.Println("The user isn't authorized!")
		errorResponce(w, "The user isn't authorized!", http.StatusUnauthorized)
		return nil
	}
	return cookie
}

func errorResponce(w http.ResponseWriter,  errorMessage string, errorStatus int) {
	err := &JSONError{errorMessage, errorStatus}
	result, _ := json.Marshal(err)
	http.Error(w, string(result), errorStatus)
	w.WriteHeader(errorStatus)
}