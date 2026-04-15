# ssh-wars

A [Charm](https://charm.sh) / [Bubble Tea](https://github.com/charmbracelet/bubbletea) ASCII animation player served over SSH. Connect and watch the classic Star Wars ASCII art right in your terminal.

Inspired by [towel.blinkenlights.nl](http://towel.blinkenlights.nl).

## Quick Start

### SSH (hosted)

```bash
ssh -p 2222 localhost
```

### Local (no SSH)

```bash
go run ./bin
```

### Docker

```bash
docker build -t ssh-wars .
docker run -p 2222:2222 ssh-wars
```

## Controls

| Key | Action |
|---|---|
| `Space` | Play / pause |
| `←` / `h` | Back one frame |
| `→` / `l` | Forward one frame |
| `↑` / `k` | Increase speed |
| `↓` / `j` | Decrease speed |
| `0`–`9` | Jump to position (0% – 90%) |
| `G` | Jump to end |
| `?` | Toggle help |
| `q` / `Ctrl+C` | Quit |

## Server Flags

```
-host string   Host to listen on (default "0.0.0.0")
-port int      Port to listen on (default 2222)
```

## License

[0BSD](LICENSE)
