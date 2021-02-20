<!--
order: 3
-->

# Events

## BeginBlocker

| Type            | Attribute Key      | Attribute Value |
| :-------------- | :----------------- | :-------------- |
| generate_random | request_id         | {requestID}     |
| generate_random | random             | {random}        |
| request_service | request_id         | {hashLock}      |
| request_service | request_context_id | {hashLock}      |

## Handlers

### MsgRequestRandom

| Type           | Attribute Key   | Attribute Value   |
| :------------- | :-------------- | :---------------- |
| request_random | request_id      | {requestID}       |
| request_random | consumer        | {consumerAddress} |
| request_random | generate_height | {generateHeight}  |
| request_random | oracle          | {useOracleOrNot}  |
| message        | module          | random            |
| message        | sender          | {senderAddress}   |

## Callbacks

| Type            | Attribute Key | Attribute Value |
| :-------------- | :------------ | :-------------- |
| generate_random | request_id    | {requestID}     |
| generate_random | random        | {random}        |
