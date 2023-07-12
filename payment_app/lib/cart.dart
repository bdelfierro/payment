import 'package:flutter/material.dart';
import 'cartItems.dart';
import 'providers.dart';
import 'package:provider/provider.dart';


class CartPage extends StatelessWidget {

  const CartPage({super.key});

  @override
  Widget build(BuildContext context) {
    print('refresh CartPage');
    print('cartID from CartDetailsNotifier: ${context.read<CartDetailsNotifier>().getCartDetails().cartID}');
    var userID = context.watch<CartDetailsNotifier>().getCartDetails().userID;
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: Text(
          'My Shopping Cart',
          style: Theme.of(context).textTheme.headlineMedium,
        ),
        actions: <Widget>[
          Padding(
              padding: const EdgeInsets.all(10.0),
            child: Badge(
              label: Text('${context.watch<CartDetailsNotifier>().getCartDetails().totalCount}'),
              child: IconButton(
                onPressed: (){
                  print('cart pressed');
                },
                icon: const Icon(Icons.shopping_cart),
              ),
            ),
          ),
        ],
      ),
      body: CartItemList(userID: userID),
    );
  }
}

class Count extends StatelessWidget {
  const Count({Key? key}) : super(key: key);
  
  @override
  Widget build(BuildContext context) {
    return Text(
      /// Calls `context.watch` to make [Count] rebuild when [Counter] changes.
      '${context.watch<CartDetailsNotifier>().getCartDetails().totalCount}',
    );
  }
}