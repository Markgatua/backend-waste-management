package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type CountriesSeeder struct{}

func (countriesSeeder CountriesSeeder) Run(q *queries.DBQuerier) {
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

				q.CreateCountry(context.Background(), queries.CreateCountryParams{
					Name:            pgtype.Varchar{String: fmt.Sprint(name),Status: pgtype.Present},
					CurrencyCode:    pgtype.Varchar{String: fmt.Sprint(currencyCode),Status: pgtype.Present},
					Capital:         pgtype.Varchar{String: fmt.Sprint(capital),Status: pgtype.Present},
					Citizenship:     pgtype.Varchar{String: fmt.Sprint(citizenship),Status: pgtype.Present},
					CountryCode:     pgtype.Varchar{String: fmt.Sprint(countryCode),Status: pgtype.Present},
					Currency:        pgtype.Varchar{String: fmt.Sprint(currency),Status: pgtype.Present},
					CurrencySubUnit: pgtype.Varchar{String: fmt.Sprint(currencySubUnit),Status: pgtype.Present},
					CurrencySymbol:  pgtype.Varchar{String: fmt.Sprint(currencySymbol),Status: pgtype.Present},
					FullName:        pgtype.Varchar{String: fmt.Sprint(fullName),Status: pgtype.Present},
					Iso31662:        pgtype.Varchar{String: fmt.Sprint(iso_3166_2),Status: pgtype.Present},
					Iso31663:        pgtype.Varchar{String: fmt.Sprint(iso_3166_3),Status: pgtype.Present},
					RegionCode:      pgtype.Varchar{String: fmt.Sprint(regionCode),Status: pgtype.Present},
					SubRegionCode:   pgtype.Varchar{String: fmt.Sprint(subRegionCode),Status: pgtype.Present},
					//Eea:              eea.(pgtype.Int2),
					CallingCode: pgtype.Varchar{String: callingCode.(string),Status: pgtype.Present},
					Flag:        pgtype.Varchar{String: fmt.Sprint(flag),Status: pgtype.Present},
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
