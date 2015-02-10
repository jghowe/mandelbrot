package main

import (
  "github.com/codegangsta/negroni"
  "net/http"
  "image"
  "image/png"
  "fmt"
  "strconv"
  "strings"

  "github.com/jghowe/mandelbrot"
)

func complexFromString(vals string) complex128 {
  fmt.Println("vals: ", vals)
  parts := strings.Split(vals,",")
  x, _ := strconv.ParseFloat(parts[0], 64)
  y, _ := strconv.ParseFloat(parts[1], 64)
  return complex(x,y)
}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/render", func(w http.ResponseWriter, req *http.Request) {
    query := req.URL.Query()

    width, _ := strconv.Atoi(req.FormValue("width"))
    height, _ := strconv.Atoi(req.FormValue("height"))
    center := complexFromString(req.FormValue("center"))
    zoom, _ := strconv.ParseFloat(req.FormValue("zoom"), 64)
    colors := strings.Split(req.FormValue("colors"), ",")
    spacing, _ := strconv.Atoi(query["colorSpacing"][0])

    dimensions := image.Pt(width,height)

    centerX :=real(center)
    centerY := imag(center)

    m := mandelbrot.Render(&dimensions, centerX, centerY, zoom, colors, spacing)

    header := w.Header()

    header["Content-Type"] = []string{"image/png"}

    w.WriteHeader(http.StatusOK)

    if err := png.Encode(w, m); err != nil {
      fmt.Println("image sent")
    } else {
      fmt.Println("error: ", err)
    }
  })

  static := negroni.NewStatic(http.Dir("./static"))

  n := negroni.Classic()
  n.Use(static)
  n.UseHandler(mux)
  n.Run(":3000")
}
