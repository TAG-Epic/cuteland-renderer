package main

import (
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"log"
	"net/http"
	"os"
	"strings"
)

var listen = "0.0.0.0:5050"
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)
var (
	spriteCache = make(map[string]image.Image)
)

func init() {
	// Loggers
	InfoLogger = log.New(os.Stdout, "INFO: ", 0)
	WarningLogger = log.New(os.Stdout, "INFO: ", 0)
	ErrorLogger = log.New(os.Stdout, "INFO: ", 0)
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
	// Loading data
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ErrorLogger.Printf("Failed decoding JSON %s", err)
	}

	// Constants
	size := 1000
	margin := size / 4
	tileCount := 5
	tileSize := size / tileCount / 2

	context := gg.NewContext(size, size)

	InfoLogger.Printf("Drawing board with map type: %s", data["world"])

	// Background
	context.DrawImage(getSprite(data["world"].(string) + "-world"), 0, 0)



	for x := 0; x < tileCount; x++ {
		for y := 0; y < tileCount; y++ {
			drawX := float64((tileSize * x) + margin)
			drawY := float64((tileSize * y) + margin)
			currentTile := data["tiles"].([]interface{})[x].([]interface{})[y]

			InfoLogger.Printf("Tile pos: %f %f", drawX, drawY)
			InfoLogger.Printf("Tile at %d %d is %s", x, y, currentTile)
			if currentTile == nil {
				continue
			}
			context.DrawImage(getTile(currentTile.(string)), int(drawX+1), int(drawY+1))

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

func getTile(tileName string) image.Image {
	tileName = strings.Replace(tileName, ":", "_", -1)
	return getSprite(tileName)
}
