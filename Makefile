CMD="main.go"
PATH_MAIN="./${CMD}"
CONFIG_FILE="./config/config.ENV.json"
CONFIG_TEST_FILE="./config/config.TEST.json"
PATH_TO_DIR=$(pwd)

DOCKER_USER="falconr"
APP_NAME="yutube"
DOCKER_IMAGE="${DOCKER_USER}/${APP_NAME}"
CONTAINER_PORT="6543"
FORWARD_PORT="6543"
IMAGE_TAG="latest"

all: test build 


.PHONY: dep 

dep:
	if ! [ -f go.mod ]; then go mod init; fi
	go get -u ./... 

run: dep build 

build: 
	if ! [ -d ${LOG_DIR} ]; then mkdir ${LOG_DIR}; fi
	CONFIG_PATH=${CONFIG_FILE} go run ${PATH_MAIN} 

test:
	CONFIG_PATH=${HOME}/go/src/github.com/Remydeme/yutube/${CONFIG_TEST_FILE} go test ./...

.PHONY: docker 

docker: clean 
	docker build -t "${DOCKER_IMAGE}":${IMAGE_TAG} .

docker-run:
	docker run -d  -P --name "${APP_NAME}"  -p "${FORWARD_PORT}:${CONTAINER_PORT}" "${DOCKER_IMAGE}"; 
	docker container logs ${APP_NAME}

docker-stop:
	docker stop ${APP_NAME}

.PHONY: docker-compose

docker-compose:
	docker-compose up -d
	docker container logs ${APP_NAME}

docker-clean:
	docker image rm $(docker image ls -a -q)

clean:
	rm -rf Gopkg.*
	rm -rf vendor 
	rm -rf go.*
	
