package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fogleman/gg"
)

var listen = "0.0.0.0:5050"
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)
var (
	backgroundSprite image.Image
	rocksSprite      image.Image
)
var (
	grid_border_color string = "#707070"
	grid_color        string = "#e4e4a1"
	grid_id_color     string = "#000"
)
var (
	alphabet = []string{"a", "b", "c", "d", "e"}
)

func init() {
	// Loggers
	InfoLogger = log.New(os.Stdout, "INFO: ", 0)
	WarningLogger = log.New(os.Stdout, "INFO: ", 0)
	ErrorLogger = log.New(os.Stdout, "INFO: ", 0)

	// Sprite loading
	InfoLogger.Print("Loading sprites")
	loadSprite("background", &backgroundSprite)
	loadSprite("rocks", &rocksSprite)
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
	context.DrawImage(rocksSprite, 0, 0)

	for x := 0; x < tileCount; x++ {
		for y := 0; y < tileCount; y++ {
			drawX := float64((tileSize * x) + margin)
			drawY := float64((tileSize * y) + margin)

			// Main tile
			context.DrawRectangle(drawX+1, drawY+1, float64(tileSize-1), float64(tileSize-1))
			context.SetHexColor(grid_color)
			context.Fill()

			// Borders

			// Horizontal
			context.DrawLine(drawX, drawY, drawX+float64(tileSize), drawY)
			context.SetLineWidth(2)
			context.SetHexColor(grid_border_color)
			context.DrawLine(drawX, drawY+float64(tileSize), drawX+float64(tileSize), drawY+float64(tileSize))
			context.SetLineWidth(2)
			context.SetHexColor(grid_border_color)

			// Vertical
			context.DrawLine(drawX, drawY, drawX, drawY+float64(tileSize))
			context.SetLineWidth(2)
			context.SetHexColor(grid_border_color)
			context.DrawLine(drawX+float64(tileSize), drawY, drawX+float64(tileSize), drawY+float64(tileSize))
			context.SetLineWidth(2)
			context.SetHexColor(grid_border_color)

			context.Stroke()

			// Grid id
			context.SetHexColor(grid_id_color)
			context.DrawStringAnchored(alphabet[y]+strconv.FormatInt(int64(x), 10), drawX+10, drawY+float64(tileSize)-10, 1, 1)

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
