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
curl localhost:8080/games
```
Response:
```
{}
```

## Create a new game
Request:
```
curl -XPOST localhost:8080/game/new
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Techno shit",
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
curl localhost:8080/game/123
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Techno shit",
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
curl -XPOST localhost:8080/game/123/join
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Techno shit",
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
curl -XPOST localhost:8080/game/123/players/abc/ready
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Techno shit",
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
curl -XPOST localhost:8080/game/123/start
```
Response:
```
{
    "id": "123",
    "host": "abc",
    "song": "Techno shit",
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