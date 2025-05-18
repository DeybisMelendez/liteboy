package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deybismelendez/liteboy/bus"
	"github.com/deybismelendez/liteboy/cartridge"
	"github.com/deybismelendez/liteboy/cpu"
	"github.com/deybismelendez/liteboy/ppu"
	"github.com/deybismelendez/liteboy/timer"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 160
	ScreenHeight = 144
	Scale        = 4
	TargetCycle  = 70224
)

type Game struct {
	cpu       *cpu.CPU
	ppu       *ppu.PPU
	bus       *bus.Bus
	lastCycle int
	image     *ebiten.Image
}

func NewGame(cpu *cpu.CPU, ppu *ppu.PPU, bus *bus.Bus) *Game {
	return &Game{
		cpu:   cpu,
		ppu:   ppu,
		bus:   bus,
		image: ebiten.NewImage(ScreenWidth, ScreenHeight),
	}
}

func (g *Game) Update() error {
	// Ejecutamos ciclos hasta llegar al target para refrescar la pantalla (simil a frame)
	g.lastCycle += g.cpu.Step()
	if g.lastCycle > TargetCycle {
		g.lastCycle -= TargetCycle
		g.image.WritePixels(g.ppu.Framebuffer)
	}
	// Leer el valor del registro P1 (0xFF00)
	p1 := g.bus.Read(0xFF00)

	// Bit 4: dirección (0=activado), Bit 5: botones
	directionKeys := (p1 & (1 << 4)) == 0
	buttonKeys := (p1 & (1 << 5)) == 0

	var input byte = 0x0F // bits 0-3: todos presionados (1 = no presionado)

	if directionKeys {
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			input &= ^byte(1 << 0) // Bit 0 - Derecha
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			input &= ^byte(1 << 1) // Bit 1 - Izquierda
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			input &= ^byte(1 << 2) // Bit 2 - Arriba
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			input &= ^byte(1 << 3) // Bit 3 - Abajo
		}
	}

	if buttonKeys {
		if ebiten.IsKeyPressed(ebiten.KeyZ) {
			input &= ^byte(1 << 0) // Bit 0 - A
		}
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			input &= ^byte(1 << 1) // Bit 1 - B
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			input &= ^byte(1 << 2) // Bit 2 - Select
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			input &= ^byte(1 << 3) // Bit 3 - Start
		}
	}

	// Escribir bits 0-3 en el registro FF00 sin tocar bits 4-7
	g.bus.Write(0xFF00, (p1&0xF0)|input)
	// Actualizar la imagen con el framebuffer RGBA (suponiendo que ppu.Framebuffer es []byte RGBA8888)
	// ebiten espera un slice []byte con pixels en formato RGBA8888

	// Asegurarse que el framebuffer tenga el tamaño correcto
	if len(g.ppu.Framebuffer) != ScreenWidth*ScreenHeight*4 {
		return fmt.Errorf("framebuffer size incorrecta")
	}

	// Copy pixels directo a ebiten.Image
	// ebiten.Image.WritePixels espera []byte en formato RGBA8888

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Escalar la imagen a la ventana (multiplicando el tamaño)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(Scale, Scale)
	screen.DrawImage(g.image, op)

	// Mostrar FPS en pantalla
	msg := fmt.Sprintf("LiteBoy Emulator - Press ESC to quit\nFPS: %.2f", ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth * Scale, ScreenHeight * Scale
}

func loadROM(path string) *cartridge.Cartridge {
	cart := cartridge.NewCartridge(path)
	if cart == nil {
		log.Fatalf("Error al cargar ROM: %s", path)
	}
	return cart
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <path_a_la_rom.gb>")
		return
	}

	romPath := os.Args[1]
	if len(os.Args) == 3 && os.Args[2] == "--info" {
		cart := loadROM(romPath)
		cart.PrintHeaderInfo()
		os.Exit(0)
	}

	cart := loadROM(romPath)
	gameBus := bus.NewBus(cart)
	gamePPU := ppu.NewPPU(gameBus)
	gameTimer := timer.NewTimer(gameBus)
	gameCPU := cpu.NewCPU(gameBus, gameTimer, gamePPU)

	game := NewGame(gameCPU, gamePPU, gameBus)

	// Configurar ventana y correr el loop de Ebiten
	ebiten.SetWindowSize(ScreenWidth*Scale, ScreenHeight*Scale)
	ebiten.SetWindowTitle("LiteBoy Emulator")
	ebiten.SetTPS(800_000)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
