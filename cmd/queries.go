package cmd

// language=sql
const (
	getDatabaseInfo = `
with dependencies as (SELECT tc.table_name   as in_table,
                             kcu.column_name as column_in,
                             ccu.table_name  AS out_table,
                             ccu.column_name AS column_name_out
                      FROM information_schema.table_constraints AS tc
                               JOIN information_schema.key_column_usage AS kcu
                                    ON tc.constraint_name = kcu.constraint_name
                               JOIN information_schema.constraint_column_usage AS ccu
                                    ON ccu.constraint_name = tc.constraint_name
                      where tc.constraint_schema = 'public'
                        and tc.table_name != ccu.table_name)


 select
     table_catalog,
     table_name,
     column_name,
     column_default,
     data_type,
     is_nullable,
     character_maximum_length,
     dependencies.out_table as dependency_table_name,
     dependencies.column_name_out as dependency_column_name
 from INFORMATION_SCHEMA.COLUMNS
 left join dependencies on dependencies.in_table = table_name and dependencies.column_in = column_name
 where table_schema = 'public'
`
)
