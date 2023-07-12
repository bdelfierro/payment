## Shopping cart payment system using stripe checkout, Flutter mobile app and Golang backend services.

---

### Payment System Architecture using Google Cloud Platform services

![Alternate Text](payment%20diagram.jpg)


1. When user hits checkout cart button, Flutter app calls payment-checkout service to get stripe session passing the list of products for payment
2. payment-checkout service calls stripe api to create new session and returns session data to Flutter app
3. Flutter app calls stripe session url link and displays the stripe checkout page via webview widget
4. stripe processes payment (e.g. credit card payment) and calls webhook endpoint configured in your account
5. payment-webhook endpoint will publish the checkout.session.completed event received from stripe to PubSub
6. PubSub message will trigger payment workflow that will call payment-cart service to update cart status and publish event to notification PubSub topic which can be used to send email, sms notification to user. This can also be extended to track delivery.
7. payment-cart service updates cart status in postgreSQL DB.

---

### Screenshots of Flutter app

<div style="position:relative;">
<img src="images/product_list.png" alt="Product list" width="350" height="620" style="margin-right:50px; margin-top: 50px" />
<img src="images/shopping_cart.png" alt="Shopping cart" width="350" height="620" style="margin-left:50px; margin-top: 50px"/>
</div>

---

### Checkout page hosted by stripe 
    Displayed in flutter app via webview widget

<div style="display:flex;">
    <img src="images/checkout_p1.png" alt="checkout page1" width="350" height="620" style="margin-right:20px; margin-top: 20px" />
    <img src="images/checkout_p2.png" alt="checkout page2" width="350" height="620" style="margin-right:20px; margin-top: 20px" />
</div>

---

<div style="display:flex;">
    <img src="images/checkout_p3.png" alt="checkout page3" width="350" height="620" style="margin-right:50px; margin-top: 50px" />
</div>

---

### Payment transactions can be viewed from your stripe account 

<div style="display:flex;">
<img src="images/stripe_payment_trans.png" alt="checkout page1" width="981" height="900" style="margin-right:20px; margin-top: 20px; margin-bottom: 20px;" />
</div>