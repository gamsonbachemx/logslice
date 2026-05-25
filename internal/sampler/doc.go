// Package sampler provides probabilistic sampling for log streams.
//
// Use New to create a Sampler with a given rate and random seed.
// Call Keep on each incoming line to decide whether it should be
// forwarded to output. A rate of 1.0 disables sampling entirely.
package sampler
