package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strconv"

	"github.com/fogleman/gg"
)

var worldName string
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
	alphabet = []string{"A", "B", "C", "D", "E"}
)

func init() {
	// Loggers
	InfoLogger = log.New(os.Stdout, "INFO: ", 0)
	WarningLogger = log.New(os.Stdout, "INFO: ", 0)
	ErrorLogger = log.New(os.Stdout, "INFO: ", 0)

}

func main() {
	var world string
	InfoLogger.Printf("Please input a world name")
	_, _ = fmt.Scanln(&world)
	worldName = world

	InfoLogger.Printf("Please input a background sprite name")
	var backgroundSpriteName string
	_, _ = fmt.Scanln(&backgroundSpriteName)
	loadSprite(backgroundSpriteName, &backgroundSprite)

	InfoLogger.Printf("Please input a border sprite name")
	var borderSpriteName string
	_, _ = fmt.Scanln(&borderSpriteName)
	loadSprite(borderSpriteName, &rocksSprite)

	renderBoard()
}

func renderBoard() {
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
			gridX := drawX + 20
			gridY := drawY + float64(tileSize-1) - 20

			context.SetHexColor(grid_id_color)
			context.DrawStringAnchored(alphabet[y]+strconv.FormatInt(int64(x+1), 10), gridX, gridY, 1, 1)

		}
	}

	// Export
	_ = context.SavePNG(fmt.Sprintf("../../sprites/%s-world.png", worldName))
}

func loadSprite(spriteName string, sprite *image.Image) {
	InfoLogger.Printf("Loading sprite %s", spriteName)
	spriteImage, err := gg.LoadImage(fmt.Sprintf("../../sprites/%s.png", spriteName))
	if err != nil {
		panic(err)
	}
	*sprite = spriteImage
}
