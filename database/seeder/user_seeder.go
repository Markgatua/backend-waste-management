package seeder

import (
	"context"
	"database/sql"
	_ "fmt"
	"time"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/helpers"
	"ttnmwastemanagementsystem/logger"

	"gopkg.in/guregu/null.v3"
)

type UserSeeder struct{}

func (userSeeder UserSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/USER SEEDER]", "=======Seeding users======", logger.LOG_LEVEL_INFO)

	err := q.CreateMainOrganizationAdmin(context.Background(), gen.CreateMainOrganizationAdminParams{
		FirstName: null.StringFrom("Super").NullString,
		LastName:  null.StringFrom("Admin").NullString,
		Email:     null.StringFrom("superadmin@admin.com").NullString,
		RoleID:    sql.NullInt32{Int32: 12, Valid: true},
		Provider:  null.StringFrom("email").NullString,
		IsMainOrganizationUser: true,
		ConfirmedAt: sql.NullTime{Time: time.Now(),Valid: true},
		Password:  null.StringFrom(helpers.Functions{}.HashPassword("%$#TYEWY")).NullString,
	})
	if err != nil {
		//logger.Log("[SEEDER/USER SEEDER]", fmt.Sprint(err.Error()), logger.LOG_LEVEL_ERROR)
	}
}
