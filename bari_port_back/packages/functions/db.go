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
	ProjectName string `dynamo:"name"`
	Description string `dynamo:"description"`
	TestUrl     string `dynamo:"testUrl"`
	ChatRoomId  string `dynamo:"chatRoomId"`
}

// messagesテーブル
type Message struct {
	Id         string `dynamo:"id"`
	UserId     string `dynamo:"userId"`
	CompanyId  string `dynamo:"companyId"`
	ChatRoomId string `dynamo:"chatRoomId"`
	Text       string `dynamo:"text"`
	ImgUrl     string `dynamo:"imgUrl"`
	IsMine     string `dynamo:"isMine"`
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

// reviewsテーブル
type Review struct {
	Id              string `dynamo:"id"`
	CompanyId       string `dynamo:"companyId"`
	UserId          string `dynamo:"userId"`
	EvaluationScore int    `dynamo:"evaluationScore"`
	ImgUrl          string `dynamo:"imgUrl"`
	Description     string `dynamo:"description"`
	SendAt          string `dynamo:"sendAt"`
}

// chat_room_participantsテーブル
type ChatRoomParticipant struct {
	Id         string `dynamo:"id"`
	ChatRoomId string `dynamo:"chatRoomId"`
	UserId     string `dynamo:"userId"`
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

func GetProject(projectId string) (Project, error) {
	table := connect(os.Getenv("SST_Table_tableName_projects"))
	var project Project
	err := table.Get("id", projectId).One(&project)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return project, errors.New("project not found")
		}
		return project, err
	}
	return project, nil
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

func GetChatRoom(chatRoomId string) (ChatRoom, error) {
	table := connect(os.Getenv("SST_Table_tableName_chatRooms"))
	var chatRoom ChatRoom
	err := table.Get("id", chatRoomId).One(&chatRoom)
	if err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return chatRoom, errors.New("chatRoom not found")
		}
		return chatRoom, err
	}
	return chatRoom, nil
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

func PostChatRoomParticipant(chatRoomParticipant ChatRoomParticipant) (bool, error) {
	table := connect(os.Getenv("SST_Table_tableName_chat_room_participants"))

	err := table.Put(chatRoomParticipant).Run()
	if err != nil {
		return false, err
	}

	return true, nil
}

func PostMessage(message Message) (bool, error) {
	table := connect(os.Getenv("SST_Table_tableName_messages"))

	err := table.Put(message).Run()
	if err != nil {
		return false, err
	}

	return true, nil
}
