package main

import (
  "github.com/1Macho/neuralnetworking"
  "github.com/1Macho/geometry"
  "github.com/1Macho/physics"
  "github.com/veandco/go-sdl2/sdl"
  "math/rand"
  "sync"
)

type Simulation struct {
  Loop
  Cars []Car
  Frames int
  Generation int
}

func ObtainTrack () Loop {
  return BuildSimplexTrack(4000, 900, 1200, 100, 25)
}

func CreateSimulation (SampleSize int) Simulation {
  loop := ObtainTrack()
  cars := make([]Car, SampleSize)
  for i := 0; i < SampleSize; i++ {
    carNetwork := neuralnetworking.BuildRandomNetwork([]int{7,5,2})
    carParticle := physics.Particle{
      loop.Start,
      geometry.Point{0,0},
      geometry.Point{0,0},
      1.0}
    r := 25 + uint8(rand.Float64() * 230)
    g := 25 + uint8(rand.Float64() * 230)
    b := 25 + uint8(rand.Float64() * 230)
    newCar := Car{carNetwork, Drivable{carParticle, 0, loop.StartAngle}, loop, true, false, 0, r, g, b}
    cars[i] = newCar
  }
  return Simulation{loop, cars, 0, 0}
}

func (s *Simulation) Tick () {
  var waitgroup sync.WaitGroup
  for i := 0; i < len(s.Cars); i++ {
    waitgroup.Add(1)
    go s.Cars[i].Tick(&waitgroup)
  }
  waitgroup.Wait()
  allDead := true
  oneWon := false
  for i := 0; i < len(s.Cars); i++ {
    oneWon = oneWon || s.Cars[i].Finished
    allDead = allDead && !s.Cars[i].Alive
  }
  if(allDead || s.Frames >= 1600 || oneWon) {
    s.Loop = ObtainTrack()
    bestCar := s.Cars[0]
    for i := 0; i < len(s.Cars); i++ {
      if (s.Cars[i].Finished) {
        bestCar = s.Cars[i]
        break
      }
      if (s.Cars[i].Fitness() > bestCar.Fitness()) {
        bestCar = s.Cars[i]
      }
    }
    newCars := make([]Car, len(s.Cars))
    baseNetwork := bestCar.Network
    for i := 0; i < len(s.Cars); i++ {
      thisNetwork := baseNetwork.Mutate(0.1, 0.05, 0.08, 0.005, []int{7,5,2})
      carParticle := physics.Particle{
        s.Loop.Start,
        geometry.Point{0,0},
        geometry.Point{0,0},
        1.0}
      newCar := Car{thisNetwork, Drivable{carParticle, 0, s.Loop.StartAngle}, s.Loop, true, false, 0, bestCar.R, bestCar.G, bestCar.B}
      newCars[i] = newCar
    }
    s.Generation += 1
    s.Cars = newCars
    s.Frames = 0
  }
  s.Frames += 1
}

func (s *Simulation) Draw (renderer *sdl.Renderer) {
  count := 0
  xOffset := 0.0
  yOffset := 0.0
  for i := 0; i < len(s.Cars); i++ {
    if (s.Cars[i].Alive) {
      count++
      xOffset += s.Cars[i].Drivable.Particle.Position.X
      yOffset += s.Cars[i].Drivable.Particle.Position.Y
    }
  }
  xOffset = xOffset / float64(count)
  yOffset = yOffset / float64(count)
  offset := geometry.Point{xOffset,yOffset}.Inverse().Add(geometry.Point{600, 500})
  renderer.SetDrawColor(0,0,0,0)
  renderer.Clear()
  s.Loop.Draw(renderer, offset)
  for i := 0; i < len(s.Cars); i++ {
    s.Cars[i].Draw(renderer, offset)
  }
}
