package main

import (
	"strings"
)

func ContainString(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}

	return false
}

func RemoveParams(key string, sourceURL string) string {
	// split to get the url before query
	rtn := strings.Split(sourceURL, "?")[0]
	queryString := ""

	// get the query url
	if strings.Contains(sourceURL, "?") {
		queryString = strings.Split(sourceURL, "?")[1]
	}

	if queryString != "" {
		// split querys into sub query
		params_arr := strings.Split(queryString, "&")

		// delete query related to key
		for i := len(params_arr) - 1; 0 <= i; i-- {
			param := strings.Split(params_arr[i], "=")[0]
			if param == key {
				params_arr = params_arr[0:i]
			}
		}

		// then add query not related to key into url
		if len(params_arr) != 0 {
			rtn += "?" + strings.Join(params_arr, "&")
		}
	}

	return rtn
}
