package sqlfmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_DirectSQLStatement(t *testing.T) {
	p := NewParser(`SELECT model_num FROM phones AS p WHERE p.release_date > '2014-09-30';`)
	l, err := p.DirectSQLStatement()
	assert.NoError(t, err)
	assert.Equal(t, Juxtaposition{
		Concentrated: true,
		Left: Stack{
			Up: Juxtaposition{
				Concentrated: false,
				Left:         Atom("SELECT"),
				Right:        Atom("model_num"),
			},
			Down: Stack{
				Up: Juxtaposition{
					Concentrated: false,
					Left:         Atom("FROM"),
					Right: Juxtaposition{
						Concentrated: false,
						Left:         Atom("phones"),
						Right: Juxtaposition{
							Concentrated: false,
							Left:         Atom("AS"),
							Right:        Atom("p"),
						},
					},
				},
				Down: Juxtaposition{
					Concentrated: false,
					Left:         Atom("WHERE"),
					Right: Juxtaposition{
						Concentrated: false,
						Left: Juxtaposition{
							Concentrated: true,
							Left:         Atom("p"),
							Right: Juxtaposition{
								Concentrated: true,
								Left:         Atom("."),
								Right:        Atom("release_date"),
							},
						},
						Right: Juxtaposition{
							Concentrated: false,
							Left:         Atom(">"),
							Right:        Atom("'2014-09-30'"),
						},
					},
				},
			},
		},
		Right: Atom(";"),
	}, l)
}
