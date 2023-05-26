// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"admin/model/model"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newTemplate(db *gorm.DB, opts ...gen.DOOption) template {
	_template := template{}

	_template.templateDo.UseDB(db, opts...)
	_template.templateDo.UseModel(&model.Template{})

	tableName := _template.templateDo.TableName()
	_template.ALL = field.NewAsterisk(tableName)
	_template.TemplateID = field.NewInt64(tableName, "template_id")
	_template.ImageName = field.NewString(tableName, "image_name")
	_template.TemplateLabel = field.NewString(tableName, "template_label")
	_template.FileName = field.NewString(tableName, "file_name")

	_template.fillFieldMap()

	return _template
}

type template struct {
	templateDo

	ALL           field.Asterisk
	TemplateID    field.Int64  // 主键id
	ImageName     field.String // 模版对应镜像定位符号
	TemplateLabel field.String // 模板标签
	FileName      field.String // 工程模板文件名

	fieldMap map[string]field.Expr
}

func (t template) Table(newTableName string) *template {
	t.templateDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t template) As(alias string) *template {
	t.templateDo.DO = *(t.templateDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *template) updateTableName(table string) *template {
	t.ALL = field.NewAsterisk(table)
	t.TemplateID = field.NewInt64(table, "template_id")
	t.ImageName = field.NewString(table, "image_name")
	t.TemplateLabel = field.NewString(table, "template_label")
	t.FileName = field.NewString(table, "file_name")

	t.fillFieldMap()

	return t
}

func (t *template) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *template) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 4)
	t.fieldMap["template_id"] = t.TemplateID
	t.fieldMap["image_name"] = t.ImageName
	t.fieldMap["template_label"] = t.TemplateLabel
	t.fieldMap["file_name"] = t.FileName
}

func (t template) clone(db *gorm.DB) template {
	t.templateDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t template) replaceDB(db *gorm.DB) template {
	t.templateDo.ReplaceDB(db)
	return t
}

type templateDo struct{ gen.DO }

type ITemplateDo interface {
	gen.SubQuery
	Debug() ITemplateDo
	WithContext(ctx context.Context) ITemplateDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITemplateDo
	WriteDB() ITemplateDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITemplateDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITemplateDo
	Not(conds ...gen.Condition) ITemplateDo
	Or(conds ...gen.Condition) ITemplateDo
	Select(conds ...field.Expr) ITemplateDo
	Where(conds ...gen.Condition) ITemplateDo
	Order(conds ...field.Expr) ITemplateDo
	Distinct(cols ...field.Expr) ITemplateDo
	Omit(cols ...field.Expr) ITemplateDo
	Join(table schema.Tabler, on ...field.Expr) ITemplateDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITemplateDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITemplateDo
	Group(cols ...field.Expr) ITemplateDo
	Having(conds ...gen.Condition) ITemplateDo
	Limit(limit int) ITemplateDo
	Offset(offset int) ITemplateDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITemplateDo
	Unscoped() ITemplateDo
	Create(values ...*model.Template) error
	CreateInBatches(values []*model.Template, batchSize int) error
	Save(values ...*model.Template) error
	First() (*model.Template, error)
	Take() (*model.Template, error)
	Last() (*model.Template, error)
	Find() ([]*model.Template, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Template, err error)
	FindInBatches(result *[]*model.Template, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Template) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITemplateDo
	Assign(attrs ...field.AssignExpr) ITemplateDo
	Joins(fields ...field.RelationField) ITemplateDo
	Preload(fields ...field.RelationField) ITemplateDo
	FirstOrInit() (*model.Template, error)
	FirstOrCreate() (*model.Template, error)
	FindByPage(offset int, limit int) (result []*model.Template, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITemplateDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t templateDo) Debug() ITemplateDo {
	return t.withDO(t.DO.Debug())
}

func (t templateDo) WithContext(ctx context.Context) ITemplateDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t templateDo) ReadDB() ITemplateDo {
	return t.Clauses(dbresolver.Read)
}

func (t templateDo) WriteDB() ITemplateDo {
	return t.Clauses(dbresolver.Write)
}

func (t templateDo) Session(config *gorm.Session) ITemplateDo {
	return t.withDO(t.DO.Session(config))
}

func (t templateDo) Clauses(conds ...clause.Expression) ITemplateDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t templateDo) Returning(value interface{}, columns ...string) ITemplateDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t templateDo) Not(conds ...gen.Condition) ITemplateDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t templateDo) Or(conds ...gen.Condition) ITemplateDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t templateDo) Select(conds ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t templateDo) Where(conds ...gen.Condition) ITemplateDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t templateDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) ITemplateDo {
	return t.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (t templateDo) Order(conds ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t templateDo) Distinct(cols ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t templateDo) Omit(cols ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t templateDo) Join(table schema.Tabler, on ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t templateDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t templateDo) RightJoin(table schema.Tabler, on ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t templateDo) Group(cols ...field.Expr) ITemplateDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t templateDo) Having(conds ...gen.Condition) ITemplateDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t templateDo) Limit(limit int) ITemplateDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t templateDo) Offset(offset int) ITemplateDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t templateDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITemplateDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t templateDo) Unscoped() ITemplateDo {
	return t.withDO(t.DO.Unscoped())
}

func (t templateDo) Create(values ...*model.Template) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t templateDo) CreateInBatches(values []*model.Template, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t templateDo) Save(values ...*model.Template) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t templateDo) First() (*model.Template, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Template), nil
	}
}

func (t templateDo) Take() (*model.Template, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Template), nil
	}
}

func (t templateDo) Last() (*model.Template, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Template), nil
	}
}

func (t templateDo) Find() ([]*model.Template, error) {
	result, err := t.DO.Find()
	return result.([]*model.Template), err
}

func (t templateDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Template, err error) {
	buf := make([]*model.Template, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t templateDo) FindInBatches(result *[]*model.Template, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t templateDo) Attrs(attrs ...field.AssignExpr) ITemplateDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t templateDo) Assign(attrs ...field.AssignExpr) ITemplateDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t templateDo) Joins(fields ...field.RelationField) ITemplateDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t templateDo) Preload(fields ...field.RelationField) ITemplateDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t templateDo) FirstOrInit() (*model.Template, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Template), nil
	}
}

func (t templateDo) FirstOrCreate() (*model.Template, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Template), nil
	}
}

func (t templateDo) FindByPage(offset int, limit int) (result []*model.Template, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t templateDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t templateDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t templateDo) Delete(models ...*model.Template) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *templateDo) withDO(do gen.Dao) *templateDo {
	t.DO = *do.(*gen.DO)
	return t
}
