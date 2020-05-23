package sqlfmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_Next(t *testing.T) {
	t.Run("character string literal", func(t *testing.T) {
		l := NewLexer(`'foo'`)

		assert.Equal(t, token{typ: characterString, val: "'foo'"}, l.Next())
		assert.Equal(t, token{typ: eos}, l.Next())
	})

	t.Run("simple select", func(t *testing.T) {
		l := NewLexer("SELECT * FROM Customers;")

		assert.Equal(t, token{typ: keyword, val: "SELECT"}, l.Next())
		assert.Equal(t, token{typ: asterisk, val: "*"}, l.Next())
		assert.Equal(t, token{typ: keyword, val: "FROM"}, l.Next())
		assert.Equal(t, token{typ: identifier, val: "Customers"}, l.Next())
		assert.Equal(t, token{typ: semicolon, val: ";"}, l.Next())
		assert.Equal(t, token{typ: eos}, l.Next())
	})

	t.Run("create table dept", func(t *testing.T) {
		l := NewLexer(`
create table dept(  
  deptno     integer primary key,
  dname      text,
  loc        text,
)
`)

		assert.Equal(t, token{typ: keyword, val: "CREATE"}, l.Next())
		assert.Equal(t, token{typ: keyword, val: "TABLE"}, l.Next())
		assert.Equal(t, token{typ: identifier, val: "dept"}, l.Next())
		assert.Equal(t, token{typ: leftParen, val: "("}, l.Next())
		assert.Equal(t, token{typ: identifier, val: "deptno"}, l.Next())
		assert.Equal(t, token{typ: keyword, val: "INTEGER"}, l.Next())
		assert.Equal(t, token{typ: keyword, val: "PRIMARY"}, l.Next())
		assert.Equal(t, token{typ: keyword, val: "KEY"}, l.Next())
		assert.Equal(t, token{typ: comma, val: ","}, l.Next())
		assert.Equal(t, token{typ: identifier, val: "dname"}, l.Next())
		assert.Equal(t, token{typ: keyword, val: "TEXT"}, l.Next())
		assert.Equal(t, token{typ: comma, val: ","}, l.Next())
		assert.Equal(t, token{typ: identifier, val: "loc"}, l.Next())
		assert.Equal(t, token{typ: keyword, val: "TEXT"}, l.Next())
		assert.Equal(t, token{typ: comma, val: ","}, l.Next())
		assert.Equal(t, token{typ: rightParen, val: ")"}, l.Next())
		assert.Equal(t, token{typ: eos}, l.Next())
	})

	t.Run("simple insert", func(t *testing.T) {
		assert := assert.New(t)

		l := NewLexer(`
insert into DEPT (DEPTNO, DNAME, LOC)
values(10, 'ACCOUNTING', 'NEW YORK');
`)

		assert.Equal(token{typ: keyword, val: "INSERT"}, l.Next())
		assert.Equal(token{typ: keyword, val: "INTO"}, l.Next())
		assert.Equal(token{typ: identifier, val: "DEPT"}, l.Next())
		assert.Equal(token{typ: leftParen, val: "("}, l.Next())
		assert.Equal(token{typ: identifier, val: "DEPTNO"}, l.Next())
		assert.Equal(token{typ: comma, val: ","}, l.Next())
		assert.Equal(token{typ: identifier, val: "DNAME"}, l.Next())
		assert.Equal(token{typ: comma, val: ","}, l.Next())
		assert.Equal(token{typ: identifier, val: "LOC"}, l.Next())
		assert.Equal(token{typ: rightParen, val: ")"}, l.Next())
		assert.Equal(token{typ: keyword, val: "VALUES"}, l.Next())
		assert.Equal(token{typ: leftParen, val: "("}, l.Next())
		assert.Equal(token{typ: unsignedNumeric, val: "10"}, l.Next())
		assert.Equal(token{typ: comma, val: ","}, l.Next())
		assert.Equal(token{typ: characterString, val: "'ACCOUNTING'"}, l.Next())
		assert.Equal(token{typ: comma, val: ","}, l.Next())
		assert.Equal(token{typ: characterString, val: "'NEW YORK'"}, l.Next())
		assert.Equal(token{typ: rightParen, val: ")"}, l.Next())
		assert.Equal(token{typ: semicolon, val: ";"}, l.Next())
		assert.Equal(token{typ: eos}, l.Next())
	})

	t.Run("simple query", func(t *testing.T) {
		assert := assert.New(t)

		l := NewLexer(`SELECT model_num FROM phones AS p WHERE p.release_date > '2014-09-30';`)
		assert.Equal(token{typ: keyword, val: "SELECT"}, l.Next())
		assert.Equal(token{typ: identifier, val: "model_num"}, l.Next())
		assert.Equal(token{typ: keyword, val: "FROM"}, l.Next())
		assert.Equal(token{typ: identifier, val: "phones"}, l.Next())
		assert.Equal(token{typ: keyword, val: "AS"}, l.Next())
		assert.Equal(token{typ: identifier, val: "p"}, l.Next())
		assert.Equal(token{typ: keyword, val: "WHERE"}, l.Next())
		assert.Equal(token{typ: identifier, val: "p"}, l.Next())
		assert.Equal(token{typ: period, val: "."}, l.Next())
		assert.Equal(token{typ: identifier, val: "release_date"}, l.Next())
		assert.Equal(token{typ: greaterThanOperator, val: ">"}, l.Next())
		assert.Equal(token{typ: characterString, val: "'2014-09-30'"}, l.Next())
		assert.Equal(token{typ: semicolon, val: ";"}, l.Next())
		assert.Equal(token{typ: eos}, l.Next())
	})
}
