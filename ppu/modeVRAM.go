package ppu

func (ppu *PPU) runVRAM() {
	// Calculamos el número de ciclos basados en la posición de SCX y los sprites
	ppu.pixelFIFO = ppu.pixelFIFO[:0]
	if ppu.cycles < 172 {
		return
	}
	ppu.cycles -= 172
	ppu.setMode(ModeHBlank)

	// Si ya se dibujaron todas las líneas pasa a HBlank
	ly := ppu.bus.Read(LYRegister)
	if ly >= ScreenHeight {
		return
	}

	// Procedemos con el renderizado de una línea...
	scx := ppu.bus.Read(SCXRegister)
	scy := ppu.bus.Read(SCYRegister)
	lcdc := ppu.bus.Read(LCDCRegister)
	wx := ppu.bus.Read(WXRegister)
	wy := ppu.bus.Read(WYRegister)
	drawWindow := (lcdc&LCDCFlagWindowEnable) != 0 && int(ly) >= int(wy)

	bgTileMapAddr := uint16(0x9800)
	if lcdc&LCDCFlagBGTileMap != 0 {
		bgTileMapAddr = 0x9C00
	}

	//tileDataAddr := uint16(0x8800)
	useSigned := true
	if lcdc&LCDCFlagBGTileData != 0 {
		//tileDataAddr = 0x8000
		useSigned = false
	}

	for x := 0; x < ScreenWidth; x++ {
		var scrollX, scrollY uint16
		var tileMapAddr uint16

		if drawWindow && x >= int(wx)-7 {
			// Dibujamos Window
			tileMapAddr = 0x9800
			if lcdc&LCDCFlagWindowTileMap != 0 {
				tileMapAddr = 0x9C00
			}

			windowY := uint16(ly) - uint16(wy)
			windowX := uint16(x) - (uint16(wx) - 7)

			tileX := windowX / 8
			tileY := windowY / 8
			tileIndexOffset := tileY*32 + tileX
			tileIndex := ppu.bus.Read(tileMapAddr + tileIndexOffset)

			var tileAddr uint16
			if useSigned {
				signedTileIndex := int8(tileIndex)
				if signedTileIndex >= 0 {
					tileAddr = 0x9000 + uint16(signedTileIndex)*16
				} else {
					tileAddr = 0x9000 - uint16(signedTileIndex*-1)*16
				}
			} else {
				tileAddr = 0x8000 + uint16(tileIndex)*16
			}

			row := (windowY % 8) * 2
			byte1 := ppu.bus.Read(tileAddr + uint16(row))
			byte2 := ppu.bus.Read(tileAddr + uint16(row) + 1)
			bit := 7 - (scrollX % 8)

			colorID := (((byte2 >> bit) & 1) << 1) | ((byte1 >> bit) & 1)
			palette := ppu.bus.Read(0xFF47)
			color := (palette >> (colorID * 2)) & 0x03
			ppu.addPixelToFIFO(getColorFromPalette(color))

		} else {
			// Dibujamos Background
			scrollX = (uint16(x) + uint16(scx)) & 0xFF
			scrollY = (uint16(ly) + uint16(scy)) & 0xFF

			tileX := (scrollX) / 8
			tileY := (scrollY) / 8

			tileIndexOffset := tileY*32 + tileX
			tileIndex := ppu.bus.Read(bgTileMapAddr + tileIndexOffset)

			var tileAddr uint16
			if useSigned {
				tileAddr = 0x9000 + uint16(int8(tileIndex))*16
			} else {
				tileAddr = 0x8000 + uint16(tileIndex)*16
			}

			row := (scrollY % 8) * 2
			byte1 := ppu.bus.Read(tileAddr + uint16(row))
			byte2 := ppu.bus.Read(tileAddr + uint16(row) + 1)
			bit := 7 - (scrollX % 8)

			colorID := (((byte2 >> bit) & 1) << 1) | ((byte1 >> bit) & 1)
			palette := ppu.bus.Read(0xFF47)
			color := (palette >> (colorID * 2)) & 0x03
			ppu.addPixelToFIFO(getColorFromPalette(color))
		}
	}

	// Transferimos los píxeles de la FIFO al framebuffer
	for x := 0; x < ScreenWidth; x++ {
		ppu.popPixelFromFIFO(x, int(ly))
	}
	ppu.renderSprites()
}

func (ppu *PPU) renderSprites() {
	spriteHeight := byte(8)
	if !ppu.isObj8x8() {
		spriteHeight = 16
	}

	ly := ppu.bus.Read(LYRegister)

	for _, sprite := range ppu.spritesOnCurrentLine {
		spriteY := int(sprite.Y) - 16
		spriteX := int(sprite.X) - 8
		line := int(ly) - spriteY

		if sprite.Atributes&0x40 != 0 { // Y flip
			line = int(spriteHeight) - 1 - line
		}

		tileIndex := sprite.TileIndex
		if spriteHeight == 16 {
			tileIndex &= 0xFE // Ignorar bit 0 en modo 8x16
		}

		tileAddr := 0x8000 + uint16(tileIndex)*16 + uint16(line)*2
		byte1 := ppu.bus.Read(tileAddr)
		byte2 := ppu.bus.Read(tileAddr + 1)

		for x := 0; x < 8; x++ {
			bit := 7 - x
			if sprite.Atributes&0x20 != 0 { // X flip
				bit = x
			}

			colorID := (((byte2 >> bit) & 1) << 1) | ((byte1 >> bit) & 1)
			if colorID == 0 {
				continue // Transparente
			}

			var paletteAddr uint16 = 0xFF48
			if sprite.Atributes&0x10 != 0 {
				paletteAddr = 0xFF49
			}

			palette := ppu.bus.Read(paletteAddr)
			color := (palette >> (colorID * 2)) & 0x03
			screenX := spriteX + x

			if screenX < 0 || screenX >= ScreenWidth {
				continue
			}

			// Prioridad: fondo (bit 7)
			bgPriority := sprite.Atributes&0x80 != 0
			if bgPriority {
				// Omitimos dibujar si fondo no es color 0
				bgPixel := ppu.Framebuffer[getFramebufferIndex(screenX, int(ly))]
				if bgPixel != getColorFromPalette(0) {
					continue
				}
			}
			//ppu.addPixelToFIFO(getColorFromPalette(color))
			ppu.Framebuffer[getFramebufferIndex(screenX, int(ly))] = getColorFromPalette(color)
		}
	}

}
