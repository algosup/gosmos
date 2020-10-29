package main

import (
	"github.com/algosup/game"
	"github.com/algosup/game/color"
	"github.com/algosup/game/key"
)

const width = 640
const height = 480
const imageSize = 40
const alienGap = 10
const alienCount = 10
const alienVerticalSpeed = 5
const missileWidth = 3
const missileHeight = 8
const missileSpeed = 5

var base game.Bitmap
var alienBitmaps [6]game.Bitmap

var xBase = (width - imageSize) / 2
var xAlien = 10
var yAlien = 10
var deltaAlien = 1
var xMissile = 0
var yMissile = height
var missileOn = false

var isAlienAlive = [6][alienCount]bool{
	{true, true, true, true, true, true, true, true, true, true},
	{true, true, true, true, true, true, true, true, true, true},
	{true, true, true, true, true, true, true, true, true, true},
	{true, true, true, true, true, true, true, true, true, true},
	{true, true, true, true, true, true, true, true, true, true},
	{true, true, true, true, true, true, true, true, true, true},
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
func intersect(start1, end1, start2, end2 int) bool {
	start := max(start1, start2)
	end := min(end1, end2)
	return end > start
}

func intersectRect(
	left1,
	top1,
	right1,
	bottom1,
	left2,
	top2,
	right2,
	bottom2 int) bool {
	return intersect(left1, right1, left2, right2) &&
		intersect(top1, bottom1, top2, bottom2)
}

func checkCollision() {
	for j := 0; j < len(alienBitmaps); j++ {
		for i := 0; i < alienCount; i++ {
			left := xAlien + (imageSize+alienGap)*i
			right := left + imageSize
			top := yAlien + (imageSize+alienGap)*j
			bottom := top + imageSize
			if isAlienAlive[j][i] && intersectRect(
				left,
				top,
				right,
				bottom,
				xMissile,
				yMissile,
				xMissile+missileWidth,
				yMissile+missileHeight) {
				isAlienAlive[j][i] = false
				missileOn = false
			}
		}
	}
}
func drawAliens(surface game.Surface) {
	for j := 0; j < len(alienBitmaps); j++ {
		for i := 0; i < alienCount; i++ {
			if isAlienAlive[j][i] {

				game.DrawBitmap(
					surface,
					xAlien+(imageSize+alienGap)*i,
					yAlien+(imageSize+alienGap)*j,
					alienBitmaps[j])
			}
		}
	}
}

func isGameOver() bool {
	for j := 0; j < len(alienBitmaps); j++ {
		for i := 0; i < alienCount; i++ {
			if isAlienAlive[j][i] {
				y := yAlien + (imageSize+alienGap)*j
				if y > height-2*imageSize {
					return true
				}
			}
		}
	}

	return false
}

func move() {
	xAlien = xAlien + deltaAlien
	rightBorder := xAlien + (imageSize+alienGap)*alienCount - alienGap
	if rightBorder > width || xAlien < 0 {
		deltaAlien = -deltaAlien
		yAlien = yAlien + alienVerticalSpeed
	}
	if key.IsPressed(key.Left) {
		xBase--
		if xBase < (-imageSize) {
			xBase = width
		}
	}
	if key.IsPressed(key.Right) {
		xBase++
		if xBase > width {
			xBase = -imageSize
		}
	}
	if missileOn {
		checkCollision()
		yMissile = yMissile - missileSpeed
		if yMissile < (-missileHeight) {
			missileOn = false
		}

	} else {
		if key.IsPressed(key.Space) {
			missileOn = true
			xMissile = xBase +
				(imageSize-missileWidth)/2
			yMissile = height - imageSize/2
		}
	}
}
func draw(surface game.Surface) {
	drawAliens(surface)
	game.DrawBitmap(surface, xBase, height-imageSize, base)
	if missileOn {
		game.DrawRect(
			surface,
			xMissile,
			yMissile,
			missileWidth,
			missileHeight,
			color.White)
	}
	if isGameOver() {
		game.DrawText(surface, "GAME OVER", width/2-40, 20)
	} else {
		move()
	}

}

func main() {
	var e error
	base, e = game.LoadBitmap("base.png")
	if e != nil {
		panic(e)
	}
	files := []string{
		"blue",
		"green",
		"orange",
		"red",
		"violet",
		"yellow"}
	for i, v := range files {
		alienBitmaps[i], e = game.LoadBitmap(v + ".png")
		if e != nil {
			panic(e)
		}
	}
	game.Run("Gosmos", width, height, draw)
}
