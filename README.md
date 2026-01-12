# Go Auto Clicker

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Wails](https://img.shields.io/badge/Wails-v2-DF5620?style=flat&logo=wails)](https://wails.io/)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A modern, cross-platform auto-clicker application built with Go and React. Features a clean interface for managing multiple click points with customizable delays and an emergency stop hotkey.

## Features

- **Multiple Click Points** - Set up multiple click locations that run simultaneously
- **Visual Coordinate Picker** - Click anywhere on screen to record coordinates
- **Customizable Delays** - Set individual delay intervals for each click point
- **Emergency Stop Hotkey** - Press **F10** to instantly stop clicking (only active during clicking)
- **Real-time Countdown** - Visual feedback showing time until next click
- **Clean, Modern UI** - Built with React and TypeScript

## Requirements

- Go 1.24 or higher
- Node.js 18 or higher
- Wails CLI v2
- Linux with X11 (for mouse/keyboard control)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Go-autoclicker
```

2. Install Wails CLI (if not already installed):
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

3. Install frontend dependencies:
```bash
cd frontend
npm install
cd ..
```

## Usage

### Running in Development Mode

```bash
wails dev
```

This starts the app in development mode with hot-reload for frontend changes.

### Building for Production

```bash
wails build
```

The compiled binary will be in the `build/bin` directory.

## How to Use

1. **Add Click Points**
   - Click the "+ Add Point" button to create a new click point
   - Each point starts at coordinates (0, 0)

2. **Set Coordinates**
   - Click "Set Coords" for a point
   - The app window will hide
   - Click anywhere on your screen to capture those coordinates
   - The window will reappear with the coordinates saved

3. **Adjust Delays**
   - Use the delay control to set the interval between clicks (in milliseconds)
   - Each point can have its own delay

4. **Start Clicking**
   - Click the "Start" button to begin auto-clicking
   - All points will click simultaneously at their respective intervals
   - The UI shows real-time countdowns for each point

5. **Stop Clicking**
   - Click the "Stop" button in the UI, OR
   - Press **F10** on your keyboard for emergency stop

## Tech Stack

### Backend
- **Go 1.24** - Core application logic
- **Wails v2** - Go + Web frontend framework
- **robotgo** - Mouse/keyboard control
- **gohook** - Global keyboard event listener

### Frontend
- **React 18** - UI framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Build tool and dev server

## Project Structure

```
.
├── main.go                      # Application entry point
├── app.go                       # Main app logic and Wails bindings
├── internal/
│   ├── domain/                  # Domain models
│   ├── service/
│   │   └── clicker.go          # Auto-clicker service with F10 hotkey
│   └── input/
│       └── recorder.go         # Coordinate recording service
├── frontend/
│   └── src/
│       ├── App.tsx             # Main React component
│       └── components/         # UI components
└── wails.json                  # Wails configuration
```

## Development

### Live Development with Hot Reload

```bash
wails dev
```

### Building

For your current platform:
```bash
wails build
```

For production with optimization:
```bash
wails build -clean -production
```

## Key Implementation Details

### Global Hotkey (F10)
- Only active when clicker is running
- Automatically starts with clicker, stops when clicker stops
- Prevents conflicts with coordinate recording
- Emits event to frontend for UI synchronization

### Multiple Click Points
- Each point runs in its own goroutine
- Independent timing for each point
- Thread-safe state management

### Coordinate Recording
- Uses gohook to detect mouse clicks
- Temporarily hides app window for unobstructed view
- Records exact mouse position on click

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Wails](https://wails.io/) - Go + Web framework
- [robotgo](https://github.com/go-vgo/robotgo) - Cross-platform automation
- [gohook](https://github.com/robotn/gohook) - Global event hooks