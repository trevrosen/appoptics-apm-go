package reporter

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/appoptics/appoptics-apm-go/v1/ao/internal/config"
	"github.com/stretchr/testify/assert"
)

// SanitizeTestCase defines the sanitizer input and output
type SanitizeTestCase struct {
	mode         int
	dbType       string
	sql          string
	sanitizedSQL string
}

func TestSQLSanitize(t *testing.T) {
	cases := []SanitizeTestCase{
		// Disabled
		{
			Disabled,
			MySQL,
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
		},
		{
			Disabled,
			Oracle,
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
		},
		{
			Disabled,
			Sybase,
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
		},
		{
			Disabled,
			SQLServer,
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
		},
		{
			Disabled,
			PostgreSQL,
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
		},
		{
			Disabled,
			Default,
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
		},
		// EnabledAuto
		{
			EnabledAuto,
			Default,
			"select * from schema.tbl where name = '';",
			"select * from schema.tbl where name = ?;",
		},
		{
			EnabledDropDoubleQuoted,
			Default,
			"select * from tbl where name = 'private';",
			"select * from tbl where name = ?;",
		},
		{
			EnabledKeepDoubleQuoted,
			Default,
			"select * from tbl where name = 'private' order by age;",
			"select * from tbl where name = ? order by age;",
		},
		{
			EnabledAuto,
			Default,
			"select ssn from accounts where password = \"mypass\" group by dept order by age;",
			"select ssn from accounts where password = ? group by dept order by age;",
		},
		{
			EnabledDropDoubleQuoted,
			Default,
			"select ssn from accounts where password = \"mypass\";",
			"select ssn from accounts where password = ?;",
		},
		{
			EnabledKeepDoubleQuoted,
			Default,
			"select ssn from accounts where password = \"mypass\";",
			"select ssn from accounts where password = \"mypass\";",
		},
		{
			EnabledAuto,
			Default,
			"select ssn from accounts where name = 'Chris O''Corner'",
			"select ssn from accounts where name = ?",
		},
		{
			EnabledAuto,
			Default,
			"SELECT name FROM employees WHERE age = 37 AND firstName = 'Eric'",
			"SELECT name FROM employees WHERE age = ? AND firstName = ?",
		},
		{
			EnabledAuto,
			Default,
			"SELECT name FROM employees WHERE name IN ('Eric', 'Tom')",
			"SELECT name FROM employees WHERE name IN (?, ?)",
		},
		{
			EnabledAuto,
			Default,
			"SELECT TOP 10 FROM employees WHERE age > 28",
			"SELECT TOP ? FROM employees WHERE age > ?",
		},
		{
			EnabledAuto,
			Default,
			"UPDATE Customers SET zip='V3B 6Z6', phone='000-000-0000'",
			"UPDATE Customers SET zip=?, phone=?",
		},
		{
			EnabledAuto,
			Default,
			"SELECT id FROM employees WHERE date BETWEEN '01/01/2019' AND '05/30/2019'",
			"SELECT id FROM employees WHERE date BETWEEN ? AND ?",
		},
		{
			EnabledAuto,
			Default,
			`SELECT name FROM employees 
			 WHERE EXISTS (
				SELECT eid FROM orders WHERE employees.id = orders.eid AND price > 1000
			 )`,
			`SELECT name FROM employees 
			 WHERE EXISTS (
				SELECT eid FROM orders WHERE employees.id = orders.eid AND price > ?
			 )`,
		},
		{
			EnabledAuto,
			Default,
			"SELECT COUNT(id), team FROM employees GROUP BY team HAVING MIN(age) > 30",
			"SELECT COUNT(id), team FROM employees GROUP BY team HAVING MIN(age) > ?",
		},
		{
			EnabledAuto,
			Default,
			`WITH tmp AS (SELECT * FROM employees WHERE team = 'IT') 
			 SELECT * FROM tmp WHERE name = 'Tom'`,
			`WITH tmp AS (SELECT * FROM employees WHERE team = ?) 
			 SELECT * FROM tmp WHERE name = ?`,
		},
		{
			EnabledKeepDoubleQuoted,
			Oracle,
			`WITH tmp AS (SELECT * FROM \"Employees\" WHERE team = 'IT') 
			 SELECT * FROM tmp WHERE name = 'Tom'`,
			`WITH tmp AS (SELECT * FROM \"Employees\" WHERE team = ?) 
			 SELECT * FROM tmp WHERE name = ?`,
		},
		{
			EnabledDropDoubleQuoted,
			Default,
			`WITH tmp AS (SELECT * FROM employees WHERE team = "IT") 
			 SELECT * FROM tmp WHERE name = 'Tom' 
			 LEFT JOIN tickets 
			 WHERE eid = id AND last_update BETWEEN '01/01/2019' AND '05/30/2019'`,
			`WITH tmp AS (SELECT * FROM employees WHERE team = ?) 
			 SELECT * FROM tmp WHERE name = ? 
			 LEFT JOIN tickets 
			 WHERE eid = id AND last_update BETWEEN ? AND ?`,
		},
	}

	for _, c := range cases {
		_ = os.Setenv("APPOPTICS_SQL_SANITIZE", strconv.Itoa(c.mode))
		assert.Nil(t, config.Load())
		ss := initSanitizersMap()

		assert.Equal(t, c.sanitizedSQL, sqlSanitizeInternal(ss, c.dbType, c.sql),
			fmt.Sprintf("Test case: %+v", c))
	}
}
