package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Moinsen Welt")

	req, err := http.NewRequest("GET", "https://localhost:5001/api/search", nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	query := req.URL.Query()
	query.Add("blNumber", "123")
	req.URL.RawQuery = query.Encode()

	fmt.Println(req.URL.String())
	
	res, err := http.Get(req.URL.String())
	if err != nil {
		log.Fatal(err)	
	}

	fmt.Println("Worked")
	fmt.Println(res.Status)

}
