timestamp := `date +%s`

semver := "1.0.2"
commit := `git show -s --format=%h`
version := semver + "+" + commit

coverage_profile_log := "coverage_profile.txt"

# print available targets
default:
    @just --list --justfile {{justfile()}}

# evaluate and print all just variables
evaluate:
    @just --evaluate

# print system information such as OS and architecture
system-info:
  @echo "architecture: {{arch()}}"
  @echo "os: {{os()}}"
  @echo "os family: {{os_family()}}"

# print supported architectures for release builds (GOOS)
supported-architectures:
    go tool dist list

# format source code
format:
    @echo "Formatting source code ..."
    gofmt -l -s -w .

# detect outdated modules (requires https://github.com/psampaz/go-mod-outdated)
outdated:
    # `-mod=readonly` is required when using a vendored setup (like we have),
    # see https://go.dev/ref/mod#vendoring
    go list -mod=readonly -u -m -json all | go-mod-outdated -update

# detect known vulnerabilities
audit: vulnerabilities-govulncheck vulnerabilities-nancy
    go vet ./...

# detect known vulnerabilities (requires https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
vulnerabilities-govulncheck:
    govulncheck ./...

# detect known vulnerabilities (requires https://github.com/sonatype-nexus-community/nancy)
vulnerabilities-nancy:
    # `-mod=readonly` is required when using a vendored setup (like we have),
    # see https://go.dev/ref/mod#vendoring
    go list -mod=readonly -json -m all | nancy sleuth --loud

# run linters (requires https://github.com/dominikh/go-tools)
lint:
    staticcheck -f stylish ./... || \
        (echo "\nRun \`just explain <LintIdentifier, e.g. SA1006>\` for details." && \
        exit 1)

# explain lint identifier (e.g., "SA1006")
explain lint-identifier:
    staticcheck -explain {{lint-identifier}}

# add missing module requirements for imported packages, removes requirements that aren't used anymore
tidy:
    go mod tidy

# show dependencies
deps:
    go mod graph

# run tests with colorized output (requires https://github.com/kyoh86/richgo)
test *FLAGS:
    richgo test -cover {{FLAGS}} ./...

# run tests (vanilla), used for CI workflow
test-vanilla *FLAGS:
    go test -cover {{FLAGS}} ./...

# show test coverage
coverage:
    go test -coverprofile={{coverage_profile_log}} ./...
    go tool cover -html={{coverage_profile_log}}

# run executable for local OS
run:
    @echo "Running sauber with defaults ..."
    go run -ldflags="-X 'main.Version={{version}}'" cmd/sauber/main.go

# build executable for local OS
build: test-vanilla
    @echo "Building executable for local OS ..."
    go build -trimpath -ldflags="-X 'main.Version={{version}}'" -o sauber cmd/sauber/main.go

# build release executables for all supported platforms
release: test-vanilla
    @echo "Building release executables (incl. cross compilation) ..."
    # `go tool dist list` shows supported architectures (GOOS)
    GOOS=darwin GOARCH=arm64 \
        go build -trimpath -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_macos-arm64 cmd/sauber/main.go
    GOOS=linux  GOARCH=386 \
        go build -trimpath -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_linux-386   cmd/sauber/main.go
    GOOS=linux  GOARCH=amd64 \
        go build -trimpath -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_linux-amd64 cmd/sauber/main.go
    GOOS=linux  GOARCH=arm \
        go build -trimpath -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_linux-arm   cmd/sauber/main.go
    GOOS=linux  GOARCH=arm64 \
        go build -trimpath -ldflags "-X 'main.Version={{version}}' -s -w" -o sauber_linux-arm64 cmd/sauber/main.go

# vendor (https://go.dev/ref/mod#vendoring)
vendor:
    go mod vendor

# run build when sources change (requires https://github.com/watchexec/watchexec)
watch:
    # Watch all go files in the current directory and all subdirectories for
    # changes.  If something changed, re-run the build.
    @watchexec --clear -exts go -- just build

# run tests when sources change (requires https://github.com/watchexec/watchexec)
watch-test:
    # Watch all go files in the current directory and all subdirectories for
    # changes.  If something changed, re-run the build.
    @watchexec --clear --exts go -- just test
