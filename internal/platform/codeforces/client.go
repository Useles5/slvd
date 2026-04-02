package codeforces

import (
	"fmt"
	"io"
	"net/http"
)

func FetchRecent(handle string) {
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=10", handle)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to fetch codeforces data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %v\n", err)
		return
	}
	fmt.Println(string(body))
}
