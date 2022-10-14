package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple static webserver:
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/home/tonycark/Desktop/acronis/find"))))
	//change /home/tonycark/Desktop/acronis/find to the directory the files are located
}
