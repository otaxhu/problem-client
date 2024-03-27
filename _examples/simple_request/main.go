package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	problem_client "github.com/otaxhu/problem-client"
)

func main() {
	res, err := http.Get("https://your-api-that-returns-problem-json.com")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	prob, extra, err := problem_client.ProblemResponse(res)
	if err != nil {
		panic(err)
	}
	if prob == nil {
		// There is no problem details, we can also check that the API returns a OK status code
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			bussinessLogic()
		} else {
			fmt.Printf("API returned %d status code but there is no problem details, got following response:\n", res.StatusCode)

			// Writing to the terminal the body
			_, err := io.Copy(os.Stdout, res.Body)
			if err != nil {
				panic(err)
			}
		}
	} else {
		// Here we got a Problem details structure and possibly a Extension members map
		fmt.Printf("Problem details: %v\nExtension members: %v\n", prob, extra)
	}
}

func bussinessLogic() {
	fmt.Println("bussiness logic")
}
