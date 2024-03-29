# Limestone
[![Go Report Card](https://goreportcard.com/badge/github.com/dxbednarczyk/limestone)](https://goreportcard.com/report/github.com/dxbednarczyk/limestone)
[![Maintainability](https://api.codeclimate.com/v1/badges/173a1d2d8003627a7261/maintainability)](https://codeclimate.com/github/dxbednarczyk/limestone/maintainability)

CLI app for downloading music from Slav Art.

## Usage

All commands have help/usage information available, just pass `--help` or `-h` as a flag.

### Divolt 

You can use `limestone` to download zipped albums/tracks from the Slav Art server, provided you have an account on [Divolt](https://divolt.xyz). `limestone` does not currently support multi-factor authentication.

```bash
$ limestone login bob@bob.com
Enter the password for bob@bob.com:
Logging in... login successful.
$ limestone divolt <url>
```

You can specify the quality of the download, according to the table on [the FAQ](https://rentry.org/slavart):
```bash
$ limestone divolt -q 3 <url>
```

### Web

`limestone` also supports [the website](https://slavart.gamesdrive.net)'s API for individual tracks from Qobuz. You do not need to authenticate to use it. This download method only downloads the highest quality available.

```bash
$ limestone web "the police"
Getting results for query "the police"...

# Searching for a track in the TUI... found one!

# Fancy progress bar...

Downloaded to /home/dxbednarczyk/Downloads/The Police - Every Breath You Take.flac
```

You can also use the `-c` flag to automatically download the closest match to your query.

```bash
$ limestone web -c "cherry bomb tyler the creator"
Getting results for query "cherry bomb tyler the creator"...

# Fancy progress bar...

Downloaded to /home/dxbednarczyk/Downloads/Tyler, The Creator - CHERRY BOMB.flac
```

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
Install Go using `winget`, or alternatively download the latest MSI package from [the download page](https://go.dev/dl/). Preferably, use Powershell 7 or higher.

```powershell
> git clone https://github.com/dxbednarczyk/limestone
> cd limestone
> go get
> ./make
```

## Installing

### Using `go install`

```bash
$ go install github.com/dxbednarczyk/limestone@latest
```

### Using `make`

Installs to `~/.local/bin` by default, make sure this directory is somewhere on your PATH.

```bash
$ make && make install
$ which limestone
/home/damian/.local/bin/limestone
$ make uninstall
Do you want to remove all saved logins and configuration? [y/n] n
rm -f /home/damian/.local/bin/limestone
```

## Comparison
See [similar tools.](https://github.com/tywil04/slavartdl?tab=readme-ov-file#similar-tools)
