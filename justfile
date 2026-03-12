set dotenv-load

[private]
default:
    just --list

# setup python workspace
[group: 'setup']
py-setup:
    cd sdk/python && uv venv && uv pip install -e .

# run linting
[group: 'python']
py-lint:
    cd sdk/python && uv run ruff check src/

# run formatter check
[group: 'python']
py-format:
    cd sdk/python && uv run ruff format --check

# static typing
[group: 'python']
py-typecheck:
    cd sdk/python && uv run ty check

# run tests
[group: 'python']
py-test:
    cd sdk/python && uv run pytest

# build server
[group: 'go']
go-build:
    go build -o bin/tempura ./cmd/tempura

# install tempura binary locally
[group: 'go']
go-install:
    go install ./cmd/tempura

# run linting
[group: 'go']
go-lint:
    golangci-lint run ./...

# run formatter check
[group: 'go']
go-format:
    test -z "$(gofmt -l .)"

# run tests
[group: 'go']
go-test:
    go test ./...

# creates a new git branch after checking out main (type: branch-name)
[group: 'git']
new branch:
    git stash
    git checkout main
    git pull
    git checkout -b {{branch}}
