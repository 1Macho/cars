package main

import (
  "github.com/1Macho/physics"
  "github.com/1Macho/geometry"
  "github.com/veandco/go-sdl2/sdl"
)

type Drivable struct {
  physics.Particle
  Acceleration float64
  Direction geometry.Angle
}

func (d Drivable) Draw (renderer *sdl.Renderer, offset geometry.Point) {
  rightAngle := d.Direction.Add(90 * 0)
  topAngle := d.Direction.Add(90 * 1)
  leftAngle := d.Direction.Add(90 * 2)
  bottomAngle := d.Direction.Add(90 * 3)
  rightPoint := d.Particle.Position.Translate(8, rightAngle)
  topPoint := d.Particle.Position.Translate(8, topAngle)
  leftPoint := d.Particle.Position.Translate(8, leftAngle)
  bottomPoint := d.Particle.Position.Translate(8, bottomAngle)
  renderer.SetDrawColor(255,255,255,255)
  renderer.DrawLine(int32(rightPoint.X + offset.X), int32(rightPoint.Y + offset.Y), int32(leftPoint.X + offset.X), int32(leftPoint.Y + offset.Y))
  renderer.DrawLine(int32(topPoint.X + offset.X), int32(topPoint.Y + offset.Y), int32(bottomPoint.X + offset.X), int32(bottomPoint.Y + offset.Y))
}

func (d *Drivable) Tick () {
  frictionForce := geometry.Point{d.Particle.Velocity.X,d.Particle.Velocity.Y}
  frictionForce = frictionForce.Inverse()
  frictionForce = frictionForce.Multiply(0.0175)
  d.Particle.ApplyForce(frictionForce)
  accelerationForce := geometry.Point{0,0}
  accelerationForce = accelerationForce.Translate(d.Acceleration / 100, d.Direction)
  d.Particle.ApplyForce(accelerationForce)
  d.Particle.Tick()
  //println(d.Particle.Position.X)
}
