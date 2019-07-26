package db_template

type RepositoryTemplate interface {
	GetTemplate() string
}

type Mysql struct {
}

type TemplateParams struct {
	RepositoryName   string
	EntityType       string
	EntityScans      string
	EntityIdType     string
	FindById         string
	ModelPackage     string
	DestinatePackage string
}

func NewMysqlTemplate() RepositoryTemplate {
	return Mysql{}
}

func (template Mysql) GetTemplate() string {
	return repoTemplate
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
