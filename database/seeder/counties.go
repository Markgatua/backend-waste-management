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

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	
	// if len(count) > 0{
	// 	return;
	// }
	// Prepare and execute SQL statements to insert records into the PostgreSQL table

	for _, record := range records {
		name := record[0] // Assuming the name column is at index 0 in the CSV file

		count, err := gen.REPO.ViewCounties(context.Background())
		fmt.Println(count)
		q.InsertCounties(context.Background(), name)
		// ("INSERT INTO counties (name) VALUES ($1)", name)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data inserted successfully.")
}
