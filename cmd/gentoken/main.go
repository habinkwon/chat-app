package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	secret := flag.String("secret", "", "")
	payload := flag.String("payload", `{"sub":"user@example.com","userID":1,"role":"USER"}`, "")
	flag.Parse()

	var claims jwt.MapClaims
	err := json.Unmarshal([]byte(*payload), &claims)
	if err != nil {
		log.Fatal(err)
	}

	t := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	s, err := t.SignedString([]byte(*secret))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}
