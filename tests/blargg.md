# Tests de blargg

## cpu_instrs

- cpu_instrs OK

### individual
- 01 Passed
- 02 Passed
- 03 Passed
- 04 Passed
- 05 Passed
- 06 Passed
- 07 Passed
- 08 Passed
- 09 Passed
- 10 Passed
- 11 Passed

## mem_timing

### individual

- 01 Failed
- 02 Failed
- 03 Failed

## mem_timing-2

### rom_singles

- 01 Failed
- 02 Failed
- 03 Failed

## oam_bug

### rom_singles

- 01 Turning LCD on starts too early in scanline Failed #3
- 02 LD DE, $FE00: INC DE Failed #2
- 03 Passed
- 04 INC DE at first corruption Failed #3
- 05 Should corrupt at beginning of first scanline Failed #2
- 06 Passed
- 07 00000000 Failed
- 08 INC/DEC rp pattern is wrong Failed #2

## halt_bug
- Failed