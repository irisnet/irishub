module.exports = {
  title: 'IRISnet Document',
  description: '',
  base: "/docs/",
  themeConfig: {
    displayAllHeaders: false,
    nav: [
      { text: 'Back to IRISnet', link: 'https://www.irisnet.org' }
    ],
    sidebar: [
      {
        title: 'Introduction',
        collapsable: false,
        children: [
          ['/introduction/Whitepaper.md', 'Whitepaper - English'],
          ['/introduction/Whitepaper_CN.md', 'Whitepaper - 中文']
        ]
      },{
        title: 'Getting Started',
        collapsable: false,
        children: [
          ['/get-started/', 'Join the Testnet'],
          ['/get-started/Install-Iris.md', 'Install'],
          ['/get-started/Full-Node.md', 'Run a Full Node'],
          ['/get-started/Validator-Node.md', 'Run a Validator Node'],
          ['/get-started/Genesis-Generation-Process.md', 'Genesis Generation'],
          ['/get-started/Bech32-on-IRISnet.md', 'Bech32 on IRISnet'],
        ]
      },{
        title: 'Modules',
        collapsable: false,
        children: [
          // ['/modules/coin/README.md', 'Coin Type'],
          ['/modules/fee-token/', 'Fee Token']
          // ['/modules/gov/README.md', 'Governance']
        ]
      },{
        title: 'Tools',
        collapsable: false,
        children: [
          ['/tools/Deploy-IRIS-Monitor.md', 'Monitor']
        ]
      },{
        title: 'Validators',
        collapsable: false,
        children: [
          ['/validators/', 'Overview'],
          ['/validators/Setup-Sentry-Node.md', 'Sentry Node'],
          ['/validators/How-to-participate-in-onchain-governance.md', 'Onchain Governance'],
          ['/validators/FAQ.md', 'FAQ']
        ]
      }
    ]
  }
}

