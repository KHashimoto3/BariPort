import { Api, StackContext, Table } from "sst/constructs";

export function BariPortStack({ stack }: StackContext) {
  //DynamoDB table
  const users_table = new Table(stack, "Users", {
    fields: {
      id: "string",
      displayName: "string",
      apiKey: "string",
    },
    primaryIndex: { partitionKey: "id" },
  });

  const messages_table = new Table(stack, "messages", {
    fields: {
      id: "string",
      userId: "string",
      chatRoomId: "string",
      text: "string",
    },
    primaryIndex: { partitionKey: "id" },
  });

  const chat_room_participants_table = new Table(
    stack,
    "chat_room_participants",
    {
      fields: {
        id: "string",
        userId: "string",
        chatRoomId: "string",
      },
      primaryIndex: { partitionKey: "id" },
    }
  );

  const chat_rooms_table = new Table(stack, "chat_rooms", {
    fields: {
      id: "string",
      name: "string",
      type: "number",
      projectId: "string",
    },
    primaryIndex: { partitionKey: "id" },
  });

  const reviews_table = new Table(stack, "reviews", {
    fields: {
      id: "string",
      companyId: "string",
      userId: "string",
      evaluationScore: "number",
      description: "string",
    },
    primaryIndex: { partitionKey: "id" },
  });

  const companies_table = new Table(stack, "companies", {
    fields: {
      id: "string",
      name: "string",
    },
    primaryIndex: { partitionKey: "id" },
  });

  const projects_table = new Table(stack, "projects", {
    fields: {
      id: "string",
      companyId: "string",
      name: "string",
      description: "string",
    },
    primaryIndex: { partitionKey: "id" },
  });

  //API Gateway
  const api = new Api(stack, "api", {
    defaults: {
      function: {
        bind: [
          users_table,
          messages_table,
          chat_room_participants_table,
          chat_rooms_table,
          reviews_table,
          companies_table,
          projects_table,
        ],
        runtime: "go",
      },
    },
    //API Endpoints
    routes: {
      "GET /hello": "packages/functions/handlers/hello.go",
    },
  });

  stack.addOutputs({
    ApiEndpoint: api.url,
  });
}
