"use strict";

async function loadcartquantity() {
  // get information
  const token = localStorage.getItem("token");
  const cart_quantity = document.getElementById("cart-quantity");

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + token);

  const requestOptions = {
    method: "Get",
    headers: myHeaders,
  };

  const res = await fetch(
    "http://localhost:4001/api/cart",
    requestOptions
  ).then((response) => {
    return response.json();
  });

  console.log(res);

  var quantity = 0;
  if (res.cart != null) {
    for (var cart of res.cart) {
      quantity += parseInt(cart.quantity);
    }
  }
  cart_quantity.textContent = quantity;

  return Promise.resolve();
}

async function loadcart() {
  // get information
  const token = localStorage.getItem("token");
  const cart_quantity = document.getElementById("cart-quantity");

  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + token);

  const requestOptions = {
    method: "Get",
    headers: myHeaders,
  };

  const res = await fetch(
    "http://localhost:4001/api/cart",
    requestOptions
  ).then((response) => {
    return response.json();
  });

  const cartelement = document.getElementById("cart");

  // neu khong co san pham trong gio hang
  if (res.cart == null) {
    cartelement.innerHTML = ` <div class="text-center border py-3">
            <h3>Your cart is empty!</h3>
        </div>`;
    return;
  }

  // neu co san pham trong gio hang
  var html = `        <div class="row">
            <div class="col-md-12">
                <div class="table-responsive">
                    <table class="table table-bordered">
                        <thead class="thead-dark">
                            <tr>
                                <th>Image</th>
                                <th>Name</th>
                                <th>Price</th>
                                <th>Quantity</th>
                                <th>Total</th>
                                <th>Remove</th>
                            </tr>
                        </thead>
                        <tbody class="align-middle">`;

  // xu ly html cua checkout
  var subtotal = 0,
    shipping = 0,
    total = 0;

  for (var order of res.cart) {
    subtotal += order.total;
    html += `     <tr id="product${order.product_id}">
                                <td><a href="/products/${order.product_id}"><img src="${order.image_path}" alt="${order.quantity}"></a></td>
                                <td><a href="/products/${order.product_id}">${order.name}</a></td>
                                <td>${order.price}</td>
                                <td>
                                    <div class="qty">
                                        <button class="btn-minus" onclick="updateCart(${order.product_id}, parseInt(document.getElementById('quantity${order.product_id}').value) - 1).then(loadcartquantity).then(loadcart)"><i class="fa fa-minus"></i></button>
                                        <input type="number" value="${order.quantity}" readonly id="quantity${order.product_id}">
                                        <button class="btn-plus" onclick="updateCart(${order.product_id}, parseInt(document.getElementById('quantity${order.product_id}').value) + 1).then(loadcartquantity).then(loadcart)"><i class="fa fa-plus"></i></button>
                                    </div>
                                </td>
                                <td id="total${order.product_id}">${order.total}</td>
                                <td><button onclick="removeCart(${order.product_id}).then(loadcartquantity).then(loadcart)"><i class="fa fa-trash"></i></button></td>
                            </tr>`;
  }

  // handle subtotal, shipping, total
  total = subtotal;
  html += `</tbody>
                    </table>
                </div>
            </div>
        </div>
        <div class="row">
            <div class="col-md-6">
                <div class="coupon">
                    <input type="text" placeholder="Coupon Code">
                    <button>Apply Code</button>
                </div>
            </div>
            <div class="col-md-6">
                <div class="cart-summary">
                    <div class="cart-content">
                        <h3>Cart Summary</h3>
                        <p>Sub Total<span id="subtotal">${subtotal}</span></p>
                        <p>Shipping Cost<span>${shipping}</span></p>
                        <h4>Grand Total<span id="total">${total}</span></h4>
                    </div>
                    <div class="cart-btn">
                        <button onclick="ClearCart().then(loadcartquantity).then(loadcart)">Clear Cart</button>
                        <button onclick="location.href='/users/checkout'">Checkout</button>
                    </div>
                </div>
            </div>
        </div>`;

  document.getElementById("cart").innerHTML = html;

  return Promise.resolve();
}

async function addCart(id, quantity, callback) {
  if (localStorage.getItem("token") === null) {
    window.location.href = "/login";
    return;
  }

  // create header with token to authenticate
  let token = localStorage.getItem("token");
  const myHeaders = new Headers();
  myHeaders.append("Content-Type", "application/json");
  myHeaders.append("Authorization", "Bearer " + token);

  const payLoad = {
    ID: id,
    Quantity: quantity,
  };

  const requestOptions = {
    method: "put",
    headers: myHeaders,
    body: JSON.stringify(payLoad),
  };

  let status;
  const res = await fetch(
    "http://localhost:4001/api/cart/add",
    requestOptions
  ).then((response) => {
    status = response.status;
    return response.json();
  });

  console.log(res);

  callback();

  // let json = await res.json();
  // document.getElementById("cart-quantity").innerText = `(${json.quantity})`;
}

async function updateCart(id, quantity) {
  if (quantity > 0) {
    const token = localStorage.getItem("token");
    const cart_quantity = document.getElementById("cart-quantity");

    // tao headers
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + token);

    const requestOptions = {
      method: "post",
      headers: myHeaders,
      body: JSON.stringify({ id, quantity }),
    };

    const res = await fetch(
      "http://localhost:4001/api/cart/update",
      requestOptions
    );

    if (res.status == 200) {
      return Promise.resolve();
    }
  } else {
    removeCart(id);
  }
}

async function removeCart(id) {
  if (confirm("Do you really want to remove this product?")) {
    const token = localStorage.getItem("token");
    const cart_quantity = document.getElementById("cart-quantity");

    // tao headers
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + token);

    const requestOptions = {
      method: "delete",
      headers: myHeaders,
      body: JSON.stringify({ id }),
    };

    const res = await fetch(
      "http://localhost:4001/api/cart/remove",
      requestOptions
    );

    if (res.status == 200) {
      return Promise.resolve();
    }
  }
}

async function ClearCart() {
  if (confirm("Do you really want to clear all cart?")) {
    const token = localStorage.getItem("token");
    const cart_quantity = document.getElementById("cart-quantity");

    // tao headers
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", "Bearer " + token);

    const requestOptions = {
      method: "delete",
      headers: myHeaders,
    };

    const res = await fetch(
      "http://localhost:4001/api/cart/delete",
      requestOptions
    );

    if (res.status == 200) {
      return Promise.resolve();
    }
  }
}

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
