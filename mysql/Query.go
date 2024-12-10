package mysql

func (db *Database) Query(sql string, parameters ...any) ([]Record, error) {

	allRecords := make([]Record, 0)

	DatabaseConnection, err := getConnection()
	if err != nil {
		return allRecords, err
	}

	rows, err := DatabaseConnection.Query(sql, parameters...)

	if err != nil {
		return allRecords, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		_ = rows.Scan(valuePtrs...)

		out := Record{}

		for i, col := range columns {
			val := values[i]

			// TODO: Implement All the Types!

			// nolint:gosimple
			switch val.(type) {
			case int64:
				// fmt.Printf("Int: %v\n", val)
				out[col] = Field{Value: val}
			case float64:
				// fmt.Printf("Float64: %v\n", val)
				out[col] = Field{Value: val}

			case []uint8:
				b, _ := val.([]byte)
				// fmt.Printf("String: %s\n", string(b))
				// l.INFO("Type: %T", val)
				out[col] = Field{Value: string(b)}

			case interface{}:
				// l.ERROR("Unknown Type: %T", val)
				// If the Record is NULL
				out[col] = Field{Value: val}

			default:
				// l.ERROR("Unknown Type: %T", val)
				out[col] = Field{Value: val}
			}
		}
		allRecords = append(allRecords, out)
	}

	return allRecords, nil
}

// QueryV2 is used to run SELECT queries and it returns the outputs. [V2 => To handle multiple databases]
func (db *Database) QueryV2(sql string, parameters ...any) ([]Record, error) {

	allRecords := make([]Record, 0)

	DatabaseConnection, err := db.getConnectionV2()
	if err != nil {
		return allRecords, err
	}

	rows, err := DatabaseConnection.Query(sql, parameters...)

	if err != nil {
		return allRecords, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		_ = rows.Scan(valuePtrs...)

		out := Record{}

		for i, col := range columns {
			val := values[i]

			// TODO: Implement All the Types!

			// nolint:gosimple
			switch val.(type) {
			case int64:
				// fmt.Printf("Int: %v\n", val)
				out[col] = Field{Value: val}
			case float64:
				// fmt.Printf("Float64: %v\n", val)
				out[col] = Field{Value: val}

			case []uint8:
				b, _ := val.([]byte)
				// fmt.Printf("String: %s\n", string(b))
				// l.INFO("Type: %T", val)
				out[col] = Field{Value: string(b)}

			case interface{}:
				// l.ERROR("Unknown Type: %T", val)
				// If the Record is NULL
				out[col] = Field{Value: val}

			default:
				// l.ERROR("Unknown Type: %T", val)
				out[col] = Field{Value: val}
			}
		}
		allRecords = append(allRecords, out)
	}

	return allRecords, nil
}
