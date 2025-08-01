# Counter Strafe Trainer

A terminal-based counter-strafe mechanics trainer built with Go and Bubble Tea, designed to help improve your movement skills for tactical FPS games like Counter-Strike 2 and Valorant.

## 🎯 What is Counter-Strafing?

Counter-strafing is a fundamental movement technique in tactical FPS games where you quickly tap the opposite movement key to instantly stop your momentum, allowing for precise shooting. This trainer simulates that mechanic with a 300ms accuracy window.

## 🚀 Features

- **Realistic Physics**: Simulates CS2/Valorant movement with acceleration, deceleration, and momentum
- **Counter-Strafe Mechanics**: Get a 100ms shooting window when you counter-strafe
- **Visual Feedback**: Targets turn green when you can accurately shoot
- **Score Tracking**: Track your successful hits
- **Terminal-Based**: Runs entirely in your terminal with smooth animations

## 📋 Prerequisites

- Go 1.19 or higher
- Terminal that supports mouse input and colors

## 🛠️ Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/counter-strafe-trainer.git
cd counter-strafe-trainer
```

2. Initialize Go module and install dependencies:
```bash
go mod init counter-strafe-trainer
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
```

3. Run the trainer:
```bash
go run main.go
```

## 🎮 How to Play

### Controls
- **A Key**: Move left
- **D Key**: Move right  
- **Left Click**: Shoot (only when target is green)
- **Q**: Quit

### Gameplay
1. A red target block spawns randomly on the horizontal line
2. Use A/D keys to move your crosshair (`+`) toward the target
3. When your crosshair is on the target, it turns **green** only when you can shoot accurately:
   - During the 300ms counter-strafe window (when you press opposite direction)
   - When you're naturally at standstill
4. Left-click to shoot and score points
5. A new target spawns after each successful hit

### Counter-Strafe Example
```
Target spawned to the right →
Press D to move right →
When crosshair reaches target, press A →
100ms window opens (target turns green) →
Left-click to score →
After 100ms, you start moving left
```

## 🎯 Training Tips

1. **Practice the Timing**: The 100ms window simulates real CS2/Valorant timing
2. **Counter-Strafe Everything**: Always counter-strafe when approaching targets
3. **Don't Spray and Pray**: Only shoot when the target is green
4. **Build Muscle Memory**: Focus on smooth A→D→A or D→A→D sequences

## ⚙️ Configuration

You can modify the physics constants in `main.go`:

```go
const (
    acceleration = 0.4          // How fast you accelerate
    deceleration = 0.94         // How fast you slow down
    maxVelocity = 1.2          // Maximum movement speed
    counterStrafeWindow = 300   // Shooting window in milliseconds
)
```

## 🎨 Interface

```
████████████████████████████████████████████████████████████████████████████████
█                                  +     ██                                    █
████████████████████████████████████████████████████████████████████████████████

Yujon's Counter Strafe Trainer:
Score: 5
Position: 42.3 | Velocity: 0.00 | Status: STANDSTILL

Controls:
  A/D - Move left/right
  Left Click - Shoot (score only when crosshair on target)
  Q - Quit

Tip: Counter-strafe (press opposite direction)
Drifting happens when you a noob and cant counter strafe!
Make sure you're holding 'a' or 'd' at all times and only shooting when its green!
```

## 🔧 Technical Details

- **Physics Engine**: Custom movement physics with realistic acceleration/deceleration
- **Counter-Strafe Detection**: Automatically detects opposite key presses and creates shooting windows
- **Real-time Updates**: 35ms tick rate for smooth gameplay
- **Color System**: Red targets (can't shoot) vs Green targets (can shoot)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- Inspired by Counter-Strike 2 and Valorant movement mechanics

## 📈 Roadmap

- [ ] Add difficulty levels (faster targets, smaller windows)
- [ ] Statistics tracking (accuracy percentage, reaction times)
- [ ] Different target patterns and spawning algorithms
- [ ] Sound effects for successful hits
- [ ] Leaderboard system
- [ ] Training scenarios (pre-fire positions, common angles)

---

**Happy Counter-Strafing!** 🎯

*Improve your aim, perfect your timing, dominate the server.*
