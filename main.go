package main

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	width  = 80
	height = 3 // Single line height like the block

	// Physics constants
	acceleration        = 1    // how fast you speed up(set to 1.5 for max acceleration)
	deceleration        = 0.75 // How fast you slow down (0.71 = lose 29% speed each tick)
	maxVelocity         = 1    // Maximum movement speed
	standstillThreshold = 0.15 // Below this speed = "not moving"
	counterStrafeWindow = 100  //100 milliseconds to shoot after counter strafing
)

type model struct {
	//this is the positin and movement
	crosshairX float64
	crosshairY float64
	velocityX  float64

	//this is the hit block position
	targetX int
	targetY int

	//this is the input vals
	movingLeft  bool //is A being held?
	movingRight bool //is D being held?

	// Key state tracking
	leftPressed   bool      // Is A currently pressed?
	rightPressed  bool      // Is D currently pressed?
	lastLeftTime  time.Time // When was A last pressed?
	lastRightTime time.Time // When was D last pressed?

	// Counter-strafe timing
	counterStrafeTime time.Time // When did counter-strafe start?
	inCounterStrafe   bool      //are we in counter strafe window?

	score int
}

type tickMsg time.Time //a custom message type

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*35, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	rand.Seed(time.Now().UnixNano())

	return model{
		crosshairX: float64(width / 2),
		crosshairY: 1, // Fixed to center line
		velocityX:  0,
		targetX:    rand.Intn(width-4) + 2,
		targetY:    1, // Fixed to center line
		score:      0,
	}
}

// initial function
func (m model) Init() tea.Cmd {
	return tick()
}

// the update function i.e the heart of this app
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// Handle key presses - update timestamp on each press
		switch msg.String() {
		case "a":
			m.lastLeftTime = time.Now()
			if !m.leftPressed {
				m.leftPressed = true
				m.movingLeft = true
				// Counter-strafe: if moving right, enter counter-strafe state
				if m.velocityX > standstillThreshold {
					m.velocityX = 0
					m.inCounterStrafe = true
					m.counterStrafeTime = time.Now()
				}
			}
		case "d":
			m.lastRightTime = time.Now()
			if !m.rightPressed {
				m.rightPressed = true
				m.movingRight = true
				// Counter-strafe: if moving left, enter counter-strafe state
				if m.velocityX < -standstillThreshold {
					m.velocityX = 0
					m.inCounterStrafe = true
					m.counterStrafeTime = time.Now()
				}
			}
		}

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// Check if crosshair is on target
			crossX := int(m.crosshairX)
			crossY := int(m.crosshairY)
			isOnTarget := crossX >= m.targetX && crossX <= m.targetX+1 &&
				crossY == m.targetY

			canShoot := m.inCounterStrafe || abs(m.velocityX) < standstillThreshold

			// Only reset target position if we hit it AND can shoot
			if isOnTarget && canShoot {
				m.score++
				// Spawn new target only when we successfully hit the current one
				m.targetX = rand.Intn(width-4) + 2
				m.targetY = 1 // Keep on center line
			}
		}

	case tickMsg:
		// Update physics
		m = m.updatePhysics()
		return m, tick()
	}

	return m, nil
}

func (m model) updatePhysics() model {
	// Check for key releases based on time since last press (simulate key release)
	keyReleaseThreshold := 80 * time.Millisecond

	if m.leftPressed && time.Since(m.lastLeftTime) > keyReleaseThreshold {
		m.leftPressed = false
		m.movingLeft = false
	}

	if m.rightPressed && time.Since(m.lastRightTime) > keyReleaseThreshold {
		m.rightPressed = false
		m.movingRight = false
	}

	// Check if counter-strafe window has expired
	if m.inCounterStrafe {
		if time.Since(m.counterStrafeTime) > time.Duration(counterStrafeWindow)*time.Millisecond {
			m.inCounterStrafe = false
		}
	}

	// Apply input acceleration (can move even during counter-strafe window end)
	if m.movingLeft && !m.movingRight {
		if !m.inCounterStrafe {
			m.velocityX -= acceleration // Speed up leftward
		}
	} else if m.movingRight && !m.movingLeft {
		if !m.inCounterStrafe {
			m.velocityX += acceleration // Speed up rightward
		}
	} else {
		// Apply deceleration when no input
		if !m.inCounterStrafe && abs(m.velocityX) > 0.1 {
			m.velocityX *= deceleration
		} else if !m.inCounterStrafe {
			m.velocityX = 0
		}
	}

	// Clamp velocity
	if m.velocityX > maxVelocity {
		m.velocityX = maxVelocity
	} else if m.velocityX < -maxVelocity {
		m.velocityX = -maxVelocity
	}

	// Update position (always update position, even during counter-strafe)
	if !m.inCounterStrafe {
		m.crosshairX += m.velocityX
	}

	// Keep crosshair in bounds - FIXED to prevent overwriting right border
	if m.crosshairX < 1 {
		m.crosshairX = 1
		m.velocityX = 0
	} else if m.crosshairX >= float64(width-1) {
		m.crosshairX = float64(width - 2)
		m.velocityX = 0
	}

	return m
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func (m model) View() string {
	// Create the game field
	field := make([][]rune, height)
	for i := range field {
		field[i] = make([]rune, width)
		for j := range field[i] {
			if i == 0 || i == height-1 || j == 0 || j == width-1 {
				field[i][j] = '█'
			} else {
				field[i][j] = ' '
			}
		}
	}

	// Place crosshair
	crossX := int(m.crosshairX)
	crossY := int(m.crosshairY)
	if crossX >= 0 && crossX < width && crossY >= 0 && crossY < height {
		field[crossY][crossX] = '+'
	}

	// Determine if we can shoot and are on target
	targetOnCrosshair := crossX >= m.targetX && crossX <= m.targetX+1 &&
		crossY == m.targetY
	canShoot := m.inCounterStrafe || abs(m.velocityX) < standstillThreshold
	frameIsGreen := targetOnCrosshair && canShoot

	// Place target (2x2 block) - but only on center line
	for dx := 0; dx < 2; dx++ {
		if m.targetY >= 0 && m.targetY < height &&
			m.targetX+dx >= 0 && m.targetX+dx < width {
			field[m.targetY][m.targetX+dx] = '█'
		}
	}

	// Convert field to string
	var result string
	for i, row := range field {
		for j, cell := range row {
			if cell == '█' && (i == 0 || i == height-1 || j == 0 || j == width-1) {
				// Color the frame based on shoot status
				if frameIsGreen {
					result += lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render(string(cell))
				} else {
					result += lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(string(cell))
				}
			} else if cell == '█' {
				// Target blocks - keep them uncolored/default
				result += string(cell)
			} else if cell == '+' {
				// Color the crosshair
				result += lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(string(cell))
			} else {
				result += string(cell)
			}
		}
		result += "\n"
	}

	// Add UI info
	result += "\n\nYujon's Counter Strafe Trainer:"

	result += fmt.Sprintf("\nScore: %d", m.score)
	result += fmt.Sprintf("\nPosition: %.1f", m.crosshairX)
	result += fmt.Sprintf(" | Velocity: %.2f", m.velocityX)
	if m.inCounterStrafe {
		timeLeft := counterStrafeWindow - int(time.Since(m.counterStrafeTime).Milliseconds())
		result += fmt.Sprintf(" | COUNTER-STRAFE WINDOW: %dms", timeLeft)
	} else if abs(m.velocityX) < standstillThreshold {
		result += " | Status: STANDSTILL"
	} else {
		result += " | Status: MOVING"
	}
	result += "\n\nControls:"
	result += "\n  A/D - Move left/right"
	result += "\n  Left Click - Shoot (score only when crosshair on target)"
	result += "\n  Q - Quit"
	result += "\n\nTip: Counter-strafe (press opposite direction)"
	result += "\nDrifting happens when you a noob and cant counter strafe!"
	result += "\nMake sure you're holding 'a' or 'd' at all times and only shooting when its green!"
	return result
}

func isTargetPosition(charIndex, rowWidth, targetX, targetY, fieldWidth int) bool {
	row := charIndex / (rowWidth + 1) // +1 for newline
	col := charIndex % (rowWidth + 1)

	return row == targetY && col >= targetX && col < targetX+2
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
