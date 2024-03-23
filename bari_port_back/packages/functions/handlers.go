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
	Id		string `json:"id"`
	Name 	string `json:"name"`
}

type GetProjectsResponse struct {
	CompanyName string `json:"companyName"`
	ProjectName 	string `json:"projectName"`
	Description string `json:"description"`
	TestUrl 	string `json:"testUrl"`
	ChatRoomId 	string `json:"chatRoomId"`
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

//プロジェクトの一覧を、企業情報と共に返すLambdaハンドラ
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
			Id: companyRes.Id,
			Name: companyRes.Name,
		}

		res = append(res, GetProjectsResponse{
			CompanyName: company.Name,
			ProjectName: "事業名HOGEHOGE",
			Description: project.Description,
			TestUrl: 	project.TestUrl,
			ChatRoomId: project.ChatRoomId,
		})
	}

	body, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}