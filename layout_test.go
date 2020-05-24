package sqlfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayout(t *testing.T) {
	t.Run("atom", func(t *testing.T) {
		l := Atom("Lorem ipsum dolor")

		var b bytes.Buffer
		assert.NoError(t, l.Write(&b, 0))
		assert.Equal(t, `Lorem ipsum dolor`, b.String())
	})

	t.Run("stack", func(t *testing.T) {
		l := Stack{
			Atom("Lorem ipsum dolor"),
			Atom("consectetur adipiscing elit"),
		}

		var b bytes.Buffer
		assert.NoError(t, l.Write(&b, 0))
		assert.Equal(t, `Lorem ipsum dolor
consectetur adipiscing elit`, b.String())
	})

	t.Run("juxtaposition", func(t *testing.T) {
		l := Juxtaposition{
			Stack{
				Atom("Lorem ipsum dolor"),
				Atom("consectetur adipiscing elit"),
			},
			Stack{
				Atom("Aliquam erat volutpat"),
				Atom("condimentum vitae leo sit"),
			},
		}

		var b bytes.Buffer
		assert.NoError(t, l.Write(&b, 0))
		assert.Equal(t, `Lorem ipsum dolor
consectetur adipiscing elit Aliquam erat volutpat
                            condimentum vitae leo sit`, b.String())
	})

	t.Run("if", func(t *testing.T) {
		l := Concatenation{
			Stack{
				Atom("if (voltage[t] < LOW_THRESHOLD)"),
				Atom("    "),
			},
			Atom("LogLowVoltage(voltage[t])"),
		}

		var b bytes.Buffer
		assert.NoError(t, l.Write(&b, 0))
		assert.Equal(t, `if (voltage[t] < LOW_THRESHOLD)
    LogLowVoltage(voltage[t])`, b.String())
	})
}
