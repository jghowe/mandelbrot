# mandelbrot

An app to explore and generate fractal images from the Mandelbrot set.

### Dependencies

* [go](https://golang.org)
* [npm](http://npmjs.org)
* [bower](http://bower.io)
* [grunt](http://gruntjs.com)

### Installation

To install, make sure your GOPATH and go workspace is setup correctly. Then run the following.

```bash
go get github.com/jghowe/mandelbrot
npm install
```

### Running

To run, set `GOMAXPROCS` to the number of cores available on your machine and then use `npm start`.

```bash
GOMAXPROCS=4 npm start
```

This will start a local web server at http://localhost:3000.  You can use the mouse wheel or trackpad to zoom in/out and click and drag to pan.
Click **Render** in the top right corner to render a new image at the current location.
