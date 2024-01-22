package seeder

import (
	"context"
	"fmt"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
)

type MainOrganizationSeeder struct{}

func (mainOrganizationSeeder MainOrganizationSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/MainOrganizationSeeder SEEDER]", "======= Seeding MainOrganizationSeeder======", logger.LOG_LEVEL_INFO)

	organizations, err := q.GetMainOrganization(context.Background(), "TTNM_1")
	if err == nil {
		if len(organizations) == 0 {
			err := q.InsertMainOrganization(context.Background(), gen.InsertMainOrganizationParams{
				Name:                   "TakaTaka Ni Mali",
				AboutUs:                "Taka Taka ni Mali is an organisation dedicated to innovation, growth and strengthening of connections in the waste management ecosystem in Kenya, turning it into a circular economy.",
				LogoPath:               "https://takanimali.org/takatakaicon.png",
				TagLine:                "Waste is Wealth",
				WebsiteUrl:             "https://takanimali.org/",
				City:                   "Nairobi",
				Country:                "Kenya",
				Zip:                    "00100",
				State:                  "Nairobi",
				OrganizationID:         "TTNM_1",
				AppAppstoreLink:        "",
				AppGooglePlaystoreLink: "",
			})
			if err != nil {
				fmt.Printf(err.Error())
			}
		} else {
			q.UpdateMainOrganizationProfile(context.Background(), gen.UpdateMainOrganizationProfileParams{
				Name:                   "TakaTaka Ni Mali",
				AboutUs:                "Taka Taka ni Mali is an organisation dedicated to innovation, growth and strengthening of connections in the waste management ecosystem in Kenya, turning it into a circular economy.",
				LogoPath:               "https://takanimali.org/takatakaicon.png",
				TagLine:                "Waste is Wealth",
				WebsiteUrl:             "https://takanimali.org/",
				City:                   "Nairobi",
				Country:                "Kenya",
				Zip:                    "00100",
				State:                  "Nairobi",
				OrganizationID:         "TTNM_1",
				AppAppstoreLink:        "",
				AppGooglePlaystoreLink: "",
			})
		}
	}
}
