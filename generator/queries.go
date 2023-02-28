package generator

const (
	getEnumData = `SELECT  e.enumlabel as value
  FROM pg_enum e
  JOIN pg_type t ON e.enumtypid = t.oid
where typname = $1`
	w
)
