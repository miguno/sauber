timestamp := `date +%s`

semver := "1.1.1-alpha"
commit := `git show -s --format=%h`
version := semver + "+" + commit

coverage_profile_log := "coverage_profile.txt"

# print available targets
[group('project-agnostic')]
default:
    @just --list --justfile {{justfile()}}

# evaluate and print all just variables
[group('project-agnostic')]
evaluate:
    @just --evaluate

# print system information such as OS and architecture
[group('project-agnostic')]
system-info:
  @echo "architecture: {{arch()}}"
  @echo "os: {{os()}}"
  @echo "os family: {{os_family()}}"

# detect issues and known vulnerabilities
[group('security')]
audit: lint vulnerabilities

# build executable for local OS
[group('development')]
build: test-vanilla
    @echo "Building executable for local OS ..."
    go build -trimpath -ldflags="-X 'main.Version={{version}}'" -o sauber cmd/sauber/main.go

# show test coverage
[group('development')]
coverage:
    go test -coverprofile={{coverage_profile_log}} ./...
    go tool cover -html={{coverage_profile_log}}

# show dependencies
[group('development')]
deps:
    go mod graph

# explain lint identifier (e.g., "SA1006")
[group('development')]
explain lint-identifier:
    staticcheck -explain {{lint-identifier}}

# format source code
[group('development')]
format:
    @echo "Formatting source code ..."
    gofmt -l -s -w $(find . -name '*.go' -not -path './vendor/*')

# run all linters
[group('security')]
lint: lint-vet lint-staticcheck lint-golangci-lint

# run golangci-lint linter (requires https://github.com/golangci/golangci-lint)
[group('security')]
lint-golangci-lint:
    @golangci-lint run

# run staticcheck linter (requires https://github.com/dominikh/go-tools)
[group('security')]
lint-staticcheck:
    @staticcheck -f stylish ./... || \
        (echo "\nRun \`just explain <LintIdentifier, e.g. SA1006>\` for details." && \
        exit 1)

# alias for 'vet'
[group('security')]
lint-vet: vet

# detect outdated modules (requires https://github.com/psampaz/go-mod-outdated)
[group('development')]
outdated:
    # `-mod=readonly` is required when using a vendored setup (like we have),
    # see https://go.dev/ref/mod#vendoring
    go list -mod=readonly -u -m -json all | go-mod-outdated -update

# build release executables for all supported platforms
[group('development')]
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

# run executable for local OS
[group('development')]
run:
    @echo "Running sauber with defaults ..."
    go run -ldflags="-X 'main.Version={{version}}'" cmd/sauber/main.go

# print supported architectures for release builds (GOOS)
[group('development')]
supported-architectures:
    go tool dist list

# run tests with colorized output (requires https://github.com/kyoh86/richgo)
[group('development')]
test *FLAGS:
    richgo test -cover -race {{FLAGS}} ./...

# run tests (vanilla), used for CI workflow
[group('development')]
test-vanilla *FLAGS:
    go test -cover -race {{FLAGS}} ./...

# add missing module requirements for imported packages, removes requirements that aren't used anymore
[group('development')]
tidy:
    go mod tidy

# vendor (https://go.dev/ref/mod#vendoring)
[group('development')]
vendor:
    go mod vendor

# runs all vulnerability checks
[group('security')]
vulnerabilities: vulnerabilities-gosec vulnerabilities-govulncheck vulnerabilities-nancy

# analyze sources for security problems (requires https://github.com/securego/gosec/)
[group('security')]
vulnerabilities-gosec:
    gosec -tests ./...

# detect known vulnerabilities (requires https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
[group('security')]
vulnerabilities-govulncheck:
    govulncheck ./...

# detect known vulnerabilities (requires https://github.com/sonatype-nexus-community/nancy)
[group('security')]
vulnerabilities-nancy:
    # `-mod=readonly` is required when using a vendored setup (like we have),
    # see https://go.dev/ref/mod#vendoring
    go list -mod=readonly -json -m all | nancy sleuth --loud

# run build when sources change (requires https://github.com/watchexec/watchexec)
[group('development')]
watch:
    # Watch all go files in the current directory and all subdirectories for
    # changes.  If something changed, re-run the build.
    @watchexec --clear --exts go -- just build

# run tests when sources change (requires https://github.com/watchexec/watchexec)
[group('development')]
watch-test:
    # Watch all go files in the current directory and all subdirectories for
    # changes.  If something changed, re-run the build.
    @watchexec --clear --exts go -- just test

# vet the sources
[group('security')]
vet:
    go vet ./...
