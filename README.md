# manga-cli

A powerful command-line tool for searching, downloading, and reading manga directly in your terminal. Built with Go, this tool provides a seamless experience for manga enthusiasts who prefer working in the terminal.

## Features

- 🔍 Search manga titles with fuzzy matching
- 📚 List available chapters for any manga
- ⬇️ Download manga chapters as images
- 📖 Built-in terminal manga reader with image viewer support
- ⚙️ Configurable settings for image viewer dimensions
- 🎯 Support for chapter ranges and batch operations

## Installation

### Prerequisites

- Go 1.16 or higher
- [viu](https://github.com/atanunq/viu) (recommended image viewer for terminal)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/TaichiKarna/manga-cli.git
cd manga-cli

# Build the project
go build -o manga-cli

# Install globally (optional)
sudo mv manga-cli /usr/local/bin/
```

## Usage

### Search and Read Manga

```bash
# Search for a manga and read it
manga-cli search --title "One Piece"

# Specify image viewer dimensions
manga-cli search --title "One Piece" --width 100 --height 50
```

### Download Manga Chapters

```bash
# Download specific chapters
manga-cli download --title "One Piece" --from 1 --to 10

# Download a single chapter
manga-cli download --title "One Piece" --chapter 1
```

### List Available Chapters

```bash
# List downloaded manga
manga-cli list 

# List chapters downloaded for a downloaded mangaa
manga-cli list --title "One Piece"
```

### Read Downloaded Manga

```bash
# Read downloaded manga
manga-cli read --title "One Piece" --chapter 1 --width 100 --height 50
```

### Config

The `config` command allows you to manage your manga-cli settings. You can view, set, and list configuration options.

```bash
# List all configuration options
manga-cli config list

# Get a specific configuration value
manga-cli config get path

# Set a configuration value
manga-cli config set path ~/manga
manga-cli config set width 100
manga-cli config set height 50
manga-cli config set viewer viu
```

Available configuration options:
- `path`: Directory where manga chapters are downloaded
- `width`: Width of the image viewer in characters
- `height`: Height of the image viewer in characters
- `viewer`: Terminal image viewer to use (currently supports viu)

Configuration file location: `~/.config/manga-cli/config.json`

Example configuration:
```json
{
    "width": 100,
    "height": 50,
    "path": "~/manga",
    "viewer": "viu"
}
```

## Configuration

The tool uses a configuration file to store user preferences. You can configure:

- Default image viewer dimensions (width and height)
- Download directory for manga chapters
- Image viewer settings (currently supports viu)

Configuration file location: `~/.config/manga-cli/config.json`

Example configuration:
```json
{
    "width": 100,
    "height": 50,
    "path": "~/manga",
    "viwer": "viu"
}
```

## Future Implementation Plans

### Reading Progress Tracking
- Save reading progress for each manga
- Track last read chapter and page
- Resume reading from last position

### Enhanced Download Features
- Concurrent chapter downloads for faster batch downloads
- Download queue management
- Download progress indicators
- Automatic retry on failed downloads
- Support for downloading entire manga series

### User Experience Improvements
- Bookmark system for favorite chapters
- Reading lists and collections
- Search filters (by genre, status, year, etc.)

### Technical Enhancements
- Multiple source support (not just MangaDex)
- Caching system for faster loading
- Better error handling and recovery
- Automated updates for downloaded manga
- Export/import reading progress and settings
- Backup and restore functionality

### Community Features
- User ratings and reviews
- Share reading lists
- Sync reading progress across devices
- Community recommendations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please:

1. Check the [Issues](https://github.com/taichkarna/manga-cli/issues) page
2. Create a new issue if your problem hasn't been reported

## Acknowledgments

- Thanks to all contributors who have helped shape this project
- Special thanks to the [viu](https://github.com/atanunq/viu) project for terminal image viewing capabilities
- Powered by the [MangaDex API](https://api.mangadex.org/docs.html) - a free and open manga API 