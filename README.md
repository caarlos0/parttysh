# parttysh

Party Parrot over SSH.

Based on https://github.com/hugomd/parrot.live

Try it out with:

```sh
ssh -p 2225 ssh.caarlos0.dev
```

## Why?

`¯\_(ツ)_/¯`

## How?

It uses a couple of Charm's packages:

- [bubbletea](https://github.com/charmbracelet/bubbletea): TUI framework
- [bubbles](https://github.com/charmbracelet/bubbles): TUI components
- [lipgloss](https://github.com/charmbracelet/lipgloss): Styling
- [wish](https://github.com/charmbracelet/wish): SSH apps
- [promwish](https://github.com/charmbracelet/promwish): prometheus middleware for wish
