package DA

import "fmt"

const (
	host     = "db"
	port     = 5432
	user     = "it_berries"
	password = "root"
	dbname   = "it_berries"

	connectError = "An error occurred while trying to connect to the database."
	pignError    = "An error occurred while trying to ping database"
	executeError = "An error occurred while trying to execute query."
	readRowError = "An error occurred while trying to read row"
	idExtracionError ="Id extraction error"
)

var psqlInfo string

func init() {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db := connect()
	defer db.Close()
	err := db.Ping()
	errorCheck(err, pignError)
	fmt.Println("Successfully connected to database.")
}