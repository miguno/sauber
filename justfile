timestamp := `date +%s`

semver := "1.0.0-alpha"
commit := `git show -s --format=%h`
version := semver + "+" + commit

# print available targets
default:
    just --list

# format source code
format:
    @echo "Formatting source code ..."
    gofmt -l -s -w .

# detect outdated modules (requires https://github.com/psampaz/go-mod-outdated)
outdated:
    go list -u -m -json all | go-mod-outdated -update

# detect known vulnerabilities (requires https://github.com/sonatype-nexus-community/nancy)
audit:
    go list -json -m all | nancy sleuth --loud

# run linters (requires https://github.com/dominikh/go-tools)
lint:
    staticcheck -f stylish ./... || \
        echo "\nRun \`just explain <LintIdentifier, e.g. SA1006>\` for details." && \
        exit 1

# explain lint identifier (e.g., "SA1006")
explain lint-identifier:
    staticcheck -explain {{lint-identifier}}

# add missing module requirements for imported packages, removes requirements that aren't used anymore
tidy:
    go mod tidy

# show dependencies
deps:
    go mod graph

# run tests
test:
    go test ./...

# run executable for local OS
run:
    @echo "Running sauber with defaults ..."
    go run -ldflags="-X 'main.Version={{version}}'" cmd/sauber/main.go

# build executable for local OS
build:
    @echo "Building executable for local OS ..."
    go build -ldflags="-X 'main.Version={{version}}'" -o sauber cmd/sauber/main.go

# build release executables for all supported platforms
release: test
    @echo "Building release executables (incl. cross compilation) ..."
    GOOS=darwin  GOARCH=arm64 go build -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_macos-arm64       cmd/sauber/main.go
    GOOS=linux   GOARCH=386   go build -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_linux-386         cmd/sauber/main.go
    GOOS=linux   GOARCH=amd64 go build -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_linux-amd64       cmd/sauber/main.go
