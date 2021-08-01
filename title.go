package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (a *app) updateTitle() error {
	// Press any key to start
	if len(inpututil.AppendPressedKeys(nil)) > 0 {
		a.gameState = "game"
		a.tick = 0
	}
	return nil
}

func (a *app) drawTitle(screen *ebiten.Image) {
	str := "EBITEN SHOOTING"
	// 中央揃えする
	r := text.BoundString(a.mplusBigFont, str)
	x := vw/2 - r.Dx()/2
	y := vh/2 + r.Dy()/2
	text.Draw(screen, str, a.mplusBigFont, x, y, color.White)

	// タイトルの少し下に中央揃えして描画
	str = "PRESS ANY KEY"
	r = text.BoundString(a.mplusNormalFont, str)
	x = vw/2 - r.Dx()/2
	y = vh/2 + 100 + r.Dy()/2
	text.Draw(screen, str, a.mplusNormalFont, x, y, color.White)
}
