# Contributing

Thanks for your interest in improving SamaSalaire.

## Workflow

1. Fork the repository and create a feature branch off `main`.
2. Run `make fmt vet test` before pushing.
3. Open a pull request with a clear description of the change.

## Commit style

Prefer Conventional Commits (e.g. `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`).
Keep the subject under 72 characters; explain *why* in the body when useful.

## Code style

- `gofmt -s` clean.
- Errors wrap with `%w` to preserve cause chains.
- Handlers return early on validation failures; avoid deeply nested code.

## Reporting issues

Open a GitHub issue with reproduction steps, expected vs actual behavior,
and the relevant log output.
