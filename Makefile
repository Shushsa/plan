all: lint install

install: go.sum
		GO111MODULE=on go install ./cmd/pland
		GO111MODULE=on go install ./cmd/plancli

		mkdir -p ~/.pland/config

		cp -r ./installation/genesis.json ~/.pland/config/
		cp -r ./installation/config.toml ~/.pland/config/

		plancli config chain-id plan
		plancli config output json
		plancli config indent true
		plancli config node tcp://localhost:26657
		plancli config trust-node true

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify