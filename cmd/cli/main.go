package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/amettod/hourly-meter/internal/app"
)

func main() {
	filename := flag.String("filename", "", "filename *.html")
	contract := flag.String("contract", "", "contract number")
	companyName := flag.String("name", "", "company name")
	meter := flag.String("meter", "", "electronic meter serial number")
	coefficient := flag.Float64("coefficient", 1, "power factory")

	flag.Parse()

	app, err := app.New(*filename, *contract, *companyName, *meter, *coefficient)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal()
	}

	fmt.Printf("total:\t%.2f kWh\nvalues:\t%d\n", app.Total, len(app.Rows))
}
