import { SSTConfig } from "sst";
import { ExampleStack } from "./stacks/ExampleStack";

export default {
  config(_input) {
    return {
      name: "bari-port-back",
      region: "ap-northeast-1", // Tokyo
    };
  },
  stacks(app) {
    app.stack(ExampleStack);
  },
} satisfies SSTConfig;
