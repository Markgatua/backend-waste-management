SELECT 
    cr.id AS request_id,
    cr.producer_id,
    cr.collector_id,
    cr.request_date,
    cr.pickup_time_stamp_id,
    cr.location,
    cr.administrative_level_1_location,
    cr.lat,
    cr.lng,
    cr.pickup_date,
    cr.status,
    cr.first_contact_person,
    cr.second_contact_person,
    cr.created_at AS request_created_at,
    cri.id AS waste_item_id,
    cri.waste_type_id,
    cri.collector_id AS waste_collector_id,
    cri.weight,
    cri.created_at AS waste_item_created_at
FROM 
    collection_requests cr
LEFT JOIN 
    collection_request_waste_items cri ON cr.id = cri.collection_request_id
LEFT JOIN 
    companies AS champion ON champion.id = cr.producer_id
LEFT JOIN 
    companies AS collector ON collector.id = cr.collector_id
WHERE 
    cr.producer_id = $1
    AND cr.status IN (5);