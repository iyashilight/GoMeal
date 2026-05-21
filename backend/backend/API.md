# 外卖订餐系统 API 文档

## 基本信息

- **Base URL**: `http://localhost:8080/api`
- **Content-Type**: `application/json`
- **认证方式**: JWT Token，通过 `Authorization: Bearer <token>` 请求头传递

## 通用响应格式

### 成功响应

```json
{
    "code": 0,
    "message": "success",
    "data": { ... }
}
```

### 错误响应

```json
{
    "code": 400,
    "message": "错误描述",
    "data": null
}
```

| code | 含义 |
|------|------|
| 0    | 成功 |
| 400  | 请求参数错误 |
| 401  | 未认证 |
| 404  | 资源不存在 |
| 500  | 服务器内部错误 |

---

## 一、用户模块

### 1.1 用户注册

```
POST /register
```

**请求体：**

```json
{
    "phone": "13800138000",
    "password": "123456"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| password | string | 是 | 密码，至少6位 |

**响应数据：**

```json
{
    "user": {
        "id": 1,
        "phone": "13800138000",
        "nickname": "用户8000",
        "avatar": "",
        "user_type": 0
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| user.id | uint | 用户ID |
| user.phone | string | 手机号 |
| user.nickname | string | 昵称 |
| user.avatar | string | 头像URL |
| user.user_type | int | 0-普通用户，1-商家 |
| token | string | JWT令牌 |

---

### 1.2 用户登录

```
POST /login
```

**请求体：**

```json
{
    "phone": "13800138000",
    "password": "123456"
}
```

**响应数据：** 同注册

---

### 1.3 获取用户信息

```
GET /user/info
```

**请求头：** `Authorization: Bearer <token>`

**响应数据：**

```json
{
    "id": 1,
    "phone": "13800138000",
    "nickname": "用户8000",
    "avatar": "",
    "user_type": 0
}
```

---

## 二、商家模块

### 2.1 获取商家列表

```
GET /merchants
```

**响应数据：**

```json
[
    {
        "id": 1,
        "name": "必胜客（国贸店）",
        "logo": "https://via.placeholder.com/100x100/ff6b6b/fff?text=Pizza",
        "notice": "满50减10，满100减25",
        "phone": "400-123-4567",
        "address": "国贸大厦1层",
        "min_price": 20,
        "delivery_fee": 5,
        "status": 1,
        "rating": 4.8,
        "sales": 1234
    }
]
```

| 字段 | 类型 | 说明 |
|------|------|------|
| min_price | float | 起送价 |
| delivery_fee | float | 配送费 |
| status | int | 0-休息中，1-营业中 |
| rating | float | 评分（0-5） |
| sales | int | 月销量 |

---

### 2.2 获取商家详情

```
GET /merchants/:id
```

**响应数据：**

```json
{
    "id": 1,
    "name": "必胜客（国贸店）",
    "logo": "https://...",
    "notice": "满50减10",
    "phone": "400-123-4567",
    "address": "国贸大厦1层",
    "min_price": 20,
    "delivery_fee": 5,
    "status": 1,
    "rating": 4.8,
    "sales": 1234,
    "categories": [
        {
            "id": 1,
            "name": "热销推荐",
            "foods": [
                {
                    "id": 1,
                    "name": "超级至尊披萨",
                    "description": "经典口味，芝士浓郁",
                    "image": "https://...",
                    "price": 59,
                    "old_price": 79,
                    "stock": 9999,
                    "sales": 520
                }
            ]
        }
    ]
}
```

---

### 2.3 获取商品详情

```
GET /foods/:id
```

**响应数据：**

```json
{
    "id": 1,
    "name": "超级至尊披萨",
    "description": "经典口味，芝士浓郁",
    "image": "https://...",
    "price": 59,
    "old_price": 79,
    "stock": 9999,
    "sales": 520
}
```

---

## 三、购物车模块

> 以下接口需要登录（请求头携带 `Authorization: Bearer <token>`）

### 3.1 获取购物车

```
GET /cart
```

**响应数据：**

```json
{
    "merchant_id": 1,
    "merchant_name": "必胜客（国贸店）",
    "delivery_fee": 5,
    "min_price": 20,
    "items": [
        {
            "id": 1,
            "food_id": 1,
            "food_name": "超级至尊披萨",
            "food_image": "https://...",
            "price": 59,
            "quantity": 2
        }
    ],
    "total_amount": 118,
    "total_quantity": 2
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| total_amount | float | 总计金额 |
| total_quantity | int | 总计数量 |

---

### 3.2 添加到购物车

```
POST /cart
```

**请求体：**

```json
{
    "merchant_id": 1,
    "food_id": 1,
    "quantity": 1
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| merchant_id | uint | 是 | 商家ID |
| food_id | uint | 是 | 商品ID |
| quantity | int | 是 | 数量，至少1 |

**说明：** 购物车只允许同一家商家的商品。如果添加不同商家的商品，会清空原购物车。

**响应数据：** `null`

---

### 3.3 更新购物车项数量

```
PUT /cart/:id
```

**请求体：**

```json
{
    "quantity": 3
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| quantity | int | 是 | 新数量 |

**响应数据：** `null`

---

### 3.4 删除购物车项

```
DELETE /cart/:id
```

**响应数据：** `null`

---

### 3.5 清空购物车

```
DELETE /cart
```

**响应数据：** `null`

---

## 四、地址模块

> 以下接口需要登录

### 4.1 获取地址列表

```
GET /addresses
```

**响应数据：**

```json
[
    {
        "id": 1,
        "name": "张三",
        "phone": "13800138000",
        "address": "北京市朝阳区...",
        "tag": "家",
        "is_default": true
    }
]
```

---

### 4.2 创建地址

```
POST /addresses
```

**请求体：**

```json
{
    "name": "张三",
    "phone": "13800138000",
    "address": "北京市朝阳区...",
    "tag": "家",
    "is_default": true
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 联系人 |
| phone | string | 是 | 电话 |
| address | string | 是 | 详细地址 |
| tag | string | 否 | 标签（家/公司/学校） |
| is_default | bool | 否 | 是否设为默认 |

**响应数据：** `null`

---

### 4.3 获取地址详情

```
GET /addresses/:id
```

**响应数据：**

```json
{
    "id": 1,
    "name": "张三",
    "phone": "13800138000",
    "address": "北京市朝阳区...",
    "tag": "家",
    "is_default": true
}
```

---

### 4.4 更新地址

```
PUT /addresses/:id
```

**请求体：** 同创建地址

**响应数据：** `null`

---

### 4.5 删除地址

```
DELETE /addresses/:id
```

**响应数据：** `null`

---

## 五、订单模块

> 以下接口需要登录

### 5.1 创建订单

```
POST /orders
```

**请求体：**

```json
{
    "address_id": 1,
    "remark": "少放辣"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| address_id | uint | 是 | 收货地址ID |
| remark | string | 否 | 备注 |

**说明：** 订单从当前购物车创建，创建成功后自动清空购物车。

**响应数据：**

```json
{
    "id": 1,
    "order_no": "2025051412345000001",
    "merchant_id": 1,
    "merchant": { ... },
    "address": { ... },
    "items": [
        {
            "food_id": 1,
            "food_name": "超级至尊披萨",
            "food_image": "https://...",
            "price": 59,
            "quantity": 2
        }
    ],
    "total_amount": 118,
    "delivery_fee": 5,
    "remark": "少放辣",
    "status": 0,
    "status_text": "待支付",
    "created_at": "2025-05-14T12:34:50+08:00"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| order_no | string | 订单号 |
| status | int | 订单状态码 |
| status_text | string | 订单状态中文 |

**订单状态：**

| status | status_text |
|--------|-------------|
| 0 | 待支付 |
| 1 | 已支付 |
| 2 | 已接单 |
| 3 | 配送中 |
| 4 | 已完成 |
| 5 | 已取消 |
| 6 | 已退款 |

---

### 5.2 获取订单列表

```
GET /orders?status=0&page=1&size=10
```

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | int | 否 | 筛选状态，不传查全部 |
| page | int | 否 | 页码，默认1 |
| size | int | 否 | 每页条数，默认10 |

**响应数据：**

```json
[
    {
        "id": 1,
        "order_no": "2025051412345000001",
        "merchant_id": 1,
        "merchant": { ... },
        "items": [ ... ],
        "total_amount": 118,
        "delivery_fee": 5,
        "status": 0,
        "status_text": "待支付",
        "created_at": "2025-05-14T12:34:50+08:00"
    }
]
```

---

### 5.3 获取订单详情

```
GET /orders/:id
```

**响应数据：** 同创建订单的响应

---

### 5.4 支付订单

```
POST /orders/:id/pay
```

**说明：** 仅限待支付状态的订单。模拟支付，成功后更新商家销量。

**响应数据：** `null`

---

### 5.5 取消订单

```
POST /orders/:id/cancel
```

**说明：** 仅限待支付状态的订单。

**响应数据：** `null`

---

## 六、初始化数据

### 6.1 初始化测试数据

```
POST /init
```

**说明：** 创建测试商家和商品。如果已有数据，则跳过。

**响应数据：**

```json
{
    "code": 0,
    "message": "初始化成功",
    "data": null
}
```
