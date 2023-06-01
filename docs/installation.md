---
title: Installation
order: 20
---

## System requirements

Rolens can run on the following operating systems:

* Windows 10/11 amd64/arm64
* Linux amd64/arm64
* macOS 10.13+ amd64 (Intel)
* macOS 11.0+ arm64 (Apple Silicon)

## Pre-compiled binaries

You can obtain a pre-compiled Rolens binary for macOS or installer for Windows from the [release page](https://github.com/garraflavatra/rolens/releases/latest).

## From source

Rolens is open-source software, which means that you can compile it from source on your own machine by cloning [the repository](https://github.com/garraflavatra/rolens).

### Prerequisites

Rolens is written in Go, so you should download the Go compiler from [the download page](https://go.dev/dl/). The minimum version required is 1.18. You can confirm whether it's installed correctly by running `go version` and checking that it outputs something similar to `go1.18.2`.

Furthermore, you need to have [Wails ^3.1](https://wails.io/docs/gettingstarted/installation) installed: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.

In order to compile the frontend, [Node.js](https://nodejs.org/en/download) ^16.0 and the [npm](https://npmjs.com) package manager ^8.0 (included in Node.js) are required. To confirm the installed versions of those tools, execute `node -v` and `npm -v`.

### Download source

To obtain a copy of the source code, do either of the following:

* Download a tarball or zip archive from the [release page](https://github.com/garraflavatra/rolens/releases/latest). Make sure you download the source archive, and not a pre-compiled binary.
* Or clone [the Git repository](https://github.com/garraflavatra/rolens): `git clone https://github.com/garraflavatra/rolens.git`.

### Compile

`cd` into the root directory of the source code and run either:

* `wails build` to generate an executable for your platform.
* `wails build -nsis` to generate an [NSIS installer](https://nsis.sourceforge.io/Main_Page). This requires that you have NSIS installed on your machine.

The generated binary will live in `build/bin`. You may want to run the installer (Windows) or move the app to the Applications folder (Mac).
