package bariport

import (
	"context"
	"encoding/json"
	"time"

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

type GetChatRoomResponse struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	CompanyName         string `json:"companyName"`
	ProjectName         string `json:"projectName"`
	ImgUrl              string `json:"imgUrl"`
	LatestMessage       string `json:"latestMessage"`
	LatestMessageSendAt string `json:"latestMessageSendAt"`
}

type PostChatRoomParticipantRequest struct {
	Id         string `json:"id"`
	ChatRoomId string `json:"chatRoomId"`
	UserId     string `json:"userId"`
}

type PostMessageRequest struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	CompanyId  string `json:"companyId"`
	ChatRoomId string `json:"chatRoomId"`
	Text       string `json:"text"`
	ImgUrl     string `json:"imgUrl"`
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

// チャットルーム一覧を返すLambdaハンドラ
func HandlerChatRooms(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := request.QueryStringParameters["userId"]

	// パラメータで受け取った userId が該当する chatRoomParticipants を取得
	chatRoomParticipants, err := GetChatRoomParticipantsByUserId(userId)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// chatRoomParticipants が所属する chatRoomIds を返す
	var chatRoomIds []string
	for _, participant := range chatRoomParticipants {
		chatRoomIds = append(chatRoomIds, participant.ChatRoomId)
	}

	// chatRoomIds に紐づく chatRooms を取得
	chatRooms, err := GetChatRoomsByChatRoomParticipants(chatRoomIds)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// 以上で取得した chatRooms を起点に必要な API を返す
	var res []GetChatRoomResponse
	for _, chatRoom := range chatRooms {
		companyRes, err := GetCompany(chatRoom.CompanyId)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		projectRes, err := GetProject(chatRoom.ProjectId)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		res = append(res, GetChatRoomResponse{
			Id:                  chatRoom.Id,
			Name:                chatRoom.Name,
			CompanyName:         companyRes.Name,
			ProjectName:         projectRes.ProjectName,
			ImgUrl:              chatRoom.ImgUrl,
			LatestMessage:       chatRoom.LatestMessage,
			LatestMessageSendAt: chatRoom.LatestMessageSendAt,
		})
	}

	body, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

// チャット参加者を送信
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
		Body:       string(body),
		StatusCode: 201,
	}, nil
}

func HandlerPostMessage(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestBody := request.Body

	var messageRequest PostMessageRequest
	err := json.Unmarshal([]byte(requestBody), &messageRequest)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//UTC時間を取得
	localTime := time.Now()
	utcTime := localTime.UTC()
	sendTime := utcTime.Format("2006-01-02T15:04:05Z")

	message := Messages{
		Id:         messageRequest.Id,
		UserId:     messageRequest.UserId,
		CompanyId:  messageRequest.CompanyId,
		ChatRoomId: messageRequest.ChatRoomId,
		Text:       messageRequest.Text,
		ImgUrl:     messageRequest.ImgUrl,
		IsMine:     "false",
		SendAt:     sendTime,
	}

	_, err = PostMessage(message)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body, _ := json.Marshal(SuccessResponse{
		Message: "success",
	})

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 201,
	}, nil

}
