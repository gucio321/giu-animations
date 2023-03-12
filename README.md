[![Go Report Card](https://goreportcard.com/badge/github.com/gucio321/giu-animations)](https://goreportcard.com/report/github.com/gucio321/giu-animations)
[![GoDoc](https://pkg.go.dev/badge/github.com/gucio321/giu-animations?utm_source=godoc)](https://pkg.go.dev/mod/github.com/gucio321/giu-animations)

# GIU ANIMATIONS

This is a module for [giu](https://github.com/AllenDang/giu) providing an
animation system.

# Documentation

## How to use?

For complete code, please check out [examples](./_examples)

### important note

Please make shure that you're using the same version of giu
as this project (technically, you need to use giu version
that uses the same imgui-go version as yours)

### Defining an animation

At the moment, there are three implementations of animations:

- [Transition](#transition) - a smooth transition between two windows/sets of windows e.t.c.
  **NOTE** applying this animation to single widgets is not implemented yet and may
  not work as expected.
- [Color Flow](#color-flow) - you can apply this animation to any widget
  You can configure this animation to make your button hover smoother or change it into a rainbow!
- [Movement](#move) - moves DrawCursor emulating moving an object (aka `giu.Widget`).

Lets shortly discuss particular types of animations:

#### Transition

Lets look at the API:

```go
func Transition(renderers ...func(starter func())) *TransitionAnimation {...}
```

`renderers` are just [key frames](key-frame) of trasition.
In each stage appropiate renderer is called.
The argument to the renderers is a pointer to Animator.Start (see later)
so that you can call it to play the animation.

#### Color flow

```go
func ColorFlow(
        widget giu.Widget,
        applying []giu.StyleColorID,
        colors ...func() color.RGBA,
) *ColorFlowAnimation {...}
```

- The first argument is a **function** producing a **widget** that the animation should apply to
- next is the list of style-color identifiers. Color changes applies to all of them.
- and the last list of arguments are [key frames](#key-frame) of color flow.

There is also a variant of the above method called `ColorFlowStyle`, which does not need
colors list. These colors are obtained
by function like this:

```go
func() color.RGBA {
    return imgui.CurrentStyle().GetStyleColor(styleID)
}
```

#### Move

```go
func Move(w func(starter StarterFunc) giu.Widget, steps ...*MoveStep) *MoveAnimation {...}
```

This will move `w` around the steps.
Lets take a closer look on steps now:

- You create a step with `Step` or `StepVec` methods.
- You have two options of specifying position:
  - you can make it relative to the previous step. This way the system will
    take position and add it to the position of the previous step
    (and do it until it reaches first step or any step with absolute position)
  - After calling `Absolute()` method of the MoveStep, its position becomes
    absolut so that it does not relay on any previous step.
- An additional feature of Steps is a Bezier Curve implementation.
  In order to enable it, simply call `Bezier` method and specify as meny points as you wish.

One more important thing to mention about is the first step.
By default, position of the first step you specify **will be treated
absolute, even thouth it wasn't set to be.** To change this
there are two additional methods of `MoveAnimation`.

- the first one is called `StartPos` and takes one argument of the following type:
  `func(startPos imgui.Vec2) *MoveStep`. It is expected to return non-nil MoveStep.
  `startPos` argument is the position of drawing cursor at the moment **of first call** of
  `Animator`.
- another method is tu simply call `DefaultStartPos` method. It takes no arguments and acts
  like most users would like to use `StartPos` - it returns `Step(startPos)`.

### Easing

There are some extra ways of playing animation flow:

```go
const (
        EasingAlgNone EasingAlgorithmType = iota
        EasingAlgInSine
        EasingAlgOutSine
        EasingAlgInOutSine
        EasingAlgInBack
        EasingAlgOutBack
        EasingAlgInOutBack
        EasingAlgInElastic
        EasingAlgOutElastic
        EasingAlgInOutElastic
        EasingAlgInBounce
        EasingAlgOutBounce
        EasingAlgInOutBounce
)
```

for further reference, see https://easings.net

### Using animator

After constructing an animation, you need to create a special type of giu widget
called `AnimatorWidget`.

You may want to store it in a temporary variable, but, as you'll see later,
animator's api is designed so that you don't need to do so every time.

As an argument to `Animator(...)` constuctor, you pass perviously created animation.

Animator has some useful methods:

- `Duration` allows you to specify animation's duration (default is 0.25 s)
- `FPS` sets Frames per second value for animation playback (default is 60)
  **NOTE** it is not real application's FPS! It just describes how often
  animation's status is updated.
- `Start` - this method you can use to invoke animation play.
- `IsRunning` returns true, if animation is being plaid right now.

### Key Frame

Key frames are states of animation with a specified animation states.
All other states between them are calculated on-go.
Key frames system in this module is not much experienced, but it should
suit needs of most users. For more information about implementation
of this system in particular animation types, see above.

## Creating your own animation

You can use this API to create your own animation.
To do soo, lets take a look on `Animation` interface.

```go
type Animation interface {
        Init()
        Reset()
        KeyFrames() int

        BuildNormal(currentKeyFrame KeyFrame, starter func())
        BuildAnimation(animationPercentage, animationPurePercentage float32, startKeyFrame, destinationKeyFrame KeyFrame, starter func())
}
```

_This is a copy from animation.go, but I've removed comments for clearity_

### Init

init is called once, during first call of Animator.Build
you can put some initialization here.

### Reset

Reset is called along with `(*Animator).Start`

### KeyFrames

Returns a number of key frames the animation implements.
This number determines behaviour of Animator while calling Start\*

### BuildNormal

is called when `!(*Animator).IsRunning()`
It takes a pointer to `(*Animator).Start` as an argument
so you can easily start animation from there.

### BuildAnimation

is called instead of BuildNormal when playing an animation.
Along with pointer to `(*Animator).Start`, it also receives
current animation progress in percents (0 >= currentPercentage <= 1)
You can do some calculations there.

# Contribution

If you implement something interessting, find any bugs, or
improvements and would be so kind to open a PR,
your contribution is welcome!

# Motivation

For now, this system is used in one of [The Greater Heptavirate's](https://github.com/TheGraterHeptavirate) projects.
But (as I'm an author of that system) I've decided to share it for public - feel free to use if you can find any use case.

# License

This project is shared under (attached) [MIT License](LICENSE).
