[![Go Report Card](https://goreportcard.com/badge/github.com/softleader/dockerctl)](https://goreportcard.com/report/github.com/softleader/dockerctl)
![stability-stable](https://img.shields.io/badge/stability-stable-green.svg)
[![license](https://img.shields.io/github/license/softleader/dockerctl.svg)](./LICENSE)
[![release](https://img.shields.io/github/release/softleader/dockerctl.svg)](https://github.com/softleader/dockerctl/releases)

dockerctl is a command line tool that wraps `docker` in order to extend it with extra features and commands.

## Usage

``` sh
$ dockerctl --help 
```

dockerctl can be safely [aliased](#aliasing) as `docker` so you can type `$ docker COMMAND` in the shell and get all the usual `dockerctl` features.

## Installation

Download the latest [compiled binaries](https://github.com/softleader/dockerctl/releases) and put it anywhere in your executable path.

## Aliasing

Some dockerctl features feel best when it's aliased as `docker`. This is not dangerous; your _normal docker commands will all work_. dockerctl merely adds some sugar.

`dockerctl alias` displays instructions for the current shell. With the `-s` flag, it
outputs a script suitable for `eval`.

You should place this command in your `.bash_profile` or other startup script:

``` sh
eval "$(dockerctl alias)"
```

### Shell tab-completion

Place this command in your `.bash_profile` or other startup script, e.g.:

``` sh
echo "source <(slctl completion bash)" >> ~/.bashrc  # for bash users
echo "source <(slctl completion zsh)" >> ~/.zshrc    # for zsh users
```
