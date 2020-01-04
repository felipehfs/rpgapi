# RPGAPI

My project using the TDD methodology to build an API.

## Routes 

|method | endpoint       | description           |
|-------|----------------|-----------------------|
|GET    | /api/characters| Display all characters|
|POST   | /api/characters| Create new character  |

## Requirements

- Go 1.13
- PostgresSQL

## Instalation

1. Create manually the databases for development and test.
2. Setup environment variables
```
    RPGAPI_TEST_USER=""
	RPGAPI_TEST_PASSWORD=""
    RPGAPI_TEST_DATABASE_NAME=""

	RPGAPI_DEVELOPMENT_USER=""
	RPGAPI_DEVELOPMENT_PASSWORD=""
    RPGAPI_DEVELOPMENT_DATABASE_NAME=""
```
3. You can run _go run main.go_ and listen on port 8083
## Runing Test 

```bash
    go run -v ./...
```