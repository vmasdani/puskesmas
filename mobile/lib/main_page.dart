import 'dart:convert';
import 'dart:io';

import 'package:android/model.dart';
import 'package:android/state.dart';
import 'package:http/http.dart' as http;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class MainPage extends StatefulWidget {
  const MainPage({Key? key}) : super(key: key);

  @override
  _MainPageState createState() => _MainPageState();
}

class _MainPageState extends State<MainPage> {
  UserBody? _userBody;

  @override
  void initState() {
    fetchData();
  }

  Future<void> fetchData() async {
    try {
      final state = context.read<AppState>();
      final res = await http.get(Uri.parse('${state.baseUrl}/users-jwt'),
          headers: {'authorization': state.apiKey ?? ''});

      if (res.statusCode != HttpStatus.ok) throw res.body;

      setState(() {
        _userBody = UserBody.fromJson(jsonDecode(res.body));
      });
    } catch (e) {
      print('[ERROR FETCHING USERBODY] $e');
    }
  }

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
              margin: const EdgeInsets.only(top: 10),
              child: Consumer<AppState>(
                builder: (ctx, state, child) {
                  return Text(
                      'Hello, ${_userBody?.name} (Username ${_userBody?.username})');
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
                    (await SharedPreferences.getInstance()).remove('apiKey');
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
