package generator

const (
	getEnumData = `SELECT typname as type, e.enumlabel as value
  FROM pg_enum e
  JOIN pg_type t ON e.enumtypid = t.oid`
)
