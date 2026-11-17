# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Google Cloud Function (2nd gen, Go, Functions Framework) that receives error log messages over HTTP from frontend web clients and writes them to stdout as JSON, where Cloud Logging picks them up. Module: `github.com/Informasjonsforvaltning/fdk-logging-function`, package `fdk_logging_function`.

The only caller is fdk-portal (`FDK_LOGGING_FUNCTION_URI` in its `src/env.ts`), which is being replaced by fdk-frontend. fdk-frontend does not use this function, so it is expected to be retired together with fdk-portal.

## Commands

```sh
go build ./...
go vet ./...
go test ./...
go test -run 'TestHandlers/error_valid' ./...   # single subtest
```

## Architecture

`fdk_logging_function_error.go` registers the entry point `ErrorLogging` in `init()` through `functions.HTTP(...)`. The handler:
- sets `Access-Control-Allow-Origin: *` on every response. `OPTIONS` returns 204 with CORS preflight headers, `POST` is handled, and any other method gets 405.
- decodes an `ErrorLogEntry` from a body capped at 64 KB (`maxBodyBytes`) and returns 400 unless `message`, `severity`, `application` and `image` are set. On success it returns 204 without echoing the input.
- logs the entry with `fmt.Println`, which calls its `String()` method (JSON marshal). The JSON line on stdout is the log record, so struct field names and JSON tags decide the shape of the Cloud Logging entry. `severity` is interpreted by Cloud Logging, and `namespace` (prod/staging/demo/development) identifies the calling environment.

## Deployment

GitHub Actions calls the shared reusable workflow `Informasjonsforvaltning/workflows/.github/workflows/deploy-cloud-function.yaml@main`. It uses `deploy-cloud-functions` v4, which only deploys 2nd gen functions, and runs `go-test-cmd` before deploying.
- **Dev** (`deploy-dev.yaml`): runs on PRs to `main` and deploys `cloud-logging-dev`, meant for fdk-portal's staging, demo and local development.
- **Prod** (`deploy-prod.yaml`): runs on push to `main` and deploys `cloud-logging`, meant for fdk-portal prod only.

The function URL is `https://europe-north1-digdir-cloud-functions.cloudfunctions.net/<name>`. fdk-portal reads it from the `FDK_LOGGING_FUNCTION_URI` key in the Kubernetes secret `commonurl`, with a fallback in `src/env.ts` for local development.

Both run in `digdir-cloud-functions`, region `europe-north1`. The Go version is set in `go.mod` and in each workflow (`go-version`, `runtime`). Update all of them together.
