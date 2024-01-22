-- collection_requests.sql

-- name: InsertNewCollectionRequest :exec
insert into collection_requests( producer_id,collector_id,request_date ) values ($1, $2, $3) returning *;

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