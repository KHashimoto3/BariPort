import { Api, StackContext, Table } from "sst/constructs";

export function ExampleStack({ stack }: StackContext) {
  //DynamoDB table
  const table = new Table(stack, "ExampleTable", {
    fields: {
      id: "string",
      name: "string",
    },
    primaryIndex: { partitionKey: "id" },
  });

  //API Gateway
  const api = new Api(stack, "api", {
    defaults: {
      function: {
        bind: [table],
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
