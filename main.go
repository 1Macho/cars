package main

import (
  "github.com/veandco/go-sdl2/sdl"
  "github.com/1Macho/geometry"
  "github.com/1Macho/physics"
  "github.com/1Macho/raycasting"
  "time"
)

func main() {
  if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
    panic(err)
  }
  defer sdl.Quit()

  window, err := sdl.CreateWindow("Cars", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		1200,1000, sdl.WINDOW_SHOWN)
  if err != nil {
    panic(err)
  }
  defer window.Destroy()

  renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.Clear()

  //testAngle := geometry.Angle{90}
  //println(testAngle.Slope())
  //testLine := geometry.Line{geometry.Angle{90+45},geometry.Point{0,1}}
  //println(testLine.TrueOrigin().X)
  //println(testLine.TrueOrigin().Y)
  //testRay := geometry.Ray{geometry.Line{geometry.Angle{47},geometry.Point{0,0}}}
  //println(geometry.InterceptRayLine(testRay,testLine).X)
  //println(geometry.InterceptRayLine(testRay,testLine).Y)


  testLoop := BuildCircularTrack(1000, 1400, 16)

  testScene := raycasting.Scene{testLoop.Walls}

  testParticle := physics.Particle{
    testLoop.Start,
    geometry.Point{0,0},
    geometry.Point{0,0},
    1.0}

  testCar := Car{Drivable{testParticle, 0, geometry.Angle{90}}, 255, 255, 0}



  running := true
  for running {
    offset := testCar.Drivable.Particle.Position.Inverse().Add(geometry.Point{600, 500})
    renderer.SetDrawColor(0,0,0,0)
    renderer.Clear()
    renderer.SetDrawColor(200,200,0,255)
    testLoop.Draw(renderer, offset)
    testCar.Draw(renderer, offset)
    testRay := geometry.Ray{
      geometry.Line{
        testCar.Drivable.Direction,
        testCar.Drivable.Particle.Position}}

    renderer.SetDrawColor(0,190,190,255)
    hit, cast := testScene.ClosestRaycast(testRay)
    cast = cast.Add(offset)
    if(hit) {
      renderer.DrawLine(int32(testCar.Drivable.Particle.Position.X + offset.X), int32(testCar.Drivable.Particle.Position.Y + offset.Y), int32(cast.X), int32(cast.Y))
    }
    testRay.Line.Direction = testRay.Line.Direction.Add(15)
    hit, cast = testScene.ClosestRaycast(testRay)
    cast = cast.Add(offset)
    if(hit) {
      renderer.DrawLine(int32(testCar.Drivable.Particle.Position.X + offset.X), int32(testCar.Drivable.Particle.Position.Y + offset.Y), int32(cast.X), int32(cast.Y))
    }
    testRay.Line.Direction = testRay.Line.Direction.Add(-30)
    hit, cast = testScene.ClosestRaycast(testRay)
    cast = cast.Add(offset)
    if(hit) {
      renderer.DrawLine(int32(testCar.Drivable.Particle.Position.X + offset.X), int32(testCar.Drivable.Particle.Position.Y + offset.Y), int32(cast.X), int32(cast.Y))
    }
    testCar.Drivable.Tick()



//    println(testCar.Drivable.Direction.Value)
    //testDrivable.Draw(renderer, geometry.Point{0,0})
    renderer.Present()
    //testParticle.Tick()
    testCar.Drivable.Acceleration = 0
    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
      switch t := event.(type) {
      case *sdl.KeyboardEvent:
        if (sdl.K_UP == t.Keysym.Sym) {
          testCar.Drivable.Acceleration = 1
        }
        if (sdl.K_DOWN == t.Keysym.Sym) {
          testCar.Drivable.Acceleration = -1
        }
        if (sdl.K_LEFT == t.Keysym.Sym) {
          testCar.Drivable.Direction = testCar.Drivable.Direction.Add(-5)
        }
        if (sdl.K_RIGHT == t.Keysym.Sym) {
          testCar.Drivable.Direction = testCar.Drivable.Direction.Add(5)
        }
        break
      case *sdl.QuitEvent:
        println("Quit")
        running = false
        break
      }
    }
    time.Sleep(1000000000/60)
  }
}
