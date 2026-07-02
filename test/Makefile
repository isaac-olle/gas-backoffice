.PHONY: e2e e2e-verbose

ENV_FILE := .test.env

# Run E2E tests (compact output)
e2e:
	set -a && . ./$(ENV_FILE) && set +a && go test -tags e2e -count=1 ./e2e/

# Run E2E tests with verbose output (-v shows each subtest)
e2e-verbose:
	set -a && . ./$(ENV_FILE) && set +a && go test -tags e2e -count=1 -v ./e2e/
