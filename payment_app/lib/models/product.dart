class Product {
  final String productID;
  final String name;
  final double price;
  final String currency;
  final String description;
  final String imageUrl;

  Product({
    required this.productID,
    required this.name,
    required this.price,
    required this.currency,
    required this.description,
    required this.imageUrl
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
        productID: json["productID"],
        name: json["name"],
        price: json['price'],
        currency: json['currency'],
        description: json['description'],
        imageUrl: json['imageUrl']
    );
  }
}