package generator

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/abyu/sqlgogen/db_template"
)

type (
	Table struct {
		TableName string
		IdColumn  Column
		Fields    []Column
	}

	Column struct {
		DbFieldName  string
		ObjFieldName string
		Type         string
	}
)

func GenerateStruct(destPath string, structName string, typeInfo TypeInfo, repoTemplate db_template.RepositoryTemplate) {
	dir, _ := os.Stat(destPath)
	destPackage := dir.Name()
	table := convertToTable(typeInfo)
	data := db_template.TemplateParams{
		DestinatePackage: destPackage,
		ModelPackage:     typeInfo.PkgPath,
		StoreName:        fmt.Sprintf("%sStore", structName),
		EntityType:       typeInfo.Name,
		EntityScans:      generateScan(*table),
		EntityIdType:     table.IdColumn.Type,
		FindById:         generateById(*table),
		FindAll:          generateAll(*table),
	}
	t := template.Must(template.New("repository").Parse(repoTemplate.GetTemplate()))

	f, _ := os.Create(fmt.Sprintf("%s/%s_store.go", destPath, structName))

	t.Execute(f, data)
}

func convertToTable(typeInfo TypeInfo) *Table {
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

func generateAll(table Table) string {

	selectTemplate := "SELECT %s FROM %s"

	dbFields := []string{}
	for _, column := range table.Fields {
		dbFields = append(dbFields, column.DbFieldName)
	}

	fieldsQuery := strings.Join(dbFields, ",")
	sqlQuery := fmt.Sprintf(selectTemplate, fieldsQuery, table.TableName)

	return sqlQuery
}
