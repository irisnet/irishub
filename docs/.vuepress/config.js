import glob from "glob";
import markdownIt from "markdown-it";
import meta from "markdown-it-meta";
import fs from "fs";
import _ from "lodash";
import { searchPlugin } from "@vuepress/plugin-search";
const sidebar = (directory, array) => {
  return array.map((i) => {
    const children = _.sortBy(
      glob
        .sync(`./${directory}/${i[1]}/*.md`)
        .map((path) => {
          const md = new markdownIt();
          const file = fs.readFileSync(path, "utf8");
          md.use(meta);
          md.render(file);
          const order = md.meta.order;
          return { path, order };
        })
        .filter((f) => f.order !== false),
      ["order", "path"]
    )
      .map((f) => f.path.slice(1))
      .filter((f) => !f.match("README"));
    return {
      text: i[0],
      children,
      collapsible: true,
    };
  });
};
import { defineUserConfig } from "vuepress";
import { defaultTheme } from "vuepress";
export default defineUserConfig({
  base: "/docs/",
  plugins: [
    searchPlugin({
      locales: {
        "/": {
          placeholder: "Search",
        },
      },
    }),
  ],
  locales: {
    "/": {
      lang: "en-US",
      title: "IRISnet Documents",
      description: "IRISnet Documents",
    },
  },
  theme: defaultTheme({
    repo: "irisnet/irishub",
    docsDir: "docs",
    editLinks: true,
    contributors: false,
    docsBranch: "main",
    editLinkText: "Help us improve this page!",
    locales: {
      "/": {
        editLinkText: "Help us improve this page!",
        navbar: [
          {
            text: "Back to IRISnet",
            link: "https://www.irisnet.org",
          },
        ],
        sidebar: sidebar("", [
          ["Getting Started", "get-started"],
          ["Concepts", "concepts"],
          ["Features", "features"],
          ["Daemon", "daemon"],
          ["CLI Client", "cli-client"],
          ["Endpoints", "endpoints"],
          ["ChainIDE for IRISnet", "chainide-for-irisnet"],
          ["Tools", "tools"],
          ["Migration", "migration"],
          ["Resources", "resources"],
        ]),
      },
    },
  }),
});
