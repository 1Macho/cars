package cars

import (
  "github.com/1Macho/geometry"
  "github.com/veandco/go-sdl2/sdl"
)

type Car struct {
  Drivable
  R uint8
  G uint8
  B uint8
}

func (c *Car) Draw (renderer *sdl.Renderer, offset geometry.Point) {
  c.Drivable.Draw(renderer, offset)
  renderer.SetDrawColor(c.R,c.G,c.B,255)
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
  renderer.DrawLine(int32(topRightPoint.X),int32(topRightPoint.Y),int32(topLeftPoint.X),int32(topLeftPoint.Y))
  renderer.DrawLine(int32(topLeftPoint.X),int32(topLeftPoint.Y),int32(bottomLeftPoint.X),int32(bottomLeftPoint.Y))
  renderer.DrawLine(int32(bottomLeftPoint.X),int32(bottomLeftPoint.Y),int32(bottomRightPoint.X),int32(bottomRightPoint.Y))
  renderer.DrawLine(int32(topRightPoint.X),int32(topRightPoint.Y),int32(bottomRightPoint.X),int32(bottomRightPoint.Y))
}
