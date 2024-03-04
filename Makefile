# This will build the docker image and save it to the bin directory
build-docker:
	rm -rdf bin
	mkdir bin
	GOOS=linux GOARCH=amd64 go build -o bin/openleagues
	docker build -t openleagues .
	docker save openleagues > openleagues.tar
	gzip openleagues.tar
	mv openleagues.tar.gz bin/openleagues.tar.gz

# This will upload the docker image to your server and restart the server container
# Make sure to set the environment variable SERVER_ADDR
upload-to-server:
	scp ./bin/openleagues.tar.gz ${SERVER_ADDR}:/root
	ssh ${SERVER_ADDR} "docker stop openleagues; docker rm openleagues; docker rmi openleagues; docker load < /root/openleagues.tar.gz; docker run -d --name openleagues -p 8081:8080 openleagues"
