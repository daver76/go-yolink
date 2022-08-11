

# Example query:

    select time,
           dev_id,
           (payload->'data'->'temperature')::float as t,
           (payload->'data'->'humidity')::float as h
    from mqtt_messages
    order by time;

