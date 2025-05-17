package bus

func (b *Bus) doDMATransfer(value byte) {
	source := uint16(value) << 8

	if b.DMAIsActive && b.dmaDelay > 0 {
		// DMA ya activa → almacenamos solicitud pendiente
		b.pendingDMASource = &source
	} else {
		// Nueva DMA → comienza después de 2 ciclos
		b.dmaDelay = 2
		b.dmaSource = source
		b.dmaIndex = 0
		b.dmaCyclesLeft = 160
		b.enableDMA = true
	}
}

// Actualiza 4 tcycles la transferencia DMA OAM
func (b *Bus) TickDMA() {
	b.Client = ClientDMA

	// Etapa de retardo: M=0 y M=1
	if b.enableDMA && b.dmaDelay > 0 {
		b.dmaDelay--
		if b.dmaDelay == 0 {
			b.DMAIsActive = true
			b.enableDMA = false
		}
		return
	}

	// Transferencia activa
	if !b.DMAIsActive {
		return
	}

	b.OAM[b.dmaIndex] = b.Read(b.dmaSource + b.dmaIndex)
	b.dmaIndex++
	b.dmaCyclesLeft--

	if b.dmaCyclesLeft <= 0 {
		b.DMAIsActive = false

		// Si hubo una DMA reiniciada mientras esta corría
		if b.pendingDMASource != nil {
			b.dmaSource = *b.pendingDMASource
			b.dmaIndex = 0
			b.dmaCyclesLeft = 160
			b.dmaDelay = 2
			b.enableDMA = true
			b.pendingDMASource = nil
		}
	}
}
