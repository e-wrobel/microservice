# Demo microservice based on Golang
* Go 1.18
* Configuration: viper
* Database: sqlite
* Rest: net/http

## Service is intended to upload large json-like file to the Rest service.
To avoid high memory usage, json file is parsed and send in chunks.
File contains regular structures showed below:
```json
"AEAJM": {
    "name": "Ajman",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
```

First parameter is called as identifier, and the rest of the object is represented by json object.
Its represented by structure below:

```go
type PortDetails struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

type PortEntity struct {
	Identifier  string       `json:"identifier"`
	PortDetails *PortDetails `json:"details"`
}
```

Client is able to iterate through the input file, parse it and send to API.
## Remark
* Client and Server package reading configuration file: **config.yaml**. 
* File resides in **local_config** directory.

## Client
* To build: **go build pkg/cmd/main.go**.
* To start: **./main**.

## Server
* Server is waiting for connection at **http:/0.0.0.0:8080/** endpoint.
* This endpoint is responsible for creating and updating json files represented by
  **PortEntity** objects.
* To start server use: **docker-compose up**.

## Testing solution
Once server is started, we can execute Client and upload file.
From the container perspective (server) we can check if database is
populated, by inspecting:
* sqlite3 /var/tmp/database.sql 
* select * from files;

Then we will find records similar to:

```sql
ZWBUQ|{"name":"Bulawayo","city":"Bulawayo","country":"Zimbabwe","alias":[],"regions":[],"coordinates":[28.626479,-20.1325066],"province":"Bulawayo","timezone":"Africa/Harare","unlocs":["ZWBUQ"],"code":""}
ZWHRE|{"name":"Harare","city":"Harare","country":"Zimbabwe","alias":[],"regions":[],"coordinates":[31.03351,-17.8251657],"province":"Harare","timezone":"Africa/Harare","unlocs":["ZWHRE"],"code":""}
ZWUTA|{"name":"Mutare","city":"Mutare","country":"Zimbabwe","alias":[],"regions":[],"coordinates":[32.650351,-18.9757714],"province":"Manicaland","timezone":"Africa/Harare","unlocs":["ZWUTA"],"code":""}
```

Where **identifier** is used as **PRIMARY KEY** in database.