import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';
import 'package:webview_flutter_android/webview_flutter_android.dart';
import 'package:webview_flutter_wkwebview/webview_flutter_wkwebview.dart';
import 'home.dart';
import 'models/cartdetails.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'models/checkoutcart.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';



class CheckoutPage extends StatefulWidget {

  final String sessionUrl;

  const CheckoutPage({Key? key, required this.sessionUrl}) : super(key: key);


  @override
  State<CheckoutPage> createState() => _CheckoutPageState();
}

class _CheckoutPageState extends State<CheckoutPage> {
  late final WebViewController _controller;

  @override
  void initState() {
    super.initState();

    late final PlatformWebViewControllerCreationParams params;
    params = const PlatformWebViewControllerCreationParams();

    final WebViewController controller =
    WebViewController.fromPlatformCreationParams(params);

    controller
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setBackgroundColor(const Color(0x00000000))
      ..setNavigationDelegate(
        NavigationDelegate(
          onProgress: (int progress) {
            debugPrint('WebView is loading (progress : $progress%)');
          },
          onPageStarted: (String url) {
            debugPrint('Page started loading: $url');
          },
          onPageFinished: (String url) {
            debugPrint('Page finished loading: $url');
          },
          onWebResourceError: (WebResourceError error) {
            debugPrint('''
              Page resource error:
              code: ${error.errorCode}
              description: ${error.description}
              errorType: ${error.errorType}
              isForMainFrame: ${error.isForMainFrame}
          ''');
          },
          onNavigationRequest: (NavigationRequest request) {
            if (request.url.contains('success')) {
              Navigator.push(
                  context,
                  MaterialPageRoute(
                      builder: (context) =>
                      const Success()));
            }
            debugPrint('allowing navigation to ${request.url}');
            return NavigationDecision.navigate;
          },
          onUrlChange: (UrlChange change) {
            debugPrint('url change to ${change.url}');
          },
        ),
      )
      ..loadRequest(Uri.parse(widget.sessionUrl));

    if (controller.platform is AndroidWebViewController) {
      AndroidWebViewController.enableDebugging(true);
      (controller.platform as AndroidWebViewController)
          .setMediaPlaybackRequiresUserGesture(false);
    }

    _controller = controller;
  }


  @override
  Widget build(BuildContext context) {
    return Scaffold(
      resizeToAvoidBottomInset: false,
      body: WebViewWidget(controller: _controller),
    );
  }
}


Future<String> getSessionURL(CartDetails cart, String checkoutHost) async {

  List<CheckoutItem> items = [];

  for (var item in cart.lineItems) {
    CheckoutItem cartItem = CheckoutItem(
      productName: item.productDetails.name,
      description: item.productDetails.description,
      imageUrl: item.productDetails.imageUrl,
      currency: item.productDetails.currency,
      quantity: item.quantity,
      unitPrice: item.productDetails.price,
    );
    items.add(cartItem);
  }

  CheckoutCart cartItems = CheckoutCart(
      cartID: cart.cartID,
      lineItems: items,
  );

  final response = await http.post(
    Uri.parse('$checkoutHost/create-checkout-session'),
    body: jsonEncode(cartItems.toJson()),
  );

  if (response.statusCode == 200) {
    Map<String, dynamic> jsonResponse = json.decode(response.body);
    var checkoutSession = CheckoutSession.fromJson(jsonResponse);
    print(checkoutSession);
    return checkoutSession.sessionURL;
  } else {
    throw Exception('Failed to load products');
  }
}

class Success extends StatelessWidget {
  const Success({super.key});

  @override
  Widget build(BuildContext context) {

    final ThemeData theme = Theme.of(context);

    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            Text(
                'payment successful',
              style: theme.textTheme.labelLarge,
            ),
            const SizedBox(
              height: 5.0,
            ),
            ElevatedButton(
                onPressed: (){

                  Navigator.push(
                      context,
                      MaterialPageRoute(
                          builder: (context) => HomePage(userID: '2bcea833-a927-4847-990b-862fa16dfa74'),
                      ));
                },
                child: const Text('continue shopping'),
                style: ElevatedButton.styleFrom(
                  foregroundColor: Colors.black54,
                  backgroundColor: Colors.amber,
                  elevation: 8.0,
                ),
            ),
          ],
        )
      ),
    );
  }
}
