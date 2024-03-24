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
	UserName     string     `json:"userName"`
	Text     string             `json:"text"`
	ImgUrl   string             `json:"imgUrl"`
	SendAt   string             `json:"sendAt"`
	IsMine   string             `json:"isMine"`
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

// メッセージ一覧を返すLambdaハンドラ
func HandlerGetMessages(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	charRoomId := request.QueryStringParameters["chatRoomId"]
	loginUserId := request.QueryStringParameters["loginUserId"]

	//すべてのmessageを取得
	messages, err := GetMessages()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//chatRoomIdに一致するmessageに対してuserIdを取得し、userNameが一致したらisMineをtrueにする
	var res []GetMessagesForResponse
	for _, message := range messages {
		if message.ChatRoomId == charRoomId {
			user, err := GetUser(message.UserId)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}

			isMine := "false"

			if user.Id == loginUserId {
				isMine = "true"
			}

			res = append(res, GetMessagesForResponse{
				Id:       message.Id,
				UserName: user.DisplayName,
				Text:     message.Text,
				ImgUrl:   message.ImgUrl,
				SendAt:   message.SendAt,
				IsMine:   isMine,
			})
		}
	}

	//resの中身を、sendAtが早い順に並び替える
	for i := 0; i < len(res); i++ {
		for j := i + 1; j < len(res); j++ {
			if res[i].SendAt > res[j].SendAt {
				tmp := res[i]
				res[i] = res[j]
				res[j] = tmp
			}
		}
	}

	body, _ := json.Marshal(GetMessagesResponse{
		Messages: res,
		TestUrl:  "https://apps.apple.com/jp/app/testflight/id899247664",
	})

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
func HandlerChatRooms(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	chatRooms, err := GetChatRooms()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

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

	message := Message{
		Id: messageRequest.Id,
		UserId: messageRequest.UserId,
		CompanyId: messageRequest.CompanyId,
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
