package main

import (
  "github.com/1Macho/neuralnetworking"
  "github.com/1Macho/geometry"
  "github.com/1Macho/physics"
  "github.com/veandco/go-sdl2/sdl"
  "math/rand"
)

type Simulation struct {
  Loop
  Cars []Car
}

func CreateSimulation (loop Loop, SampleSize int) Simulation {
  cars := make([]Car, SampleSize)
  for i := 0; i < SampleSize; i++ {
    carNetwork := neuralnetworking.BuildRandomNetwork([]int{7,5,4,2})
    carParticle := physics.Particle{
      loop.Start,
      geometry.Point{0,0},
      geometry.Point{0,0},
      1.0}
    r := 150 + uint8(rand.Float64() * 105)
    g := 150 + uint8(rand.Float64() * 105)
    b := 150 + uint8(rand.Float64() * 105)
    newCar := Car{carNetwork, Drivable{carParticle, 0, geometry.Angle{90}}, loop, true, false, 0, r, g, b}
    cars[i] = newCar
  }
  return Simulation{loop, cars}
}

func (s *Simulation) Tick () {
  for i := 0; i < len(s.Cars); i++ {
    s.Cars[i].Tick()
  }
}

func (s *Simulation) Draw (renderer *sdl.Renderer) {
  offset := s.Cars[0].Drivable.Particle.Position.Inverse().Add(geometry.Point{600, 500})
  renderer.SetDrawColor(0,0,0,0)
  renderer.Clear()
  s.Loop.Draw(renderer, offset)
  for i := 0; i < len(s.Cars); i++ {
    s.Cars[i].Draw(renderer, offset)
  }
}
