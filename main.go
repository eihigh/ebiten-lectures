package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/ebiten/emoji"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	// 適当な画面サイズにする
	vw, vh = 800, 600

	speed       = 6
	bulletSpeed = 16
	emojiSize   = 128
)

var (
	emptyImage = ebiten.NewImage(3, 3) // vector描画用
)

type point struct {
	x, y float64
}

type enemy struct {
	pos  point
	kind enemyKind
	dead bool
}

type enemyKind int

const (
	normalEnemy enemyKind = iota
	homingEnemy
)

type app struct {
	// まずプレイヤーを矢印キーで動かせるようにする
	// 最初はx, y別の変数を使う形で説明して、次に構造体を作る
	playerPos point
	life      int

	// プレイヤーが動かせるようになったら、弾を撃てるようにする
	bulletPoses []*point

	// 弾が撃てるようになったら、敵を作る
	// 最初はnewAppで敵を作って、次のステップで時間を計測して敵を出すようにする
	enemies []*enemy

	// 敵ができたら、ゲームらしくしていく
	// 経過フレーム数を見ていろいろする
	tick int

	// シーンと呼ぶとUnityっぽすぎるので敢えて呼び方を変える
	gameState string

	mplusBigFont    font.Face
	mplusNormalFont font.Face
}

func newApp() (*app, error) {
	emptyImage.Fill(color.White)
	a := &app{}

	// フォントロード
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return nil, err
	}
	const dpi = 72
	a.mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	a.mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    72,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	// プレイヤー初期化
	a.playerPos.x = vw / 2
	a.playerPos.y = vh - 50
	a.life = 2

	// 敵の初期化

	// その他初期化
	a.gameState = "title"

	return a, nil
}

func (a *app) Update() error {
	a.tick++

	switch a.gameState {
	case "title":
		return a.updateTitle()
	case "gameover":
		return a.updateGameover()
	}

	// プレイヤーの移動
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		a.playerPos.x -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		a.playerPos.x += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		a.playerPos.y -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		a.playerPos.y += speed
	}

	// 弾を撃つ
	// TODO: 球数制限
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		p := &point{a.playerPos.x, a.playerPos.y}
		a.bulletPoses = append(a.bulletPoses, p)
	}
	// 弾を移動する
	for _, p := range a.bulletPoses {
		p.y -= bulletSpeed
	}

	// 画面外に出た弾を削除する
	// 分かれば簡単だけど、初学者向けじゃないかもなので、悩みどころ
	// （リストリンクの方がいいかもだけどGoなので使わない）
	next := make([]*point, 0, len(a.bulletPoses))
	for _, p := range a.bulletPoses {
		if p.y < -50 {
			continue
		}
		next = append(next, p)
	}
	a.bulletPoses = next

	// 敵との当たり判定、ライフ減少処理
	for _, e := range a.enemies {
		d2 := math.Pow(e.pos.x-a.playerPos.x, 2) + math.Pow(e.pos.y-a.playerPos.y, 2)
		r2 := math.Pow(32+32, 2)
		if d2 < r2 {
			// 当たった
			// 敵がいると無限に死に続けるので、敵も殺す
			e.dead = true
			a.life--
		}
	}

	// 死亡判定
	if a.life <= 0 {
		a.gameState = "gameover"
		return nil
	}

	// 敵をスポーンする
	if a.tick == 30 {
		// 30f: 上部に敵を３体スポーン
		for i := 0; i < 3; i++ {
			x := 200 + 200*i
			y := 50
			kind := normalEnemy
			a.enemies = append(a.enemies, &enemy{
				pos:  point{float64(x), float64(y)},
				kind: kind,
			})
		}
	}

	if a.tick == 180 {
		// 3s: 追尾する敵をスポーン
		x := 400
		y := 50
		kind := homingEnemy
		a.enemies = append(a.enemies, &enemy{
			pos:  point{float64(x), float64(y)},
			kind: kind,
		})
	}

	// 敵を移動する
	// ここで基本的なswitch文を覚える
	// 多態性はここではまだ使わない
	for _, e := range a.enemies {
		switch e.kind {
		case normalEnemy:
			// ふつうの敵
			// ただ下に移動する
			e.pos.y += 2

		case homingEnemy:
			// 追尾する敵
			vx := a.playerPos.x - e.pos.x
			vy := a.playerPos.y - e.pos.y
			l := math.Sqrt(vx*vx + vy*vy)
			spd := 2.2
			e.pos.x += vx / l * spd
			e.pos.y += vy / l * spd
		}
	}

	// 当たり判定する
	for _, b := range a.bulletPoses {
		for _, e := range a.enemies {
			d2 := math.Pow(b.x-e.pos.x, 2) + math.Pow(b.y-e.pos.y, 2)
			r2 := math.Pow(32+8, 2)
			if d2 < r2 {
				e.dead = true
			}
		}
	}

	// 死んだ敵を削除する
	enemies := make([]*enemy, 0, len(a.enemies))
	for _, e := range a.enemies {
		if e.dead {
			continue
		}
		enemies = append(enemies, e)
	}
	a.enemies = enemies

	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	switch a.gameState {
	case "title":
		a.drawTitle(screen)
		return
	case "gameover":
		a.drawGameover(screen)
		return
	}
	// 絵文字の画像サイズは128x128で、64x64に縮小して描画

	// emoji パッケージを使ってアセットを用意する手間を削減
	img := emoji.Image("😃")
	// 行列をどこまで説明するかは悩む・・・
	op := &ebiten.DrawImageOptions{}
	// 中央揃えの方法は先んじて説明しておく
	w, h := img.Size()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(a.playerPos.x, a.playerPos.y)
	screen.DrawImage(img, op)

	// 弾のグラフィックをどうするかは諸説ある、vector使う？
	for _, p := range a.bulletPoses {
		path := &vector.Path{}
		path.MoveTo(float32(p.x), float32(p.y))
		r := 16.0
		path.Arc(float32(p.x), float32(p.y), float32(r), 0, math.Pi*2, vector.Clockwise)
		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].SrcX = 1
			vs[i].SrcY = 1
			vs[i].ColorR = 0xff
			vs[i].ColorG = 0xff
			vs[i].ColorB = 0xff
		}
		screen.DrawTriangles(vs, is, emptyImage, nil)
	}

	// 敵を描画する
	for _, e := range a.enemies {
		img = emoji.Image("😈")
		w, h = img.Size()
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(e.pos.x, e.pos.y)
		screen.DrawImage(img, op)
	}

	// ライフを描画する
	img = emoji.Image("❤")
	for i := 0; i < a.life; i++ {
		x := float64(i) * 64
		y := 0.0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Translate(x, y)
		screen.DrawImage(img, op)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("life: %d", a.life))
}

func (a *app) Layout(ow, oh int) (int, int) {
	return vw, vh
}

func main() {
	ebiten.SetWindowSize(vw, vh)
	ebiten.SetWindowTitle("[capture]")
	a, err := newApp()
	if err != nil {
		panic(err)
	}
	if err := ebiten.RunGame(a); err != nil {
		panic(err)
	}
}
