run:
	echo "running.."
	go run ./cmd/shortener -d "host=localhost port=5432 user=postgres password=1073849 dbname=shortner sslmode=disable"

build:
	echo "building.."
	go build -a -o ./cmd/shortener ./cmd/shortener

test:
	echo "testing.."
	go test -v -cover ./...

.DEFAULT_GOAL := run