-- collection_requests.sql

-- name: InsertNewCollectionRequest :exec
insert into collection_requests( producer_id,collector_id,request_date,location,lat,lng,administrative_level_1_location,first_contact_person ) values ($1, $2, $3, $4, $5, $6, $7, $8) returning *;

-- name: UpdateCollectionRequest :exec
update collection_requests
set
    pickup_date = $1,
    status = $2
where id = $3;

-- name: ConfirmCollectionRequest :exec
update collection_requests set confirmed = $1 where id = $2;

-- name: CancelCollectionRequest :exec
update collection_requests set cancelled = $1 where id = $2;


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
WHERE collection_requests.status=$1;

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
WHERE collection_requests.cancelled=$1;

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
WHERE collection_requests.confirmed=$1;

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
WHERE collection_requests.confirmed=$1;

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
select sum(waste_items.weight) as total_weight,waste_types.name from waste_items 
inner join waste_types on waste_types.id=waste_items.waste_type_id 
inner join collection_requests on collection_requests.id=waste_items.collection_request_id
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
    waste_items AS totals ON totals.collection_request_id = collection_requests.id
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
    waste_items AS totals ON totals.collection_request_id = collection_requests.id
WHERE
    collection_requests.producer_id = $1
GROUP BY
    collection_requests.id;


-- name: GetWasteItemsProducerData :many
SELECT
    CAST(SUM(waste_items.weight) AS DECIMAL(10,2)) AS total_weight,
    waste.name AS waste_name,
    collections.status AS collection_status
FROM
    waste_items
JOIN
    waste_types AS waste ON waste_items.waste_type_id = waste.id
LEFT JOIN
    collection_requests AS collections ON collections.id = waste_items.collection_request_id
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
    waste_items AS totals ON totals.collection_request_id = collection_requests.id
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
    waste_items AS totals ON totals.collection_request_id = collection_requests.id
WHERE
    collection_requests.producer_id = $1 AND collection_requests.status = false
GROUP BY
    collection_requests.id, collector.name;


-- name: GetAggregatorNewRequests :many
SELECT
    collection_requests.*,
    producer.name AS producer_name,
    producer.location AS producer_location
FROM
    collection_requests
LEFT JOIN
    companies AS producer ON producer.id = collection_requests.producer_id
WHERE
    collection_requests.collector_id = $1 AND collection_requests.status = false
GROUP BY
    collection_requests.id, producer.name, producer.location;