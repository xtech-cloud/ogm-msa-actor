APP_NAME := ogm-actor
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )
UUID := $(shell cat /tmp/ogm-actor-uuid)

.PHONY: build
build: 
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: run-f
run-fs:
	MSA_CONFIG_DEFINE='{"source":"file","prefix":"/etc/ogm/","key":"actor.yml"}' ./bin/${APP_NAME}

.PHONY: run-e
run-cs:
	MSA_CONFIG_DEFINE='{"source":"etcd","prefix":"/xtc/ogm/config","key":"actor.yml"}' ./bin/${APP_NAME}

.PHONY: call
call:
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Healthy.Echo '{"msg":"hello"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Domain.Create '{"name":"test"}'
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Domain.List '{"offset":0, "count":9999}'
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Device.List '{"offset":0, "count":9999}'
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Domain.Execute '{"uuid":"${UUID}", "command":"reboot", "device":["0001", "0002"], "parameter":""}'
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Sync.Push '{"domain":"${UUID}", "device":{"serialNumber":"00002", "name":"G2", "operatingSystem":"Windows 11", "systemVersion":"21H2", "shape":"PC", "program":{"myapp":"1.0.0"}}, "upProperty":{"key1":"value1", "key2":"value2"}, "downProperty":["key1"]}'
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Sync.Pull '{"domain":"${UUID}", "downProperty":["key1"]}'
	gomu --registry=etcd --client=grpc call xtc.ogm.actor Domain.FetchDevice '{"uuid":"${UUID}"}'
	#gomu --registry=etcd --client=grpc call xtc.ogm.actor Domain.Delete '{"uuid":"${UUID}"}'

.PHONY: post
post:
	curl -X POST -H 'Content-Type:application/json' -d '{"msg":"hello"}' localhost/ogm/actor/Healthy/Echo

.PHONY: benchmark
benchmark:
	python3 ./benchmark.py

.PHONY: dist
dist:
	rm -rf ./dist
	mkdir ./dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}

.PHONY: docker
docker:
	docker build -t xtechcloud/${APP_NAME}:${BUILD_VERSION} .
	docker rm -f ${APP_NAME}
	docker run --restart=always --name=${APP_NAME} --net=host -v /data/${APP_NAME}:/ogm -e MSA_REGISTRY_ADDRESS='localhost:2379' -e MSA_CONFIG_DEFINE='{"source":"file","prefix":"/ogm/config","key":"${APP_NAME}.yaml"}' -d xtechcloud/${APP_NAME}:${BUILD_VERSION}
	docker logs -f ${APP_NAME}
