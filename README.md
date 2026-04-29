# Liturgy of the Hours CLI

A terminal UI for reading the [Liturgy of the Hours](https://divineoffice.org/). Browse and read today's, yesterday's, and tomorrow's Divine Office prayers without leaving your terminal.

<!-- screenshot: full app overview -->
![App Overview](docs/screenshots/overview.gif)

---

## Installation/Usage

### With Go

```sh
go install github.com/tonye4/liturgyOfTheHoursCLI@latest
```
The binary will be available as `liturgyOfTheHoursCLI` anywhere in your terminal.

```sh
liturgyOfTheHoursCLI
```
### From Source
Run the following inside of the `LiturgyOfTheHoursCLI` directory.
```sh
go run .
```

### Updating
To update to the latest version, re-run the same command:

```sh
go install github.com/tonye4/liturgyOfTheHoursCLI@latest
```
### Navigation

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up in the prayer list |
| `↓` / `j` | Move down in the prayer list |
| `←` / `h` | Switch to previous day |
| `→` / `l` | Switch to next day |
| `Enter` / `Space` | Open selected prayer |
| `Esc` / `Backspace` / `q` | Go back to prayer list |
| `q` / `Ctrl+C` | Quit |

---

## Requirements

- A terminal with true color support
- Go 1.25+ (only needed if installing via `go install`)
