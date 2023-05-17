package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"outdoorsy/daos"
	"strconv"
	"strings"

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
	queryParams := r.URL.Query()

	var err error

	priceMin := 0
	priceMax := 0
	limit := 0
	offset := 0
	var idsSlice []int
	var nearSlice []float64

	priceMinStr := queryParams.Get("price_min")
	if priceMinStr != "" {
		priceMin, err = strconv.Atoi(priceMinStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing price_min: "+err.Error())
			return
		}

	}

	priceMaxStr := queryParams.Get("price_max")
	if priceMaxStr != "" {
		priceMax, err = strconv.Atoi(priceMaxStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing price_max: "+err.Error())
			return
		}
	}

	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing limit: "+err.Error())
			return
		}
	}

	offsetStr := queryParams.Get("offset")
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing offset: "+err.Error())
			return
		}
	}

	ids := queryParams.Get("ids")
	if len(ids) > 0 {
		idsStrSlice := strings.Split(ids, ",")
		fmt.Println("ids split", idsStrSlice)
		if len(idsStrSlice) > 0 {
			for _, idStr := range idsStrSlice {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid value for ids: %s", idStr))
					return
				}
				idsSlice = append(idsSlice, id)
			}
		}

	}

	nearValues := queryParams.Get("near")
	if len(nearValues) > 0 {
		nearStrSlice := strings.Split(nearValues, ",")
		if len(nearStrSlice) == 2 {
			for _, nearStr := range nearStrSlice {
				nearFloat, err := strconv.ParseFloat(nearStr, 64)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid value for near: %s", nearStr))
					return
				}
				nearSlice = append(nearSlice, nearFloat)
			}
		} else {
			respondWithError(w, http.StatusBadRequest, "wrong number of coordinates in near param")
			return
		}
	}

	sort := queryParams.Get("sort")
	if sort != "" {
		if _, ok := allowedSortValues[sort]; !ok {
			respondWithError(w, http.StatusBadRequest, "invalid sort value: "+sort)
			return
		}
	}

	rentals, err := dao.GetRentals(priceMin, priceMax, limit, offset, idsSlice, nearSlice, sort)
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
