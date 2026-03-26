# QRGen

A simple QR code generator written in Go.
Generates QR codes from text input with configurable size, error correction, and version.

## Features

* Generate QR codes from any string input

* Adjustable scale (image size)
* Custom output file name
* Error correction levels (L, M, Q, H)
* Optional manual QR version override
* Verbose mode for debugging

## Installation

Clone the repository:

```bash
git clone https://github.com/aboutblank0/qr-code-generator.git
cd qrgen
```

Build the binary:

```bash
go build -o qrgen ./cmd/qrgen/
```

## Usage

```bash

qrgen [options] <content>
```

### Example

```bash
qrgen -scale 8 -output hello.png "Hello world"
```


## Options

| Flag       | Description                                        |
| ---------- | -------------------------------------------------- |
| `-help`    | Display help information                           |
| `-scale`   | Scale factor for the generated image (default: 10) |
| `-output`  | Output file name (default: `qrcode.png`)           |
| `-version` | Override QR version (1–40, auto if omitted)        |
| `-ec`      | Error correction level: L, M, Q, H (default: M)    |
| `-verbose` | Enable verbose output                              |

## Error Correction Levels

| Level | Description     |
| ----- | --------------- |
| L     | Low (~7%)       |
| M     | Medium (~15%)   |
| Q     | Quartile (~25%) |
| H     | High (~30%)     |

## License

MIT License

