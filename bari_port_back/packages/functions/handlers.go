package bariport

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type HelloResponse struct {
	Message string `json:"message"`
}

type GetProjectCompany struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetProjectsResponse struct {
<<<<<<< HEAD
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Company     GetProjectCompany `json:"company"`
}

type GetMessageChatRoom struct {
	Id     string `dynamo:"id"`
	Name   string `dynamo:"text"`
	ImgUrl string `dynamo:"imgUrl"`
	IsMine bool   `dynamo:"isMine"`
	SendAt string `dynamo:"sendAt"`
}

type GetMessageUser struct {
	Id          string `dynamo:"id"`
	DisplayName string `dynamo:"displayName"`
	ApiKey      string `dynamo:"apiKey"`
}

type GetMessagesResponse struct {
	Id       string             `dynamo:"id"`
	User     GetMessageUser     `dynamo:"userId"`
	ChatRoom GetMessageChatRoom `dynamo:"chatRoomId"`
	Text     string             `dynamo:"text"`
	ImgUrl   string             `dynamo:"imgUrl"`
	SendAt   string             `dynamo:"sendAt"`
=======
	CompanyName string `json:"companyName"`
	ProjectName string `json:"projectName"`
	Description string `json:"description"`
	TestUrl     string `json:"testUrl"`
	ChatRoomId  string `json:"chatRoomId"`
}

type GetReviewCompany struct {
	Name string `json:"name"`
}

type GetReviewsResponse struct {
	Id              string `json:"id"`
	SendAt          string `json:"sendAt"`
	ImgUrl          string `json:"imgUrl"`
	CompanyName     string `json:"companyName"`
	EvaluationScore int    `json:"evaluationScore"`
	Description     string `json:"description"`
>>>>>>> fbcabf00f0b295fe419fcd1cde6ca83f9b2225df
}

func HandlerHello(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(HelloResponse{
		Message: "Hello, World!",
	})
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

// プロジェクトの一覧を、企業情報と共に返すLambdaハンドラ
func HandlerGetProjects(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	projects, err := GetProjects()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var res []GetProjectsResponse
	for _, project := range projects {

		companyRes, err := GetCompany(project.CompanyId)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		company := GetProjectCompany{
			Id:   companyRes.Id,
			Name: companyRes.Name,
		}

		res = append(res, GetProjectsResponse{
<<<<<<< HEAD
			Id:          project.Id,
			Name:        project.Name,
			Description: project.Description,
			Company:     company,
=======
			CompanyName: company.Name,
			ProjectName: "事業名HOGEHOGE",
			Description: project.Description,
			TestUrl:     project.TestUrl,
			ChatRoomId:  project.ChatRoomId,
>>>>>>> fbcabf00f0b295fe419fcd1cde6ca83f9b2225df
		})
	}

	body, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

// メッセージ一覧を返すLambdaハンドラ
func HandlerGetMessages(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	chatRoomId := request.QueryStringParameters["chatRoomId"]

	allMessages, err := GetMessages()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var messages []Message
	for _, msg := range allMessages {
		if msg.ChatRoomId == chatRoomId {
			messages = append(messages, msg)
		}
	}

	// チャットルームを取得
	chatRoomRes, err := GetChatRooms(chatRoomId)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	chatRoom := GetMessageChatRoom{
		Id:   chatRoomRes.Id,
		Name: chatRoomRes.Name,
	}

	var res []GetMessagesResponse
	for _, message := range messages {

		// メッセージに紐づくユーザーを取得
		userRes, err := GetUser(message.UserId)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		user := GetMessageUser{
			Id:          userRes.Id,
			DisplayName: userRes.DisplayName,
			ApiKey:      userRes.ApiKey,
		}

		res = append(res, GetMessagesResponse{
			Id:       message.Id,
			User:     user,
			ChatRoom: chatRoom,
			Text:     message.Text,
			ImgUrl:   message.ImgUrl,
			SendAt:   message.SendAt,
		})
	}

	body, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

// 口コミの一覧を、企業情報と共に返すLambdaハンドラ
func HandlerGetReviews(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	reviews, err := GetReviews()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var res []GetReviewsResponse
	for _, review := range reviews {

		companyRes, err := GetCompany(review.CompanyId)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		company := GetReviewCompany{
			Name: companyRes.Name,
		}

		res = append(res, GetReviewsResponse{
			Id:              review.Id,
			SendAt:          review.SendAt,
			ImgUrl:          review.ImgUrl,
			CompanyName:     company.Name,
			EvaluationScore: review.EvaluationScore,
			Description:     review.Description,
		})
	}

	body, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}
