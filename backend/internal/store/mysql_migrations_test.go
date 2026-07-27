package store

import "testing"

func TestMySQLMigrationsAreValid(t *testing.T) {
	if err := validateMySQLMigrations(); err != nil {
		t.Fatal(err)
	}
	if len(mysqlMigrations) == 0 {
		t.Fatal("expected at least one MySQL migration")
	}
}
