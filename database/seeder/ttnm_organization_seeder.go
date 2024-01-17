package seeder

import (
	"context"
	"fmt"
	"ttnmwastemanagementsystem/gen"
	"ttnmwastemanagementsystem/logger"
)

type TtnmOrganizationSeeder struct{}



func (ttnmOrganizationSeeder TtnmOrganizationSeeder) Run(q *gen.Queries) {
	logger.Log("[SEEDER/TtnmOrganizationSeeder SEEDER]", "======= Seeding TtnmOrganizationSeeder======", logger.LOG_LEVEL_INFO)

	organizations,err :=q.GetTTNMOrganizations(context.Background(),"TTNM_1");
	if err==nil{
		if len(organizations)==0{
			err:=q.InsertTTNMOrganization(context.Background(),gen.InsertTTNMOrganizationParams{
				Name: "TakaTaka Ni Mali",
				AboutUs: "Taka Taka ni Mali is an organisation dedicated to innovation, growth and strengthening of connections in the waste management ecosystem in Kenya, turning it into a circular economy.",
				LogoPath: "https://takanimali.org/takatakaicon.png",
				TagLine: "Waste is Wealth",
				WebsiteUrl: "https://takanimali.org/",
				City: "Nairobi",
				Country: "Kenya",
				Zip: "00100",
				State: "Nairobi",
				OrganizationID: "TTNM_1",
				AppAppstoreLink: "",
				AppGooglePlaystoreLink: "",
			})
			if(err != nil){
				fmt.Printf(err.Error())
			}
		}else{
			q.UpdateTtnmOrganizationProfile(context.Background(),gen.UpdateTtnmOrganizationProfileParams{
				Name: "TakaTaka Ni Mali",
				AboutUs: "Taka Taka ni Mali is an organisation dedicated to innovation, growth and strengthening of connections in the waste management ecosystem in Kenya, turning it into a circular economy.",
				LogoPath: "https://takanimali.org/takatakaicon.png",
				TagLine: "Waste is Wealth",
				WebsiteUrl: "https://takanimali.org/",
				City: "Nairobi",
				Country: "Kenya",
				Zip: "00100",
				State: "Nairobi",
				OrganizationID: "TTNM_1",
				AppAppstoreLink: "",
				AppGooglePlaystoreLink: "",
			})
		}
	}
}