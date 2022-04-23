# This project is shifted at https://git.mills.io/saltyim/app #

## Salty UI

A cross-platform, decentralized, chat app based on [SaltyIM](https://git.mills.io/saltyim/saltyim) for functionality and 
[GioUI](https://gioui.org/) for UI

## Supported Platform Status

- [x] WebAssembly (Modern Browsers)
- [x] Linux
- [x] Windows 
- [x] macOS
- [x] Android (Incomplete)
- [x] iOS / tvOS (Incomplete)

## Prerequisites

Before continuing, please make sure you satisfy prerequisites from the following:

* [Go](https://go.dev/)
* [SaltyIM](https://git.mills.io/saltyim/saltyim)
* [GioUI](https://gioui.org/)
* [GoGio](https://pkg.go.dev/gioui.org/cmd/gogio)

Install `gogio` with:

```#!console
go install gioui.org/cmd/gogio@latest
```

* [Android Studio](https://developer.android.com/studio) for android development.

### Local Development

Run (_from the terminal, inside the  root directory of this project, where [main.go](/main.go) file resides_):

```#!console
go run .
```

### Android Debug Development

Run the following to generate apk:

```#!console
gogio -target android .
```

* The above will generate saltyui.apk.

You can then install apk to the emulator or real device using:

```#!console
adb install saltyui.apk
```

### iOS Debug Development

```#!console
gogio -o saltyui.app -target ios .
```
Startup an iOS sim ( and wait for eternity )
```#!console
xcrun simctl install booted saltyui.app
```

### WebServer

Run the following to build the Web assets into `./web`:

```#!console
gogio -target js -o ./web .
```

## Troubleshooting

### Wasm
* Enable debugging in chrome [https://developer.chrome.com/blog/wasm-debugging-2020/](https://developer.chrome.com/blog/wasm-debugging-2020/)

### Weird Issues
* Inside was 