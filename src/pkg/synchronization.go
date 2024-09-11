package pkg

import (
	"database/sql"
	"fmt"
	"log"
)

var config DBConfig

func ConnectToDB(config DBConfig) (*sql.DB, error) {
	var dsn string
	switch config.DBType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/", config.DBUserName, config.DBPassword, config.DBHost, config.DBPort)
	case "Postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/?sslmode=disable", config.DBUserName, config.DBPassword, config.DBHost, config.DBPort)
	case "MSSQL":
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s", config.DBUserName, config.DBPassword, config.DBHost, config.DBPort)
	case "Oracle":
		dsn = fmt.Sprintf("oracle://%s:%s@%s:%s", config.DBUserName, config.DBPassword, config.DBHost, config.DBPort)
	default:
		return nil, fmt.Errorf("unsupported database type")
	}

	db, err := sql.Open(config.DBType, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func FetchDatabaseStatus(db *sql.DB, dbType string) error {
	var query string

	switch dbType {
	case "mysql":
		query = "SHOW STATUS LIKE 'Uptime'"
	case "Postgres":
		query = "SELECT pg_is_in_recovery() AS is_in_recovery"
	case "MSSQL":
		query = "SELECT state_desc FROM sys.databases WHERE name = DB_NAME()"
	case "Oracle":
		query = "SELECT open_mode FROM v$database"
	default:
		return fmt.Errorf("unsupported database type")
	}

	// Execute query for database status
	row := db.QueryRow(query)

	var status string
	switch dbType {
	case "mysql":
		var uptime int
		if err := row.Scan(&status, &uptime); err != nil {
			return err
		}
		log.Printf("Database Uptime: %d seconds", uptime)
	case "Postgres":
		var isInRecovery bool
		if err := row.Scan(&isInRecovery); err != nil {
			return err
		}
		if isInRecovery {
			status = "In Recovery"
		} else {
			status = "Active"
		}
		log.Printf("Database Status: %s", status)
	case "MSSQL":
		if err := row.Scan(&status); err != nil {
			return err
		}
		log.Printf("Database Status: %s", status)
	case "Oracle":
		if err := row.Scan(&status); err != nil {
			return err
		}
		log.Printf("Database Status: %s", status)
	}

	return nil
}

func FetchTablePrivileges(db *sql.DB, dbType, dbName string) error {
	var query string

	switch dbType {
	case "mysql":
		query = `
			SELECT 
				t.table_schema AS database_name, 
				t.table_name AS table_name, 
				SUBSTRING_INDEX(p.grantee, '@', 1) AS username,
				SUBSTRING_INDEX(p.grantee, '@', -1) AS host,
				p.privilege_type
			FROM 
				information_schema.table_privileges p
			JOIN 
				information_schema.tables t 
			ON 
				p.table_schema = t.table_schema 
				AND p.table_name = t.table_name
			WHERE 
				t.table_schema = ?`
	case "Postgres":
		query = `
			SELECT 
				n.nspname AS database_name, 
				c.relname AS table_name, 
				u.usename AS username,
				'pghost' AS host, -- Change as per your requirement
				COALESCE(array_agg(DISTINCT p.privilege_type), '{}') AS privileges
			FROM 
				information_schema.role_table_grants p
			JOIN 
				information_schema.tables t ON p.table_name = t.table_name
			JOIN 
				pg_catalog.pg_class c ON t.table_name = c.relname
			JOIN 
				pg_catalog.pg_namespace n ON c.relnamespace = n.oid
			JOIN 
				pg_catalog.pg_user u ON p.grantee = u.usename
			WHERE 
				n.nspname NOT IN ('pg_catalog', 'information_schema') 
				AND n.nspname = ? 
			GROUP BY 
				n.nspname, c.relname, u.usename`
	case "MSSQL":
		query = `
			SELECT 
				s.name AS database_name, 
				t.name AS table_name, 
				p.name AS username,
				s.name AS host,
				priv.permission_name AS privilege_type
			FROM 
				sys.database_permissions priv
			JOIN 
				sys.objects t ON priv.major_id = t.object_id
			JOIN 
				sys.schemas s ON t.schema_id = s.schema_id
			JOIN 
				sys.database_principals p ON priv.grantee_principal_id = p.principal_id
			WHERE 
				s.name = ?`
	case "Oracle":
		query = `
			SELECT 
				t.owner AS database_name, 
				t.table_name AS table_name, 
				p.grantee AS username,
				'pghost' AS host, -- Change as per your requirement
				p.privilege AS privilege_type
			FROM 
				dba_tab_privs p
			JOIN 
				dba_tables t ON p.table_name = t.table_name
			WHERE 
				t.owner = ?`
	default:
		return fmt.Errorf("unsupported database type")
	}

	// Execute query for table privileges
	rows, err := db.Query(query, dbName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var dbName, tableName, username, host, privilegeType string
		if err := rows.Scan(&dbName, &tableName, &username, &host, &privilegeType); err != nil {
			return err
		}
		log.Printf("Database: %s, Table: %s, User: %s, Host: %s, Privilege: %s", dbName, tableName, username, host, privilegeType)
	}

	return nil
}

func FetchDatabaseDetails(db *sql.DB, dbType string) error {
	var databases []string

	// Fetch database names
	databasesQuery := ""
	switch dbType {
	case "mysql":
		databasesQuery = "SHOW DATABASES"
	case "Postgres":
		databasesQuery = "SELECT datname FROM pg_database WHERE datistemplate = false"
	case "MSSQL":
		databasesQuery = "SELECT name FROM sys.databases"
	case "Oracle":
		databasesQuery = "SELECT name FROM v$database"
	default:
		return fmt.Errorf("unsupported database type")
	}

	// Execute query for database names
	rows, err := db.Query(databasesQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return err
		}
		databases = append(databases, dbName)

		// Fetch database status
		err = FetchDatabaseStatus(db, dbType)
		if err != nil {
			log.Printf("Failed to fetch status for database %s: %v", dbName, err)
		}

		// Fetch tables and privileges for each database
		err = FetchTablePrivileges(db, dbType, dbName)
		if err != nil {
			log.Printf("Failed to fetch table privileges for database %s: %v", dbName, err)
		}
	}

	return nil
}
