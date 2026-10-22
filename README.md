# Morning

A Go project.

## Development

This repository uses [mise](https://mise.jdx.dev/) to manage the toolchain and [hk](https://hk.jdx.dev/) for quality checks and Git hooks.

```sh
mise install
hk install
hk check
```

The pre-push hook runs the full suite, including `gosec` and `govulncheck`. To run it manually:

```sh
mise exec -- hk run pre-push
```

GitHub Actions runs `hk check --all` on pushes and pull requests.
