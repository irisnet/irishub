<!--
order: 3
-->

# Events

The farm module emits the following events:

## Handlers

### MsgCreatePool

| Type        | Attribute Key | Attribute Value |
| :---------- | :------------ | :-------------- |
| create_pool | creator       | {creator}       |
| create_pool | pool_id       | {create_pool}   |
| message     | module        | farm            |
| message     | sender        | {senderAddress} |

### MsgDestroyPool

| Type         | Attribute Key | Attribute Value |
| :----------- | :------------ | :-------------- |
| destroy_pool | creator       | {creator}       |
| destroy_pool | pool_id       | {pool_id}       |
| message      | module        | farm            |
| message      | sender        | {senderAddress} |

### MsgAdjustPool

| Type          | Attribute Key | Attribute Value |
| :------------ | :------------ | :-------------- |
| append_reward | creator       | {creator}       |
| append_reward | pool_id       | {pool_id}       |
| message       | module        | farm            |
| message       | sender        | {senderAddress} |

### MsgStake

| Type    | Attribute Key | Attribute Value |
| :------ | :------------ | :-------------- |
| stake   | creator       | {creator}       |
| stake   | pool_id       | {pool_id}       |
| stake   | amount        | {amount}        |
| stake   | reward        | {reward}        |
| message | module        | farm            |
| message | sender        | {senderAddress} |

### MsgUnstake

| Type    | Attribute Key | Attribute Value |
| :------ | :------------ | :-------------- |
| unstake | creator       | {creator}       |
| unstake | pool_id       | {pool_id}       |
| unstake | amount        | {amount}        |
| unstake | reward        | {reward}        |
| message | module        | farm            |
| message | sender        | {senderAddress} |

### MsgHarvest

| Type    | Attribute Key | Attribute Value |
| :------ | :------------ | :-------------- |
| harvest | creator       | {creator}       |
| harvest | pool_id       | {pool_id}       |
| harvest | reward        | {reward}        |
| message | module        | farm            |
| message | sender        | {senderAddress} |
