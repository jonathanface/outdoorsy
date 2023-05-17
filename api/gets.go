package api

import (
	"database/sql"
	"net/http"
	"net/url"
	"outdoorsy/daos"
	"strconv"

	"github.com/gorilla/mux"
)

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
