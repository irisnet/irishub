module.exports = {
    title: 'IRISnet Document',
    description: '',
    base: "/docs/",
    themeConfig: {
        displayAllHeaders: false,
        nav: [
            {text: 'Back to IRISnet', link: 'https://www.irisnet.org'},
            {text: 'Introduction', link: '/introduction/'},
            {text: 'Software', link: '/software/'},
            {text: 'Getting Started', link: '/get-started/'},
            {text: 'Features', link: '/features/'},
            {text: 'CLI Client', link: '/cli-client/'},
            {text: 'Light Client', link: '/light-client/'},
            {text: 'Resources', link: '/resources/'},
        ],
        sidebar: {
            '/introduction/': [{
                title: 'The IRIS Hub',
                collapsable: false,
                children: [
                    'The-IRIS-Hub/Proof-of-Stake.md',
                    'The-IRIS-Hub/IRIS-Tokens.md',
                    'The-IRIS-Hub/Validators.md',
                    'The-IRIS-Hub/Delegators.md'
                ]
            },
                {
                    title: 'The IRIS Service',
                    collapsable: false,
                    children: [
                        'The-IRIS-Service/Lifecycle.md',
                        'The-IRIS-Service/Providers.md',
                        'The-IRIS-Service/Consumers.md',
                    ]
                },
                {
                    title: 'The IRIS Network',
                    collapsable: false,
                    children: [
                        'The-IRIS-Network/'
                    ]
                }],
            '/software/': [
                {
                    title: 'Node',
                    collapsable: false,
                    children: [
                        ['node.md', 'Node']
                    ]
                }, {
                    title: 'CLI Client',
                    collapsable: false,
                    children: [
                        ['cli-client.md', 'CLI Client']
                    ]
                }, {
                    title: 'Light Client',
                    collapsable: false,
                    children: [
                        ['light-client.md', 'Light Client']
                    ]
                }, {
                    title: 'Monitor',
                    collapsable: false,
                    children: [
                        ['monitor.md', 'Monitor']
                    ]
                }],
            '/get-started/': [{
                title: 'Getting Started',
                collapsable: false,
                children: [
                    ['Download-Rainbow.md', 'Download Rainbow'],
                    ['Install-the-Software.md', 'Install the Software'],
                    ['Join-the-Testnet.md', 'Join the Testnet']
                ]
            }],
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
            '/light-client/': [{
                title: 'Light Client',
                collapsable: false,
                children: [
                    ['', 'Light Client']
                ]
            }],
            '/resources/': [{
                title: 'Resources',
                collapsable: false,
                children: [
                    ['validator-faq.md', 'Validator FAQ'],
                    ['delegator-faq.md', 'Delegator FAQ'],
                    ['whitepaper-zh.md', 'Whitepaper ZH'],
                    ['whitepaper-en.md', 'Whitepaper EN'],
                    ['whitepaper-kr.md', 'Whitepaper KR'],
                ]
            }]
        }
    }
}