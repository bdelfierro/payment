import 'package:flutter/material.dart';

import 'product.dart';

class CartDetails {
  String cartID;
  String userID;
  String status;
  double  totalAmount;
  int totalCount = 0;
  List<CartItem> lineItems;

  CartDetails(
      {required this.cartID,
      required this.userID,
      required this.status,
      required this.totalAmount,
      required this.totalCount,
      required this.lineItems});

  factory CartDetails.fromJson(Map<String, dynamic> json) {
    return CartDetails(
      cartID: json["cartID"],
      userID: json["userID"],
      status: json['status'],
      totalAmount: (json['totalAmount'] as num).toDouble(),
      totalCount: (json['totalCount'] as num).toInt(),
      lineItems: (json['lineItems'] as List<dynamic>)
          .map((item) => CartItem.fromJson(item as Map<String, dynamic>))
          .toList(),
    );
  }
}

class CartItem {
  late Product productDetails;
  int quantity;
  double totalPrice;

  CartItem(
      {required this.productDetails,
      required this.quantity,
      required this.totalPrice});

  factory CartItem.fromJson(Map<String, dynamic> json) {
    return CartItem(
        productDetails: Product.fromJson(json['productDetails']),
        quantity: (json['quantity'] as num).toInt(),
        totalPrice: (json['totalPrice'] as num).toDouble(),
    );
  }
}
