package repository

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/abyu/sqlxx/generator"
)

// type Repository interface {
// 	Save(interface{}) (interface{}, error)
// 	FindById(interface{}) (interface{}, error)
// 	FindAll() (interface{}, error)
// }

type RepositoryParams struct {
	RepositoryName   string
	EntityType       string
	EntityScans      string
	EntityIdType     string
	FindById         string
	ModelPackage     string
	DestinatePackage string
}

type Table struct {
	TableName string
	IdColumn  Column
	Fields    []Column
}

type Column struct {
	DbFieldName  string
	ObjFieldName string
	Type         string
}

func convertToTable(typeInfo generator.TypeInfo) *Table {
	var tableName string
	var idColumn Column
	fields := []Column{}

	for _, f := range typeInfo.Fields {

		if table, ok := f.Tags["table"]; ok {
			tableName = table
		}
		if column, ok := f.Tags["column"]; ok {
			col := Column{
				DbFieldName:  column,
				ObjFieldName: f.Name,
			}
			fields = append(fields, col)
		} else if id, ok := f.Tags["id"]; ok {
			idColumn = Column{
				DbFieldName:  id,
				ObjFieldName: f.Name,
				Type:         f.Type,
			}
			fields = append(fields, idColumn)
		} else {
			col := Column{
				DbFieldName:  f.Name,
				ObjFieldName: f.Name,
			}
			fields = append(fields, col)
		}
	}

	return &Table{
		TableName: tableName,
		IdColumn:  idColumn,
		Fields:    fields,
	}
}

func extractTable(entityType reflect.Type) *Table {

	var tableName string
	var idColumn Column
	fields := []Column{}

	for i := 0; i < entityType.NumField(); i++ {
		f := entityType.Field(i)
		if table, ok := f.Tag.Lookup("table"); ok {
			tableName = table
		}
		if column, ok := f.Tag.Lookup("column"); ok {
			col := Column{
				DbFieldName:  column,
				ObjFieldName: f.Name,
			}
			fields = append(fields, col)
		} else if id, ok := f.Tag.Lookup("id"); ok {
			idColumn = Column{
				DbFieldName:  id,
				ObjFieldName: f.Name,
				Type:         f.Type.Name(),
			}
			fields = append(fields, idColumn)
		} else {
			col := Column{
				DbFieldName:  f.Name,
				ObjFieldName: f.Name,
			}
			fields = append(fields, col)
		}
	}

	return &Table{
		TableName: tableName,
		IdColumn:  idColumn,
		Fields:    fields,
	}
}

func GenerateStructReflect(structName string, entity interface{}) {

	s := reflect.ValueOf(entity).Type()
	fmt.Println(s.PkgPath())
	table := extractTable(s)
	data := RepositoryParams{
		ModelPackage:   s.PkgPath(),
		RepositoryName: structName,
		EntityType:     s.Name(),
		EntityScans:    generateScan(*table),
		EntityIdType:   table.IdColumn.Type,
		FindById:       generateById(*table),
	}
	t := template.Must(template.New("repository").Parse(repoTemplate))

	f, err := os.Create("repository/afile.go")

	fmt.Println(err)
	t.Execute(f, data)

}

func GenerateStruct(destPath string, structName string, typeInfo generator.TypeInfo) {
	dir, _ := os.Stat(destPath)
	destPackage := dir.Name()
	table := convertToTable(typeInfo)
	data := RepositoryParams{
		DestinatePackage: destPackage,
		ModelPackage:     typeInfo.PkgPath,
		RepositoryName:   fmt.Sprintf("%sRepository", structName),
		EntityType:       typeInfo.Name,
		EntityScans:      generateScan(*table),
		EntityIdType:     table.IdColumn.Type,
		FindById:         generateById(*table),
	}
	t := template.Must(template.New("repository").Parse(repoTemplate))

	f, err := os.Create(fmt.Sprintf("%s/%s_repository.go", destPath, structName))

	fmt.Println(err)
	t.Execute(f, data)

}
func generateById(table Table) string {

	selectTemplate := "SELECT %s FROM %s %s"
	whereClause := "WHERE %s = ?"

	dbFields := []string{}
	for _, column := range table.Fields {
		dbFields = append(dbFields, column.DbFieldName)
	}

	fieldsQuery := strings.Join(dbFields, ",")
	whereQuery := fmt.Sprintf(whereClause, table.IdColumn.DbFieldName)
	sqlQuery := fmt.Sprintf(selectTemplate, fieldsQuery, table.TableName, whereQuery)

	return sqlQuery
}

func generateScan(table Table) string {
	scanFields := []string{}
	for _, field := range table.Fields {
		scanFields = append(scanFields, fmt.Sprintf("&entity.%s", field.ObjFieldName))
	}

	return strings.Join(scanFields, ",")
}

var repoTemplate = `
package {{.DestinatePackage}}

import (
	"database/sql"

	models "{{.ModelPackage}}"
)

type {{.RepositoryName}} struct {
	DB sql.DB
}

func (repo *{{.RepositoryName}}) Save(entity models.{{.EntityType}}) (*models.{{.EntityType}}, error) {

	
	return nil, nil
}

func (repo *{{.RepositoryName}}) FindById(id {{.EntityIdType}}) (*models.{{.EntityType}}, error) {

	row := repo.DB.QueryRow("{{.FindById}}", id)

	return ScanTo{{.EntityType}}(row)
}

func (repo *{{.RepositoryName}}) FindAll() (*models.{{.EntityType}}, error) {

	return nil, nil
}

func ScanTo{{.EntityType}}(row *sql.Row) (*models.{{.EntityType}}, error) {
	var entity models.{{.EntityType}}
	err := row.Scan({{.EntityScans}})
	
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

`
