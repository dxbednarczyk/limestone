# Limestone
[![Go Report Card](https://goreportcard.com/badge/github.com/dxbednarczyk/limestone)](https://goreportcard.com/report/github.com/dxbednarczyk/limestone)

CLI app for downloading music from Slav Art on [Divolt](https://divolt.xyz).

## Building

### Linux
Install Go through your system's package manager or from [the download page](https://go.dev/dl/).
Alternatively, I recommend using [go-installer](https://github.com/kerolloz/go-installer).

```bash
$ git clone https://github.com/dxbednarczyk/limestone
$ cd limestone
$ go get
$ make
```

### Windows
Install Go using `winget`, or alternatively download the latest MSI package from [go.dev/dl](https://go.dev/dl/). Preferably, use Powershell 7 or higher.

```powershell
> git clone https://github.com/dxbednarczyk/limestone
> cd limestone
> go get
> ./make
```

## Installing

### Linux

Installs to `~/.local/bin` by default, make sure this directory is somewhere on your PATH.

```bash
$ make && make install
$ which limestone
/home/damian/.local/bin/limestone
$ make uninstall
Do you want to remove all saved logins and configuration? [y/n] n
rm -f /home/damian/.local/bin/limestone
```

### Windows (planned)

*crickets\*