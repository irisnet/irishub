#!/bin/bash
rm -rf ~/.iris*
iris init --chain-id=fuxi-develop --name=init1
curl https://raw.githubusercontent.com/irisnet/irishub/develop/testnets/develop/config/genesis.json -o ~/.iris/config/genesis.json
curl https://raw.githubusercontent.com/irisnet/irishub/develop/testnets/develop/config/node_key.json -o ~/.iris/config/node_key.json
curl https://raw.githubusercontent.com/irisnet/irishub/develop/testnets/develop/config/priv_validator.json -o ~/.iris/config/priv_validator.json
SP=$(curl https://raw.githubusercontent.com/irisnet/irishub/develop/testnets/develop/seed_phrase)
command="iriscli keys add init1 --recover"
expect -c "
    spawn $command;
    expect {
        \"override the existing name*\" {send \"y\r\";
        expect {
           \"Enter a passphrase for your key:\" {send \"1234567890\r\";
           expect {
               \"Repeat the passphrase:\" {send \"1234567890\r\"; expect {
                 \"Enter your recovery seed phrase:\" {send \"$SP\r\"; exp_continue}
               }}}
           }}}
        }}

    };

"
iris start