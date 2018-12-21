module.exports = {
    base: "/docs/",
    locales: {
        '/': {
            lang: 'English',
            title: 'IRISnet Document'
        },
        '/zh/': {
            lang: '简体中文',
            title: 'IRISnet 文档'
        }
    },
    themeConfig: {
        displayAllHeaders: false,
        locales: {
            '/': {
                selectText: 'Languages',
                nav: [
                    {text: 'Back to IRISnet', link: 'https://www.irisnet.org'},
                    {text: 'Introduction', link: '/introduction/'},
                    {text: 'Software', link: '/software/node.md'},
                    {text: 'Getting Started', link: '/get-started/'},
                    {text: 'Features', link: '/features/basic-concepts/coin-type.md'},
                    {text: 'CLI Client', link: '/cli-client/'},
                    {text: 'Light Client', link: '/light-client/'},
                    {text: 'Resources', link: '/resources/'}
                ],
                sidebar: {
                    '/software/': [
                        ['node.md', 'Node'],
                        ['cli-client.md', 'CLI Client'],
                        ['light-client.md', 'Light Client'],
                        ['tool.md', 'Tool']
                    ],
                    '/get-started/': [
                        ['Download-Rainbow.md', 'Download Rainbow'],
                        ['Install-the-Software.md', 'Install the Software'],
                        ['Join-the-Testnet.md', 'Join the Testnet']
                    ],
                    '/features/': [{
                        title: 'Basic Concepts',
                        collapsable: false,
                        children: [
                            ["basic-concepts/coin-type.md", 'Coin Type'],
                            ["basic-concepts/fee.md", 'Fee'],
                            ["basic-concepts/inflation.md", 'Infation'],
                            ["basic-concepts/bech32-prefix.md", 'Bech32 Prefix'],
                            ["basic-concepts/genesis-file.md", 'Genesis File'],
                            ["basic-concepts/gov-params.md", 'Gov Params']
                        ]
                    },{
                        title: 'Modules',
                        collapsable: false,
                        children: [
                            ['bank.md', 'Bank'],
                            ['stake.md', 'Stake'],
                            ['service.md', 'Service'],
                            ['record.md', 'Record'],
                            ['governance.md', 'Governance'],
                            ['upgrade.md', 'Upgrade'],
                            ['distribution.md', 'Distribution'],
                            ['guardian.md', 'Guardian'],
                        ]
                    }],
                    '/cli-client/': [{
                        title: 'Status',
                        collapsable: false,
                        children: [
                            ['status/', 'iriscli status']
                        ]
                    },
                    {
                        title: 'Tendermint',
                        collapsable: false,
                        children: [
                            ['tendermint/', 'iriscli tendermint']
                        ]
                    },
                    {
                        title: 'Keys',
                        collapsable: false,
                        children: [
                            ['keys/', 'iriscli keys']
                        ]
                    },
                    {
                        title: 'Bank',
                        collapsable: false,
                        children: [
                            ['bank/', 'iriscli bank']
                        ]
                    },
                    {
                        title: 'Stake',
                        collapsable: false,
                        children: [
                            ['stake/', 'iriscli stake']
                        ]
                    },
                    {
                        title: 'Gov',
                        collapsable: false,
                        children: [
                            ['gov/', 'iriscli gov']
                        ]
                    },
                    {
                        title: 'Record',
                        collapsable: false,
                        children: [
                            ['record/', 'iriscli record']
                        ]
                    },
                    {
                        title: 'Upgrade',
                        collapsable: false,
                        children: [
                            ['upgrade/', 'iriscli upgrade']
                        ]
                    },
                    {
                        title: 'Service',
                        collapsable: false,
                        children: [
                            ['service/', 'iriscli service']
                        ]
                    }],
                    '/resources/': [
                        ['validator-faq.md', 'Validator FAQ'],
                        ['delegator-faq.md', 'Delegator FAQ'],
                        ['whitepaper-zh.md', 'Whitepaper ZH'],
                        ['whitepaper-en.md', 'Whitepaper EN'],
                        ['whitepaper-kr.md', 'Whitepaper KR']
                    ]
                }
            },
            '/zh/': {
                selectText: '选择语言',
                nav: [
                    {text: '返回官网', link: 'https://www.irisnet.org'},
                    {text: '简介', link: '/zh/introduction/'},
                    {text: '软件', link: '/zh/software/node.md'},
                    {text: '开始', link: '/zh/get-started/'},
                    {text: '功能', link: '/zh/features/basic-concepts/coin-type.md'},
                    {text: '命令行', link: '/zh/cli-client/'},
                    {text: '轻客户端', link: '/zh/light-client/'},
                    {text: '资源', link: '/zh/resources/'}
                ],
                sidebar: {
                    '/zh/software/': [
                        ['node.md', 'Node'],
                        ['cli-client.md', 'CLI Client'],
                        ['light-client.md', 'Light Client'],
                        ['tool.md', 'Tool']
                    ],
                    '/zh/get-started/': [
                        ['Download-Rainbow.md', 'Download Rainbow'],
                        ['Install-the-Software.md', 'Install the Software'],
                        ['Join-the-Testnet.md', 'Join the Testnet']
                    ],
                    '/zh/features/': [{
                        title: 'Basic Concepts',
                        collapsable: false,
                        children: [
                            ["basic-concepts/coin-type.md", 'Coin Type'],
                            ["basic-concepts/fee.md", 'Fee'],
                            ["basic-concepts/inflation.md", 'Infation'],
                            ["basic-concepts/bech32-prefix.md", 'Bech32 Prefix'],
                            ["basic-concepts/genesis-file.md", 'Genesis File'],
                            ["basic-concepts/gov-params.md", 'Gov Params']
                        ]
                    },{
                        title: 'Modules',
                        collapsable: false,
                        children: [
                            ['bank.md', 'Bank'],
                            ['stake.md', 'Stake'],
                            ['service.md', 'Service'],
                            ['record.md', 'Record'],
                            ['governance.md', 'Governance'],
                            ['upgrade.md', 'Upgrade'],
                            ['distribution.md', 'Distribution'],
                            ['guardian.md', 'Guardian'],
                        ]
                    }],
                    '/zh/cli-client/': [{
                        title: 'Status',
                        collapsable: false,
                        children: [
                            ['status/', 'iriscli status']
                        ]
                    },
                    {
                        title: 'Tendermint',
                        collapsable: false,
                        children: [
                            ['tendermint/', 'iriscli tendermint']
                        ]
                    },
                    {
                        title: 'Keys',
                        collapsable: false,
                        children: [
                            ['keys/', 'iriscli keys']
                        ]
                    },
                    {
                        title: 'Bank',
                        collapsable: false,
                        children: [
                            ['bank/', 'iriscli bank']
                        ]
                    },
                    {
                        title: 'Stake',
                        collapsable: false,
                        children: [
                            ['stake/', 'iriscli stake']
                        ]
                    },
                    {
                        title: 'Gov',
                        collapsable: false,
                        children: [
                            ['gov/', 'iriscli gov']
                        ]
                    },
                    {
                        title: 'Record',
                        collapsable: false,
                        children: [
                            ['record/', 'iriscli record']
                        ]
                    },
                    {
                        title: 'Upgrade',
                        collapsable: false,
                        children: [
                            ['upgrade/', 'iriscli upgrade']
                        ]
                    },
                    {
                        title: 'Service',
                        collapsable: false,
                        children: [
                            ['service/', 'iriscli service']
                        ]
                    }],
                    '/zh/resources/': [
                        ['validator-faq.md', 'Validator FAQ'],
                        ['delegator-faq.md', 'Delegator FAQ'],
                        ['whitepaper-zh.md', 'Whitepaper ZH'],
                        ['whitepaper-en.md', 'Whitepaper EN'],
                        ['whitepaper-kr.md', 'Whitepaper KR']
                    ]
                }
            }
        }
    }
}
