package main

import (
	"fmt"
	"os"

	"github.com/daver76/go-yolink"
)

func main() {
	client, err := yolink.NewAPIClient(
		os.Getenv("YOLINK_UAID"),
		os.Getenv("YOLINK_SECRET_KEY"),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("client:", client)

	home_id, err := client.GetHomeId()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("home_id:", home_id)
}
