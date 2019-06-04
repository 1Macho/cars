package main

import (
  "github.com/1Macho/geometry"
  "github.com/1Macho/raycasting"
  "github.com/1Macho/neuralnetworking"
  "github.com/veandco/go-sdl2/sdl"
  "sync"
)

type Car struct {
  neuralnetworking.Network
  Drivable
  Loop
  Alive bool
  Finished bool
  Stage int
  R uint8
  G uint8
  B uint8
}

func (c *Car) CalculateBoundaries (offset geometry.Point) (geometry.Point, geometry.Point, geometry.Point, geometry.Point){
  topRightPointB := geometry.Point{75, 40}
  topLeftPointB := geometry.Point{75, -40}
  bottomRightPointB := geometry.Point{-75, 40}
  bottomLeftPointB := geometry.Point{-75, -40}
  topRightPoint := geometry.Point{0,0}.Translate(topRightPointB.Magnitude(), topRightPointB.Heading().AddAngle(c.Drivable.Direction))
  topLeftPoint := geometry.Point{0,0}.Translate(topLeftPointB.Magnitude(), topLeftPointB.Heading().AddAngle(c.Drivable.Direction))
  bottomRightPoint := geometry.Point{0,0}.Translate(bottomRightPointB.Magnitude(), bottomRightPointB.Heading().AddAngle(c.Drivable.Direction))
  bottomLeftPoint := geometry.Point{0,0}.Translate(bottomLeftPointB.Magnitude(), bottomLeftPointB.Heading().AddAngle(c.Drivable.Direction))
  topRightPoint = topRightPoint.Add(offset).Add(c.Drivable.Particle.Position)
  topLeftPoint = topLeftPoint.Add(offset).Add(c.Drivable.Particle.Position)
  bottomRightPoint = bottomRightPoint.Add(offset).Add(c.Drivable.Particle.Position)
  bottomLeftPoint = bottomLeftPoint.Add(offset).Add(c.Drivable.Particle.Position)
  return topRightPoint, topLeftPoint, bottomLeftPoint, bottomRightPoint
}

func DrawLineRelative (a geometry.Point, b geometry.Point, offset geometry.Point, renderer *sdl.Renderer) {
  trueA := a.Add(offset)
  trueB := b.Add(offset)
  renderer.DrawLine(int32(trueA.X),int32(trueA.Y),int32(trueB.X),int32(trueB.Y))
}

func (c *Car) DistancesMultiCast () []float64 {
  multiCastResult := c.MultiCastFromCar()
  result := make([]float64, len(multiCastResult))
  for i := 0; i < len(multiCastResult); i++ {
    result[i] = c.Drivable.Particle.Position.Distance(multiCastResult[i]) / 1200
  }
  return result
}

func (c *Car) MultiCastFromCar () []geometry.Point {
  changePerTurn := 180.0 / 6.0
  baseAngle := -90.0
  result := make([]geometry.Point, 7)
  for i := 0; i < 7; i++ {
    result[i] = c.RayCastFromCar(baseAngle)
    baseAngle += changePerTurn
  }
  return result
}

func (c *Car) RayCastFromCar (offset float64) geometry.Point {
  newAngle := c.Drivable.Direction.Add(offset)
  newRay := geometry.Ray{geometry.Line{newAngle, c.Drivable.Particle.Position}}
  castScene := raycasting.Scene{c.Loop.Walls}
  _, result := castScene.ClosestRaycast(newRay)
  return result
}

func (c *Car) CollisionRayCast (scene raycasting.Scene, point geometry.Point) bool {
  raycastScene := scene
  ray := geometry.RayFromPoints(c.Drivable.Particle.Position, point)
  hit, result := raycastScene.ClosestRaycast(ray)
  if (hit) {
    checkDistance := c.Drivable.Particle.Position.Distance(point)
    castDistance := c.Drivable.Particle.Position.Distance(result)
    return castDistance < checkDistance
  }
  return false
}

func (c *Car) CollisionDetect (scene raycasting.Scene) bool {
  j, k, l, m := c.CalculateBoundaries(geometry.Point{0,0})
  collided := c.CollisionRayCast(scene, j)
  collided = collided || c.CollisionRayCast(scene, k)
  collided = collided || c.CollisionRayCast(scene, l)
  collided = collided || c.CollisionRayCast(scene, m)
  return collided
}

func (c *Car) Draw (renderer *sdl.Renderer, offset geometry.Point) {
  if(c.Alive) {
    c.Drivable.Draw(renderer, offset)
    renderer.SetDrawColor(c.R,c.G,c.B,255)
    j, k, l, m := c.CalculateBoundaries(offset)
    renderer.DrawLine(int32(j.X),int32(j.Y),int32(k.X),int32(k.Y))
    renderer.DrawLine(int32(k.X),int32(k.Y),int32(l.X),int32(l.Y))
    renderer.DrawLine(int32(l.X),int32(l.Y),int32(m.X),int32(m.Y))
    renderer.DrawLine(int32(j.X),int32(j.Y),int32(m.X),int32(m.Y))
    /*
    multiCastResult := c.MultiCastFromCar()
    renderer.SetDrawColor(0,20,20,255)
    for i := 0; i < len(multiCastResult); i++ {
      DrawLineRelative(c.Drivable.Particle.Position, multiCastResult[i], offset, renderer)
    }
    */
  }
}

func (c *Car) Rotate (degrees float64) {
  if (c.Alive && !c.Finished) {
    if (degrees > 5) {
      c.Drivable.Direction = c.Drivable.Direction.Add(5)
    } else if (degrees < -5) {
      c.Drivable.Direction = c.Drivable.Direction.Add(-5)
    } else {
      c.Drivable.Direction = c.Drivable.Direction.Add(degrees)
    }
  }
}

func (c *Car) Accelerate (acceleration float64) {
  if (c.Alive && !c.Finished) {
    if (acceleration > 2) {
      c.Drivable.Acceleration = 2
    } else if (acceleration < -2) {
      c.Drivable.Acceleration = -2
    } else {
      c.Drivable.Acceleration = acceleration
    }
  }
}

func (c *Car) DistanceFromLastCheckPoint () float64 {
  if (c.Stage == 0) {
    return 0
  }
  lastCheckPoint := c.Loop.CheckPoints[c.Stage - 1]
  checkPointA := lastCheckPoint.Line.Origin
  checkPointB := lastCheckPoint.EndPoint()
  checkPointMiddle := geometry.Point{
    (checkPointA.X + checkPointB.X)/2,
    (checkPointA.Y + checkPointB.Y)/2}
  return c.Drivable.Particle.Position.Distance(checkPointMiddle)
}

func (c *Car) Fitness () float64 {
  return float64(1000 * c.Stage) + c.DistanceFromLastCheckPoint()
}

func (c *Car) Tick (waitgroup *sync.WaitGroup) {
  if (c.Alive && !c.Finished) {
    multiCastResult := c.DistancesMultiCast()
    networkOutput := c.Network.CalculateOutput(multiCastResult)
    c.Rotate(networkOutput[0] * 5)
    c.Accelerate(networkOutput[1] * 2)
    c.Drivable.Tick()
    if (c.Stage < len(c.Loop.CheckPoints)) {
      nextCheckPoint := []geometry.Segment{c.Loop.CheckPoints[c.Stage]}
      collidedWithCheckPoint := c.CollisionDetect(raycasting.Scene{nextCheckPoint})
      if (collidedWithCheckPoint) {
        c.Stage += 1
      }
    }
    if (c.Stage == len(c.Loop.CheckPoints)) {
      finishLine := []geometry.Segment{c.Loop.FinishLine}
      collidedWithFinishLine := c.CollisionDetect(raycasting.Scene{finishLine})
      if (collidedWithFinishLine) {
        c.Finished = true
      }
    }
    collidedWithWalls := c.CollisionDetect(raycasting.Scene{c.Loop.Walls})
    if (collidedWithWalls) {
      c.Alive = false
    }
  }
  waitgroup.Done()
}
