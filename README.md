# lc0_mitm

Extremely crude hack of [Nibbler](https://github.com/fohristiwhirl/nibbler) to run via a browser with a Golang backend that actually runs the engine.

```
go get github.com/gorilla/websocket
go build lc0_mitm.go
```

Then run the Go backend in the same file as `lc0.exe`, and launch the HTML file in a browser. Tested on Firefox only, for now.

In the event that you're on Linux, you probably have `lc0` instead of `lc0.exe` so you can either rename it or change that bit of code in the .go file.
