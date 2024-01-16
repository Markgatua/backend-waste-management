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
		Name:        "Aggregator Admin",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          4,
		Name:        "Aggregator System User",
		GuardName:   "Web",
		IsActive: true,
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          5,
		Name:        "Aggregator Collector",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          6,
		Name:        "External Collector",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          7,
		Name:        "Green Champion",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          8,
		Name:        "Corporate Global Admin",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          9,
		Name:        "Corporate Admin",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          10,
		Name:        "Corporate System Admin",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
	q.CreateRole(context.Background(), gen.CreateRoleParams{
		ID:          11,
		Name:        "Green Regulator",
		GuardName:   "Web",
		Description: null.StringFrom("").NullString,
	})
}