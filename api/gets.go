package api

import (
	"database/sql"
	"net/http"
	"net/url"
	"outdoorsy/daos"
	"strconv"

	"github.com/gorilla/mux"
)

var allowedSortValues = map[string]bool{
	"price_asc":    true,
	"price_desc":   true,
	"year_asc":     true,
	"year_desc":    true,
	"make_asc":     true,
	"make_desc":    true,
	"type_asc":     true,
	"type_desc":    true,
	"created_asc":  true,
	"created_desc": true,
	"updated_asc":  true,
	"updated_desc": true,
}

func RentalEndPoint(w http.ResponseWriter, r *http.Request) {
	dao, ok := r.Context().Value("dao").(daos.DaoInterface)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "unable to parse or retrieve dao from context")
		return
	}
	rentalIDStr, err := url.PathUnescape(mux.Vars(r)["rentalID"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing rental ID")
		return
	}
	rentalID, err := strconv.Atoi(rentalIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid rental ID")
		return
	}

	rental, err := dao.GetRentalByID(rentalID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "no rental found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, rental)
}

func MultiRentalEndPoint(w http.ResponseWriter, r *http.Request) {
	dao, ok := r.Context().Value("dao").(daos.DaoInterface)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "unable to parse or retrieve dao from context")
		return
	}
	// Get the query parameters
	queryParams, err := validateQuerySortParameters(r.URL.Query())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	rentals, err := dao.GetRentals(queryParams.PriceMin, queryParams.PriceMax, queryParams.Limit,
		queryParams.Offset, queryParams.IdsSlice, queryParams.NearSlice,
		queryParams.Sort)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "no rental found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, rentals)
}
