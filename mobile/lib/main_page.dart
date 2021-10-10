import 'package:android/state.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

class MainPage extends StatefulWidget {
  const MainPage({Key? key}) : super(key: key);

  @override
  _MainPageState createState() => _MainPageState();
}

class _MainPageState extends State<MainPage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Puskesmas'),
      ),
      body: Container(
        margin: const EdgeInsets.only(left: 10, right: 10),
        child: ListView(
          children: [
            Container(
              margin: const EdgeInsets.only(),
              child: Consumer<AppState>(
                builder: (ctx, state, child) {
                  return Text('Hello, ${state.apiKey}');
                },
              ),
            ),
            Container(
              margin: const EdgeInsets.only(top: 10),
              child: Container(
                alignment: Alignment.center,
                child: MaterialButton(
                  onPressed: () async {
                    final state = context.read<AppState>();
                    state.setApiKey(null);
                  },
                  color: Colors.blue,
                  textColor: Colors.white,
                  child: const Text('Logout'),
                ),
              ),
            )
          ],
        ),
      ),
    );
  }
}
