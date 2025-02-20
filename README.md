Configuration:
set database connection parameters in database/datagase.go file

PREPARATION
Database is up to user to be set up from the list of migrations, provided in data/database_migration_log.txt
After the database is ready and the application successfully manages to connect to it, you can write a populate subcommand.
That command is just a script and does not require the main command to start up the service.
Onca everything is ready, you can start up the service and perform curl commands. (It is not a necessity to populate the dbs for this to work although)

NOTE
I have had limited time, writing this program, so some functionality is limited, particularly, signal management has been added to "TODO" box
Same goes for unit tests, if i had more time, I would have written them too.

INFO
data/data.json has been formatted, but has not been changed from the original file.
data/database_migration_log.txt has the log of migrations, done to the database (table creations)
data folder contains two resulting .csv files: signal_id.csv and track.csv

Usage:
1) DB population command. It will copy the data from a json file and add it to the database. It will skip (fail to add in the logs)
entries, if they are already present in the database.
C:\Users\maxga\go\go124\go\bin\go run cross_tech populate

2) This command will start a server, listening on a provided port.
C:\Users\maxga\go\go124\go\bin\go run cross_tech --server-port 8080

Testing:
1) GET
curl -Method Get http://localhost:8080/track?id=0
2) POST
curl -Method Post http://localhost:8080/track?entry='{"track_id": 0, "source": "custom_source", "target": "custom_target"}'
3) DELETE
curl -Method Delete http://localhost:8080/track?id=0
DELETE of a track will cause all the signals, referencing that particular track to be deleted too
4) PUT
curl -Method Put "http://localhost:8080/track?id=0&typ=source&value=new_custom_source"

// TODO
1) implement api to interact with signal_id
2) create a set of automated tests