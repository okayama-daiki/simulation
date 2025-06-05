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

	D  = 0.25 // Diffusion coefficient
	dt = 1    // Time step
	dx = 1    // Width of a cell
	dy = 1    // Height of a cell

	screenWidth  = M * cellSize
	screenHeight = N * cellSize

	fps = 1 // Frames per second
)

type Simulation struct {
	frameCount int
	rho        [][]float32
}

func New() Simulation {
	var rho [][]float32
	for range N {
		rho = append(rho, make([]float32, M))
	}

	for i := N/2 - 5; i < N/2+5; i++ {
		for j := M/2 - 5; j < M/2+5; j++ {
			rho[i][j] = 1.0 // Initial condition: set a blob in the center
		}
	}

	return Simulation{
		frameCount: 0,
		rho:        rho,
	}
}

func (s *Simulation) Update() error {
	s.frameCount++

	if s.frameCount%fps != 0 {
		return nil
	}

	// Update the grid based on following formula:
	// rho^{n+1}_{i,j} = rho^{n}_{i,j} + D * dt / dx^2 * (rho^{n}_{i+1,j} + rho^{n}_{i-1,j} + rho^{n}_{i,j+1} + rho^{n}_{i,j-1} - 4 * rho^{n}_{i,j})

	rhoNext := make([][]float32, N)
	for i := range N {
		rhoNext[i] = make([]float32, M)
	}

	for i := 1; i < N-1; i++ {
		for j := 1; j < M-1; j++ {
			rhoNext[i][j] = s.rho[i][j] + (D*dt/(dx*dx))*(s.rho[i+1][j]+s.rho[i-1][j]+s.rho[i][j+1]+s.rho[i][j-1]-4*s.rho[i][j])
		}
	}

	// Boundary conditions
	for i := range N {
		rhoNext[i][0] = rhoNext[i][1]     // Left boundary
		rhoNext[i][M-1] = rhoNext[i][M-2] // Right boundary
	}
	for j := range M {
		rhoNext[0][j] = rhoNext[1][j]     // Top boundary
		rhoNext[N-1][j] = rhoNext[N-2][j] // Bottom boundary
	}

	for i := range N {
		for j := range M {
			s.rho[i][j] = rhoNext[i][j]
		}
	}

	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	for y := range N {
		for x := range M {
			geom := ebiten.GeoM{}
			geom.Translate(float64(x*cellSize), float64(y*cellSize))

			colorValue := min(uint8(255*s.rho[y][x]), 255)
			color := color.RGBA{colorValue, 0, 0, 255} // Red channel based on rho value
			vector.DrawFilledRect(screen, float32(x*cellSize), float32(y*cellSize), float32(cellSize), float32(cellSize), color, false)
		}
	}
}

func (s *Simulation) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("2D Diffusion Simulation")

	simulation := New()
	if err := ebiten.RunGame(&simulation); err != nil {
		panic(err)
	}
}
