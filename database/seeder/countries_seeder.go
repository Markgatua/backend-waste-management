package seeder

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
)

type CountriesSeeder struct{}

func (countriesSeeder CountriesSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/COUNTRIES SEEDER]", "=======Seeding countries======", logger.LOG_LEVEL_INFO)

	jsonFile, err := os.Open("assets/data/countries.json")
	if err == nil {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var result map[string]map[string]interface{}
		unmarshalError := json.Unmarshal(byteValue, &result)
		if unmarshalError == nil {
			for _, v := range result {
				capital := v["capital"]
				citizenship := v["citizenship"]
				countryCode := v["country-code"]
				currency := v["currency"]
				currencyCode := v["currency_code"]
				currencySubUnit := v["currency_sub_unit"]
				fullName := v["full_name"]
				iso_3166_2 := v["iso_3166_2"]
				iso_3166_3 := v["iso_3166_3"]
				name := v["name"]
				regionCode := v["region-code"]
				subRegionCode := v["sub-region-code"]

				//eea := v["eea"]
				callingCode := v["calling_code"]
				currencySymbol := v["currency_symbol"]
				//currencyDecimals := v["currency_decimals"]
				flag := v["flag"]

				q.CreateCountry(context.Background(), gen.CreateCountryParams{
					Name:            fmt.Sprint(name),
					CurrencyCode:    sql.NullString{String: fmt.Sprint(currencyCode), Valid: true},
					Capital:         sql.NullString{String: fmt.Sprint(capital), Valid: true},
					Citizenship:     fmt.Sprint(citizenship),
					CountryCode:     fmt.Sprint(countryCode),
					Currency:        sql.NullString{String: fmt.Sprint(currency), Valid: true},
					CurrencySubUnit: sql.NullString{String: fmt.Sprint(currencySubUnit), Valid: true},
					CurrencySymbol:  sql.NullString{String: fmt.Sprint(currencySymbol), Valid: true},
					FullName:        sql.NullString{String: fmt.Sprint(fullName), Valid: true},
					Iso31662:        fmt.Sprint(iso_3166_2),
					Iso31663:        fmt.Sprint(iso_3166_3),
					RegionCode:      fmt.Sprint(regionCode),
					SubRegionCode:   fmt.Sprint(subRegionCode),
					CallingCode:     sql.NullString{String: callingCode.(string), Valid: true},
					Flag:            sql.NullString{String: fmt.Sprint(flag), Valid: true},
				})
				//fmt.Println(err.Error())
			}
		} else {
			fmt.Println(unmarshalError.Error())
		}
	} else {
		fmt.Println("Error reading from json file -- ", err.Error())
	}

	defer jsonFile.Close()
}
