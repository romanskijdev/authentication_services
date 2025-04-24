package dbmigration

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func CreateEnumType(db *gorm.DB, typeName string, values []string, comment string) error {
	valuesStr := "'" + strings.Join(values, "', '") + "'"
	createQuery := fmt.Sprintf(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
				CREATE TYPE %s AS ENUM (%s);
				COMMENT ON TYPE %s IS '%s';
			END IF;
		END $$;
	`, typeName, typeName, valuesStr, typeName, comment)
	return db.Exec(createQuery).Error
}
