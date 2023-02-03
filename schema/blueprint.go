package schema

import "github.com/goal-web/supports/utils"

type Blueprint struct {
	Columns []*ColumnDefinition
	Table   string
}

func (bp *Blueprint) AddColumn(typeValue, name string) *ColumnDefinition {
	column := &ColumnDefinition{TypeValue: typeValue, Name: name, Table: bp.Table}
	bp.Columns = append(bp.Columns, column)
	return column
}

func (bp *Blueprint) Id(name ...string) *ColumnDefinition {
	return bp.BigInteger(utils.DefaultString(name, "id")).Unsigned().AutoIncrement()
}

func (bp *Blueprint) BigInteger(name string) *ColumnDefinition {
	return bp.AddColumn("bigInteger", name)
}

func (bp *Blueprint) Integer(name string) *ColumnDefinition {
	return bp.AddColumn("integer", name)
}

func (bp *Blueprint) TinyInteger(name string) *ColumnDefinition {
	return bp.AddColumn("tinyInteger", name)
}

func (bp *Blueprint) SmallInteger(name string) *ColumnDefinition {
	return bp.AddColumn("smallInteger", name)
}

func (bp *Blueprint) MediumInteger(name string) *ColumnDefinition {
	return bp.AddColumn("mediumInteger", name)
}

func (bp *Blueprint) Float(name string, args ...int) *ColumnDefinition {
	var total, places = 8, 2
	switch len(args) {
	case 1:
		total = args[0]
	case 2:
		total = args[0]
		places = args[1]
	}
	return bp.AddColumn("float", name).Float(total, places)
}

func (bp *Blueprint) Double(name string, args ...int) *ColumnDefinition {
	var total, places = 8, 2
	switch len(args) {
	case 1:
		total = args[0]
	case 2:
		total = args[0]
		places = args[1]
	}
	return bp.AddColumn("double", name).Float(total, places)
}

func (bp *Blueprint) Decimal(name string, args ...int) *ColumnDefinition {
	var total, places = 8, 2
	switch len(args) {
	case 1:
		total = args[0]
	case 2:
		total = args[0]
		places = args[1]
	}
	return bp.AddColumn("decimal", name).Float(total, places)
}

func (bp *Blueprint) Boolean(name string) *ColumnDefinition {
	return bp.AddColumn("boolean", name)
}

func (bp *Blueprint) Json(name string) *ColumnDefinition {
	return bp.AddColumn("json", name)
}

func (bp *Blueprint) Jsonb(name string) *ColumnDefinition {
	return bp.AddColumn("jsonb", name)
}

func (bp *Blueprint) Date(name string) *ColumnDefinition {
	return bp.AddColumn("date", name)
}

func (bp *Blueprint) Datetime(name string, precision ...int) *ColumnDefinition {
	return bp.AddColumn("dateTime", name).Precision(utils.DefaultInt(precision))
}

func (bp *Blueprint) DatetimeTz(name string, precision ...int) *ColumnDefinition {
	return bp.AddColumn("dateTimeTz", name).Precision(utils.DefaultInt(precision))
}

func (bp *Blueprint) Time(name string, precision ...int) *ColumnDefinition {
	return bp.AddColumn("time", name).Precision(utils.DefaultInt(precision))
}

func (bp *Blueprint) TimeTz(name string, precision ...int) *ColumnDefinition {
	return bp.AddColumn("timeTz", name).Precision(utils.DefaultInt(precision))
}

func (bp *Blueprint) Timestamp(name string, precision ...int) *ColumnDefinition {
	return bp.AddColumn("timestamp", name).Precision(utils.DefaultInt(precision))
}

func (bp *Blueprint) TimestampTz(name string, precision ...int) *ColumnDefinition {
	return bp.AddColumn("timestampTz", name).Precision(utils.DefaultInt(precision))
}

func (bp *Blueprint) Timestamps(precision ...int) {
	bp.Timestamp("created_at").Precision(utils.DefaultInt(precision)).Nullable()
	bp.Timestamp("updated_at").Precision(utils.DefaultInt(precision)).Nullable()
}

func (bp *Blueprint) TimestampsTz(precision ...int) {
	bp.TimestampTz("created_at").Precision(utils.DefaultInt(precision)).Nullable()
	bp.TimestampTz("updated_at").Precision(utils.DefaultInt(precision)).Nullable()
}

func (bp *Blueprint) SoftDeletes(name ...string) *ColumnDefinition {
	return bp.Timestamp(utils.DefaultString(name, "deleted_at")).Nullable()
}

func (bp *Blueprint) SoftDeletesTz(name ...string) *ColumnDefinition {
	return bp.TimestampTz(utils.DefaultString(name, "deleted_at")).Nullable()
}

func (bp *Blueprint) Year(name string) *ColumnDefinition {
	return bp.AddColumn("year", name)
}

func (bp *Blueprint) Binary(name string) *ColumnDefinition {
	return bp.AddColumn("binary", name)
}

func (bp *Blueprint) Uuid(name string) *ColumnDefinition {
	return bp.AddColumn("uuid", name)
}

func (bp *Blueprint) IpAddress(name ...string) *ColumnDefinition {
	return bp.AddColumn("ipAddress", utils.DefaultString(name, "ip_address"))
}

func (bp *Blueprint) MacAddress(name ...string) *ColumnDefinition {
	return bp.AddColumn("macAddress", utils.DefaultString(name, "mac_address"))
}

func (bp *Blueprint) Geometry(name string) *ColumnDefinition {
	return bp.AddColumn("geometry", name)
}
