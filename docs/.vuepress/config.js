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
                    {text: 'Getting Started', link: '/get-started/'},
                    {text: 'Software', link: '/software/node.md'},
                    {text: 'Features', link: '/features/basic-concepts/coin-type.md'},
                    {text: 'CLI Client', link: '/cli-client/'},
                    {text: 'Light Client', link: '/light-client/'},
                    {text: 'Resources', link: '/resources/'}
                ],
                sidebar: {
                    '/software/': [
                        ['node.md', 'Node'],
                        ['How-to-install-irishub.md', 'Install'],
                        ['cli-client.md', 'CLI Client'],
                        ['light-client.md', 'Light Client'],
                        ['export.md', 'Export'],
                        ['sentry.md', 'Sentry'],
                        ['tool.md', 'Tool'],
                        ['monitor.md', 'Monitor'],
                        ['ledger.md', 'Ledger'],
                        ['kms/kms.md', 'Kms']
                    ],
                    '/get-started/': [
                        ['Join-the-Testnet.md', 'Use the Testnet'],
                        ['Join-the-Mainnet.md', 'Join the Mainnet'],
                        ['Full-Node.md', 'Full Node'],
                        ['Validator-Node.md', 'Validator Node']
                    ],
                    '/features/': [{
                        title: 'Basic Concepts',
                        collapsable: false,
                        children: [
                            ["basic-concepts/coin-type.md", 'Coin Type'],
                            ["basic-concepts/fee.md", 'Fee'],
                            ["basic-concepts/bech32-prefix.md", 'Bech32 Prefix'],
                            ["basic-concepts/genesis-file.md", 'Genesis File'],
                            ["basic-concepts/gov-params.md", 'Gov Params'],
                            ["basic-concepts/key.md", 'Key']
                        ]
                    },{
                        title: 'Modules',
                        collapsable: false,
                        children: [
                            ['bank.md', 'Bank'],
                            ['stake.md', 'Stake'],
                            ['slashing.md', 'Slashing'],
                            ['service.md', 'Service'],
                            ['governance.md', 'Governance'],
                            ['upgrade.md', 'Upgrade'],
                            ['distribution.md', 'Distribution'],
                            ['guardian.md', 'Guardian'],
                            ['mint.md', 'Mint']
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
                        title: 'Distribution',
                        collapsable: false,
                        children: [
                            ['distribution/', 'iriscli distribution']
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
                    {text: '开始', link: '/zh/get-started/'},
                    {text: '软件', link: '/zh/software/node.md'},
                    {text: '功能', link: '/zh/features/basic-concepts/coin-type.md'},
                    {text: '命令行', link: '/zh/cli-client/'},
                    {text: '轻客户端', link: '/zh/light-client/'},
                    {text: '资源', link: '/zh/resources/'}
                ],
                sidebar: {
                    '/zh/software/': [
                        ['node.md', '节点'],
                        ['How-to-install-irishub.md', '安装'],
                        ['cli-client.md', '命令行客户端'],
                        ['light-client.md', '轻节点客户端(LCD)'],
                        ['export.md', '导出区块链状态'],
                        ['sentry.md', '哨兵节点'],
                        ['tool.md', '调试工具'],
                        ['monitor.md', '监控'],
                        ['ledger.md', 'Ledger硬件钱包'],
                        ['kms/kms.md', 'Kms']
                    ],
                    '/zh/get-started/': [
                        ['Join-the-Testnet.md', '使用测试网'],
                        ['Join-the-Mainnet.md', '加入主网'],
                        ['Full-Node.md', '全节点'],
                        ['Validator-Node.md', '验证人节点']
                    ],
                    '/zh/features/': [{
                        title: '基础概念',
                        collapsable: false,
                        children: [
                            ["basic-concepts/coin-type.md", '代币单位'],
                            ["basic-concepts/fee.md", '交易费'],
                            ["basic-concepts/bech32-prefix.md", 'Bech32地址前缀'],
                            ["basic-concepts/genesis-file.md", 'Genesis创世文件'],
                            ["basic-concepts/gov-params.md", '链上治理参数'],
                            ["basic-concepts/key.md", '账户钱包']
                        ]
                    },{
                        title: '模块',
                        collapsable: false,
                        children: [
                            ['bank.md', '转账、查询'],
                            ['stake.md', '委托、股权'],
                            ['slashing.md', '惩罚机制'],
                            ['service.md', 'iService服务'],
                            ['governance.md', '链上治理'],
                            ['upgrade.md', '升级'],
                            ['distribution.md', '收益分配'],
                            ['guardian.md', '系统用户'],
                            ['mint.md', '通胀']
                        ]
                    }],
                    '/zh/cli-client/': [{
                        title: '状态',
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
                        title: '钱包',
                        collapsable: false,
                        children: [
                            ['keys/', 'iriscli keys']
                        ]
                    },
                    {
                        title: '转账、查询',
                        collapsable: false,
                        children: [
                            ['bank/', 'iriscli bank']
                        ]
                    },
                    {
                        title: '委托、股权',
                        collapsable: false,
                        children: [
                            ['stake/', 'iriscli stake']
                        ]
                    },
                    {
                        title: '收益分配',
                        collapsable: false,
                        children: [
                            ['distribution/', 'iriscli distribution']
                        ]
                    },
                    {
                        title: '链上治理',
                        collapsable: false,
                        children: [
                            ['gov/', 'iriscli gov']
                        ]
                    },
                    {
                        title: '升级',
                        collapsable: false,
                        children: [
                            ['upgrade/', 'iriscli upgrade']
                        ]
                    },
                    {
                        title: 'iService服务',
                        collapsable: false,
                        children: [
                            ['service/', 'iriscli service']
                        ]
                    }],
                    '/zh/resources/': [
                        ['validator-faq.md', '验证人 FAQ'],
                        ['delegator-faq.md', '委托人 FAQ'],
                        ['whitepaper-zh.md', '白皮书 ZH'],
                        ['whitepaper-en.md', '白皮书 EN'],
                        ['whitepaper-kr.md', '白皮书 KR']
                    ]
                }
            }
        }
    }
}
