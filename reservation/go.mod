module github.com/SergeyBogomolovv/restaurant/reservation

go 1.23.2

replace github.com/SergeyBogomolovv/restaurant/common => ../common

require (
	github.com/SergeyBogomolovv/restaurant/common v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.4.0
	google.golang.org/grpc v1.68.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
