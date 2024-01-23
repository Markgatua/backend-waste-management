package seeder

import (
	"context"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"

	"gopkg.in/guregu/null.v3"
)

type RoleSeeder struct{}

func (roleSeeder RoleSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/ROLE SEEDER]", "=======Seeding roles======", logger.LOG_LEVEL_INFO)
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          1,
		Name:        "TTNM Admin",
		GuardName:   "super_admin_web",
		IsActive: true,
		Description: null.StringFrom("Super Admin can perform most actions in the system").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          2,
		Name:        "Aggregator Global Admin",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})	
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          3,
		Name:        "Aggregator Super Admin",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          4,
		Name:        "Aggregator admin",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          5,
		Name:        "Aggregator system user",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})


	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          6,
		Name:        "Green champion global admin",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})	
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          7,
		Name:        "Green champion super admin",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          8,
		Name:        "Green champion admin",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})

	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          9,
		Name:        "Green champion system user",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
	
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          10,
		Name:        "External Collector",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})

	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          12,
		Name:        "TTNM super Admin",
		GuardName:   "super_admin_web",
		IsActive: true,
		Description: null.StringFrom("Super Admin can perform most actions in the system").NullString,
	})
}