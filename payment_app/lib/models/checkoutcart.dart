class CheckoutItem {
  final String productName;
  final String description;
  final String imageUrl;
  final String currency;
  final int quantity;
  final double unitPrice;

  CheckoutItem(
      {required this.productName,
      required this.description,
      required this.imageUrl,
      required this.currency,
      required this.quantity,
      required this.unitPrice});

  Map<String, dynamic> toJson() {
    return {
      "productName": productName,
      "description": description,
      "imageUrl": imageUrl,
      "currency": currency,
      "quantity": quantity,
      "unitPrice": unitPrice,
    };
  }
}

class CheckoutCart {
  final String cartID;
  final List<CheckoutItem> lineItems;

  CheckoutCart({required this.cartID, required this.lineItems});

  Map<String, dynamic> toJson() {
    return {
      "cartID": cartID,
      "lineItems": lineItems.map((item) => item.toJson()).toList(),
    };
  }
}

class CheckoutSession {
  final String cartID;
  final String sessionID;
  final String sessionURL;

  CheckoutSession(
      {required this.cartID,
      required this.sessionID,
      required this.sessionURL});

  Map<String, dynamic> toJson() {
    return {
      "cartID": cartID,
      "sessionID": sessionID,
      "sessionURL": sessionURL,
    };
  }

  factory CheckoutSession.fromJson(Map<String, dynamic> json) {
    return CheckoutSession(
      cartID: json['cartID'],
      sessionID: json['sessionID'],
      sessionURL: json['sessionUrl'],
    );
  }
}
