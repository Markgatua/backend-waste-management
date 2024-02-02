-- name: CreateBuyer :one
insert into
    buyers (
        company_id, company, first_name, last_name, calling_code, phone, administrative_level_1_location, location, is_active, lat, lng, created_at, updated_at,region
    )
values (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,$14
    ) returning *;

-- name: UpdateBuyer :exec
update buyers set company_id=$1,company=$2,first_name=$3,last_name=$4,calling_code=$5,phone=$6,administrative_level_1_location=$7,location=$8,is_active=$9,lat=$10,lng=$11 ,region=$12 where id=$13;

-- name: DeleteBuyer :exec
delete from buyers where id = $1;


-- name: CreateSale :one
insert into sales(ref,company_id,buyer_id,total_amount_of_waste,total_amount,dump) VALUES ($1,$2,$3,$4,$5,$6) returning *;

-- name: CreateSaleItem :one
insert into sale_items(company_id,sale_id,waste_type_id,amount_of_waste,cost_per_kg,total_amount) VALUES($1,$2,$3,$4,$5,$6) returning *;

-- name: MakeCashPayment :one
insert into sale_transactions(ref,sale_id,company_id,payment_method,transaction_date) VALUES($1,$2,$3,"CASH",$4) returning *;
