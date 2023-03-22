# Limestone

CLI app for downloading music from Slav Art on [Divolt](https://divolt.xyz).

## Building

Install Go through your system's package manager or from [the download page](https://go.dev/dl/).
Alternatively, I recommend using [go-installer](https://github.com/kerolloz/go-installer).

```bash
$ git clone https://github.com/dxbednarczyk/limestone
$ cd limestone
$ go get
$ make
```

## Installing from source

Installs to `~/.local/bin` by default, make sure this directory is somewhere on your PATH.

```bash
$ make && make install
$ which limestone
/home/damian/.local/bin/limestone
$ make uninstall
# This would silently remove $HOME/.config/limestone
Do you want to remove all saved logins and configuration? [y/n] n
rm -f /home/damian/.local/bin/limestone
```