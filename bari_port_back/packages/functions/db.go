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

// messagesテーブル
type Message struct {
	Id         string `dynamo:"id"`
	UserId     string `dynamo:"userId"`
	CompanyId  string `dynamo:"companyId"`
	ChatRoomId string `dynamo:"chatRoomId"`
	Text       string `dynamo:"text"`
	ImgUrl     string `dynamo:"imgUrl"`
	IsMine     bool   `dynamo:"isMine"`
	SendAt     string `dynamo:"sendAt"`
}

// chatRoomsテーブル
type ChatRoom struct {
	Id        string `dynamo:"id"`
	Name      string `dynamo:"displayName"`
	Type      string `dynamo:"apiKey"`
	ProjectId string `dynamo:"projectId"`
}

// usersテーブル
type User struct {
	Id          string `dynamo:"id"`
	DisplayName string `dynamo:"displayName"`
	ApiKey      string `dynamo:"apiKey"`
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

func GetMessages() ([]Message, error) {
	table := connect(os.Getenv("SST_Table_tableName_messages"))
	var messages []Message
	err := table.Scan().All(&messages)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return messages, errors.New("messages not found")
		}
		return messages, err
	}
	return messages, nil
}

func GetChatRooms() ([]ChatRoom, error) {
	table := connect(os.Getenv("SST_Table_tableName_users"))
	var chatRooms []ChatRoom
	err := table.Scan().All(&chatRooms)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return chatRooms, errors.New("chatRooms not found")
		}
		return chatRooms, err
	}
	return chatRooms, nil
}

func GetUsers() ([]User, error) {
	table := connect(os.Getenv("SST_Table_tableName_users"))
	var users []User
	err := table.Scan().All(&users)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return users, errors.New("users not found")
		}
		return users, err
	}
	return users, nil
}

func GetUser(userId string) (User, error) {
	table := connect(os.Getenv("SST_Table_tableName_users"))
	var user User
	err := table.Get("id", userId).One(&user)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}
