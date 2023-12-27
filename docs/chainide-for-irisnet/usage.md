---
order: 2
---

# How to Use IRISnet EVM IDE

> The process of writing, compiling, deploying, interacting, and querying with IRISnet EVM smart contracts

## Write a contract

Upon entering the project, the README.md file included in the folder will be automatically previewed.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FQgBGZS91QA5GP2AXMTIl%2Fimage.png?alt=media&token=2f54c253-7113-457e-918c-6d8922c5be95)

In the Explorer panel, you can create new files (or folders), refresh the directory, and download files. You can also directly click on the files that come with the template.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FE2cSwSZdbKLZK5YN7XGK%2Fimage.png?alt=media&token=47ea9493-d238-474e-ab62-17a9d1f83dd3)

Click on a contract file to edit the code.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2Fc3NCif0cmQb6idg1jY4z%2Fimage.png?alt=media&token=a3bfae68-88a5-480d-b296-de50afd1b019)

## Compile a contract

Once your contract code is written, click on the "Compiler" button in the right-side menu to open the compilation module. Choose the compiler version and decide whether to enable optimization, then click "Compile \*\*\*.sol" to initiate the compilation.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FIOilLElJDiaICoy6AUpF%2Fimage.png?alt=media&token=26ef30aa-6746-417e-a634-4970b64bb53d)

After successful compilation, the ABI and BYTE CODE will be displayed below, and there will be a confirmation message in the console stating "Compile contract success."

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FczE48oDmrmEi2hNWpGds%2Fimage.png?alt=media&token=3f51d9b4-063e-48dd-a60e-b13b54bc0781)

## Connect to IRISnet EVM

Before deploying a contract, you need to click on "Connect Wallet" in the upper right corner and select connect to JavaScript VM (for testing, implemented in JavaScript) or Metamask (for deployment on the IRISnet blockchain).

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FLRhasXHh6XnfX6AsxOko%2Fimage.png?alt=media&token=81c389b6-c8b0-43c2-a26e-7f5ee01d60c9)

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2F57zu9JFBQunvzHoVKsvC%2Fimage.png?alt=media&token=5d255844-fec0-4287-a493-44d90c8d88f3)

## Deploy a contract

Click the "Deploy & Interaction" button on the right-hand side, which will bring up the deployment and interaction pages. Select the compiled contract and click "Deploy" to initiate the deployment (then confirm in Metamask). After successful contract deployment, the console will display the contract deployment result and relevant information.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2F6rkcZA6xBoCtRpYVM2cp%2Fimage.png?alt=media&token=a38b25f5-c15d-4621-97d6-0fdce7c138e0)

In addition, you can click "Import Deployed Contract" to import a contract that has already been deployed for contract interactions.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FLrJQS2V3i4xlCzSpqlXJ%2Fimage.png?alt=media&token=d1465e33-034d-414b-a7e9-41b0c069508e)

## Contract interaction

After a successful contract deployment, you can interact with the contract. Click on the deployed contract, choose the corresponding interface, and click "Submit" or "Get" to perform interactions.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FqUmdHgsKhKzqCsCqotgp%2Fimage.png?alt=media&token=38ca1cc9-9c42-497b-b1e4-5e1173eace35)

## Transaction Query

Click on the transaction hash in the Output section to view the specific details of each transaction.

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2FbyVtEZJW5HOzHxtw6mx2%2Fimage.png?alt=media&token=3b456fc1-40c8-4c91-889e-09ba4dd8e950)

## IRISnet EVM Sandbox

If you prefer using the command line for development, you can open the IRISnet EVM Sandbox, which comes pre-loaded with [Hardhat](https://hardhat.org/), [Truffle](https://trufflesuite.com/), [Brownie](https://eth-brownie.readthedocs.io/en/stable/) [Ganache](https://trufflesuite.com/ganache/), [Git](https://git-scm.com/) and [Node.js V16](https://nodejs.org/en).

![img](https://3869740696-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F-MYy-lqJKjq1m0yBAX4r%2Fuploads%2F6WvZX7unWkGkXKV7OfEE%2Fimage.png?alt=media&token=1845bc06-3128-44e4-9ec3-9c3e7e18c44d)

If you've started a process on a port in the Sandbox and wish to access that port, please refer to [ChainIDE - Port Forwarding](https://chainide.gitbook.io/chainide-english-1/port-forwarding).
