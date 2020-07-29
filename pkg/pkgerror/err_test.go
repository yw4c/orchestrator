// +build framework

package pkgerror

import (
	"github.com/rotisserie/eris"
	"log"
	"testing"
)


func TestErr(t *testing.T) {

	// you can wrap multiple err before response.
	err := eris.Wrap(ErrAccountIsDisabled, "origin errs")
	err2 := eris.Wrap(err, "wrap2")
	err3 := eris.Wrap(err2, "wrap3")

	format := eris.NewDefaultStringFormat(eris.FormatOptions{
		InvertOutput: true, // flag that inverts the error output (wrap errors shown first)
		WithTrace: true,    // flag that enables stack trace output
		InvertTrace: true,  // flag that inverts the stack trace output (top of call stack shown first)
	})

	formattedStr := eris.ToCustomString(err3, format)
	log.Println(formattedStr)
}

func TestConvertProtoErr(t *testing.T) {

	err := eris.Wrap(ErrInvalidInput, "param Format is wrong")
	err = convertProtoErr(err)
	log.Println(err)
}


