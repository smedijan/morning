# Morning

`morning` runs `.sh` and `.ps1` scripts directly inside `$XDG_CONFIG_HOME/morning.d` in alphanumeric filename order. Shell scripts are sourced with `sh`; PowerShell scripts are dot-sourced with `pwsh` (or `powershell` if `pwsh` is unavailable). Consecutive scripts using the same shell share shell state. A missing scripts directory is treated as an empty list; a script failure stops the run. Install the relevant interpreter for the script types you use.

When `XDG_CONFIG_HOME` is unset or empty, Morning uses `~/.config/morning.d`.

```sh
mkdir -p ~/.config/morning.d
cat > ~/.config/morning.d/10-update.sh <<'EOF'
# Runs as part of `morning`
echo "Good morning"
EOF
morning
```

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
