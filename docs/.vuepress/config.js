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
          ['/get-started/Install-Iris.md', 'Install'],
          ['/get-started/Full-Node.md', 'Run a Full Node'],
          ['/get-started/Validator-Node.md', 'Run a Validator Node']
        ]
      },{
        title: 'Modules',
        collapsable: false,
        children: [
          // ['/modules/coin/README.md', 'Coin Type'],
          ['/modules/fee-token/README.md', 'Fee Token'],
          // ['/modules/gov/README.md', 'Governance']
        ]
      },{
        title: 'Validators',
        collapsable: false,
        children: [
          ['/validators/README.md', 'Overview'],
          ['/validators/Setup-A-Sentry-Node.md', 'Setup a Sentry Node'],
          ['/validators/FAQ.md', 'FAQ']
        ]
      }
    ]
  }
}
