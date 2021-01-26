package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"net/url"
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
	borderSprite     image.Image
)
var (
	grid_border_color string = "#707070"
	grid_color        string = "#e4e4a1"
	grid_id_color     string = "#000"
)
var (
	alphabet    = []string{"A", "B", "C", "D", "E"}
	spriteCache = make(map[string]image.Image)
)

func init() {
	// Loggers
	InfoLogger = log.New(os.Stdout, "INFO: ", 0)
	WarningLogger = log.New(os.Stdout, "INFO: ", 0)
	ErrorLogger = log.New(os.Stdout, "INFO: ", 0)

	// Sprite loading
	InfoLogger.Print("Loading sprites")
	loadSprite("background", &backgroundSprite)
	loadSprite("default-border", &borderSprite)
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
	world_name = r.URL.Query()
	size := 1000
	margin := size / 4
	tileCount := 5
	tileSize := size / tileCount / 2

	context := gg.NewContext(size, size)

	// Background
	context.DrawImage(backgroundSprite, 0, 0)
	context.DrawImage(borderSprite, 0, 0)

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
			gridX := drawX + 20
			gridY := drawY + float64(tileSize-1) - 20

			context.SetHexColor(grid_id_color)
			context.DrawStringAnchored(alphabet[y]+strconv.FormatInt(int64(x+1), 10), gridX, gridY, 1, 1)

		}
	}

	// Export
	_ = context.EncodePNG(w)
}

func loadSprite(spriteName string, sprite *image.Image) bool {
	InfoLogger.Printf("Loading sprite %s", spriteName)
	spriteImage, err := gg.LoadImage(fmt.Sprintf("../sprites/%s.png", spriteName))
	if err != nil {
		return false
	}
	*sprite = spriteImage
	return true
}

func getSprite(spriteName string) image.Image {
	sprite, ok := spriteCache[spriteName]
	if !ok {
		ok = loadSprite(spriteName, &sprite)
		if !ok {
			return nil
		}
		spriteCache[spriteName] = sprite
		return sprite
	}
	return sprite
}
