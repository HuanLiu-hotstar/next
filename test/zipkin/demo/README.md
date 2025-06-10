# Gateway Demo

## gateway

- addr: localhost:8080/playback

## pc

- addr: localhost:8083/pc

## um

- addr : localhost:8084/um

## ratelimit

- add: localhost:8082/ratelimit


protoc --go_out=. --go_opt=paths=test/zipkin/demo/proto --go-grpc_out=. --go-grpc_opt=paths=test/zipkin/demo/proto  ratelimit/ratelimit.proto