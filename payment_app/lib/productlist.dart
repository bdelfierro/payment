import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'providers.dart';
import 'package:provider/provider.dart';
import 'dart:convert';
import 'models/cartdetails.dart';
import 'models/product.dart';
import 'package:intl/intl.dart';
import 'cartItems.dart';

class ProductList extends StatelessWidget {
  final Future<List<Product>> products;

  const ProductList({super.key, required this.products});

  @override
  Widget build(BuildContext context) {
    print('refresh ProductList');
    final ThemeData theme = Theme.of(context);
    final NumberFormat formatter = NumberFormat.simpleCurrency(
        locale: Localizations.localeOf(context).toString());

    return FutureBuilder<List<Product>>(
      future: products,
      builder: (BuildContext context, AsyncSnapshot<List<Product>> snapshot) {
        if (snapshot.hasData) {
          return ListView.builder(
            shrinkWrap: true,
            itemCount: snapshot.data!.length,
            itemBuilder: (context, index) {
              print('ProductList ListView.builder snapshot : ${snapshot.data![index].description}');
              return Card(
                clipBehavior: Clip.antiAlias,
                elevation: 0.5,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  mainAxisSize: MainAxisSize.max,
                  children: <Widget>[
                    Image.network(
                      snapshot.data![index].imageUrl,
                      height: 150,
                      width: 100,
                    ),
                    SizedBox(
                      width: 130,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const SizedBox(
                            height: 5.0,
                          ),
                          Text(
                            snapshot.data![index].description,
                            style: theme.textTheme.labelLarge,
                            softWrap: true,
                            overflow: TextOverflow.ellipsis,
                            maxLines: 3,
                          ),
                          const SizedBox(
                            height: 4.0,
                          ),
                          Text(
                            formatter.format(snapshot.data![index].price),
                            style: theme.textTheme.bodySmall,
                          ),
                          const SizedBox(
                            height: 4.0,
                          ),
                        ],
                      ),
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      mainAxisSize: MainAxisSize.max,
                      children: [
                        ElevatedButton(
                          onPressed: () {
                            print('added to cart ${snapshot.data![index].name}, ${snapshot.data![index].productID}');
                            addItem(context, snapshot.data![index]);
                          },
                          child: Text('Add to Cart'),
                          style: ElevatedButton.styleFrom(
                            foregroundColor: Colors.black54,
                            backgroundColor: Colors.amber,
                            elevation: 8.0,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              );
            },
          );
        } else if (snapshot.hasError) {
          return Text("${snapshot.error}");
        }
        return const CircularProgressIndicator();
      },
    );
  }
}

Future<List<Product>> fetchProducts() async {
  final response = await http.post(
    Uri.parse('http://34.143.248.70:80/productlist'),
    // Uri.parse('http://localhost:8080/productlist'),
    body: jsonEncode(<String, dynamic>{
      "action": "",
    }),
  );
  print('http code: ${response.statusCode}');
  if (response.statusCode == 200) {
    Map<String, dynamic> jsonResponse = json.decode(response.body);
    List<dynamic> list = jsonResponse["products"];
    List<Product> productList =
        list.map((item) => Product.fromJson(item)).toList();
    return productList;
  } else {
    throw Exception('Failed to load products');
  }
}

Future<CartDetails> addItem(BuildContext context, Product item) async {
  final response = await http.post(
    Uri.parse('http://34.126.92.60:80/addItem'),
    // Uri.parse('http://localhost:8081/addItem'),
    body: jsonEncode(<String, dynamic>{
      "action": "",
      "userID": "2bcea833-a927-4847-990b-862fa16dfa74",
      "lineItems": [
        {
          "quantity": 1,
          "productDetails": {
            "productID": item.productID,
            "name": item.name,
            "description": item.description,
            "imageUrl": item.imageUrl,
            "currency": item.currency,
            "price": item.price
          }
        }
      ]
    }),
  );
  if (response.statusCode == 200) {
    // print('Response status: ${response.statusCode}');
    // print('Response body: ${response.body}');
    Map<String, dynamic> jsonResponse = json.decode(response.body);
    var cartDetails = CartDetails.fromJson(jsonResponse);
    context.read<CartDetailsNotifier>().cartDetails = cartDetails;
    return cartDetails;
  } else {
    print(response.statusCode);
    throw Exception('Failed to load products');
  }
}
