package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (a *app) updateGameover() error {
	// Press any key to restart
	if len(inpututil.AppendPressedKeys(nil)) > 0 {
		a.gameState = "title"
		a.tick = 0
		// TODO: シーン構造体をちゃんと用意してやる
		a.life = 2
	}
	return nil
}

func (a *app) drawGameover(screen *ebiten.Image) {
	str := "GAME OVER"
	// 中央揃えする
	r := text.BoundString(a.mplusBigFont, str)
	x := vw/2 - r.Dx()/2
	y := vh/2 + r.Dy()/2
	text.Draw(screen, str, a.mplusBigFont, x, y, color.White)

	// タイトルの少し下に中央揃えして描画
	str = "PRESS ANY KEY TO RESTART"
	r = text.BoundString(a.mplusNormalFont, str)
	x = vw/2 - r.Dx()/2
	y = vh/2 + 100 + r.Dy()/2
	text.Draw(screen, str, a.mplusNormalFont, x, y, color.White)
}
