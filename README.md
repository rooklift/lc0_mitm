# lc0_mitm

Extremely crude hack of [Nibbler](https://github.com/fohristiwhirl/nibbler) to run via a browser with a Golang backend that actually runs the engine.

```
go get github.com/gorilla/websocket
go build lc0_mitm.go
```

Then run the Go backend, and launch the HTML file in a browser. Tested on Firefox only, for now.
