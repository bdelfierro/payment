import 'package:flutter/material.dart';
import 'models/cartdetails.dart';
import 'models/product.dart';

class CartDetailsNotifier extends ChangeNotifier {

  late CartDetails _cartDetails;

  CartDetailsNotifier() {
    _cartDetails = CartDetails(
      cartID: "",
      userID: "",
      status: "",
      totalAmount: 0.0,
      totalCount: 0,
      lineItems: <CartItem>[],
    );
  }

  set cartDetails(CartDetails newValue) {
    _cartDetails =  newValue;
    notifyListeners();
  }

  void init(){
    _cartDetails = CartDetails(
      cartID: "",
      userID: "",
      status: "",
      totalAmount: 0.0,
      totalCount: 0,
      lineItems: <CartItem>[],
    );
  }

  void addItem(Product product) {

    var newProduct = true;
    _cartDetails.lineItems.asMap().forEach((i, item) {
      if (item.productDetails.productID == product.productID){
        _cartDetails.lineItems[i].quantity += 1;
        _cartDetails.lineItems[i].totalPrice += _cartDetails.lineItems[i].productDetails.price;
        _cartDetails.totalCount++;
        _cartDetails.totalAmount += _cartDetails.lineItems[i].productDetails.price;
        newProduct = false;
      }
    });
    if (newProduct) {
      var newItem = CartItem(productDetails: product, quantity: 1, totalPrice: product.price);
      _cartDetails.lineItems.add(newItem);
    }
    notifyListeners();
  }

  void removeItem(String productID) {
    _cartDetails.lineItems.asMap().forEach((i, item) {
      if (item.productDetails.productID == productID){
        if (_cartDetails.lineItems[i].quantity > 0) {
          _cartDetails.lineItems[i].quantity -= 1;
          _cartDetails.lineItems[i].totalPrice -= _cartDetails.lineItems[i].productDetails.price;
          _cartDetails.totalCount--;
          _cartDetails.totalAmount -= _cartDetails.lineItems[i].productDetails.price;
          notifyListeners();
        }
      }
    });
  }

  CartDetails getCartDetails() {
    return _cartDetails;
  }

}
