package daos

import (
	"fmt"
	"log"
	"os"
	"outdoorsy/models"
	"strconv"

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
	stmt, err := d.dbClient.Preparex("SELECT r.*, u.id AS sub_user_id, u.first_name, u.last_name FROM rentals r, users u WHERE r.id = $1")
	if err != nil {
		return nil, err
	}

	rental := models.Rental{}
	err = stmt.Get(&rental, rentalID)
	if err != nil {
		return nil, err
	}
	return &rental, nil
}

func (d *DAO) GetRentals(priceMin, priceMax, limit, offset int, ids []int, near []float64, sort string) (rentals []*models.Rental, err error) {
	query := "SELECT r.*, u.id AS sub_user_id, u.first_name, u.last_name FROM rentals r, users u"
	// Add WHERE clauses for filtering by price range
	if priceMin > 0 {
		query += " WHERE price >= " + strconv.Itoa(priceMin)
	}
	if priceMax > 0 {
		if priceMin > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " price <= " + strconv.Itoa(priceMax)
	}

	// Add WHERE clauses for filtering by rental IDs
	if len(ids) > 0 {
		idStr := ""
		for i, id := range ids {
			if i > 0 {
				idStr += ","
			}
			idStr += strconv.Itoa(id)
		}
		if priceMin > 0 || priceMax > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " id IN (" + idStr + ")"
	}

	// Add WHERE clauses for filtering by proximity
	if len(near) == 2 {
		lat := near[0]
		lng := near[1]
		if priceMin > 0 || priceMax > 0 || len(ids) > 0 {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += fmt.Sprintf(" ST_DWithin(ST_MakePoint(%f, %f)::geography, ST_MakePoint(lng, lat)::geography, 100 * 1609.34)", lng, lat)
	}
	if priceMin > 0 || priceMax > 0 || len(ids) > 0 || len(near) == 2 {
		query += " AND"
	} else {
		query += " WHERE"
	}
	query += " u.id = r.user_id"

	// Add ORDER BY clause for sorting
	switch sort {
	case "price_asc":
		query += " ORDER BY price_per_day ASC"
	case "price_desc":
		query += " ORDER BY price_per_day DESC"
	case "year_asc":
		query += " ORDER BY vehicle_year ASC"
	case "year_desc":
		query += " ORDER BY vehicle_year DESC"
	case "make_asc":
		query += " ORDER BY vehicle_make ASC"
	case "make_desc":
		query += " ORDER BY vehicle_make DESC"
	case "type_asc":
		query += " ORDER BY type ASC"
	case "type_desc":
		query += " ORDER BY type DESC"
	case "created_asc":
		query += " ORDER BY created ASC"
	case "created_desc":
		query += " ORDER BY created DESC"
	case "updated_asc":
		query += " ORDER BY updated ASC"
	case "updated_desc":
		query += " ORDER BY updated DESC"
	}

	// Add LIMIT and OFFSET clauses for pagination
	if limit > 0 {
		query += " LIMIT " + strconv.Itoa(limit)
	}
	if offset > 0 {
		query += " OFFSET " + strconv.Itoa(offset)
	}
	stmt, err := d.dbClient.Preparex(query)
	if err != nil {
		return nil, err
	}
	err = stmt.Select(&rentals)
	if err != nil {
		return nil, err
	}
	return rentals, nil
}
