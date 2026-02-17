# Repository Guidelines

## Project Structure & Module Organization
This repository is a single-package Go module: `github.com/akaswenwilk/jman`.

- Core library code is at the repo root (`arr.go`, `obj.go`, `equal.go`, `getter.go`, `setter.go`, etc.).
- Package docs live in `doc.go`; user-facing usage examples are in `README.md`.
- Tests are colocated with source as `*_test.go` (for example, `obj_test.go`, `arr_test.go`).
- CI configuration is in `.github/workflows/pull_request.yaml`.

When adding features, keep related logic in existing root files (or add another focused `*.go` file in root) and add matching tests nearby.

## Build, Test, and Development Commands
- `go test ./...`: run all unit tests across the module.
- `go test -run TestObj_Equal ./...`: run a focused subset while iterating.
- `go test -cover ./...`: check coverage impact of your change.
- `golangci-lint run`: run lint checks locally (matches CI linter family).
- `go fmt ./...`: format code before opening a PR.

## Coding Style & Naming Conventions
- Follow standard Go formatting (`gofmt`) and idioms.
- Use tabs/Go-default indentation (do not hand-align with spaces).
- Keep package name as `jman`; exported APIs use `CamelCase`, unexported helpers use `camelCase`.
- Preserve existing file naming style: short, responsibility-based names (`options.go`, `matchers.go`).
- Test names follow `Test<Type>_<Method>_<Scenario>` (example: `TestObj_Equal_Unequal_MissingKeyFromActual`).
- Prefer public interfaces that accept `jman.T` and fail via `t.Fatalf(...)` on invalid usage/inputs, instead of returning `error` values (unless explicitly requested otherwise).

## Testing Guidelines
- Add or update tests for every behavioral change.
- Prefer table-driven tests for multiple variants when it improves readability.
- Keep assertions explicit about JSON paths and mismatch behavior.
- Use targeted runs during development, then finish with `go test ./...`.

## Commit & Pull Request Guidelines
- Current history uses short, imperative commit messages (`add ...`, `fix ...`, `update ...`); keep the same style.
- Keep commits scoped to one logical change.
- PRs should include a brief description of behavior changes.
- PRs should link issue/context when applicable.
- PRs should include test evidence (command run and result).
- Ensure CI passes (`go test ./...` and lint checks) before requesting review.

## Required Agent Git Workflow
- Before making any file changes, pull the latest changes from `origin/master` and create a new worktree with a new branch from that updated base so your work is isolated from other parallel agents/users.
- Use descriptive branch names with a clear purpose (for example, `fix/obj-equal-missing-key` or `feat/setter-path-validation`).
- Make one or more well-scoped commits with short, imperative, descriptive messages.
- Push the branch to `origin` after committing.
- Use GitHub CLI to open a pull request from that branch (`gh pr create`), with a clear title and summary.
- Do not commit directly to `master` for agent-driven changes.
