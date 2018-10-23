package DA

import (
	"IT-Berries_Go_server/auth/models"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func FindUserByEmail(email string) []*models.User {
	db := connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE email = $1", email)
	errorCheck(err, executeError)
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(user.GetIdPoint(), user.GetAvatarPoint(), user.GetEmailPoint(), user.GetPasswordPoint(), user.GetUsernamePoint())
		errorCheck(err, readRowError)
		users = append(users, user)
	}
	errorCheck(rows.Err(), readRowError)
	return users
}

func FindUserById(userId int) []*models.User {
	db := connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users WHERE user_id = $1", userId)
	errorCheck(err, executeError)
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(user.GetIdPoint(), user.GetAvatarPoint(), user.GetEmailPoint(), user.GetPasswordPoint(), user.GetUsernamePoint())
		errorCheck(err, readRowError)
		users = append(users, user)
	}
	errorCheck(rows.Err(), readRowError)
	return users
}

func AddUser(user models.User) int64 {
	db := connect()
	defer db.Close()
	result, err := db.Query("insert into users (avatar, email, password, user_name) values ($1, $2, $3, $4) RETURNING user_id",
		user.GetAvatar(), user.GetEmail(), user.GetPassword(), user.GetUsername())
	errorCheck(err, executeError)
	defer result.Close()
	var id int64
	result.Next()
	err = result.Scan(&id)
	errorCheck(err, idExtracionError)
	_, err = db.Exec("insert into history (date_result, score, user_id) values ($1, $2, $3)",
		time.Now(), 0, id)
	errorCheck(err, executeError)
	return id
}

func SaveUser(user models.User) int64 {
	db :=connect()
	defer db.Close()
	result, err := db.Exec("UPDATE users SET avatar = $1, email = $2, password = $3, user_name = $4",
		user.GetAvatar(), user.GetEmail(), user.GetPassword(), user.GetUsername())
	errorCheck(err, executeError)
	id, err := result.RowsAffected()
	errorCheck(err, idExtracionError)
	return id
}

func errorCheck(err error, message string) {
	if err != nil {
		log.Println(message, err)
		panic(err)
	}
}

func connect() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	errorCheck(err, connectError)
	return db
}
