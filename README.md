# hops

A TUI for managing `/etc/hosts` profiles. Inspired by [hostctl](https://github.com/guumaster/hostctl).

## What it does

- Manage groups of host entries as **profiles**
- Toggle profiles on/off without touching `/etc/hosts`
- Apply enabled profiles to `/etc/hosts` with a single keystroke (sudo only when needed)
- Import host lists from URLs
- Edit entries inline

Profiles are stored locally in `~/.local/share/hops/` as plain text host files. Nothing is written to `/etc/hosts` until you explicitly apply.

## Install

```bash
go install github.com/houz42/hops@latest
```

Or build from source:

```bash
git clone https://github.com/houz42/hops.git
cd hops
go build -o hops .
```

## Usage

```bash
hops              # uses /etc/hosts by default
hops -f /path/to/hosts  # custom hosts file
```

## Keybindings

### Profile List

| Key              | Action               |
| ---------------- | -------------------- |
| `↑/k` `↓/j`     | Navigate             |
| `Enter` / `Space`| Toggle enable/disable|
| `→/l`            | View entries         |
| `a`              | Add profile          |
| `d/x`            | Delete profile       |
| `i`              | Import from URL      |
| `S`              | Apply to /etc/hosts  |
| `q`              | Quit                 |

### Detail View

| Key          | Action       |
| ------------ | ------------ |
| `↑/k` `↓/j` | Navigate     |
| `←/h` / `Esc`| Back        |
| `a`          | Add entry    |
| `d`          | Delete entry |

## How it works

```
~/.local/share/hops/
├── state.json          # which profiles are enabled
└── profiles/
    ├── work.hosts      # 10.0.0.1 api.internal
    ├── staging.hosts   # 10.0.1.1 staging.example.com
    └── adblock.hosts   # imported from URL
```

Profiles are plain `IP hostname` text files. Toggle and edit freely — no privileges required.

When you press `S`, hops writes enabled profiles into `/etc/hosts` inside a managed block:

```
# your existing entries stay untouched

# --- BEGIN hops managed ---
# profile: work
10.0.0.1 api.internal
# profile: staging
10.0.1.1 staging.example.com
# --- END hops managed ---
```

Only the managed block is ever modified. Everything outside it is preserved.

## Credits

Inspired by [guumaster/hostctl](https://github.com/guumaster/hostctl). Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## License

MIT
