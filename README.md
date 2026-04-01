# lazymake

Interactive TUI for Makefiles — like [lazygit](https://github.com/jesseduffield/lazygit), but for `make`.

Fuzzy-search your Makefile targets and run them instantly.

![Go](https://img.shields.io/github/go-mod/go-version/febriliankr/lazymake)
![License](https://img.shields.io/github/license/febriliankr/lazymake)

## Demo

```
  lazymake   Makefile

  > build_

    build         Build the binary
  > build-docker  Build Docker image
    clean         Remove build artifacts
    dev           Run in dev mode
    test          Run all tests

  ↑↓ navigate • enter select • esc quit
```

## Install

### Homebrew

```bash
brew tap febriliankr/tap
brew install lazymake
```

### Ubuntu/Debian

```bash
echo "deb [trusted=yes] https://febriliankr.github.io/lazymake/deb stable main" | sudo tee /etc/apt/sources.list.d/lazymake.list
sudo apt update
sudo apt install lazymake
```

Or download the `.deb` directly from [Releases](https://github.com/febriliankr/lazymake/releases):

```bash
sudo dpkg -i lazymake_*.deb
```

### Go

```bash
go install github.com/febriliankr/lazymake@latest
```

### From source

```bash
git clone https://github.com/febriliankr/lazymake.git
cd lazymake
make install
```

## Usage

```bash
lazymake          # pick a target from ./Makefile
lazymake -r       # recursively find Makefiles in subdirectories
lazymake -f path  # specify a Makefile
lazymake -n       # dry-run: print the command instead of running it
lazymake -v       # print version
```

### Makefile comments

lazymake picks up target descriptions from comments:

```makefile
## Build the binary
build:
	go build -o app .

test: ## Run all tests
	go test ./...
```

Both styles work — comment above the target, or inline `##` comment.

## Keybindings

| Key | Action |
|-----|--------|
| `↑` / `ctrl+k` | Move up |
| `↓` / `ctrl+j` | Move down |
| `enter` | Run selected target |
| `esc` / `ctrl+c` | Quit |
| Type anything | Fuzzy filter |

## License

MIT
