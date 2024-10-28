package client

import (
	"net/http"
	"strconv"
)

func respRetryAfter(resp *http.Response) (uint, error) {
	v, err := strconv.ParseUint(resp.Header.Values("Retry-After")[0], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
