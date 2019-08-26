# Run

```
cd database
docker build -t dda-postgres .
docker run -ti -p 5432:5432 dda-postgres
```

# Requests

## Get all games
Request:
```
curl localhost:9090/games
```
Response:
```
{}
```

## Create a new game
Request:
```
curl -XPOST localhost:9090/game/new
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Rick Astley - Never gonna give you up",
    "players": {
        "abc": {
            "id": "abc",
            "score": 0,
            "ready": false
        }
    },
    "started": 0,
    "stopped": 0
}
```

## Get the details of an existing game
Request:
```
curl localhost:9090/game/123
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Rick Astley - Never gonna give you up",
    "players": {
        "abc": {
            "id": "abc",
            "score": 0,
            "ready": false
        }
    },
    "started": 0,
    "stopped": 0
}
```

## Join an existing game
Request:
```
curl -XPOST localhost:9090/game/123/join
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Rick Astley - Never gonna give you up",
    "players": {
        "abc": {
            "id": "abc",
            "score": 0,
            "ready": false
        },
        "def": {
            "id": "def",
            "score": 0,
            "ready": false
        }
    },
    "started": 0,
    "stopped": 0
}
```

## Set the player status to ready
Request:
```
curl -XPOST localhost:9090/game/123/players/abc/ready
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Rick Astley - Never gonna give you up",
    "players": {
        "abc": {
            "id": "abc",
            "score": 0,
            "ready": true
        },
        "def": {
            "id": "def",
            "score": 0,
            "ready": false
        }
    },
    "started": 0,
    "stopped": 0
}
```

## Start an existing game
Request:
```
curl -XPOST localhost:9090/game/123/start
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Rick Astley - Never gonna give you up",
    "players": {
        "abc": {
            "id": "abc",
            "score": 0,
            "ready": true
        },
        "def": {
            "id": "def",
            "score": 0,
            "ready": false
        }
    },
    "started": 1565613892,
    "stopped": 0
}
```

## Get a list of unassigned allocations, maximum 10
Request:
```
curl localhost:9090/allocations?player=nickyjams
```
Response:
```
["77f06c2a-d2f2-9a25-0ad6-4d7e3d9fa4e7"]
```

## Stop an allocation
Request:
```
curl localhost:9090/allocations/f1ba4520-9180-056b-2e9b-ed62c593e434/stop
```
Response:
```
"f1ba4520-9180-056b-2e9b-ed62c593e434"
```