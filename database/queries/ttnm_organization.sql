-- ttnm_organization.sql

-- name: UpdateTtnmOrganizationProfile :exec
update ttnm_organization set name=$2, tag_line=$3, about_us=$4, logo_path=$5, website_url=$6, city=$7,state=$8,zip=$9,country=$10,app_appstore_link=$11,app_google_playstore_link=$12 where organization_id=$1;

-- name: GetTTNMOrganizations :many
select * from ttnm_organization where organization_id=$1;

-- name: InsertTTNMOrganization :exec
insert into ttnm_organization (name,organization_id,tag_line,about_us,logo_path,website_url,city,state,zip,country,app_appstore_link,app_google_playstore_link) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);
