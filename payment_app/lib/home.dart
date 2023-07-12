import 'package:flutter/material.dart';
import 'productlist.dart';
import 'package:provider/provider.dart';
import 'models/cartdetails.dart';
import 'providers.dart';
import 'cartItems.dart';
import 'cart.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';


class HomePage extends StatefulWidget {
  final String userID;
  final String paymentHost = dotenv.get('PAYMENT_CART_HOST', fallback: 'http://localhost:8081');

  HomePage({super.key, required this.userID});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {

  @override
  void initState(){
    super.initState();
    print('HomePage initState()');
    cartItemsCount(context, widget.userID, widget.paymentHost);
  }

  @override
  Widget build(BuildContext context) {
    print('refresh HomePage');
    print('cartID during refresh of CartPage: ${context.read<CartDetailsNotifier>().getCartDetails().cartID}');
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: Text(
          'Product List',
          style: Theme.of(context).textTheme.headlineMedium,
        ),
        actions: <Widget>[
          Padding(
            padding: const EdgeInsets.all(10.0),
            child: Badge(
              label: Text('${context.watch<CartDetailsNotifier>().getCartDetails().totalCount}'),
              child: IconButton(
                onPressed: () {
                  print('cartID before routing to CartPage: ${context.read<CartDetailsNotifier>().getCartDetails().cartID}');
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) {
                      return const CartPage();
                    }),
                  );
                },
                icon: const Icon(Icons.shopping_cart),
              ),
            ),
          ),
        ],
      ),
      body: ProductList(products: fetchProducts()),
    );
  }

}

void cartItemsCount(BuildContext context, String userID, String paymentHost)  {
  print('cartItemsCount: userID: $userID');

  fetchActiveCartItems(context, userID, paymentHost)
  .then((cartDetails) {
    print('cartID in cartItemsCount: ${context.read<CartDetailsNotifier>().getCartDetails().cartID}');
  })
  .catchError((error) {
      print('error in fetchActiveCartItems: $error');
  });
}

