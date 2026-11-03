package lab

import (
	"fmt"
	"unicode/utf8"
)

type ConversionError struct {
	Code   string
	Offset int
}

func (e *ConversionError) Error() string { return fmt.Sprintf("%s at byte %d", e.Code, e.Offset) }

type Warning struct {
	Offset int    `json:"offset"`
	Rune   string `json:"rune"`
}
type Conversion struct {
	Status   string    `json:"status"`
	Output   string    `json:"output"`
	Warnings []Warning `json:"warnings"`
}

func Convert(s, policy string, max int) (Conversion, error) {
	fail := func(code string, off int) (Conversion, error) {
		return Conversion{Status: "rejected"}, &ConversionError{code, off}
	}
	if policy != "strict" && policy != "replacement" {
		return fail("invalid_policy", 0)
	}
	if max < 0 {
		return fail("invalid_limit", 0)
	}
	if !utf8.ValidString(s) {
		return fail("invalid_utf8", 0)
	}
	r := Conversion{Status: "unchanged", Warnings: []Warning{}}
	out := make([]byte, 0)
	for off, c := range s {
		if c < 32 || c > 126 {
			if policy == "strict" {
				return fail("unsupported_character", off)
			}
			r.Warnings = append(r.Warnings, Warning{off, string(c)})
			c = '?'
			r.Status = "changed"
		}
		if len(out) >= max {
			return fail("output_overflow", off)
		}
		out = append(out, byte(c))
	}
	r.Output = string(out)
	return r, nil
}
