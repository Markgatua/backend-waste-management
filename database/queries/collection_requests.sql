-- collection_requests.sql

-- name: InsertNewCollectionRequest :exec
insert into collection_requests( producer_id,collector_id,request_date,pickup_time_stamp_id,location,lat,lng,administrative_level_1_location,first_contact_person,status ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning *;

-- name: UpdateCollectionRequest :exec
update collection_requests
set
    pickup_date = $1,
    status = $2
where id = $3;

-- name: ConfirmCollectionRequest :exec
update collection_requests set status = 2 where id = $1;

-- name: CancelCollectionRequest :exec
update collection_requests set status = 4 where id = $1;

-- name: CompleteCollectionRequest :exec
update collection_requests set status = 5 where id = $1;


-- name: GetAllCollectionRequests :many
SELECT 
    collection_requests.*,
    champion.name AS aggregator_name,
    collector.name AS champion_name,
    secondcollector.name as second_collector_name
FROM 
    collection_requests
LEFT JOIN 
    companies AS champion ON champion.id = collection_requests.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN 
    companies AS secondcollector ON secondcollector.id = collection_requests.second_collector_id;


-- name: GetAllCompletedCollectionRequests :many
SELECT 
    collection_requests.*,
    champion.name AS aggregator_name,
    collector.name AS champion_name,
    secondcollector.name as second_collector_name
FROM 
    collection_requests
LEFT JOIN 
    companies AS champion ON champion.id = collection_requests.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN 
    companies AS secondcollector ON secondcollector.id = collection_requests.second_collector_id
WHERE collection_requests.status=5;

-- name: GetAllCancelledCollectionRequests :many
SELECT 
    collection_requests.*,
    champion.name AS aggregator_name,
    collector.name AS champion_name,
    secondcollector.name as second_collector_name
FROM 
    collection_requests
LEFT JOIN 
    companies AS champion ON champion.id = collection_requests.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN 
    companies AS secondcollector ON secondcollector.id = collection_requests.second_collector_id
WHERE collection_requests.cancelled=4;

-- name: GetCollectionRequest :one
SELECT 
    collection_requests.*  
FROM 
    collection_requests
WHERE collection_requests.id=$1;

-- name: ChangeCollectionRequestStatus :exec
update collection_requests set status = $1 where id=$2;

-- name: GetAllPendingConfirmationCollectionRequests :many
SELECT 
    collection_requests.*,
    champion.name AS aggregator_name,
    collector.name AS champion_name,
    secondcollector.name as second_collector_name
FROM 
    collection_requests
LEFT JOIN 
    companies AS champion ON champion.id = collection_requests.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN 
    companies AS secondcollector ON secondcollector.id = collection_requests.second_collector_id
WHERE collection_requests.status=2;

-- name: GetAllPendingCollectionRequests :many
SELECT 
    collection_requests.*,
    champion.name AS aggregator_name,
    collector.name AS champion_name,
    secondcollector.name as second_collector_name
FROM 
    collection_requests
LEFT JOIN 
    companies AS champion ON champion.id = collection_requests.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN 
    companies AS secondcollector ON secondcollector.id = collection_requests.second_collector_id
WHERE collection_requests.confirmed=2;

-- name: GetAllCollectionRequestsForACollector :many
SELECT 
    collection_requests.*,
    champion.name AS aggregator_name,
    collector.name AS champion_name,
    secondcollector.name as second_collector_name
FROM 
    collection_requests
LEFT JOIN 
    companies AS champion ON champion.id = collection_requests.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN 
    companies AS secondcollector ON secondcollector.id = collection_requests.second_collector_id
WHERE collection_requests.collector_id=$1;


-- name: CollectionWeightTotals :one
select sum(collection_request_waste_items.weight) as total_weight,waste_types.name from collection_request_waste_items 
inner join waste_types on waste_types.id=collection_request_waste_items.waste_type_id 
inner join collection_requests on collection_requests.id=collection_request_waste_items.collection_request_id
where collection_requests.producer_id=$1 GROUP BY waste_types.name;

-- name: GetLatestCollection :one
SELECT
    collection_requests.*,
    collector.name AS collector_name,
    CAST(SUM(totals.weight) AS DECIMAL(10,2)) AS total_weight
FROM
    collection_requests
LEFT JOIN
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN
    collection_request_waste_items AS totals ON totals.collection_request_id = collection_requests.id
WHERE
    collection_requests.id = $1
GROUP BY
    collection_requests.id, collector.name;


-- name: GetProducerLatestCollectionId :one
SELECT *
FROM collection_requests
WHERE producer_id = $1
ORDER BY created_at DESC
LIMIT 1;


-- name: GetCollectionStats :many
SELECT
    collection_requests.*,
    CAST(SUM(totals.weight) AS DECIMAL(10,2)) AS total_weight
FROM
    collection_requests
LEFT JOIN
    collection_request_waste_items AS totals ON totals.collection_request_id = collection_requests.id
WHERE
    collection_requests.producer_id = $1
GROUP BY
    collection_requests.id;


-- name: GetWasteItemsProducerData :many
SELECT
    CAST(SUM(collection_request_waste_items.weight) AS DECIMAL(10,2)) AS total_weight,
    waste.name AS waste_name,
    collections.status AS collection_status
FROM
    collection_request_waste_items
JOIN
    waste_types AS waste ON collection_request_waste_items.waste_type_id = waste.id
LEFT JOIN
    collection_requests AS collections ON collections.id = collection_request_waste_items.collection_request_id
WHERE
    collections.producer_id = $1
GROUP BY
    collections.status, waste.name;
    
-- name: GetAllProducerCompletedCollectionRequests :many
SELECT
    collection_requests.*,
    collector.name AS collector_name,
    CAST(SUM(totals.weight) AS DECIMAL(10,2)) AS total_weight
FROM
    collection_requests
LEFT JOIN
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN
    collection_request_waste_items AS totals ON totals.collection_request_id = collection_requests.id
WHERE
    collection_requests.producer_id = $1 AND collection_requests.status = true
GROUP BY
    collection_requests.id, collector.name;

-- name: GetAllProducerPendingCollectionRequests :many
SELECT
    collection_requests.*,
    collector.name AS collector_name,
    CAST(SUM(totals.weight) AS DECIMAL(10,2)) AS total_weight
FROM
    collection_requests
LEFT JOIN
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN
    collection_request_waste_items AS totals ON totals.collection_request_id = collection_requests.id
WHERE
    collection_requests.producer_id = $1 AND collection_requests.status = false
GROUP BY
    collection_requests.id, collector.name;


-- name: GetAggregatorNewRequests :many
SELECT
    collection_requests.*,
    producer.name AS producer_name,
    producer.location AS producer_location,
    pickup.stamp AS pickup_stamp,
    pickup.time_range AS pickup_time_range
FROM
    collection_requests
LEFT JOIN
    companies AS producer ON producer.id = collection_requests.producer_id
LEFT JOIN
    pickup_time_stamps AS pickup ON pickup.id = collection_requests.pickup_time_stamp_id
WHERE
    collection_requests.collector_id = $1 AND collection_requests.status = 1
GROUP BY
    collection_requests.id, producer.name, producer.location,pickup.stamp,pickup.time_range;


-- name: GetMyLatestRequests :many
SELECT
    collection_requests.*,
    collector.name AS collector_name,
    collector.location AS collector_location,
    pickup.stamp AS pickup_stamp,
    pickup.time_range AS pickup_time_range
FROM
    collection_requests
LEFT JOIN
    companies AS collector ON collector.id = collection_requests.collector_id
LEFT JOIN
    pickup_time_stamps AS pickup ON pickup.id = collection_requests.pickup_time_stamp_id
WHERE
    collection_requests.producer_id = $1 AND collection_requests.status = 1 AND CAST(collection_requests.request_date AS TIMESTAMP) >= CURRENT_DATE 
GROUP BY

    collection_requests.id, collector.name, collector.location,pickup.stamp,pickup.time_range;


-- name: GetCollectionRequestsInArray :many
select collection_requests.id,collection_requests.producer_id, companies.name as champion_name,collection_requests.collector_id,
collection_requests.request_date,collection_requests.pickup_date,collection_requests.status,collection_requests.lat,
collection_requests.lng,collection_requests.created_at,collection_requests.pickup_time_stamp_id,collection_requests.id,
collection_requests.first_contact_person,collection_requests.second_contact_person,pickup_time_stamps.stamp,
pickup_time_stamps.time_range from collection_requests 
inner join pickup_time_stamps on pickup_time_stamps.id=collection_requests.pickup_time_stamp_id
inner join companies on companies.id=collection_requests.producer_id
where collection_requests.id=ANY(sqlc.arg('collectionIds')::int[]);

-- name: GetCollectionScheduleInArray :many
select champion_pickup_times.id as pickup_time_id,champion_pickup_times.champion_aggregator_assignment_id,champion_pickup_times.pickup_time_stamp_id,
champion_pickup_times.pickup_day,champion_aggregator_assignments.champion_id,champion_aggregator_assignments.collector_id,companies.lat,
companies.lng,companies.name as champion_name,companies.location,companies.contact_person1_first_name,companies.contact_person1_last_name,
companies.contact_person1_phone,companies.contact_person1_email,companies.contact_person2_email,companies.administrative_level_1_location,
companies.contact_person2_first_name,companies.contact_person2_last_name,companies.contact_person2_phone,pickup_time_stamps.stamp,
pickup_time_stamps.time_range from champion_pickup_times 
left join pickup_time_stamps on pickup_time_stamps.id=champion_pickup_times.pickup_time_stamp_id
left join champion_aggregator_assignments on champion_pickup_times.champion_aggregator_assignment_id=champion_aggregator_assignments.id
left join companies on companies.id = champion_aggregator_assignments.champion_id
where champion_pickup_times.id=ANY(sqlc.arg('pickupTimeIds')::int[]);

