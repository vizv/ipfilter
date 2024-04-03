package sync

import (
	"fmt"
	"io"
	"net/http"
)

func Download(url string) ([]byte, error) {
	client := http.Client{}

	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf(`failed to download "%s": %+v`, url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(`failed to download "%s" with status: %d`, url, res.StatusCode)
	}

	return io.ReadAll(res.Body)
}
