# regi

[![Go](https://github.com/iamharvey/regi/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/iamharvey/regi/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamharvey/regi)](https://goreportcard.com/report/github.com/iamharvey/regi) 
[![GoDoc](https://godoc.org/github.com/iamharvey/regi?status.svg)](https://godoc.org/github.com/iamharvey/regi)
<a href="https://codecov.io/gh/iamharvey/regi"><img src="https://codecov.io/gh/iamharvey/regi/main/graph/badge.svg" alt="codeCov"></a>
<a href="https://github.com/iamharvey/regi/blob/main/LICENSE.md"><img src="https://img.shields.io/github/license/iamharvey/regi" alt="License"></a>



regi is a CLI tool for managing your accessibility to multiple Docker registries.

<br>

## Features

Context management:
- [x] List all the contexts;
- [x] Add a new context (new connection settings);
- [x] Get info about a context;
- [x] Set current context;
- [x] Remove context.

Login
- [x] Login to current registry.

Image
- [x] List all the images with/without tags;
- [x] Tag & push local image to remote registry;
- [x] Pull image from remote registry;
- [x] Delete specific versions of an image from remote registry;
- [x] Delete image repository from remote registry;

more features are coming ...

<br><br>

## Install

### Install From Source

Clone the source first. Then, go the project folder and run:

```shell
go install
```

<br><br>

## User Guide

```shell
Usage:
  regi [flags]
  regi [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  context     Manage connection settings of multiple Docker registries.
  help        Help about any command
  image       Pull, push, delete and list images over Docker registry
  login       Login to current Docker registry.

Flags:
  -h, --help   help for regi
```

<br>

### Context Management

User can manage accessibility of multiple Docker registries via `context` command.
Each context denotes a connection setting of a Docker registry.

```shell
Available Commands:
  add         Add a new context.
  get         Get context info given context name.
  list        List all the contexts.
  set         Set current context with context name.

Flags:
  -h, --help   help for context
```

<br>

**Add Context**
New connection setting can be added via `context add`:

```shell
$ regi context add \
  --name=context1 \
  --server=http://192.168.0.168:5000 \
  --user=user \
  --password=123
  
Context added:
- name: context1
- server: http://192.168.0.168:5000
- insecure skip TLS verify: false
- user: user
- password: ***
```

<br>

**Get Context Info**

User can check out context info via `context get`:

```shell
$ regi context get context1

Context "context1":
- server: 192.168.0.168:5000
- insecure skip TLS verify: false
- user: user
- password: ***
```

<br>

**Set Current Context**

User can set current connection context via `context set`:

```shell
$ regi context set context1

Context switched. Current context: context1
```

<br>

**List Contexts**

User can list all the connection settings via `context list`:

```shell
$ regi context list

Contexts:
- 123.456.789.000:45678
- 123.456.789.000:12350
- 192.168.0.168:5000
- context1 <---
- context2
- context3
- context4
- context5
- context6
- context7
```

The current context can be identified by the dashed left arrow.

<br><br>

## Login

User can login to the current registry via `login`:

```shell
$ regi login

Connecting Docker registry with current context [context1](192.168.0.168:5000)
Authenticating with existing credentials...
Login Succeeded
```

Once you login to the registry, you can perform list, pull, push and delete of images and 
repositories on that registry.

<br><br>

## Image Management

User can list, pull, push and delete images on registry via `image` commands.

```shell
Usage:
  regi image [command]

Aliases:
  image, i, im, img

Available Commands:
  delete      Delete image from current registry.
  list        List images on current registry.
  pull        Pull image from current registry.
  push        Push image to current registry.

Flags:
  -h, --help   help for image

Use "regi image [command] --help" for more information about a command.
```

<br>

### List Images

User can list all the images with/without tags using `image list`:

```shell
$ regi image list

Images:
- gohash  [latest]
- golang  [1.18 1.17]
- my-first-func  [latest]
- mysql  [5.7 8.0]
- nltk  [latest]
- otp  [latest]
- print-pi  [latest]
- rand-str-gen  [latest]
- rsa  [latest]
- shc-ech-main  [1.0.0-dev]
- shc-grt-main  [1.0.0-dev]
```

<br>

### Pull Image

User can pull image from registry via `image pull`:

```shell
$ regi image pull golang 1.17

1.17: Pulling from golang                                                                                                                      15:44:34
Digest: sha256:bfb57478eb0b381f242b3ab27b373bca5516eb9d35eef98a41a0ba2742ab517d
Status: Image is up to date for 192.168.0.168:5000/golang:1.17
192.168.0.168:5000/golang:1.17
```

<br>

### Push Image

User can push image to registry via `image push`:

```shell
$ regi image push golang 1.17

The push refers to repository [192.168.0.168:5000/golang]
b208c5304d1a: Preparing
676f12fd4802: Preparing
2b09084b5ad5: Preparing
7372faf8e603: Preparing
9be7f4e74e71: Preparing
36cd374265f4: Preparing
5bdeef4a08f3: Preparing
36cd374265f4: Waiting
5bdeef4a08f3: Waiting
676f12fd4802: Layer already exists
9be7f4e74e71: Layer already exists
b208c5304d1a: Layer already exists
2b09084b5ad5: Layer already exists
7372faf8e603: Layer already exists
36cd374265f4: Layer already exists
5bdeef4a08f3: Layer already exists
1.17: digest: sha256:bfb57478eb0b381f242b3ab27b373bca5516eb9d35eef98a41a0ba2742ab517d size: 1796
```

<br>

### Delete Image

User can delete image via `image delete`:

```shell
$ regi image delete golang 1.17

image golang:1.17 is deleted
```

<br><br>

## Limitation

Regi is currently only support standard Docker registry. It is not tested with customized Docker registries (e.g., JFrog virtual Docker registry) or non Docker registries. Feel free to post issues or contribute.

<br><br>

## License

MIT License

Copyright (c) 2022 Harvey Li

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
