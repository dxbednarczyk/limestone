# Limestone
[![Go Report Card](https://goreportcard.com/badge/github.com/dxbednarczyk/limestone)](https://goreportcard.com/report/github.com/dxbednarczyk/limestone)

CLI app for downloading music from Slav Art.

## Usage

You can use `limestone` to download zipped albums/tracks from the Slav Art server, provided you have an account on [Divolt](https://divolt.xyz). As of v0.3, `limestone` does not support multi-factor authentication.

```bash
$ limestone login bob@example.com
Enter the password for bob@bob.com:
Logging in... login successful.
$ limestone divolt <url>
```

If you don't want to cache your login details, you can pass in your email and password as flags:
```
$ limestone divolt --email "bob@bob.com" --pass "bob123!" <url>
```

`limestone` also supports [the website](https://slavart.gamesdrive.net)'s API for individual tracks from Qobuz. You do not need to authenticate to use it.

All commands have help/usage information available, just pass `--help` or `-h` as a flag.

## Building

### Unix
Install Go through your system's package manager or from [the download page](https://go.dev/dl/).

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

### Unix

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
