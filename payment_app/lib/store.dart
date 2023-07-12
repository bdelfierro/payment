import 'package:flutter/material.dart';
import 'home.dart';
import 'colors.dart';


class Store extends StatelessWidget {
  const Store({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      title: 'Store',
      initialRoute: '/',
      routes: {
        '/': (BuildContext context) => HomePage(userID: '2bcea833-a927-4847-990b-862fa16dfa74'),
      },
      theme: _kStoreTheme,
    );
  }
}

final ThemeData _kStoreTheme = _buildStoreTheme();

ThemeData _buildStoreTheme() {

  final ThemeData base = ThemeData.light(useMaterial3: true);

  return base.copyWith(
    colorScheme: base.colorScheme.copyWith(
      primary: kStorePink100,
      onPrimary: kStoreBrown900,
      secondary: kStoreBrown900,
      error: kStoreErrorRed,
    ),
    textTheme: _buildStoreTextTheme(base.textTheme),
    textSelectionTheme: const TextSelectionThemeData(
      selectionColor: kStorePink100,
    ),
    appBarTheme: const AppBarTheme(
      foregroundColor: kStoreBrown900,
      backgroundColor: kStorePink100,
    ),
  );
}


TextTheme _buildStoreTextTheme(TextTheme base) {
  return base
      .copyWith(
    headlineSmall: base.headlineSmall!.copyWith(
      fontWeight: FontWeight.w500,
    ),
    titleLarge: base.titleLarge!.copyWith(
      fontSize: 18.0,
    ),
    bodySmall: base.bodySmall!.copyWith(
      fontWeight: FontWeight.w400,
      fontSize: 14.0,
    ),
    bodyLarge: base.bodyLarge!.copyWith(
      fontWeight: FontWeight.w500,
      fontSize: 16.0,
    ),
  )
      .apply(
    fontFamily: 'Rubik',
    displayColor: kStoreBrown900,
    bodyColor: kStoreBrown900,
  );
}