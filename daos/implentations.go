package daos

import (
	"fmt"
	"log"
	"os"
	"outdoorsy/models"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DAO struct {
	dbClient *sqlx.DB
}

func NewDAO() *DAO {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PGSQL_HOST"), os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_USER"), os.Getenv("PGSQL_PASS"),
		os.Getenv("PGSQL_DB"))

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	dao := DAO{
		dbClient: db,
	}
	return &dao
}

func (d *DAO) CloseDB() {
	err := d.dbClient.Close()
	if err != nil {
		log.Fatal("unable to close DB connection")
	}
}

func (d *DAO) GetRentalByID(rentalID int) (*models.Rental, error) {
	stmt, err := d.dbClient.Preparex("SELECT * FROM rentals r, users u WHERE r.id=$1 AND r.user_id=u.id")
	if err != nil {
		return &models.Rental{}, err
	}

	// Execute the prepared statement, passing in an id value for the
	// parameter whose placeholder is ?
	rental := models.Rental{}
	err = stmt.Get(&rental, rentalID)
	if err != nil {
		return &rental, err
	}
	return &rental, nil
}
