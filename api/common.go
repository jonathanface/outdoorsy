package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	var (
		response []byte
		err      error
	)
	if response, err = json.Marshal(payload); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type queryParams struct {
	PriceMin  int
	PriceMax  int
	Limit     int
	Offset    int
	IdsSlice  []int
	NearSlice []float64
	Sort      string
}

func validateQuerySortParameters(values url.Values) (params queryParams, err error) {
	priceMinStr := values.Get("price_min")
	if priceMinStr != "" {
		params.PriceMin, err = strconv.Atoi(priceMinStr)
		if err != nil {
			return params, errors.New("error parsing min price: " + err.Error())
		}

	}

	priceMaxStr := values.Get("price_max")
	if priceMaxStr != "" {
		params.PriceMax, err = strconv.Atoi(priceMaxStr)
		if err != nil {
			return params, errors.New("error parsing max price: " + err.Error())
		}
	}

	limitStr := values.Get("limit")
	if limitStr != "" {
		params.Limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return params, errors.New("error parsing limit: " + err.Error())
		}
	}

	offsetStr := values.Get("offset")
	if offsetStr != "" {
		params.Offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return params, errors.New("error parsing offset: " + err.Error())
		}
	}

	ids := values.Get("ids")
	if len(ids) > 0 {
		idsStrSlice := strings.Split(ids, ",")
		if len(idsStrSlice) > 0 {
			for _, idStr := range idsStrSlice {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return params, fmt.Errorf("Invalid value for ids: %s", idStr)
				}
				params.IdsSlice = append(params.IdsSlice, id)
			}
		}

	}

	nearValues := values.Get("near")
	if len(nearValues) > 0 {
		nearStrSlice := strings.Split(nearValues, ",")
		if len(nearStrSlice) == 2 {
			for _, nearStr := range nearStrSlice {
				nearFloat, err := strconv.ParseFloat(nearStr, 64)
				if err != nil {
					return params, fmt.Errorf("Invalid value for near: %s", nearStr)
				}
				params.NearSlice = append(params.NearSlice, nearFloat)
			}
		} else {
			return params, fmt.Errorf("wrong number of coordinates in near param")
		}
	}

	sort := values.Get("sort")
	if sort != "" {
		if _, ok := allowedSortValues[sort]; !ok {
			return params, fmt.Errorf("invalid sort value: %s", sort)
		}
	}
	return
}
