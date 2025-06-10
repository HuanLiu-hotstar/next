package main

import (
	"fmt"
	"regexp"
)

func main() {
	query := `query (
		$token: String
		$countryCode: String
		$platformCode: String
	      ) {
		widget15: ymalCollection(
		  itemId:1770005028
		  token: $token
		  countryCode: $countryCode
		  platformCode: $platformCode
		) {}
		widget12:{}
		`
	r, _ := regexp.Compile("widget[0-9]*")
	fmt.Println(r.MatchString("peach"))
	fmt.Println(r.FindAllString(query, -1))
}
