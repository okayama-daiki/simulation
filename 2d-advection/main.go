package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	M        = 100
	N        = 100
	cellSize = 2

	U_0 = 0.1 // Advection speed in x-direction
	V_0 = 0.1 // Advection speed in y-direction
	dt  = 1.0 // Time step
	dx  = 1.0 // Width of a cell
	dy  = 1.0 // Height of a cell

	screenWidth  = M * cellSize
	screenHeight = N * cellSize

	fps = 1
)

type Simulation struct {
	frameCount int
	phi        [][]float32
}

func New() Simulation {
	phi := make([][]float32, N)
	for i := range phi {
		phi[i] = make([]float32, M)
	}

	for i := N/2 - 7; i < N/2+7; i++ {
		for j := M/2 - 7; j < M/2+7; j++ {
			phi[i][j] = 1.0
		}
	}

	return Simulation{
		frameCount: 0,
		phi:        phi,
	}
}

func (s *Simulation) Update() error {
	s.frameCount++

	if s.frameCount%fps != 0 {
		return nil
	}

	// Update the grid based on the advection formula:
	// phi^{n+1}_{i,j} = phi^{n}_{i,j} - U_0 * dt / dx * (phi^{n}_{i,j} - phi^{n}_{i-1,j}) - V_0 * dt / dy * (phi^{n}_{i,j} - phi^{n}_{i,j-1})

	phiNext := make([][]float32, N)
	for i := range phiNext {
		phiNext[i] = make([]float32, M)
	}

	for i := 1; i < N; i++ {
		for j := 1; j < M; j++ {
			du := float32(U_0 * dt / dx)
			dv := float32(V_0 * dt / dy)

			phiNext[i][j] = s.phi[i][j] -
				du*(s.phi[i][j]-s.phi[i-1][j]) -
				dv*(s.phi[i][j]-s.phi[i][j-1])
		}
	}

	// Boundary condition: fixed to zero

	for i := range N {
		phiNext[i][0] = 0   // Left boundary
		phiNext[i][M-1] = 0 // Right boundary
	}

	for j := range M {
		phiNext[0][j] = 0   // Top boundary
		phiNext[N-1][j] = 0 // Bottom boundary
	}

	for i := range s.phi {
		copy(s.phi[i], phiNext[i])
	}

	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	for y := range N {
		for x := range M {
			value := s.phi[y][x]
			colorValue := min(uint8(255*value), 255) // Red channel based on phi value
			vector.DrawFilledRect(screen, float32(x*cellSize), float32(y*cellSize), float32(cellSize), float32(cellSize), color.RGBA{colorValue, 0, 0, 255}, false)
		}
	}
}

func (s *Simulation) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("2D Advection Simulation")

	simulation := New()
	if err := ebiten.RunGame(&simulation); err != nil {
		panic(err)
	}
}
