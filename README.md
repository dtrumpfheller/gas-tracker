# Gas Tracker
Gas prices only know one direction and as a consequence, I wanted to know when is the best time to refuel my iron horse.

Sadly, I couldn't find a reliable and free service available that provides current and historical prices for Canada.
Other countries are far ahead of us, but I guess we are too busy shoveling snow and have no time for such 'futuristic'
solutions. One okish source I could find is GasBuddy where members can report the prices.
Depending on how often that happens prices can still be badly outdated and thus this can't serve as a reliable source
for current prices. But it is good enough to get a somewhat accurate picture of prices over time and it gives me
some info to plan my refueling trip based on historical data.

This tracker regularly gets the prices for specific gas stations and uploads them into an InfluxDB2 database. It ensures that
already entered prices are not uploaded again until updated on the GasBuddy page.

Please ensure you have all necessary permissions and you aren't violating any terms and conditions of neither this project nor 
GasBuddy before using this software.

## Configuration
The configuration file must be a valid YAML file. Its path can be passed into the application as an argument, else **config.yml** is assumed.

Example **config.yml** file:
```
  url: http://127.0.0.1:9086
  token: "abcd"
  organization: "home"
  bucket: "gasbuddy"
  sleepDuration: 60
  stations:
    - 12345
    - 67890
```

| Name            | Description                              |
|-----------------|------------------------------------------|
| url             | address of InfluxDB2 server              |
| token           | auth token to access InfluxDB2 server    |
| organization    | organization of InfluxDB2 server         |
| bucket          | name of bucket                           |
| sleepDuration   | sleep time between gas checks in minutes |
| stations        | station ids from GasBuddy                |

## Docker
The agent was written with the intent of running it in docker. You can also run it directly if this is preferred.

### Build Image
Execute following statement, then either start via docker or docker compose.
```
docker build -t gas-tracker .
```

### Docker
```
docker run -d --restart unless-stopped --name=gas-tracker -v ./config.yml:/config.yml gas-tracker
```

### Docker Compose
```
version: "3.4"
services:
  gas-tracker:
    image: gas-tracker
    container_name: gas-tracker
    restart: unless-stopped
    volumes:
      - ./config.yml:/config.yml:ro
```