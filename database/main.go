package database
import (
	"database/sql"
	"github.com/lib/pq"
	"fmt"
)
//InitialDatabase is to initial postgresQL
func InitialDatabase(url string) (*sql.DB, error){
	pgURL, err := pq.ParseURL(url)
	if err != nil {
		fmt.Println("was not able to connect to the database")
		return nil, err
	}
	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		fmt.Println("was not able to connect to the database")
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		fmt.Println("ping did not work")
		return nil, err
	}

	return db, nil
}