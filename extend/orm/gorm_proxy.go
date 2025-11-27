package orm

//
//type GormClient struct {
//	GormDb *gorm.DB
//}
//
//func (gormDb *GormClient) GetDb() *gorm.DB {
//	return gormDb.GormDb
//}
//
//func (gormDb *GormClient) Session(config *gorm.Session) *GormClient {
//	gormDb.GormDb.Session(config)
//	return gormDb
//}
//
//func (gormDb *GormClient) WithContext(ctx context.Context) *GormClient {
//	gormDb.GormDb.WithContext(ctx)
//	return gormDb
//}
//
//func (gormDb *GormClient) Debug() *GormClient {
//	gormDb.GormDb.Debug()
//	return gormDb
//}
//
//func (gormDb *GormClient) Set(key string, value interface{}) *GormClient {
//	gormDb.GormDb.Set(key, value)
//	return gormDb
//}
//
//func (gormDb *GormClient) Get(key string) (interface{}, bool) {
//	return gormDb.GormDb.Get(key)
//}
//
//func (gormDb *GormClient) InstanceSet(key string, value interface{}) *GormClient {
//	gormDb.GormDb.InstanceSet(key, value)
//	return gormDb
//}
//
//func (gormDb *GormClient) InstanceGet(key string) (interface{}, bool) {
//	return gormDb.GormDb.InstanceGet(key)
//}
//
//func (gormDb *GormClient) AddError(err error) error {
//	return gormDb.GormDb.AddError(err)
//}
//
//func (gormDb *GormClient) DB() (*sql.DB, error) {
//	return gormDb.GormDb.DB()
//}
//
//func (gormDb *GormClient) SetupJoinTable(model interface{}, field string, joinTable interface{}) error {
//	return gormDb.GormDb.SetupJoinTable(model, field, joinTable)
//}
//
//func (gormDb *GormClient) Use(plugin gorm.Plugin) error {
//	return gormDb.GormDb.Use(plugin)
//}
//
//func (gormDb *GormClient) ToSQL(queryFn func(tx *gorm.DB) *gorm.DB) string {
//	return gormDb.GormDb.ToSQL(queryFn)
//}
//
//func (gormDb *GormClient) Association(column string) *gorm.Association {
//	return gormDb.GormDb.Association(column)
//}
//
//func (gormDb *GormClient) Create(value interface{}) *GormClient {
//	gormDb.GormDb.Create(value)
//	return gormDb
//}
//
//func (gormDb *GormClient) CreateInBatches(value interface{}, batchSize int) *GormClient {
//	gormDb.GormDb.CreateInBatches(value, batchSize)
//	return gormDb
//}
//
//func (gormDb *GormClient) Save(value interface{}) *GormClient {
//	gormDb.GormDb.Save(value)
//	return gormDb
//}
//
//func (gormDb *GormClient) First(dest interface{}, conds ...interface{}) *GormClient {
//	gormDb.GormDb.First(dest, conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Take(dest interface{}, conds ...interface{}) *GormClient {
//	gormDb.GormDb.Take(dest, conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Last(dest interface{}, conds ...interface{}) *GormClient {
//	gormDb.GormDb.Last(dest, conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Find(dest interface{}, conds ...interface{}) *GormClient {
//	gormDb.GormDb.Find(dest, conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) FindInBatches(dest interface{}, batchSize int, fc func(tx *GormClient, batch int) error) *GormClient {
//	gormDb.GormDb.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error {
//		return fc(gormDb, batch)
//	})
//	return gormDb
//}
//
//func (gormDb *GormClient) FirstOrInit(dest interface{}, conds ...interface{}) *GormClient {
//	gormDb.GormDb.FirstOrInit(dest, conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) FirstOrCreate(dest interface{}, conds ...interface{}) *GormClient {
//	gormDb.GormDb.FirstOrCreate(dest, conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Update(column string, value interface{}) *GormClient {
//	gormDb.GormDb.Update(column, value)
//	return gormDb
//}
//
//func (gormDb *GormClient) Updates(values interface{}) *GormClient {
//	gormDb.GormDb.Updates(values)
//	return gormDb
//}
//
//func (gormDb *GormClient) UpdateColumn(column string, value interface{}) *GormClient {
//	gormDb.GormDb.UpdateColumn(column, value)
//	return gormDb
//}
//
//func (gormDb *GormClient) UpdateColumns(values interface{}) *GormClient {
//	gormDb.GormDb.UpdateColumns(values)
//	return gormDb
//}
//
//func (gormDb *GormClient) Delete(value interface{}, conds ...interface{}) *GormClient {
//	gormDb.GormDb.Delete(value, conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Count(count *int64) *GormClient {
//	gormDb.GormDb.Count(count)
//	return gormDb
//}
//
//func (gormDb *GormClient) Row() *sql.Row {
//	return gormDb.GormDb.Row()
//}
//
//func (gormDb *GormClient) Rows() (*sql.Rows, error) {
//	return gormDb.GormDb.Rows()
//}
//
//func (gormDb *GormClient) Scan(dest interface{}) *GormClient {
//	gormDb.GormDb.Scan(dest)
//	return gormDb
//}
//
//func (gormDb *GormClient) Pluck(column string, dest interface{}) *GormClient {
//	gormDb.GormDb.Pluck(column, dest)
//	return gormDb
//}
//
//func (gormDb *GormClient) ScanRows(rows *sql.Rows, dest interface{}) error {
//	return gormDb.GormDb.ScanRows(rows, dest)
//}
//
//func (gormDb *GormClient) Connection(fc func(tx *GormClient) error) (err error) {
//	return gormDb.GormDb.Connection(func(tx *gorm.DB) error {
//		return fc(gormDb)
//	})
//}
//
//func (gormDb *GormClient) Transaction(fc func(tx *GormClient) error, opts ...*sql.TxOptions) (err error) {
//	return gormDb.GormDb.Transaction(func(tx *gorm.DB) error {
//		return fc(gormDb)
//	}, opts...)
//}
//
//func (gormDb *GormClient) Begin(opts ...*sql.TxOptions) *GormClient {
//	gormDb.GormDb.Begin(opts...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Commit() *GormClient {
//	gormDb.GormDb.Commit()
//	return gormDb
//}
//
//func (gormDb *GormClient) Rollback() *GormClient {
//	gormDb.GormDb.Rollback()
//	return gormDb
//}
//
//func (gormDb *GormClient) SavePoint(name string) *GormClient {
//	gormDb.GormDb.SavePoint(name)
//	return gormDb
//}
//
//func (gormDb *GormClient) RollbackTo(name string) *GormClient {
//	gormDb.GormDb.RollbackTo(name)
//	return gormDb
//}
//
//func (gormDb *GormClient) Exec(sql string, values ...interface{}) *GormClient {
//	gormDb.GormDb.Exec(sql, values...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Migrator() gorm.Migrator {
//	return gormDb.GormDb.Migrator()
//}
//
//func (gormDb *GormClient) AutoMigrate(dst ...interface{}) error {
//	return gormDb.GormDb.AutoMigrate(dst...)
//}
//
//func (gormDb *GormClient) Model(value interface{}) *GormClient {
//	gormDb.GormDb.Model(value)
//	return gormDb
//}
//
//func (gormDb *GormClient) Clauses(conds ...clause.Expression) *GormClient {
//	gormDb.GormDb.Clauses(conds...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Table(name string, args ...interface{}) *GormClient {
//	gormDb.GormDb.Table(name, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Distinct(args ...interface{}) *GormClient {
//	gormDb.GormDb.Distinct(args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Select(query interface{}, args ...interface{}) *GormClient {
//	gormDb.GormDb.Select(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Omit(columns ...string) *GormClient {
//	gormDb.GormDb.Omit(columns...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Where(query interface{}, args ...interface{}) *GormClient {
//	gormDb.GormDb.Where(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Not(query interface{}, args ...interface{}) *GormClient {
//	gormDb.GormDb.Not(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Or(query interface{}, args ...interface{}) *GormClient {
//	gormDb.GormDb.Or(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Joins(query string, args ...interface{}) *GormClient {
//	gormDb.GormDb.Joins(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) InnerJoins(query string, args ...interface{}) *GormClient {
//	gormDb.GormDb.InnerJoins(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Group(name string) *GormClient {
//	gormDb.GormDb.Group(name)
//	return gormDb
//}
//
//func (gormDb *GormClient) Having(query interface{}, args ...interface{}) *GormClient {
//	gormDb.GormDb.Having(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Order(value interface{}) *GormClient {
//	gormDb.GormDb.Or(value)
//	return gormDb
//}
//
//func (gormDb *GormClient) Limit(limit int) *GormClient {
//	gormDb.GormDb.Limit(limit)
//	return gormDb
//}
//
//func (gormDb *GormClient) Offset(offset int) *GormClient {
//	gormDb.GormDb.Offset(offset)
//	return gormDb
//}
//
//func (gormDb *GormClient) Scopes(funcs ...func(db *GormClient) *GormClient) *GormClient {
//	var funcS []func(db *gorm.DB) *gorm.DB
//	for _, fun := range funcs {
//		funcS = append(funcS, func(db *gorm.DB) *gorm.DB {
//			return fun(gormDb).GormDb
//		})
//	}
//	gormDb.GormDb.Scopes(funcS...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Preload(query string, args ...interface{}) *GormClient {
//	gormDb.GormDb.Preload(query, args...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Attrs(attrs ...interface{}) *GormClient {
//	gormDb.GormDb.Attrs(attrs...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Assign(attrs ...interface{}) *GormClient {
//	gormDb.GormDb.Assign(attrs...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Unscoped() *GormClient {
//	gormDb.GormDb.Unscoped()
//	return gormDb
//}
//
//func (gormDb *GormClient) Raw(sql string, values ...interface{}) *GormClient {
//	gormDb.GormDb.Raw(sql, values...)
//	return gormDb
//}
//
//func (gormDb *GormClient) Result() (int64, error) {
//	return gormDb.GormDb.RowsAffected, gormDb.GormDb.Error
//}
