# satellite
    `go get github.com/jamiecrisman/go-satellite`

## Intro

[![Go](https://github.com/jamiecrisman/go-satellite/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/jamiecrisman/go-satellite/actions/workflows/go.yml) [![GoDoc](https://godoc.org/github.com/jamiecrisman/go-satellite?status.svg)](https://godoc.org/github.com/jamiecrisman/go-satellite)


`go-satellite` lets you take a TLE and propagate it to a given time utilizing sgp4. It also lets you calculate the look angles to a satellite from a given location.

You should utilize [the original repo](https://github.com/joshuaferrara/go-satellite) over this one.

This fork diverges to narrow the focus of the repo/package.

- removes panics/fatal calls and replaces with errors
  - adds sentinel errors for propagation rather than error codes
- removes dependencies
- removes spacetrak/celestrak TLE fetching features, those can be implemented elsewhere
- reorganizes some code
  - APIs were ambiguous with similar signatures, but different functionality. This was cleaned up to have clearer separation of TLE and Satellite.
- pulls in some community PR fixes that weren't merged into the original repo
- adding some cmd utilities
