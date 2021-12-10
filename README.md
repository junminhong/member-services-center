## 簡介
基於Golang開發的會員服務中心

## 目標
為了實現single sign on以及加深對於微服務架構的拆分經驗，將整個架構的會員控管都交由member services center去做控管。

## 需求
1. [x] 註冊會員帳戶
2. [x] 會員驗證信件發送
3. [ ] 登入會員帳戶並回傳驗證token
4. [ ] 提供其他服務驗證會員token有效性
5. [ ] CRUD會員資訊

## API Document
- register [POST]: /api/v1/member/register?email=&password=
> requirement:
>> email
>> password
- login [POST]: /api/v1/member/login?email=&password=
> requirement:
>> token
>> email
>> password

## 技術選用
1. Golang
2. Gorm
3. gRPC
4. Gin

## 資料庫選用
1. Postgresql
2. Redis