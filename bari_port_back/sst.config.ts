import { SSTConfig } from "sst";
import { BariPortStack } from "./stacks/BariPortStack";

export default {
  config(_input) {
    return {
      name: "bari-port-back-prod",
      region: "ap-northeast-1", // Tokyo
    };
  },
  stacks(app) {
    app.stack(BariPortStack);
  },
} satisfies SSTConfig;
