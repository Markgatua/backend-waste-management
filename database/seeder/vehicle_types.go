package seeder

import (
	"context"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
)

type VehicleTypesSeeder struct{}

func (seeder VehicleTypesSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/VEHICLE TYPES SEEDER]", "=======Seeding vehicle types======", logger.LOG_LEVEL_INFO)
	_,err:=q.CreateVehicleTypes(context.Background(), gen.CreateVehicleTypesParams{
		ID:          1,
		Name:        "light_truck",
		MaxVehicleWeight: 3.5,
		MaxVehicleHeight: 3.2,
		Description: "A light truck, for example, small delivery truck or camping car",
	})
	if err!=nil{
		logger.Log("VehivleTypeSeeder",err.Error(),logger.LOG_LEVEL_ERROR)
	}
	q.CreateVehicleTypes(context.Background(), gen.CreateVehicleTypesParams{
		ID:          2,
		Name:        "medium_truck",
		MaxVehicleWeight: 7.5,
		MaxVehicleHeight: 4.1,
		Description: "A medum-size truck",
	})
	q.CreateVehicleTypes(context.Background(), gen.CreateVehicleTypesParams{
		ID:          3,
		Name:        "truck",
		MaxVehicleWeight: 22,
		MaxVehicleHeight: 4.1,
		Description: "A truck",
	})
	q.CreateVehicleTypes(context.Background(), gen.CreateVehicleTypesParams{
		ID:          4,
		Name:        "heavy_truck",
		MaxVehicleWeight: 40,
		MaxVehicleHeight: 4.1,
		Description: "A heavy truck",
	})
	q.CreateVehicleTypes(context.Background(), gen.CreateVehicleTypesParams{
		ID:          5,
		Name:        "truck_dangerous_goods",
		MaxVehicleWeight: 22,
		MaxVehicleHeight: 4.1,
		Description: "A truck carrying dangerous goods",
	})
	q.CreateVehicleTypes(context.Background(), gen.CreateVehicleTypesParams{
		ID:          6,
		Name:        "long_truck",
		MaxVehicleWeight: 22,
		MaxVehicleHeight: 4.1,
		Description: "A long truck with a maximal lenth of 34m",
	})
	

}