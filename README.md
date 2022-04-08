# Salty UI
A cross-platform, decentralized, chat app based on [SaltyIM](https://git.mills.io/saltyim/saltyim) for functionality and 
[GioUI](https://gioui.org/) for UI

## Platform Support
- [x] WebAssembly (Modern Browsers)
- [x] Linux
- [x] Windows
- [x] macOS
- [x] Android
- [x] iOS / tvOS

## Prerequisites
Before continuing, please make sure you have the prerequisites from the following
* [Go](https://go.dev/)
* [SaltyIM](https://git.mills.io/saltyim/saltyim)
* [GioUI](https://gioui.org/)

### Local Development
* Run ```go run .``` from the terminal, inside the  root directory of this project, where [main.go](/main.go) file resides.

### WebServer
```gogio -target js .```