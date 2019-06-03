package main

import (
  "github.com/1Macho/geometry"
  "github.com/1Macho/raycasting"
  "github.com/veandco/go-sdl2/sdl"
)

type Car struct {
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
  c.Drivable.Draw(renderer, offset)
  renderer.SetDrawColor(c.R,c.G,c.B,255)
  j, k, l, m := c.CalculateBoundaries(offset)
  renderer.DrawLine(int32(j.X),int32(j.Y),int32(k.X),int32(k.Y))
  renderer.DrawLine(int32(k.X),int32(k.Y),int32(l.X),int32(l.Y))
  renderer.DrawLine(int32(l.X),int32(l.Y),int32(m.X),int32(m.Y))
  renderer.DrawLine(int32(j.X),int32(j.Y),int32(m.X),int32(m.Y))
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
    if (acceleration > 1) {
      c.Drivable.Acceleration = -1
    } else if (acceleration < -1) {
      c.Drivable.Acceleration = -1
    } else {
      c.Drivable.Acceleration = acceleration
    }
  }
}

func (c *Car) Tick () {
  if (c.Alive && !c.Finished) {
    c.Drivable.Tick()
    if (c.Stage < len(c.Loop.CheckPoints)) {
      nextCheckPoint := []geometry.Segment{c.Loop.CheckPoints[c.Stage]}
      collidedWithCheckPoint := c.CollisionDetect(raycasting.Scene{nextCheckPoint})
      if (collidedWithCheckPoint) {
        c.Stage += 1
        println(c.Stage)
      }
    }
    if (c.Stage == len(c.Loop.CheckPoints)) {
      finishLine := []geometry.Segment{c.Loop.FinishLine}
      collidedWithFinishLine := c.CollisionDetect(raycasting.Scene{finishLine})
      if (collidedWithFinishLine) {
        c.Finished = true
        c.R = 255
        c.G = 255
        c.B = 255
      }
    }
    collidedWithWalls := c.CollisionDetect(raycasting.Scene{c.Loop.Walls})
    if (collidedWithWalls) {
      c.Alive = false
      c.R = 100
      c.G = 100
      c.B = 100
    }
  }
}
