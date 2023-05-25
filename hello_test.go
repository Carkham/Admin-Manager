package main

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"testing"
)

func randomFormat() string {
	// A slice of message formats.
	formats := []string{
		"Hi, %v. Welcome!",
		"Great to see you, %v!",
		"Hail, %v! Well met!",
	}

	// Return one of the message formats selected at random.
	return formats[rand.Intn(len(formats))]
}

func Hello(name string) (string, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return name, errors.New("empty name")
	}
	// Create a message using a random format.
	message := fmt.Sprintf(randomFormat(), name)
	return message, nil
}

func TestHelloName(t *testing.T) {
	name := "Dkhtn"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello("Dkhtn")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("Dkhtn") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
