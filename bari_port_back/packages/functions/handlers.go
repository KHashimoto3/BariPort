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

type GetMessageChatRoom struct {
	Id        string `json:"id"`
	ProjectId string `json:"projectId"`
	Name      string `json:"text"`
	ImgUrl    string `json:"imgUrl"`
	SendAt    string `json:"sendAt"`
}

type GetMessageUser struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	ApiKey      string `json:"apiKey"`
}

type GetMessageProject struct {
	TestUrl string `json:"testUrl"`
}

type GetMessagesForResponse struct {
	Id       string             `json:"id"`
	User     GetMessageUser     `json:"userId"`
	ChatRoom GetMessageChatRoom `json:"chatRoomId"`
	Text     string             `json:"text"`
	ImgUrl   string             `json:"imgUrl"`
	SendAt   string             `json:"sendAt"`
}

type GetMessagesResponse struct {
	Messages []GetMessagesForResponse `json:"messages"`
	TestUrl  string                   `json:"testUrl"`
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

type PostMessageRequest struct {
	Id string `json:"id"`
	UserId string `json:"userId"`
	CompanyId string `json:"companyId"`
	ChatRoomId string `json:"chatRoomId"`
	Text string `json:"text"`
	ImgUrl string `json:"imgUrl"`
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

// メッセージ一覧を返すLambdaハンドラ
func HandlerGetMessages(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	chatRoomId := request.QueryStringParameters["chatRoomId"]

	allMessages, err := GetMessages()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// チャットルームを取得
	chatRoomRes, err := GetChatRoom(chatRoomId)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// メッセージのフィルタリングとユーザー情報の取得
	var messagesInChatRoom []GetMessagesForResponse
	for _, message := range allMessages {
		if message.ChatRoomId == chatRoomId {
			userRes, err := GetUser(message.UserId)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}
			user := GetMessageUser{
				Id:          userRes.Id,
				DisplayName: userRes.DisplayName,
				ApiKey:      userRes.ApiKey,
			}

			messagesInChatRoom = append(messagesInChatRoom, GetMessagesForResponse{
				Id:   message.Id,
				User: user,
				ChatRoom: GetMessageChatRoom{
					Id:        chatRoomRes.Id,
					Name:      chatRoomRes.Name,
					ProjectId: chatRoomRes.ProjectId,
				},
				Text:   message.Text,
				ImgUrl: message.ImgUrl,
				SendAt: message.SendAt,
			})
		}
	}

	// プロジェクト情報の取得
	projectRes, err := GetProject(chatRoomId)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// 最終的なレスポンスの構築
	res := GetMessagesResponse{
		Messages: messagesInChatRoom,
		TestUrl:  projectRes.TestUrl,
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
		Id: messageRequest.Id,
		UserId: messageRequest.UserId,
		CompanyId: messageRequest.CompanyId,
		ChatRoomId: messageRequest.ChatRoomId,
		Text: messageRequest.Text,
		ImgUrl: messageRequest.ImgUrl,
		IsMine: "false",
		SendAt: sendTime,
	}

	_, err = PostMessage(message)

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