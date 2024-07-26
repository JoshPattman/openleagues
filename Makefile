build-and-deploy-local: build
	docker compose up

build-and-deploy-remote: build deploy

# This will build the docker image and save it to the bin directory
build:
	rm -rdf bin && mkdir bin
	GOOS=linux GOARCH=amd64 go build -o bin/olsrv_bin
	docker build -t olsrv_img .
	docker save olsrv_img > bin/olsrv_img.tar
	gzip bin/olsrv_img.tar

# This will upload the docker image to your server and restart the server container
# Make sure to set the environment variable SERVER_ADDR
deploy:
	ssh ${SERVER_ADDR} "mkdir -p /root/olsrv"
	ssh ${SERVER_ADDR} "cd /root/olsrv; docker compose down || true"
	scp ./bin/olsrv_img.tar.gz ${SERVER_ADDR}:/root/olsrv
	scp ./compose.yaml ${SERVER_ADDR}:/root/olsrv
	ssh ${SERVER_ADDR} "cd /root/olsrv; docker load < olsrv_img.tar.gz"
	ssh ${SERVER_ADDR} "cd /root/olsrv; docker compose up -d"