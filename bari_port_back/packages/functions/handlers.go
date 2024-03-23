package bariport

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type HelloResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type GetProjectCompany struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetProjectsResponse struct {
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
}

type PostChatRoomParticipantRequest struct {
	Id 	 string `json:"id"`
	ChatRoomId string `json:"chatRoomId"`
	UserId	 string `json:"userId"`
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
			CompanyName: company.Name,
			ProjectName: "事業名HOGEHOGE",
			Description: project.Description,
			TestUrl:     project.TestUrl,
			ChatRoomId:  project.ChatRoomId,
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

//チャット参加者を送信
func HandlerPutChatRoomParticipant(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestBody := request.Body

	var chatRoomParticipant ChatRoomParticipant

	err := json.Unmarshal([]byte(requestBody), &chatRoomParticipant)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	_, err = PostChatRoomParticipant(chatRoomParticipant)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body, _ := json.Marshal(SuccessResponse{
		Message: "success",
	})

	return events.APIGatewayProxyResponse{
		Body: string(body),
		StatusCode: 201,
	}, nil
}