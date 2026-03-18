# SQL 执行说明

## 系统（sys）
- `001_rbac_schema_and_seed.sql`

## 业务（tk_）
- `business/001_business_schema.sql`：业务建表 + 历史结构兼容升级
- `business/002_business_data.sql`：业务初始化数据 + 可选 `w_* -> tk_*` 迁移

## 推荐执行顺序
1. `001_rbac_schema_and_seed.sql`
2. `business/001_business_schema.sql`
3. `business/002_business_data.sql`
