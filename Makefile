-include ../../x-collection/tools/release-scripts/src/MakefileQA


run:
	cd cmd; go run *.go

test-api:
	cd cmd/api; go test -v *.go

test: test-api

login-user:
	curl --location --request POST 'http://localhost:10000/api/auth' \
	--header 'Content-Type: application/json' \
	--data-raw '{"username":"user","password":"user"}'

login-admin:
	curl --location --request POST 'http://localhost:10000/api/auth' \
	--header 'Content-Type: application/json' \
	--data-raw '{"username":"admin","password":"admin"}'

stop-local-db:
	docker rm -f mysql

start-local-db:
	docker run --name mysql  -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=hancock -p 3306:3306 -d mysql:latest

restart-local-db:stop-local-db start-local-db