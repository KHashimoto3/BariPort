package bariport

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// dynamo db接続関数
func connect(dbName string) dynamo.Table {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
	return db.Table(dbName)
}

// 各テーブルのデータ構造
// companiesテーブル
type Company struct {
	Id   string `dynamo:"id"`
	Name string `dynamo:"name"`
}

// projectsテーブル
type Project struct {
	Id          string `dynamo:"id"`
	CompanyId   string `dynamo:"companyId"`
	Name        string `dynamo:"name"`
	Description string `dynamo:"description"`
}

// reviewsテーブル
type Review struct {
	Id              string `dynamo:"id"`
	CompanyId       string `dynamo:"companyId"`
	UserId          string `dynamo:"userId"`
	EvaluationScore int    `dynamo:"evaluationScore"`
	ImageUrl        string `dynamo:"evaluationScore"`
	Description     string `dynamo:"description"`
	SendAt          string `dynamo:"sendAt"`
}

// companyのデータを取得
func GetCompany(companyId string) (Company, error) {
	table := connect(os.Getenv("SST_Table_tableName_companies"))
	var company Company
	err := table.Get("id", companyId).One(&company)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return company, errors.New("company not found")
		}
		return company, err
	}
	return company, nil
}

func GetProjects() ([]Project, error) {
	table := connect(os.Getenv("SST_Table_tableName_projects"))
	var projects []Project
	err := table.Scan().All(&projects)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return projects, errors.New("projects not found")
		}
		return projects, err
	}
	return projects, nil
}

func GetReviews() ([]Review, error) {
	table := connect(os.Getenv("SST_Table_tableName_reviews"))
	var reviews []Review
	err := table.Scan().All(&reviews)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return reviews, errors.New("reviews not found")
		}
		return reviews, err
	}
	return reviews, nil
}
