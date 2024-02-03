-- name: CreateBuyer :one
insert into
    buyers (
        company_id, company, first_name, last_name, calling_code, phone, administrative_level_1_location, location, is_active, lat, lng, created_at, updated_at,region
    )
values (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,$14
    ) returning *;

-- name: CreateSupplier :one
insert into
    suppliers (
        company_id, company, first_name, last_name, calling_code, phone, administrative_level_1_location, location, is_active, lat, lng, created_at, updated_at,region
    )
values (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,$14
    ) returning *;

-- name: UpdateBuyer :exec
update buyers set company_id=$1,company=$2,first_name=$3,last_name=$4,calling_code=$5,phone=$6,administrative_level_1_location=$7,location=$8,is_active=$9,lat=$10,lng=$11 ,region=$12 where id=$13;


-- name: UpdateSupplier :exec
update suppliers set company_id=$1,company=$2,first_name=$3,last_name=$4,calling_code=$5,phone=$6,administrative_level_1_location=$7,location=$8,is_active=$9,lat=$10,lng=$11 ,region=$12 where id=$13;


-- name: DeleteBuyer :exec
delete from buyers where id = $1;

-- name: DeleteSupplier :exec
delete from suppliers where id = $1;

-- name: CreateSale :one
insert into sales(ref,company_id,buyer_id,total_weight,total_amount,dump) VALUES ($1,$2,$3,$4,$5,$6) returning *;

-- name: CreateBuys :one
insert into purchases(ref,company_id,supplier_id,total_weight,total_amount,dump) VALUES ($1,$2,$3,$4,$5,$6) returning *;


-- name: DeleteSale :exec
delete from sales where id=$1;

-- name: DeletePurchase :exec
delete from purchases where id=$1;


-- name: CreateSaleItem :one
insert into sale_items(company_id,sale_id,waste_type_id,weight,cost_per_kg,total_amount) VALUES($1,$2,$3,$4,$5,$6) returning *;


-- name: CreatePurchaseItem :one
insert into purchase_items(company_id,purchase_id,waste_type_id,weight,cost_per_kg,total_amount) VALUES($1,$2,$3,$4,$5,$6) returning *;


-- name: MakeCashPayment :one
insert into sale_transactions(ref,sale_id,company_id,payment_method,amount,transaction_date) VALUES($1,$2,$3,$4,$5,$6) returning *;


-- name: MakePurchaseCashPayment :one
insert into purchase_transactions(ref,purchase_id,company_id,payment_method,amount,transaction_date) VALUES($1,$2,$3,$4,$5,$6) returning *;


-- name: InventoryItemCount :one
select count(*) from inventory where waste_type_id=$1 and company_id = $2;

-- name: GetInventoryItem :one
select * from inventory where  waste_type_id=$1 and company_id = $2;

-- name: UpdateInventoryItem :exec
update inventory set total_weight=$1 where id =$2;

-- name: InsertToInventory :exec
insert into inventory(waste_type_id,company_id,total_weight) VALUES($1,$2,$3);