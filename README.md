# YoLink Client API
## Setup

* Get credentials from YoLink mobile app:
    * Settings, Account, Advanced Settings, User Access Credentials
* Store the credentials in environment variables:
        
        export YOLINK_UAID="UAID"
        export YOLINK_SECRET_KEY="Secret Key"

* Build and run the `yl-request` app to test the API.

## Commands / Sample Apps
### yl-request

* Utility to query device list and device state.
* Example commands:

        yl-request list-devices
        yl-request get-state DEVICE_ID

### yl-msgdump

* Logs MQTT messages to stdout.

### yl-pgwriter

* Logs MQTT messages to PostgreSQL.
* See `cmd/yl-pgwriter/postgres_setup.sql` for DB setup script.
* Configure DB connection in an environment variable:

        export YOLINK_DB_URL="postgres://yolink:abcd@localhost:5432/yolink"

* Example query:

        select time,
               dev_id,
               (payload->'data'->'temperature')::float as t,
               (payload->'data'->'humidity')::float as h
        from mqtt_messages
        order by time;

