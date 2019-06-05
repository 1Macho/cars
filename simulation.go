package main

import (
  "github.com/1Macho/neuralnetworking"
  "github.com/1Macho/geometry"
  "github.com/1Macho/physics"
  "github.com/veandco/go-sdl2/sdl"
  "math/rand"
  "sync"
  "sort"
)

type byFitness []Car

func (s byFitness) Len() int {
    return len(s)
}
func (s byFitness) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s byFitness) Less(i, j int) bool {
    return s[i].Fitness() > s[j].Fitness()
}

type Simulation struct {
  Loop
  Cars []Car
  Frames int
  Generation int
  MaxFitness []float64
  AvgFitness []float64
  MedianFitness []float64
  MinFitness []float64
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
  return Simulation{loop, cars, 0, 0, []float64{},[]float64{},[]float64{},[]float64{}}
}

func (s *Simulation) NextGeneration () {
  s.Loop = ObtainTrack()
  sort.Sort(byFitness(s.Cars))
  newCars := make([]Car, len(s.Cars))
  mutantsLeft := 5
  fitnessAcum := 0.0
  for i := 0; i < len(s.Cars); i++ {

    fitnessAcum += s.Cars[i].Fitness()

    thisCar := int(rand.Float64() * (float64(len(s.Cars))/30.0))
    bestCar := s.Cars[thisCar]
    baseNetwork := bestCar.Network
    thisNetwork := baseNetwork.Mutate(0.1, 0.05, 0.08, 0.005, []int{7,5,2})
    newR := bestCar.R
    newG := bestCar.G
    newB := bestCar.B
    if (mutantsLeft > 0) {
      for j := 0; j < 100; j++{
        thisNetwork = thisNetwork.Mutate(0.1, 0.05, 0.08, 0.005, []int{7,5,2})
      }
      newR = newR + uint8(rand.Float64() * 50) - 25
      newG = newG + uint8(rand.Float64() * 50) - 25
      newB = newB + uint8(rand.Float64() * 50) - 25
      mutantsLeft--
    }
    carParticle := physics.Particle{
      s.Loop.Start,
      geometry.Point{0,0},
      geometry.Point{0,0},
      1.0}
    newCar := Car{thisNetwork, Drivable{carParticle, 0, s.Loop.StartAngle}, s.Loop, true, false, 0, newR, newG, newB}
    newCars[i] = newCar
  }

  maxFitness := s.Cars[0].Fitness()
  avgFitness := fitnessAcum / float64(len(s.Cars))
  medianFitness := s.Cars[len(s.Cars) / 2].Fitness()
  minFitness := s.Cars[len(s.Cars) - 1].Fitness()
  s.MaxFitness = append(s.MaxFitness, maxFitness)
  s.AvgFitness = append(s.AvgFitness, avgFitness)
  s.MedianFitness = append(s.MedianFitness, medianFitness)
  s.MinFitness = append(s.MinFitness, minFitness)

  s.Generation += 1
  s.Cars = newCars
  s.Frames = 0
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
  s.Frames += 1
  if(allDead || s.Frames >= 3200 || oneWon) {
    s.NextGeneration()
  }
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

  graphDivider := 50.0
  graphMultiplier := 10
  baseHeight := 990

  if(len(s.MaxFitness) >= 2) {
    for index := 0; index <= 100; index++{
      i := len(s.MaxFitness) - 2 - index
      if(i > 0) {
        maxFitness := baseHeight - int(s.MaxFitness[i] / graphDivider)
        avgFitness := baseHeight - int(s.AvgFitness[i] / graphDivider)
        medianFitness := baseHeight - int(s.MedianFitness[i] / graphDivider)
        minFitness := baseHeight - int(s.MinFitness[i] / graphDivider)
        maxFitnessN := baseHeight - int(s.MaxFitness[i+1] / graphDivider)
        avgFitnessN := baseHeight - int(s.AvgFitness[i+1] / graphDivider)
        medianFitnessN := baseHeight - int(s.MedianFitness[i+1] / graphDivider)
        minFitnessN := baseHeight - int(s.MinFitness[i+1] / graphDivider)
        thisX := ((100 - index) * graphMultiplier + 10)
        nextX := (((100 - index) + 1) * graphMultiplier + 10)
        renderer.SetDrawColor(0,255,0,255)
        renderer.DrawLine(int32(thisX),int32(maxFitness),int32(nextX),int32(maxFitnessN))
        renderer.SetDrawColor(0,0,255,255)
        renderer.DrawLine(int32(thisX),int32(avgFitness),int32(nextX),int32(avgFitnessN))
        renderer.SetDrawColor(0,255,255,255)
        renderer.DrawLine(int32(thisX),int32(medianFitness),int32(nextX),int32(medianFitnessN))
        renderer.SetDrawColor(255,0,0,255)
        renderer.DrawLine(int32(thisX),int32(minFitness),int32(nextX),int32(minFitnessN))
      } else { break; }
    }
  }

}
