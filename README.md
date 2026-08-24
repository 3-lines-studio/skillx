# Skillx

Skillx provides Agent Skills through the Unix tool convention.

```sh
skillx describe
printf '%s\n' '{}' | skillx run skills
printf '%s\n' '{"name":"example"}' | skillx run skill
```

Tools:

```text
skills
skill
```

Skillx merges `./.agents/skills` and `$HOME/.agents/skills`. Project skills take precedence when both roots contain the same name. Missing default roots are ignored. Set `SKILLX_ROOT` to replace both defaults with one explicit root. Skillx uses a separate Go `os.Root` for each root and strict skill names to confine reads.

Run the local coding stack:

```sh
AX_TOOLS="fsx bashx skillx" ax
```

## Build

```sh
go test ./...
go build -o skillx .
```

## Install

```sh
curl -fsSL https://ax.3lines.studio/install.sh | sh -s -- skillx
```
