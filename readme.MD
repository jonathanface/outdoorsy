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

5. Hit http://localhost/api/rentals or http://localhost/api/rentals/[RENTAL_ID] 

Supported parameters:
```
	price_min  int
	price_max int
	limit     int
	offset    int
	ids       comma,separated,ints
	near      comma,separated,float,pair
	Sort      string
  ```
Sorting options:
```
	price_asc
	price_desc
	year_asc
	year_desc
	make_asc
	make_desc
	type_asc
	type_desc
	created_asc
	created_desc
	updated_asc
	updated_desc
```

5. To test just run go test ./...
