# Shange

### Installation:

 - Database

docker run --name shange-postgres -e POSTGRES_USER=postuser -e POSTGRES_PASSWORD=postpass -e POSTGRES_DB=shange-db -p 5432:5432 -d postgres:12.1

 - Test database
 
docker run --name shange-test-postgres -e POSTGRES_USER=postuser -e POSTGRES_PASSWORD=postpass -e POSTGRES_DB=shange-db -p 5432:5432 -d postgres:12.1


### Deployment

go run main.go

 | Flags                    | Default          |
 |:-------------------------|------------------|
 | mdir                     | ./db/migrations  |
 | db.type                  | postgres         |
 | db.user                  | postuser         |
 | db.pwd                   | postuser         |
 | db.host                  | localhost        |
 | db.port                  | 5432             |
 | db.name                  | shange-db        |