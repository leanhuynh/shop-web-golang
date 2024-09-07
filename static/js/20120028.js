"use strict";

// load constant value
const PORT = 4000;
const HOST = "http://localhost";
const DB_DSN =
  "root:leanhuynh236@tcp(localhost:3306)/go_stripe?parseTime=true&tls=false";
const SMTP_HOST = "smtp.mailtrap.io";
const SMTP_USERNAME = "a486f2e0236b7c";
const SMTP_PASSWORD = "72825a9ab13788";
const SMTP_PORT = 587;
const SECRET_KEY = "secret";
const SESSION_KEY = "session";

// hàm thông báo lỗi
function showNotification(message) {
  // chỉ thông báo khi message khác ""
  if (message != "" && message != undefined && message != null) {
    alert(message);
  }
}

// hàm thêm sản phẩm vào giỏ hàng
function addToCart(id, quantity) {
  // tạo header của request
  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Accept", "application/json");

  // payLoad chứa thông tin request
  const payLoad = {
    ID: id, // product_id
    Quantity: quantity,
  };

  // tạo request
  const requestOptions = {
    method: "POST",
    headers: myHeaders,
    body: JSON.stringify(payLoad),
  };

  // gửi request
  fetch(
    `${HOST}:${PORT}/cart/add`, // tạo url để gửi request
    requestOptions
  ).then((response) => {
    // lấy thông tin status của response
    const status = response?.status;
    /**
     * nếu thêm thành công vào giỏ hàng --> reload lại trang
     * nếu thêm không thành công --> chỉ xuất ra thông báo không thành công
     */
    if (status === 200) {
      showNotification("Thêm sản phẩm vào giỏ hàng thành công");
      location.reload(); // reload lại page
    } else {
      showNotification("Thêm sản phẩm vào giỏ hàng không thành công");
    }
  });
}

// hàm cập nhật số lượng sản phẩm (khi người dùng đang ở giao diện giỏ hàng và cần thay dổi số lượng sản phẩm)
function updateCart(id, quantity) {
  if (quantity > 0) {
    // const token = localStorage.getItem("token");
    // const cart_quantity = document.getElementById("cart-quantity");

    // tạo header của request
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");

    const requestOptions = {
      method: "POST",
      headers: myHeaders,
      body: JSON.stringify({ id, quantity }),
    };

    // gửi yêu cầu cập nhật lại số lượng của sản phẩm trong giỏ hàng
    fetch(`${HOST}:${PORT}/cart/update`, requestOptions).then((response) => {
      const status = response.status;
      console.log(response);
      if (status === 200) {
        // showNotification("Cập nhật số lượng sản phẩm thành công");
        location.reload(); // reload lại page
      } else {
        showNotification("Cập nhật số lượng sản phẩm thất bại");
        location.reload();
      }
    });
  } else {
    removeProductFromCart(id);
  }
}

// hàm xóa sản phẩm ra khỏi giỏ hàng
function removeProductFromCart(id) {
  // cần xác nhận lại người dùng có muốn xóa sản phẩm ra khỏi giỏ hàng
  if (confirm("Bạn có muốn xóa sản phẩm ra khỏi giỏ hàng?")) {
    // tạo header cho request
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");

    // tạo request
    const requestOptions = {
      method: "DELETE",
      headers: myHeaders,
      body: JSON.stringify({ id }),
    };

    // gửi request tới server
    fetch(`${HOST}:${PORT}/cart/remove`, requestOptions).then((response) => {
      const status = response.status; // lấy thông tin status trong response

      if (status === 200) {
        // thực hiện thành công
        showNotification("Xóa sản phẩm ra khỏi giỏ hàng thành công");
        location.reload(); // reload lại page
      } else {
        // thao tác không thể thực hiện
        showNotification("Thao tác hiện không thể thực hiện được");
        location.reload();
      }
    });
  }
}

// hàm xóa giỏ hàng
function clearCart() {
  // xác nhận xóa toàn bộ sản phẩm trong giỏ hàng
  if (confirm("Bạn có muốn xóa toàn bộ giỏ hàng?")) {
    // tạo header cho request
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Accept", "application/json");

    // tạo request
    const requestOptions = {
      method: "DELETE",
      headers: myHeaders,
    };

    // gửi request tới server
    fetch(`${HOST}:${PORT}/cart/delete`, requestOptions).then((response) => {
      const status = response.status; // lấy thông tin status trong response

      if (status === 200) {
        // thành công
        showNotification("Xóa giỏ hàng thành công");
        location.reload(); // reload lại page
      } else {
        // thất bại
        showNotification("Xóa giỏ hàng thất bại");
        location.reload();
      }
    });
  }
}

// hàm lấy originalURL cho các thẻ a
function RemoveParams(key, sourceURL) {
  // split to get the url before query
  var rtn;
  var queryString;

  // get the query url
  if (sourceURL.includes("?")) {
    rtn = sourceURL.split("?")[0];
    queryString = sourceURL.split("?")[1];
  } else {
    rtn = sourceURL;
    return rtn;
  }

  if (queryString != "") {
    // split querys into sub query
    var params_arr = queryString.split("&");

    // delete query related to key
    var param;
    for (let i = params_arr.length - 1; 0 <= i; i--) {
      param = params_arr[i].split("=")[0];
      if (param == key) {
        params_arr = params_arr.slice(0, i);
      }
    }

    // then add query not related to key into url
    if (params_arr.length != 0) {
      rtn += "?" + params_arr.join("&");
    }
  }

  return rtn;
}

// hàm kiểm tra confirmPassword và password có match với nhau hay không
function checkPasswordConfirm(formId) {
  let password = document.querySelector(`#${formId} [name=password]`);
  let confirmPassword = document.querySelector(
    `#${formId} [name=confirmPassword]`
  );

  if (password.value != confirmPassword.value) {
    confirmPassword.setCustomValidity("Passwords not match");
    confirmPassword.reportValidity();
  } else {
    confirmPassword.setCustomValidity("");
  }
}

// hàm đặt hàng
function placeOrder() {
  // lấy shipping address được chọn
  const shippingAddress = document.querySelector(
    'input[name="addressId"]:checked'
  );

  if (shippingAddress === null) {
    showNotification("Vui lòng chọn địa chỉ nhận hàng");
    return;
  }

  /**
   * lấy id của shipping address
   * value rỗng (địa chỉ nhận hàng khác) --> addressId = 0
   */
  const addressId = shippingAddress.value;

  // lấy method được chọn
  const methodPay = document.querySelector('input[name="payment"]:checked');

  if (methodPay === null) {
    showNotification("Vui lòng chọn phương thức thanh toán");
    return;
  }

  // lấy phương thức thanh toán
  const method = methodPay.value;

  // tạo request
  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Accept", "application/json");

  // lấy thông tin trong form nếu người dùng chọn địa chỉ khác
  let firstName, lastName, email, mobile, address, country, city;
  if (addressId === "0") {
    firstName = document.getElementById("first-name").value;
    lastName = document.getElementById("last-name").value;
    email = document.getElementById("e-mail").value;
    mobile = document.getElementById("mobile-no").value;
    address = document.getElementById("address").value;
    country = document.getElementById("country").value;
    city = document.getElementById("city").value;
  }

  console.log(
    `firstName: ${firstName} lastName: ${lastName} email: ${email} mobile: ${mobile} address: ${address} country: ${country} city: ${city}`
  );

  console.log(`id: ${addressId} method: ${method}`);

  // tạo request
  const requestOptions = {
    method: "POST",
    headers: myHeaders,
    body: JSON.stringify({
      addressId,
      method,
      firstName,
      lastName,
      email,
      address,
      country,
      city,
    }),
  };

  // gửi request
  fetch(`${HOST}:${PORT}/user/order`, requestOptions).then((response) => {
    const status = response.status;
    response.json().then((res) => {
      const message = res.message; // lấy message trong response
      if (status === 200) {
        showNotification(message);
        window.location.href = "/cart";
      } else {
        showNotification(message);
      }
    });
  });
}

// hàm tìm kiểm sản phẩm theo keyword
function search(keyword) {
  if (keyword !== "" && keyword !== undefined) {
    location.href = `${HOST}:${PORT}/product?search=${keyword}`;
  }
}
