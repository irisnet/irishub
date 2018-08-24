module.exports = {
  title: 'IRIShub Document',
  description: '',
  base: "/",
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
          ['/get-started/install-iris.md', 'Install'],
          ['/get-started/full-node.md', 'Run a Full Node'],
          ['/get-started/validator-node.md', 'Run a Validator Node']
        ]
      },{
        title: 'Modules',
        collapsable: false,
        children: [
          ['/modules/fee-token/feeToken.md', 'Fee Token'],
          ['/modules/gov/gov_spec.md', 'Governance'],
          ['/modules/software-upgrade/software-upgrade.md', 'Software Upgrade']
        ]
      },{
        title: 'Validators',
        collapsable: false,
        children: [
          ['/validators/overview.md', 'Overview'],
          ['/validators/Setup A Sentry Node.md', 'Setup a Sentry Node']          
        ]
      }
    ]
  }
}
