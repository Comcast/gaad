package bitreader

import "errors"

var (
	// ErrExpGolombDecode represents an error in decoding a number represented in exponential-golomb coding
	ErrExpGolombDecode = errors.New("could not complete Exp-Golomb decode")
	// ErrReaderOutOfBounds is returned if an attempt is made to read from a parser that has been depleted
	ErrReaderOutOfBounds = errors.New("cannot read past end of reader bytes")
)
