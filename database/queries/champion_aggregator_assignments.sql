-- champion_aggregator_assignments.sql

-- name: GetAllChampionCollectorsAssignments :many
SELECT 
    champion_aggregator_assignments.*,
    champion.name AS aggregator_name,
    collector.name AS champion_name
FROM 
    champion_aggregator_assignments
LEFT JOIN 
    companies AS champion ON champion.id = champion_aggregator_assignments.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = champion_aggregator_assignments.collector_id;

-- name: AssignChampionToCollector :one
insert into champion_aggregator_assignments( champion_id,collector_id ) values ($1, $2) returning *;

-- name: GetAllChampionsForACollector :many
SELECT 
    champion_aggregator_assignments.*,
    champion.name AS champion_name,
    collector.name AS collector_name
FROM 
    champion_aggregator_assignments
LEFT JOIN 
    companies AS champion ON champion.id = champion_aggregator_assignments.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = champion_aggregator_assignments.collector_id
WHERE collector_id = $1;

-- name: GetCollectorsForGreenChampion :many
SELECT 
    champion_aggregator_assignments.*,
    champion.name AS champion_name,
    collector.name AS collector_name
FROM 
    champion_aggregator_assignments
LEFT JOIN 
    companies AS champion ON champion.id = champion_aggregator_assignments.champion_id
LEFT JOIN 
    companies AS collector ON collector.id = champion_aggregator_assignments.collector_id
WHERE champion_id = $1;

-- name: UpdateChampionCollector :exec
update champion_aggregator_assignments
set
    collector_id = $1
where id = $2;

-- name: DeleteChampionCollector :exec
delete from champion_aggregator_assignments where id=$1;

-- name: GetAssignedCollectorsToGreenChampion :many
select * from champion_aggregator_assignments where champion_id=$1;

-- name: RemoveAggrigatorsAssignedFromGreenChampions :exec
delete from champion_aggregator_assignments where champion_id =$1;

-- name: AssignCollectorsToGreenChampion :exec
insert into champion_aggregator_assignments(champion_id,collector_id) VALUES($1,$2);