package daos

import (
	"outdoorsy/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDAO_GetRentalByID(t *testing.T) {
	// Create a new mock DB connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create the DAO with the mock DB connection
	dao := &DAO{
		dbClient: sqlx.NewDb(db, "sqlmock"),
	}

	// Define the test cases in a table format
	testCases := []struct {
		Name           string
		RentalID       int
		PrepareMock    func()
		ExpectedResult *models.Rental
		ExpectedError  error
	}{
		{
			Name:     "Successful retrieval",
			RentalID: 1,
			PrepareMock: func() {
				// Define the expected query and result rows
				columns := []string{"id", "name", "price_per_day"}
				rows := sqlmock.NewRows(columns).AddRow(1, "Test Rental", 100.0)

				// Set up the expectations
				mock.ExpectPrepare(`SELECT r\.\*, u\.id AS sub_user_id, u\.first_name, u\.last_name FROM rentals r, users u WHERE r\.id = \$1`).
					ExpectQuery().
					WithArgs(1).
					WillReturnRows(rows)
			},
			ExpectedResult: &models.Rental{
				Id:   1,
				Name: "Test Rental",
				RentalPrice: models.RentalPrice{
					Day: 100,
				},
			},
			ExpectedError: nil,
		},
		// Add more test cases for different scenarios if needed
	}

	// Run the test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// Prepare the mock DB with the expected behavior
			testCase.PrepareMock()

			// Call the function under test
			rental, err := dao.GetRentalByID(testCase.RentalID)

			// Verify the results
			assert.Equal(t, testCase.ExpectedError, err)
			assert.Equal(t, testCase.ExpectedResult, rental)
		})
	}
}

func TestGetRentals(t *testing.T) {
	// Prepare test data
	priceMin := 50
	priceMax := 200
	limit := 10
	offset := 0
	ids := []int{1, 2, 3}
	near := []float64{40.0, -75.0}
	sort := "price_asc"

	// Create a new mock DB connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create the DAO with the mock DB connection
	dao := &DAO{
		dbClient: sqlx.NewDb(db, "sqlmock"),
	}

	// Define the test cases in a table format
	testCases := []struct {
		Name           string
		PrepareMock    func()
		ExpectedResult []*models.Rental
		ExpectedError  error
	}{
		{
			Name: "Successful retrieval",
			PrepareMock: func() {
				// Define the expected query and result rows

				columns := []string{"id", "name", "price_per_day"}
				rows := sqlmock.NewRows(columns).
					AddRow(1, "Rental 1", 100.0).
					AddRow(2, "Rental 2", 150.0).
					AddRow(3, "Rental 3", 180.0)

				// Set up the expectations
				mock.ExpectPrepare(`SELECT r\..*, u\.id AS sub_user_id, u\.first_name, u\.last_name FROM rentals r JOIN users u ON u\.id = r\.user_id WHERE price_per_day >= \$1 AND price_per_day <= \$2 AND r\.id IN \(\$3, \$4, \$5\) AND ST_DWithin\(ST_MakePoint\(-75.000000, 40.000000\)::geography, ST_MakePoint\(lng, lat\)::geography, 100.000000 \* 1609.34\) ORDER BY r.price_per_day ASC LIMIT \$6`).
					ExpectQuery().
					WithArgs(limit, priceMin, priceMax, ids[0], ids[1], ids[2]).
					WillReturnRows(rows)
			},
			ExpectedResult: []*models.Rental{
				{Id: 1, UserID: 0, Name: "Rental 1", RentalPrice: models.RentalPrice{Day: 100}},
				{Id: 2, UserID: 0, Name: "Rental 2", RentalPrice: models.RentalPrice{Day: 150}},
				{Id: 3, UserID: 0, Name: "Rental 3", RentalPrice: models.RentalPrice{Day: 180}},
			},
			ExpectedError: nil,
		},
		// Add more test cases for different scenarios if needed
	}

	// Run the test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// Prepare the mock DB with the expected behavior
			testCase.PrepareMock()

			// Call the function under test
			rentals, err := dao.GetRentals(priceMin, priceMax, limit, offset, ids, near, sort)

			// Verify the results
			assert.Equal(t, testCase.ExpectedError, err)
			assert.Equal(t, testCase.ExpectedResult, rentals)
		})
	}
}
