package daos

import "outdoorsy/models"

type DaoInterface interface {
	GetRentalByID(rentalID int) (*models.Rental, error)
}
