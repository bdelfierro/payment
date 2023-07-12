import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'models/cartdetails.dart';
import 'providers.dart';
import 'package:provider/provider.dart';
import 'checkout.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';


Future<CartDetails> fetchActiveCartItems(BuildContext context, String userID, String paymentHost) async {
  print('fetchActiveCartItems userID: $userID');

  final response = await http.post(
    Uri.parse('$paymentHost/getActiveCartDetails'),
    body: jsonEncode(<String, dynamic>{
      "userID": userID,
    }),
  );
  print(response.headers);
  // print(response.body);
  switch (response.statusCode) {
    case 200:
      Map<String, dynamic> jsonResponse = json.decode(response.body);
      // print('cartDetails jsonResponse: ${jsonResponse}');
      var cartDetails = CartDetails.fromJson(jsonResponse);
      context.read<CartDetailsNotifier>().cartDetails =  cartDetails;
      return cartDetails;
    case 404:
      context.read<CartDetailsNotifier>().init();
      return context.read<CartDetailsNotifier>().getCartDetails();
    default:
      throw Exception('Failed to load products, ${response.statusCode}');

  }
}

Future<CartDetails> addItemCart(BuildContext context, CartDetails cart, int index, String paymentHost) async {
  print('addItemCart');
  final response = await http.post(
    Uri.parse('$paymentHost/addItem'),
    body: jsonEncode(<String, dynamic>{
      "userID": cart.userID,
      "lineItems": [
        {
          "quantity": 1,
          "productDetails": {
            "productID": cart.lineItems[index].productDetails.productID
          }
        }
      ]
    }),
  );
  // print(response.body);
  if (response.statusCode == 200) {
    Map<String, dynamic> jsonResponse = json.decode(response.body);
    var cartDetails = CartDetails.fromJson(jsonResponse);
    context.read<CartDetailsNotifier>().cartDetails =  cartDetails;
    return cartDetails;
  } else {
    throw Exception('Failed to load products');
  }
}

Future<CartDetails> removeItemCart(BuildContext context, CartDetails cart, int index, String paymentHost) async {
  print('removeItemCart');
  final response = await http.post(
    Uri.parse('$paymentHost/removeItem'),
    body: jsonEncode(<String, dynamic>{
      "userID": cart.userID,
      "cartID": cart.cartID,
      "lineItems": [
        {
          "quantity": -1,
          "productDetails": {
            "productID": cart.lineItems[index].productDetails.productID
          }
        }
      ]
    }),
  );
  // print(response.body);
  if (response.statusCode == 200) {
    Map<String, dynamic> jsonResponse = json.decode(response.body);
    var cartDetails = CartDetails.fromJson(jsonResponse);
    context.read<CartDetailsNotifier>().cartDetails =  cartDetails;
    return cartDetails;
  } else {
    throw Exception('Failed to load products');
  }
}

Future<CartDetails> removeProductFromCart(BuildContext context, CartDetails cart, int index, String paymentHost) async {
  print('removeProductFromCart');
  final response = await http.post(
    Uri.parse('$paymentHost/removeProduct'),
    body: jsonEncode(<String, dynamic>{
      "userID": cart.userID,
      "cartID": cart.cartID,
      "lineItems": [
        {
          "quantity": 1,
          "productDetails": {
            "productID": cart.lineItems[index].productDetails.productID
          }
        }
      ]
    }),
  );
  // print(response.body);
  if (response.statusCode == 200) {
    Map<String, dynamic> jsonResponse = json.decode(response.body);
    var cartDetails = CartDetails.fromJson(jsonResponse);
    context.read<CartDetailsNotifier>().cartDetails =  cartDetails;
    return cartDetails;
  } else {
    throw Exception('Failed to load products');
  }
}

class cart extends StatelessWidget {
  const cart({super.key});

  @override
  Widget build(BuildContext context) {
    return const Placeholder();
  }
}


class CartItemList extends StatelessWidget {
  final String userID;
  final String paymentHost = dotenv.get('PAYMENT_CART_HOST', fallback: 'http://localhost:8081');
  final String checkoutHost = dotenv.get('PAYMENT_CHECKOUT_HOST', fallback: 'http://localhost:4242');

 CartItemList({super.key, required this.userID});

  @override
  Widget build(BuildContext context) {
    print('CartItemList refresh');
    Future<CartDetails> cartDetails = fetchActiveCartItems(context, userID, paymentHost);
    final ThemeData theme = Theme.of(context);
    final NumberFormat formatter = NumberFormat.simpleCurrency(
        locale: Localizations.localeOf(context).toString());

    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.center,
      mainAxisSize: MainAxisSize.min,
      children: [
        Center(
          child: Card(
              clipBehavior: Clip.antiAlias,
              elevation: 0.0,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'Subtotal (${context.watch<CartDetailsNotifier>().getCartDetails().totalCount} items): '
                        '${formatter.format(context.watch<CartDetailsNotifier>().getCartDetails().totalAmount)}',
                  ),
                  const SizedBox(
                    height: 10.0,
                  ),
                  ElevatedButton(
                    onPressed: () async{
                      final sessionUrl = await getSessionURL(context.read<CartDetailsNotifier>().getCartDetails(), checkoutHost);
                      Navigator.push(
                          context,
                          MaterialPageRoute(
                              builder: (context) =>
                                  CheckoutPage(sessionUrl: sessionUrl)));

                    },
                    child: Text(
                      'Check out Cart',
                      style: theme.textTheme.labelLarge,
                    ),
                    style: ButtonStyle(
                      backgroundColor: MaterialStateProperty.all<Color>(Colors.amberAccent),
                    ),
                  ),
                ],
              )
          ),
        ),
        Flexible(
            fit: FlexFit.loose,
            child: FutureBuilder<CartDetails>(
              future: cartDetails,
              builder: (context, snapshot) {
                if (snapshot.hasData) {
                  return ListView.builder(
                    shrinkWrap: true,
                    itemCount: snapshot.data!.lineItems.length,
                    itemBuilder: (context, index) {
                      print('index: $index');
                      print('index: ${snapshot.data!.lineItems}');
                      print('length: ${snapshot.data!.lineItems.length}');
                      if (index >= snapshot.data!.lineItems.length) {
                          return const Text('');
                      }else {
                        return Padding(
                          padding: EdgeInsets.all(5.0),
                          child: Card(
                            clipBehavior: Clip.antiAlias,
                            elevation: 0.5,
                            child: Padding(
                              padding: EdgeInsets.all(4.0),
                              child: Column(
                                mainAxisAlignment: MainAxisAlignment.start,
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.start,
                                    mainAxisSize: MainAxisSize.min,
                                    children: [
                                      Padding(
                                        padding: EdgeInsets.all(5.0),
                                        child: Image.network(
                                          snapshot.data!.lineItems[index].productDetails
                                              .imageUrl,
                                          height: 100,
                                          width: 80,
                                        ),
                                      ),
                                      SizedBox(
                                        width: 200,
                                        child: Column(
                                          mainAxisAlignment: MainAxisAlignment.start,
                                          crossAxisAlignment: CrossAxisAlignment.start,
                                          children: [
                                            Text(
                                              snapshot.data!.lineItems[index].productDetails.description,
                                              style: theme.textTheme.labelLarge,
                                              softWrap: true,
                                              overflow: TextOverflow.ellipsis,
                                              maxLines: 3,
                                            ),
                                            const SizedBox(
                                              height: 4.0,
                                            ),
                                            Text(formatter.format(snapshot
                                                .data!.lineItems[index].productDetails.price)),
                                          ],
                                        ),
                                      ),
                                    ],
                                  ),
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.start,
                                    mainAxisSize: MainAxisSize.min,
                                    children: <Widget>[
                                      IconButton(
                                          onPressed: () {
                                            print('remove 1 item');
                                            removeItemCart(context, snapshot.data!, index, paymentHost);
                                            // context.read<CartDetailsNotifier>().removeItem(snapshot.data!.lineItems[index].productDetails.productID);
                                          },
                                          icon: const Icon(Icons.remove)),
                                      Padding(
                                        padding: EdgeInsets.symmetric(horizontal: 10),
                                        child: Text(
                                          '${context.watch<CartDetailsNotifier>().getCartDetails().lineItems[index].quantity}',
                                          style: Theme.of(context).textTheme.titleSmall,
                                        ),
                                      ),
                                      IconButton(
                                        onPressed: () {
                                          print('add 1 item');
                                          addItemCart(context, snapshot.data!, index, paymentHost);
                                          // context.read<CartDetailsNotifier>().addItem(snapshot.data!.lineItems[index].productDetails);
                                        },
                                        icon: const Icon(Icons.add),
                                      ),
                                      const SizedBox(
                                        height: 5.0,
                                      ),
                                      IconButton(
                                        onPressed: () {
                                          print('remove product');
                                          removeProductFromCart(context, snapshot.data!, index, paymentHost);
                                        },
                                        icon: const Icon(Icons.delete),
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                          ),
                        );
                      }

                    },
                  );
                } else if (snapshot.hasError) {
                  return Text("${snapshot.error}");
                }
                // By default, show a loading spinner.
                return const CircularProgressIndicator();
              },
            ),
        ),
      ],
    );
  }
}

