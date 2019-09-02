.PHONY: protos database api

build-database:
	pushd database && \
	docker build -t eveld/da-dance-database . && \
	docker push eveld/da-dance-database

run-database:
	docker run -dit --net=host -e POSTGRES_USER=secret_user -e POSTGRES_PASSWORD=secret_password -e POSTGRES_DB=dda eveld/da-dance-database

build-api:
	docker build -t eveld/da-dance-api . && \
	docker tag eveld/da-dance-api:latest eveld/da-dance-api:${TAG}
	docker push eveld/da-dance-api:${TAG}

run-api:
	docker run -it --net=host -e NOMAD_ADDR=http://34.76.12.196:4646 -e POSTGRES_USER=secret_user -e POSTGRES_PASSWORD=secret_password -e POSTGRES_DATABASE=dda eveld/da-dance-api

run-go:
	NOMAD_ADDR=http://34.76.12.196:4646 POSTGRES_USER=secret_user POSTGRES_PASSWORD=secret_password POSTGRES_DATABASE=dda go run .

deploy-api:
	NOMAD_ADDR=http://34.76.12.196:4646 nomad run deploy/api.hcl