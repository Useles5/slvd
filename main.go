package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	userName := os.Args[1]
	url := fmt.Sprintf("https://codeforces.com/api/user.status?handle=%s&from=1&count=10", userName)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

}
