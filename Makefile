-include ../../x-collection/tools/release-scripts/src/MakefileQA


stop-local-db:
	docker rm -f mysql

start-local-db:
	docker run --name mysql  -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=hancock -p 3306:3306 -d mysql:latest

restart-local-db:stop-local-db start-local-db