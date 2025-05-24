# LiteBoy 🎮 (En desarrollo)
_Un emulador sencillo de Game Boy escrito en Go_

LiteBoy es un emulador ligero de la consola Nintendo Game Boy, desarrollado usando el lenguaje Go y el motor de videojuegos Ebitengine. Su objetivo es ser simple, rápido y fácil de entender para aprender sobre emulación y desarrollo de sistemas gráficos retro.

## 📦 Requisitos

- Go 1.24
- Ebitengine 2.8

## Uso

Ejecutar en raiz del proyecto:

go run . [path-rom]

Agrega --info para visualizar información de la rom.

Para ejecutar tests requiere descargar los test rom de Blargg y Mooneye en la carpeta roms/blargg y roms/mooneye respectivamente. Luego puedes proceder a ejecutar go test.

# Que hace bien el emulador

- Ejecuta decentemente todas las instrucciones de CPU con timings correctos
- Realiza un renderizado de imagen decente pero sin timings exactos
- Genera audio de los canales 1, 2 y 3 decentemente
- Lee cartuchos de tipo ROM ONLY, MBC1, MBC2, MBC3, MBC5, MBC7 (algunos no están completos)
- Pasa todos los tests de Blargg excepto los que prueban bugs
- Pasa casi todos los test de Mooneye excepto los de PPU
- Pasa el test de dmg-acid2

## TODO

- APU:
    - Falta generar canal 4 (ruido) correctamente
    - Falta mejorar los canales de audio
- PPU:
    - Falta mejorar timings
- Cartuchos:
    - Agregar soporte a mas tipos de cartuchos
    - Mejorar soporte de algunos tipos de cartuchos
- Otros:
    - Refactorizar y optimizar proyecto

## ¡Se busca contribución!

Si te gusta el proyecto, puedes contribuir en el desarrollo enviando pull request.