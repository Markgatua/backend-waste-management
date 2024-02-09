package seeder

import (
	"context"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"

	_ "github.com/lib/pq"
)

type PickupTimeStampsSeeder struct{}

func (pickupTimeStampsSeeder PickupTimeStampsSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/PICKUP TIME STAMPS SEEDER]", "=======Seeding pickup time stamps======", logger.LOG_LEVEL_INFO)
	q.InsertPickupTimeStsmp(context.Background(), gen.InsertPickupTimeStsmpParams{
		ID: 1,
		Stamp: "Morning",
		TimeRange: "8AM - 11Am",
		Position: 1,
	})
	q.InsertPickupTimeStsmp(context.Background(), gen.InsertPickupTimeStsmpParams{
		ID: 2,
		Stamp: "Afternoon",
		TimeRange: "12PM - 5PM",
		Position: 2,
	})
}

