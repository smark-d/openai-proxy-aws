package api

import (
	"github.com/go-redis/redis"
	"net/http"
)

func Filter(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization") // get Authorization header
	if authHeader == "" {
		return false // if Authorization header is empty, return false
	}

	customKey := authHeader[7:] // get customKey after "Bearer "

	currCount, err := GetCurrCount(customKey)

	if err != nil && err != redis.Nil {
		panic(err)
	}
	totalCount, err := GetTotalCount(customKey)
	if err != nil && err != redis.Nil {
		panic(err)
	}
	// assuming the total number of requests allowed is stored in a variable called "totalRequests"
	if currCount >= totalCount {
		return false // if request count is greater than or equal to total requests, return false
	}

	IncrCurrCount(customKey) // increment request count for customKey
	return true              // return true since request count is less than total requests

}
