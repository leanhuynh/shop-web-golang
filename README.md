# Phát triển ứng dụng WEB-SHOPPING

**Chủ đề**

- Ứng dụng mua sắm

**Các tính năng nổi bật**

- Forgot-password

- Authentication

- Tìm kiểm sản phẩm (theo keyword, tag, category, giá, ngày)

- Xác nhận mua hàng (gửi email xác nhận)

**Công nghệ sử dụng**

- Golang

- HTML, CSS

- JWT

- MariaDB

**Cách cài đặt**

- Cài đặt package go

```bash
go get
```

```bash
go mod tidy
```

```bash
go install ./...
```

- Khởi tạo database

```bash
soda migrate up
```

**Hình ảnh của WEB-SHOPPING**

- Giao diện trang chủ
  <img src="./trang_chu.png" width="" height="800">

- chức năng thêm sản phẩm
  <img src="./them_san_pahm.png" width="" height="800">

- chức năng xem giỏ hàng
  <img src="./xem_gio_hang.png" width="" height="800">

- chức năng xem thông tin đơn hàng
  <img src="./thong_tin_don_hang.png" width="" height="800">
