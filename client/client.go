package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/urlshortener/shortener"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()

	c := shortener.NewUrlServiceClient(conn)

	for {
		fmt.Println("\n*******to Create type: C orig/url\n*******to Get type: G short/url")
		fmt.Print("\n***********************Enter method, url: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if strings.HasPrefix(scanner.Text(), "C") {
			// slice first two symbols
			message := shortener.Message{Body: scanner.Text()[2:]}

			response, err := c.Create(context.Background(), &message)
			if err != nil {
				log.Fatalf("Error when calling Create: %s\n", err)
			}
			log.Printf("Response from Server: %s\n", response.Body)

		} else if strings.HasPrefix(scanner.Text(), "G") {
			// slice first two symbols
			message := shortener.Message{Body: scanner.Text()[2:]}

			response, err := c.Get(context.Background(), &message)
			if err != nil {
				log.Fatalf("Error when calling Create: %s\n", err)
			}
			log.Printf("Response from Server: %s\n", response.Body)
		}
	}
}
