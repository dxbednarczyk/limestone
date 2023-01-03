# Limestone

CLI app for downloading music from Slav Art on [Divolt](https://divolt.xyz).

## Building

Install Go through your system's package manager or from [the download page.](https://go.dev/dl/) Ensure it is somewhere on your path.

```bash
$ git clone https://github.com/dxbednarczyk/limestone
$ cd limestone
$ go get && go build
```

## Roadmap

- [x] Initial implementation
- [ ] Caching of emails and passwords (hashed, of course)
- [ ] Pretty it up using [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [ ] Ease of use for users not used to Slav Art (creating accounts, provider list)
