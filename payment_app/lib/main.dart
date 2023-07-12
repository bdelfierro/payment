import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'store.dart';
import 'providers.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

Future<void> main() async {
  await dotenv.load(fileName: ".env");
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => CartDetailsNotifier(),
        )
      ],
      child: const Store(),
    ),
  );
}