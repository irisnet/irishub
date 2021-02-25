# Record

Record 模块用于管理 IRIS Hub 上的记录

## 可用命令

| 名称                                | 描述         |
| ----------------------------------- | ------------ |
| [create](#iris-tx-record-create)    | 创建一条记录 |
| [record](#iris-query-record-record) | 查询记录     |

## iris tx record create

创建一条记录

```bash
iris tx record create [digest] [digest-algo] [flags]
```

**Flags:**

| 名称，速记 | 类型   | 必须 | 默认 | 描述                       |
| ---------- | ------ | ---- | ---- | -------------------------- |
| --uri      | string |      |      | 记录的 uri，比如 ipfs 链接 |
| --meta     | string |      |      | 记录的元数据               |

## iris query record record

通过记录 ID 查询记录

```bash
iris query record record [record-id]
```
