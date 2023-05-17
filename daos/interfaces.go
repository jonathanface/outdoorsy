package daos

import "outdoorsy/models"

type DaoInterface interface {
	GetRentalByID(rentalID int) (*models.Rental, error)
	GetRentals(priceMin, priceMax, limit, offset int, ids []int, near []float64, sort string) ([]*models.Rental, error)
}
