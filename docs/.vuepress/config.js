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
                    ["Endpoints", "endpoints"],
                    ["Tools", "tools"],
                    ["Migration", "migration"],
                    ["Resources", "resources"]
                ])
            }
        },
    }
};