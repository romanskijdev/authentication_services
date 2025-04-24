package dbmigration

import (
	"github.com/sirupsen/logrus"
	"time"
)

func (m *MigrationFunctions) FunctionMigrations() {
	logrus.Info("🚀 Start db functions creation: ", time.Now().UTC())

	// Function: update_updated_at_column
	m.gormDB.Exec(`
        CREATE OR REPLACE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = now();
            RETURN NEW;
        END;
        $$ language 'plpgsql';
    `)
}
