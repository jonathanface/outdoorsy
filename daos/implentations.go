package daos

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"outdoorsy/models"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DAO struct {
	dbClient *sqlx.DB
}

func NewDAO() (*DAO, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PGSQL_HOST"), os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_USER"), os.Getenv("PGSQL_PASS"),
		os.Getenv("PGSQL_DB"))

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	dao := DAO{
		dbClient: db,
	}
	return &dao, err
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
	// Build the placeholders and arguments for the query
	placeholders := []string{}
	args := []interface{}{}
	argIndex := 1

	if priceMin > 0 {
		placeholders = append(placeholders, "price_per_day >= $"+strconv.Itoa(argIndex))
		args = append(args, priceMin)
		argIndex++
	}
	if priceMax > 0 {
		placeholders = append(placeholders, "price_per_day <= $"+strconv.Itoa(argIndex))
		args = append(args, priceMax)
		argIndex++
	}
	if len(ids) > 0 {
		idPlaceholders := []string{}
		for range ids {
			idPlaceholders = append(idPlaceholders, "$"+strconv.Itoa(argIndex))
			argIndex++
		}
		placeholders = append(placeholders, "id IN ("+strings.Join(idPlaceholders, ", ")+")")
		idArgs := make([]interface{}, len(ids))
		for i, id := range ids {
			idArgs[i] = id
		}

		// Append idArgs to args
		args = append(args, idArgs...)
	}
	if len(near) == 2 {
		lat := near[0]
		lng := near[1]
		radius := 100.0 // in miles
		placeholders = append(placeholders, fmt.Sprintf("ST_DWithin(ST_MakePoint(%f, %f)::geography, ST_MakePoint(lng, lat)::geography, %f * 1609.34)", lng, lat, radius))
	}

	// Construct the WHERE clause from the placeholders
	whereClause := ""
	if len(placeholders) > 0 {
		whereClause = "WHERE " + strings.Join(placeholders, " AND ")
	}

	// Construct the ORDER BY clause based on the sort parameter
	sortClause := ""
	switch sort {
	case "price_asc":
		sortClause = " ORDER BY r.price_per_day ASC"
	case "price_desc":
		sortClause = " ORDER BY r.price_per_day DESC"
	case "year_asc":
		sortClause = " ORDER BY r.vehicle_year ASC"
	case "year_desc":
		sortClause = " ORDER BY r.vehicle_year DESC"
	case "make_asc":
		sortClause = " ORDER BY r.vehicle_make ASC"
	case "make_desc":
		sortClause = " ORDER BY r.vehicle_make DESC"
	case "type_asc":
		sortClause = " ORDER BY r.type ASC"
	case "type_desc":
		sortClause = " ORDER BY r.type DESC"
	case "created_asc":
		sortClause = " ORDER BY r.created ASC"
	case "created_desc":
		sortClause = " ORDER BY r.created DESC"
	case "updated_asc":
		sortClause = " ORDER BY r.updated ASC"
	case "updated_desc":
		sortClause = " ORDER BY r.updated DESC"
	// Add other cases for different sort options
	default:
		// Default sorting if no valid sort option is provided
		sortClause = "ORDER BY created DESC"
	}

	// Construct the final query
	query := fmt.Sprintf("SELECT r.*, u.id AS sub_user_id, u.first_name, u.last_name FROM rentals r JOIN users u ON u.id = r.user_id %s %s", whereClause, sortClause)

	// Add LIMIT and OFFSET clauses for pagination
	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		argIndex++
		fmt.Println("args", args)
		args = append([]interface{}{limit}, args...)
	}
	if offset > 0 {
		query += " OFFSET $" + strconv.Itoa(argIndex)
		argIndex++
		args = append([]interface{}{offset}, args...)
	}

	// Prepare the SQL statement
	stmt, err := d.dbClient.Preparex(query)
	if err != nil {
		return nil, err
	}

	// Execute the SQL statement and fetch the results
	err = stmt.Select(&rentals, args...)
	if err != nil {
		return nil, err
	}
	if len(rentals) == 0 {
		return nil, sql.ErrNoRows
	}

	return rentals, nil
}
