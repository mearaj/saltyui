# Salty UI
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
Before continuing, please make sure you satisfy prerequisites from the following
* [Go](https://go.dev/)
* [SaltyIM](https://git.mills.io/saltyim/saltyim)
* [GioUI](https://gioui.org/)
* [GoGio](https://pkg.go.dev/gioui.org/cmd/gogio)
```go install gioui.org/cmd/gogio@latest```
* [Android Studio](https://developer.android.com/studio) for android development.

### Local Development
* Run ```go run .``` from the terminal, inside the  root directory of this project, where [main.go](/main.go) file resides.

### Android Debug Development
* Run the following to generate apk<br>
```gogio -target android .```
* The above will generate saltyui.apk.
You can then install apk to the emulator or real device using<br>
```adb install slatyui.apk```

### WebServer
```gogio -target js .```