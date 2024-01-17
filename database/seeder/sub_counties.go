package seeder

import (
	"context"
	_ "database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"

	_ "github.com/lib/pq"
)

type SubCountiesSeeder struct{}

func (subCountiesSeeder SubCountiesSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/SUBCOUNTIES SEEDER]", "=======Seeding sub counties======", logger.LOG_LEVEL_INFO)

	// Open the CSV file
	file, err := os.Open("assets/files/subcounties.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	count,err := q.CheckSubCountiesDuplicate(context.Background(), "Tetu");
	if err == nil {
		if count == 0 {

	for _, record := range records {
		countyIDStr := record[0]
		name := record[1] 

		countyID, err := strconv.ParseInt(countyIDStr, 10, 32)
		q.InsertSubcounties(context.Background(), gen.InsertSubcountiesParams{
			Name: name,
			CountyID: int32(countyID),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	}else{
		return
	}

	fmt.Println("Data inserted successfully.")
}
}
