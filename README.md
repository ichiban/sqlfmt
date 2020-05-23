# sqlfmt

An opinionated SQL formatter.

```console
$ go install github.com/ichiban/sqlfmt/cmd/sqlfmt
```

```console
$ echo "SELECT model_num FROM phones AS p WHERE p.release_date > '2014-09-30';" | sqlfmt -
SELECT model_num
  FROM phones AS p
 WHERE p.release_date > '2014-09-30';
```

- https://www.sqlstyle.guide/
- https://static.googleusercontent.com/media/research.google.com/en//pubs/archive/44667.pdf
- https://jakewheat.github.io/sql-overview/sql-2011-foundation-grammar.html