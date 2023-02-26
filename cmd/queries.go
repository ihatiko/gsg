package cmd

// language=sql
const (
	getDatabaseInfoQuery = `
with dependencies as (
select distinct kcu.table_name  as in_table,
                                      kcu.column_name as column_in, ccu.table_name as out_table,
ccu.column_name AS column_name_out,
                                      itc.constraint_type
from information_schema.key_column_usage as kcu
 JOIN information_schema.constraint_column_usage AS ccu
           ON ccu.constraint_name = kcu.constraint_name
left join information_schema.table_constraints itc on itc.constraint_name = kcu.constraint_name
where kcu.table_schema = 'public' and itc.constraint_type != 'PRIMARY KEY')


 select
     table_catalog,
     table_name,
     column_name,
     column_default,
     data_type,
     case is_nullable when 'YES' then true else false end as is_nullable,
     character_maximum_length,
     dependencies.out_table as dependency_table_name,
     dependencies.column_name_out as dependency_column_name,
     dependencies.constraint_type
 from INFORMATION_SCHEMA.COLUMNS
 left join dependencies on dependencies.in_table = table_name and dependencies.column_in = column_name
 where table_schema = 'public' and not exists(select * from information_schema.views iv
where table_schema = 'public' and INFORMATION_SCHEMA.COLUMNS.table_name = iv.table_name)
`

	getDatabasesQuery = `SELECT datname as name FROM pg_database
where datistemplate = false`
)
