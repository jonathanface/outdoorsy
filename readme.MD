Rest API demo

1. git clone https://github.com/jonathanface/outdoorsy

2. Create a .env file in the root directory with the following parameters:
```
PGSQL_HOST=localhost
PGSQL_USER=root
PGSQL_PASS=root
PGSQL_PORT=5434
PGSQL_DB=testingwithrentals
VERSION=PRODUCTION_V1
```
3. From a command prompt, in your repo's root directory, docker-compose up

4. In another command prompt, in your repo's root directory, go build && sudo ./outdoorsy

5. To test just run go test ./...
