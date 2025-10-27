# DENSENDITHER

A powerful command-line image dithering tool that applies various dithering algorithms with custom color palettes. Perfect for creating retro-styled images, reducing color depth, and artistic image processing.

## Features

- **Multiple Dithering Algorithms**: Floyd-Steinberg and Ordered (Bayer) dithering
- **Custom Color Palettes**: Create, manage, and use your own color palettes
- **Image Resizing**: Optional resizing with high-quality Lanczos3 resampling
- **Easy CLI Interface**: Intuitive commands with comprehensive help
- **Batch Processing Ready**: Perfect for scripts and automation

## Installation

### Quick Install (Recommended)

```bash
git clone <your-repo-url>
cd densendither
./install.sh
```

The install script will:
1. Build the binary
2. Install it to your system PATH
3. Verify the installation

### Manual Installation

```bash
# Build the binary
go build -o densendither .

# Move to a directory in your PATH (optional)
sudo mv densendither /usr/local/bin/
```

### Requirements

- Go 1.21 or later
- Linux, macOS, or Windows

## Quick Start

```bash
# 1. Create a color palette
densendither palette add -n grayscale -c "#000000,#404040,#808080,#C0C0C0,#FFFFFF"

# 2. Dither an image
densendither dither -i input.jpg -p grayscale

# 3. View your result: dithered-grayscale-input.png
```

## Commands

### Palette Management

#### Create a Palette
```bash
# Add a new palette
densendither palette add -n pico8 -c "#000000,#1D2B53,#7E2553,#008751,#AB5236"

# Add more colors to existing palette
densendither palette add -n cyberpunk -c "#ff00ff,#00ffff,#ffff00"
```

#### List Palettes
```bash
# Show all available palettes
densendither palette list

# Show details of a specific palette
densendither palette show -n pico8
```

#### Remove Palettes
```bash
# Remove a palette
densendither palette remove -n old_palette
```

### Image Dithering

#### Basic Dithering
```bash
# Use Floyd-Steinberg dithering (default)
densendither dither -i image.png -p grayscale

# Use ordered dithering
densendither dither -i image.png -p pico8 -a ordered
```

#### With Resizing
```bash
# Resize to 800x600 before dithering
densendither dither -i large_image.jpg -p sunset -r 800x600

# Resize and use ordered dithering
densendither dither -i photo.png -p cyberpunk -r 400x300 -a ordered
```

### Image Resizing (Standalone)
```bash
# Resize image using Lanczos3 algorithm
densendither resize -i image.png -H 800 -w 600

# With custom output filename
densendither resize -i input.jpg -H 1200 -w 800 -o resized_image.png
```

## Detailed Usage

### Palette Command Reference

```bash
densendither palette [subcommand] [flags]

Subcommands:
  add     Add a new color palette
  remove  Remove an existing palette  
  list    List all available palettes
  show    Show details of a specific palette

Examples:
  densendither palette add -n "gameboy" -c "#0f380f,#306230,#8bac0f,#9bbc0f"
  densendither palette show -n gameboy
  densendither palette list
```

### Dither Command Reference

```bash
densendither dither [flags]

Required Flags:
  -i, --image     Path to input image file
  -p, --palette   Name of color palette to use

Optional Flags:
  -a, --algorithm  Dithering algorithm (floyd, ordered) [default: floyd]
  -r, --resize     Resize format: widthxheight (e.g., 800x600)

Examples:
  densendither dither -i photo.jpg -p pico8
  densendither dither -i artwork.png -p grayscale -a ordered
  densendither dither -i large.jpg -p gameboy -r 320x240
```

### Resize Command Reference

```bash
densendither resize [flags]

Required Flags:
  -i, --image     Path to input image file
  -H, --height    Target height in pixels
  -w, --width     Target width in pixels

Optional Flags:
  -o, --output    Output filename (auto-generated if not specified)
  -a, --algorithm Resampling algorithm (lanczos3) [default: lanczos3]

Examples:
  densendither resize -i large.png -H 800 -w 600
  densendither resize -i photo.jpg -H 400 -w 300 -o thumbnail.png
```

## Supported Formats

**Input**: JPEG, PNG, GIF, BMP, TIFF, WebP
**Output**: PNG (high quality, preserves transparency)

## Color Palette Format

Color palettes use hex color codes in the format `#RRGGBB` or `#RGB`.

**Valid Examples**:
- `#FF0000` (red)
- `#00FF00` (green) 
- `#0000FF` (blue)
- `#F00` (short red)

**Palette Limits**:
- Minimum: 1 color
- Maximum: 10 colors per palette
- Names must be unique and non-empty

## Dithering Algorithms

### Floyd-Steinberg (`floyd`)
- **Best for**: Photographs, complex images, smooth gradients
- **Characteristics**: Error diffusion creates organic, natural-looking patterns
- **Use when**: You want the highest quality dithering with natural textures

### Ordered/Bayer (`ordered`)
- **Best for**: Graphics, logos, geometric patterns
- **Characteristics**: Regular pattern, consistent texture
- **Use when**: You want predictable, uniform dithering patterns

## Configuration

Palettes are stored in `~/.config/densendither/conf.json`.

Example configuration:
```json
{
  "pico8": [
    "#000000", "#1D2B53", "#7E2553", "#008751",
    "#AB5236", "#5F574F", "#C2C3C7", "#FFF1E8"
  ],
  "grayscale": [
    "#000000", "#404040", "#808080", "#C0C0C0", "#FFFFFF"
  ]
}
```

## Examples & Use Cases

### Retro Gaming Art
```bash
# Create Game Boy palette
densendither palette add -n gameboy -c "#0f380f,#306230,#8bac0f,#9bbc0f"

# Process artwork
densendither dither -i character.png -p gameboy -r 160x144
```

### Web Graphics Optimization
```bash
# Create web-safe palette
densendither palette add -n websafe -c "#000000,#666666,#999999,#CCCCCC,#FFFFFF"

# Optimize image
densendither dither -i hero-image.jpg -p websafe -r 800x400
```

### Artistic Effects
```bash
# Create custom artistic palette
densendither palette add -n sunset -c "#ff6b35,#f7931e,#ffd23f,#874c62,#c98686"

# Apply artistic effect
densendither dither -i landscape.jpg -p sunset -a floyd
```

## Output Files

Files are automatically named using the pattern:
- **Dithered images**: `dithered-{palette}-{filename}.png`
- **Resized images**: `resized_{width}x{height}_{filename}.png` (or custom name)

## Tips & Best Practices

1. **Choose the right algorithm**:
   - Use Floyd-Steinberg for photos and natural images
   - Use Ordered for graphics and pixel art

2. **Optimize palette size**:
   - More colors = better quality but larger files
   - 4-8 colors often provide the best balance

3. **Consider resizing**:
   - Smaller images dither better
   - Resize before dithering for better results

4. **Test different palettes**:
   - Create multiple palettes for different moods
   - Use `palette show` to preview colors

## Troubleshooting

### Common Issues

**"palette not found"**: Use `densendither palette list` to see available palettes

**"invalid color format"**: Ensure colors use hex format (#RRGGBB or #RGB)

**"image failed to load"**: Check file path and format support

**Command not found**: Ensure densendither is in your PATH or use `./densendither`

### Getting Help

```bash
# General help
densendither --help

# Command-specific help  
densendither palette --help
densendither dither --help
densendither resize --help
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Your License Here]

## Changelog

### v1.0.0
- Initial release
- Floyd-Steinberg and Ordered dithering
- Palette management system
- Image resizing with Lanczos3
- Complete CLI interface