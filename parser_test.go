package sqlfmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_DirectSQLStatement(t *testing.T) {
	t.Run("simple query", func(t *testing.T) {
		p := NewParser(`SELECT model_num FROM phones AS p WHERE p.release_date > '2014-09-30';`)
		l, err := p.DirectSQLStatement()
		assert.NoError(t, err)
		assert.Equal(t, Concatenation{
			Stack{
				Juxtaposition{Atom("SELECT"), Atom("model_num")},
				Juxtaposition{Atom("FROM"), Atom("phones"), Atom("AS"), Atom("p")},
				Juxtaposition{Atom("WHERE"), Concatenation{Atom("p"), Atom("."), Atom("release_date")}, Atom(">"), Atom("'2014-09-30'")},
			},
			Atom(";"),
		}, l)
	})

	t.Run("union", func(t *testing.T) {
		p := NewParser(`(SELECT f.species_name,
        AVG(f.height) AS average_height, AVG(f.diameter) AS average_diameter
 FROM flora AS f
 WHERE f.species_name = 'Banksia'
    OR f.species_name = 'Sheoak'
    OR f.species_name = 'Wattle'
 GROUP BY f.species_name, f.observation_date)

UNION ALL

(SELECT b.species_name,
        AVG(b.height) AS average_height, AVG(b.diameter) AS average_diameter
 FROM botanic_garden_flora AS b
 WHERE b.species_name = 'Banksia'
    OR b.species_name = 'Sheoak'
    OR b.species_name = 'Wattle'
 GROUP BY b.species_name, b.observation_date);`)

		l, err := p.DirectSQLStatement()
		assert.NoError(t, err)
		assert.Equal(t, Concatenation{
			Stack{
				Concatenation{
					Atom("("),
					Stack{
						Juxtaposition{
							Atom("SELECT"),
							Concatenation{
								Juxtaposition{
									Concatenation{
										Atom("f"),
										Atom("."),
										Atom("species_name"),
									},
								},
								Atom(","),
							},
							Concatenation{
								Juxtaposition{
									Concatenation{
										Atom("AVG"),
										Atom("("),
										Atom("f"),
										Atom("."),
										Atom("height"),
										Atom(")"),
									},
									Atom("AS"),
									Atom("average_height"),
								},
								Atom(","),
							},
							Concatenation{
								Atom("AVG"),
								Atom("("),
								Atom("f"),
								Atom("."),
								Atom("diameter"),
								Atom(")"),
							},
							Atom("AS"),
							Atom("average_diameter"),
						},
						Juxtaposition{
							Atom("FROM"),
							Atom("flora"),
							Atom("AS"),
							Atom("f"),
						},
						Juxtaposition{
							Atom("WHERE"),
							Stack{
								Juxtaposition{
									Concatenation{
										Atom("f"),
										Atom("."),
										Atom("species_name"),
									},
									Atom("="),
									Atom("'Banksia'"),
								},
								Juxtaposition{
									Atom("OR"),
									Concatenation{
										Atom("f"),
										Atom("."),
										Atom("species_name"),
									},
									Atom("="),
									Atom("'Sheoak'"),
								},
								Juxtaposition{
									Atom("OR"),
									Concatenation{
										Atom("f"),
										Atom("."),
										Atom("species_name"),
									},
									Atom("="),
									Atom("'Wattle'"),
								},
							},
						},
						Juxtaposition{
							Atom("GROUP"),
							Atom("BY"),
							Concatenation{
								Concatenation{
									Atom("f"),
									Atom("."),
									Atom("species_name"),
								},
								Atom(","),
							},
							Concatenation{
								Atom("f"),
								Atom("."),
								Atom("observation_date"),
							},
						},
					},
					Atom(")"),
				},
				Juxtaposition{
					Atom("UNION"),
					Atom("ALL"),
				},
				Concatenation{
					Atom("("),
					Stack{
						Juxtaposition{
							Atom("SELECT"),
							Concatenation{
								Juxtaposition{
									Concatenation{
										Atom("b"),
										Atom("."),
										Atom("species_name"),
									},
								},
								Atom(","),
							},
							Concatenation{
								Juxtaposition{
									Concatenation{
										Atom("AVG"),
										Atom("("),
										Atom("b"),
										Atom("."),
										Atom("height"),
										Atom(")"),
									},
									Atom("AS"),
									Atom("average_height"),
								},
								Atom(","),
							},
							Concatenation{
								Atom("AVG"),
								Atom("("),
								Atom("b"),
								Atom("."),
								Atom("diameter"),
								Atom(")"),
							},
							Atom("AS"),
							Atom("average_diameter"),
						},
						Juxtaposition{
							Atom("FROM"),
							Atom("botanic_garden_flora"),
							Atom("AS"),
							Atom("b"),
						},
						Juxtaposition{
							Atom("WHERE"),
							Stack{
								Juxtaposition{
									Concatenation{
										Atom("b"),
										Atom("."),
										Atom("species_name"),
									},
									Atom("="),
									Atom("'Banksia'"),
								},
								Juxtaposition{
									Atom("OR"),
									Concatenation{
										Atom("b"),
										Atom("."),
										Atom("species_name"),
									},
									Atom("="),
									Atom("'Sheoak'"),
								},
								Juxtaposition{
									Atom("OR"),
									Concatenation{
										Atom("b"),
										Atom("."),
										Atom("species_name"),
									},
									Atom("="),
									Atom("'Wattle'"),
								},
							},
						},
						Juxtaposition{
							Atom("GROUP"),
							Atom("BY"),
							Concatenation{
								Concatenation{
									Atom("b"),
									Atom("."),
									Atom("species_name"),
								},
								Atom(","),
							},
							Concatenation{
								Atom("b"),
								Atom("."),
								Atom("observation_date"),
							},
						},
					},
					Atom(")"),
				},
			},
			Atom(";"),
		}, l)
	})
}
