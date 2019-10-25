const glob = require("glob");
const markdownIt = require("markdown-it");
const meta = require("markdown-it-meta");
const fs = require("fs");
const _ = require("lodash");

const sidebar = (directory, array) => {
    return array.map(i => {
        const children = _.sortBy(
            glob
                .sync(`./${directory}/${i[1]}/*.md`)
                .map(path => {
                    const md = new markdownIt();
                    const file = fs.readFileSync(path, "utf8");
                    md.use(meta);
                    md.render(file);
                    const order = md.meta.order;
                    return { path, order };
                })
                .filter(f => f.order !== false),
            ["order", "path"]
        )
            .map(f => f.path)
            .filter(f => !f.match("README"));

        return {
            title: i[0],
            children
        };
    });
};

module.exports = {
    base: "/docs/",
    plugins: [
        ['@vuepress/search', {
            searchMaxSuggestions: 10
        }]
    ],
    locales: {
        "/": {
            lang: "en-US",
            title: "IRISnet Documents",
            description: "IRISnet Documents",
        },
        "/zh/": {
            lang: "简体中文",
            title: "IRISnet 文档",
            description: "IRISnet 文档",
        }
    },
    themeConfig: {
        repo: "irisnet/irishub",
        docsDir: "docs",
        editLinks: true,
        docsBranch: "master",
        editLinkText: 'Help us improve this page!',
        locales: {
            "/": {
                selectText: 'Languages',
                label: 'English',
                editLinkText: 'Help us improve this page!',
                nav: [
                    {
                        text: 'Back to IRISnet',
                        link: 'https://www.irisnet.org'
                    }
                ],
                sidebar: sidebar("", [
                    ["Getting Started", "get-started"],
                    ["Concepts", "concepts"],
                    ["Features", "features"],
                    ["Daemon", "daemon"],
                    ["CLI Client", "cli-client"],
                    ["API Server", "light-client"],
                    ["Tools", "tools"],
                    ["Resources", "resources"]
                ])
            },
            "/zh/": {
                selectText: '选择语言',
                label: '简体中文',
                editLinkText: '帮助我们完善此文档',
                nav: [{
                    text: 'IRISnet 官网',
                    link: 'https://www.irisnet.org'
                }],
                sidebar: sidebar("", [
                    ["快速开始", "/zh/get-started"],
                    ["概念", "/zh/concepts"],
                    ["功能模块", "/zh/features"],
                    ["守护进程", "/zh/daemon"],
                    ["命令行客户端", "/zh/cli-client"],
                    ["API 服务", "/zh/light-client"],
                    ["工具", "/zh/tools"],
                    ["资源", "/zh/resources"]
                ])
            }
        },
    }
};