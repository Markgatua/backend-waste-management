package seeder

import (
	"context"
	_ "database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"

	_ "github.com/lib/pq"
)

type CountiesSeeder struct{}

func (countiesSeeder CountiesSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/COUNTIES SEEDER]", "=======Seeding counties======", logger.LOG_LEVEL_INFO)

	// Open the CSV file
	file, err := os.Open("assets/files/counties.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	count,err := q.DuplicateCounties(context.Background(),"Nyeri");
	if err == nil {
		if count == 0 {

	for _, record := range records {
		name := record[0] 

		q.InsertCounties(context.Background(), name)
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
