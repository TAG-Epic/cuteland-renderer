package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"log"
	"net/http"
	"os"
)

var listen = "0.0.0.0:5050"
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)
var (
	backgroundSprite image.Image
)

func init() {
	//
	InfoLogger = log.New(os.Stdout, "INFO: ", 0)
	WarningLogger = log.New(os.Stdout, "INFO: ", 0)
	ErrorLogger = log.New(os.Stdout, "INFO: ", 0)

	// Sprite loading
	InfoLogger.Print("Loading sprites")
	loadSprite("background", &backgroundSprite)
}

func main() {
	// Register handlers
	http.HandleFunc("/render", renderBoard)

	// Start the server
	InfoLogger.Printf("Listening on %s", listen)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		ErrorLogger.Fatal(err)
	}
}

func renderBoard(w http.ResponseWriter, r *http.Request) {
	// Constants
	size := 1000
	margin := size / 4
	tileCount := 5
	tileSize := size / tileCount / 2

	context := gg.NewContext(size, size)

	// Background
	context.DrawImage(backgroundSprite, 0, 0)

	for x := 0; x < tileCount; x++ {
		for y := 0; y < tileCount; y++ {
			drawX := float64((tileSize * x) + margin)
			drawY := float64((tileSize * y) + margin)


			// Main tile
			context.DrawRectangle(drawX, drawY, float64(tileSize), float64(tileSize))
			context.SetHexColor("#e4e4a1")
			context.Fill()

			// Borders
			context.DrawLine(drawX, drawY, drawX+float64(tileSize), drawY)
			context.SetLineWidth(2)
			context.SetHexColor("#00FFFF")

			context.Fill()
		}
	}

	// Export
	_ = context.EncodePNG(w)
}

func loadSprite(spriteName string, sprite *image.Image) {
	InfoLogger.Printf("Loading sprite %s", spriteName)
	spriteImage, err := gg.LoadImage(fmt.Sprintf("../sprites/%s.png", spriteName))
	if err != nil {
		panic(err)
	}
	*sprite = spriteImage
}
