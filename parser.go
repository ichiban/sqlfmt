package sqlfmt

import (
	"fmt"
)

type Parser struct {
	lexer   *Lexer
	current token
}

func NewParser(input string) *Parser {
	p := Parser{lexer: NewLexer(input)}
	p.current = p.lexer.Next()
	return &p
}

// beginning of syntax

func (p *Parser) DirectSQLStatement() (Layout, error) {
	s, err := p.directlyExecutableStatement()
	if err != nil {
		return nil, fmt.Errorf("directly executable statement: %w", err)
	}

	if _, err := p.accept(semicolon); err != nil {
		return nil, err
	}
	return Juxtaposition{
		Concentrated: true,
		Left:         s,
		Right:        Atom(";"),
	}, nil
}

func (p *Parser) directlyExecutableStatement() (Layout, error) {
	return p.cursorSpecification()
}

func (p *Parser) cursorSpecification() (Layout, error) {
	return p.queryExpression()
}

func (p *Parser) queryExpression() (Layout, error) {
	return p.queryExpressionBody()
}

func (p *Parser) queryExpressionBody() (Layout, error) {
	q, err := p.queryTerm()
	if err != nil {
		return nil, fmt.Errorf("query term: %w", err)
	}

	for {
		v, err := p.accept(keyword, "UNION", "EXCEPT")
		if err != nil {
			break
		}
		u := Layout(Atom(v))
		if v, err := p.accept(keyword, "ALL", "DISTINCT"); err == nil {
			u = Juxtaposition{
				Left:  u,
				Right: Atom(v),
			}
		}
		if c, err := p.correspondingSpec(); err == nil {
			u = Juxtaposition{
				Left:  u,
				Right: c,
			}
		}

		t, err := p.queryTerm()
		if err != nil {
			return nil, err
		}

		q = Stack{
			Up: q,
			Down: Stack{
				Up:   u,
				Down: t,
			},
		}
	}

	return q, nil
}

func (p *Parser) queryTerm() (Layout, error) {
	return p.queryPrimary()
}

func (p *Parser) queryPrimary() (Layout, error) {
	if s, err := p.simpleTable(); err == nil {
		return s, nil
	}

	v, err := p.accept(leftParen)
	if err != nil {
		return nil, err
	}

	l := Layout(Atom(v))

	b, err := p.queryExpressionBody()
	if err != nil {
		return nil, fmt.Errorf("query expression body: %w", err)
	}
	l = Juxtaposition{
		Concentrated: true,
		Left:         l,
		Right:        b,
	}

	// TODO: order by clause, result offset clause, fetch first clause

	v, err = p.accept(rightParen)
	if err != nil {
		return nil, err
	}

	return Juxtaposition{
		Concentrated: true,
		Left:         l,
		Right:        Atom(v),
	}, nil
}

func (p *Parser) simpleTable() (Layout, error) {
	return p.querySpecification()
}

func (p *Parser) correspondingSpec() (Layout, error) {
	v, err := p.accept(keyword, "CORRESPONDING")
	if err != nil {
		return nil, err
	}

	c := Layout(Atom(v))

	if v, err := p.accept(keyword, "BY"); err == nil {
		c = Juxtaposition{
			Left:  c,
			Right: Atom(v),
		}

		if _, err := p.accept(leftParen); err != nil {
			return nil, err
		}

		l, err := p.correspondingColumnList()
		if err != nil {
			return nil, err
		}
		c = Juxtaposition{
			Left:  c,
			Right: l,
		}

		if _, err := p.accept(rightParen); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (p *Parser) correspondingColumnList() (Layout, error) {
	return p.columnNameList()
}

func (p *Parser) querySpecification() (Layout, error) {
	if _, err := p.accept(keyword, "SELECT"); err != nil {
		return nil, err
	}

	s, err := p.selectList()
	if err != nil {
		return nil, fmt.Errorf("select list: %w", err)
	}

	t, err := p.tableExpression()
	if err != nil {
		return nil, fmt.Errorf("table list: %w", err)
	}

	return Stack{
		Up: Juxtaposition{
			Left:  Atom("SELECT"),
			Right: s,
		},
		Down: t,
	}, nil
}

func (p *Parser) selectList() (Layout, error) {
	if _, err := p.accept(asterisk); err == nil {
		return Atom("*"), nil
	}

	s, err := p.selectSublist()
	if err != nil {
		return nil, fmt.Errorf("select sublist: %w", err)
	}

	for {
		if _, err := p.accept(comma); err != nil {
			break
		}

		n, err := p.selectSublist()
		if err != nil {
			return nil, fmt.Errorf("select sublist: %w", err)
		}

		s = Juxtaposition{
			Left:  s,
			Right: n,
		}
	}

	return s, nil
}

func (p *Parser) selectSublist() (Layout, error) {
	return p.derivedColumn()
}

func (p *Parser) derivedColumn() (Layout, error) {
	v, err := p.valueExpression()
	if err != nil {
		return nil, fmt.Errorf("value expression: %w", err)
	}

	if a, err := p.asClause(); err == nil {
		v = Juxtaposition{
			Left:  v,
			Right: a,
		}
	}

	return v, nil
}

func (p *Parser) asClause() (Layout, error) {
	var l Layout
	if v, err := p.accept(keyword, "AS"); err == nil {
		l = Atom(v)
	}

	c, err := p.columnName()
	if err != nil {
		return nil, fmt.Errorf("column name: %w", err)
	}

	if l != nil {
		c = Juxtaposition{
			Left:  l,
			Right: c,
		}
	}

	return c, nil
}

func (p *Parser) valueExpression() (Layout, error) {
	return p.commonValueExpression()
}

func (p *Parser) commonValueExpression() (Layout, error) {
	if e, err := p.stringValueExpression(); err == nil {
		return e, nil
	}
	return p.referenceValueExpression()
}

func (p *Parser) stringValueExpression() (Layout, error) {
	return p.characterValueExpression()
}

func (p *Parser) characterValueExpression() (Layout, error) {
	return p.characterFactor()
}

func (p *Parser) characterFactor() (Layout, error) {
	return p.characterPrimary()
}

func (p *Parser) characterPrimary() (Layout, error) {
	return p.valueExpressionPrimary()
}

func (p *Parser) valueExpressionPrimary() (Layout, error) {
	return p.nonparenthesizedValueExpressionPrimary()
}

func (p *Parser) nonparenthesizedValueExpressionPrimary() (Layout, error) {
	if v, err := p.unsignedValueSpecification(); err == nil {
		return v, nil
	}

	if r, err := p.columnReference(); err == nil {
		return r, nil
	}

	return p.setFunctionSpecification()
}

func (p *Parser) unsignedValueSpecification() (Layout, error) {
	return p.unsignedLiteral()
}

func (p *Parser) unsignedLiteral() (Layout, error) {
	return p.generalLiteral()
}

func (p *Parser) generalLiteral() (Layout, error) {
	c, err := p.accept(characterString)
	if err != nil {
		return nil, err
	}
	return Atom(c), nil
}

func (p *Parser) columnReference() (Layout, error) {
	return p.basicIdentifierChain()
}

func (p *Parser) basicIdentifierChain() (Layout, error) {
	return p.identifierChain()
}

func (p *Parser) setFunctionSpecification() (Layout, error) {
	return p.aggregateFunction()
}

func (p *Parser) aggregateFunction() (Layout, error) {
	return p.generalSetFunction()
}

func (p *Parser) generalSetFunction() (Layout, error) {
	l, err := p.setFunctionType()
	if err != nil {
		return nil, fmt.Errorf("set function type: %w", err)
	}

	v, err := p.accept(leftParen)
	if err != nil {
		return nil, err
	}
	l = Juxtaposition{
		Concentrated: true,
		Left:         l,
		Right:        Atom(v),
	}

	e, err := p.valueExpression()
	if err != nil {
		return nil, fmt.Errorf("value expression: %w", err)
	}
	l = Juxtaposition{
		Concentrated: true,
		Left:         l,
		Right:        e,
	}

	v, err = p.accept(rightParen)
	if err != nil {
		return nil, err
	}

	return Juxtaposition{
		Concentrated: true,
		Left:         l,
		Right:        Atom(v),
	}, nil
}

func (p *Parser) setFunctionType() (Layout, error) {
	return p.computationalOperation()
}

func (p *Parser) computationalOperation() (Layout, error) {
	v, err := p.accept(keyword, "AVG", "MAX", "MIN", "SUM", "EVERY", "ANY", "SOME", "COUNT", "STDDEV_POP", "STDDEV_SAMP", "VAR_SAMP", "VAR_POP", "COLLECT", "FUSION", "INTERSECTION")
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) identifierChain() (Layout, error) {
	l, err := p.identifier()
	if err != nil {
		return nil, err
	}

	for {
		if _, err := p.accept(period); err != nil {
			break
		}

		i, err := p.identifier()
		if err != nil {
			return nil, fmt.Errorf("identifier: %w", err)
		}
		l = Juxtaposition{
			Concentrated: true,
			Left:         l,
			Right: Juxtaposition{
				Concentrated: true,
				Left:         Atom("."),
				Right:        i,
			},
		}
	}

	return l, nil
}

func (p *Parser) identifier() (Layout, error) {
	return p.actualIdentifier()
}

func (p *Parser) actualIdentifier() (Layout, error) {
	return p.regularIdentifier()
}

func (p *Parser) regularIdentifier() (Layout, error) {
	v, err := p.accept(identifier)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) referenceValueExpression() (Layout, error) {
	return p.valueExpressionPrimary()
}

func (p *Parser) tableExpression() (Layout, error) {
	f, err := p.fromClause()
	if err != nil {
		return nil, err
	}

	if w, err := p.whereClause(); err == nil {
		f = Stack{
			Up:   f,
			Down: w,
		}
	}

	if g, err := p.groupByClause(); err == nil {
		f = Stack{
			Up:   f,
			Down: g,
		}
	}

	return f, nil
}

func (p *Parser) fromClause() (Layout, error) {
	if _, err := p.accept(keyword, "FROM"); err != nil {
		return nil, err
	}
	l, err := p.tableReferenceList()
	if err != nil {
		return nil, err
	}
	return Juxtaposition{
		Left:  Atom("FROM"),
		Right: l,
	}, nil
}

func (p *Parser) tableReferenceList() (Layout, error) {
	l, err := p.tableReference()
	if err != nil {
		return nil, err
	}

	for {
		if _, err := p.accept(comma); err != nil {
			break
		}

		t, err := p.tableReference()
		if err != nil {
			return nil, err
		}

		l = Juxtaposition{
			Left:  l,
			Right: t,
		}
	}

	return l, nil
}

func (p *Parser) tableReference() (Layout, error) {
	return p.tableFactor()
}

func (p *Parser) tableFactor() (Layout, error) {
	return p.tablePrimary()
}

func (p *Parser) tablePrimary() (Layout, error) {
	n, err := p.tableOrQueryName()
	if err != nil {
		return nil, err
	}

	if _, err := p.accept(keyword, "AS"); err == nil {
		c, err := p.correlationName()
		if err != nil {
			return nil, err
		}
		n = Juxtaposition{
			Left: n,
			Right: Juxtaposition{
				Left:  Atom("AS"),
				Right: c,
			},
		}
	}

	return n, nil
}

func (p *Parser) tableOrQueryName() (Layout, error) {
	return p.tableName()
}

func (p *Parser) columnNameList() (Layout, error) {
	l, err := p.columnName()
	if err != nil {
		return nil, fmt.Errorf("column name: %w", err)
	}

	for {
		v, err := p.accept(comma)
		if err != nil {
			break
		}

		n, err := p.columnName()
		if err != nil {
			return nil, err
		}

		l = Juxtaposition{
			Concentrated: true,
			Left:         l,
			Right: Juxtaposition{
				Left:  Atom(v),
				Right: n,
			},
		}
	}

	return l, nil
}

func (p *Parser) whereClause() (Layout, error) {
	if _, err := p.accept(keyword, "WHERE"); err != nil {
		return nil, err
	}

	s, err := p.searchCondition()
	if err != nil {
		return nil, fmt.Errorf("search condition: %w", err)
	}

	return Juxtaposition{
		Left:  Atom("WHERE"),
		Right: s,
	}, nil
}

func (p *Parser) groupByClause() (Layout, error) {
	v, err := p.accept(keyword, "GROUP")
	if err != nil {
		return nil, err
	}
	l := Layout(Atom(v))

	v, err = p.accept(keyword, "BY")
	if err != nil {
		return nil, err
	}
	l = Juxtaposition{
		Left:  l,
		Right: Atom(v),
	}

	g, err := p.groupingElementList()
	if err != nil {
		return nil, fmt.Errorf("grouping element list: %w", err)
	}

	return Juxtaposition{
		Left:  l,
		Right: g,
	}, nil
}

func (p *Parser) groupingElementList() (Layout, error) {
	l, err := p.groupingElement()
	if err != nil {
		return nil, fmt.Errorf("groupoing element: %w", err)
	}

	for {
		v, err := p.accept(comma)
		if err != nil {
			break
		}

		e, err := p.groupingElement()
		if err != nil {
			return nil, fmt.Errorf("grouping element: %w", err)
		}

		l = Juxtaposition{
			Concentrated: true,
			Left:         l,
			Right: Juxtaposition{
				Left:  Atom(v),
				Right: e,
			},
		}
	}

	return l, nil
}

func (p *Parser) groupingElement() (Layout, error) {
	return p.ordinaryGroupingSet()
}

func (p *Parser) ordinaryGroupingSet() (Layout, error) {
	return p.groupingColumnReference()
}

func (p *Parser) groupingColumnReference() (Layout, error) {
	return p.columnReference()
}

func (p *Parser) searchCondition() (Layout, error) {
	return p.booleanValueExpression()
}

func (p *Parser) booleanValueExpression() (Layout, error) {
	l, err := p.booleanTerm()
	if err != nil {
		return nil, fmt.Errorf("boolean term: %w", err)
	}

	for {
		v, err := p.accept(keyword, "OR")
		if err != nil {
			break
		}

		t, err := p.booleanTerm()
		if err != nil {
			return nil, fmt.Errorf("boolean term: %w", err)
		}

		l = Stack{
			Up: l,
			Down: Juxtaposition{
				Left:  Atom(v),
				Right: t,
			},
		}
	}

	return l, nil
}

func (p *Parser) booleanTerm() (Layout, error) {
	return p.booleanFactor()
}

func (p *Parser) booleanFactor() (Layout, error) {
	var not bool
	if _, err := p.accept(keyword, "NOT"); err == nil {
		not = true
	}

	t, err := p.booleanTest()
	if err != nil {
		return nil, err
	}

	if not {
		t = Juxtaposition{
			Left:  Atom("NOT"),
			Right: t,
		}
	}

	return t, nil
}

func (p *Parser) booleanTest() (Layout, error) {
	return p.booleanPrimary()
}

func (p *Parser) booleanPrimary() (Layout, error) {
	return p.predicate()
}

func (p *Parser) predicate() (Layout, error) {
	return p.comparisonPredicate()
}

func (p *Parser) comparisonPredicate() (Layout, error) {
	l, err := p.rowValuePredicand()
	if err != nil {
		return nil, fmt.Errorf("row value predicand: %w", err)
	}

	r, err := p.comparisonPredicatePart2()
	if err != nil {
		return nil, fmt.Errorf("comparison predicate part 2: %w", err)
	}

	return Juxtaposition{
		Left:  l,
		Right: r,
	}, nil
}

func (p *Parser) comparisonPredicatePart2() (Layout, error) {
	o, err := p.compOp()
	if err != nil {
		return nil, fmt.Errorf("comp op: %w", err)
	}

	v, err := p.rowValuePredicand()
	if err != nil {
		return nil, fmt.Errorf("row value predicand: %w", err)
	}

	return Juxtaposition{
		Left:  o,
		Right: v,
	}, nil
}

func (p *Parser) compOp() (Layout, error) {
	if o, err := p.equalsOperator(); err == nil {
		return o, nil
	}
	if o, err := p.notEqualsOperator(); err == nil {
		return o, nil
	}
	if o, err := p.lessThanOperator(); err == nil {
		return o, nil
	}
	if o, err := p.greaterThanOperator(); err == nil {
		return o, nil
	}
	if o, err := p.lessThanOrEqualsOperator(); err == nil {
		return o, nil
	}
	return p.greaterThanOrEqualsOperator()
}
func (p *Parser) equalsOperator() (Layout, error) {
	v, err := p.accept(equalsOperator)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) notEqualsOperator() (Layout, error) {
	v, err := p.accept(notEqualsOperator)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) lessThanOperator() (Layout, error) {
	v, err := p.accept(lessThanOperator)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) greaterThanOperator() (Layout, error) {
	v, err := p.accept(greaterThanOperator)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) lessThanOrEqualsOperator() (Layout, error) {
	v, err := p.accept(lessThanOrEqualsOperator)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) greaterThanOrEqualsOperator() (Layout, error) {
	v, err := p.accept(greaterThanOrEqualsOperator)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) rowValuePredicand() (Layout, error) {
	if v, err := p.rowValueSpecialCase(); err == nil {
		return v, nil
	}
	return p.rowValueConstructorPredicand()
}

func (p *Parser) rowValueSpecialCase() (Layout, error) {
	return p.nonparenthesizedValueExpressionPrimary()
}

func (p *Parser) rowValueConstructorPredicand() (Layout, error) {
	return p.commonValueExpression()
}

func (p *Parser) tableName() (Layout, error) {
	return p.localOrSchemaQualifiedName()
}

func (p *Parser) localOrSchemaQualifiedName() (Layout, error) {
	return p.qualifiedIdentifier()
}

func (p *Parser) qualifiedIdentifier() (Layout, error) {
	v, err := p.accept(identifier)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) columnName() (Layout, error) {
	v, err := p.accept(identifier)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

func (p *Parser) correlationName() (Layout, error) {
	v, err := p.accept(identifier)
	if err != nil {
		return nil, err
	}
	return Atom(v), nil
}

// end of syntax

func (p *Parser) accept(t tokenType, vals ...string) (string, error) {
	v, err := p.expect(t, vals...)
	if err != nil {
		return "", err
	}
	p.current = p.lexer.Next()
	return v, nil
}

func (p *Parser) expect(t tokenType, vals ...string) (string, error) {
	if p.current.typ != t {
		return "", &ErrUnexpected{
			ExpectedType:   t,
			ExpectedValues: vals,
			Actual:         p.current,
		}
	}

	if len(vals) > 0 {
		for _, v := range vals {
			if v == p.current.val {
				return p.current.val, nil
			}
		}
		return "", &ErrUnexpected{
			ExpectedType:   t,
			ExpectedValues: vals,
			Actual:         p.current,
		}
	}

	return p.current.val, nil
}

type ErrUnexpected struct {
	ExpectedType   tokenType
	ExpectedValues []string
	Actual         token
}

func (e *ErrUnexpected) Error() string {
	return fmt.Sprintf("expected: <%s %s>, actual: %s", e.ExpectedType, e.ExpectedValues, e.Actual)
}
