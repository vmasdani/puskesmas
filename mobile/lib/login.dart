import 'dart:convert';
import 'dart:io';

import 'package:android/state.dart';
import 'package:flutter/material.dart';
import 'package:flutter/painting.dart';
import 'package:provider/provider.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

import 'helpers.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({Key? key}) : super(key: key);

  @override
  _LoginPageState createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  var username = '';
  var password = '';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        margin: const EdgeInsets.only(top: 10, left: 10, right: 10),
        child: ListView(
          children: [
            Container(
              margin: EdgeInsets.only(top: 10, bottom: 10),
              child: const Text(
                'Puskesmas Pasir Nangka',
                style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
              ),
            ),
            Container(
              margin: const EdgeInsets.only(top: 5, bottom: 5),
              child: TextField(
                decoration: const InputDecoration(
                    hintText: 'Username...',
                    border: OutlineInputBorder(),
                    isDense: true,
                    label: Text('Username')),
                controller: (() {
                  final controller = TextEditingController();

                  return controller;
                })(),
                onChanged: (u) {
                  username = u;
                },
              ),
            ),
            Container(
              margin: const EdgeInsets.only(top: 5, bottom: 5),
              child: TextField(
                decoration: const InputDecoration(
                    hintText: 'Password...',
                    border: OutlineInputBorder(),
                    isDense: true,
                    label: Text('Password')),
                controller: (() {
                  final controller = TextEditingController();

                  return controller;
                })(),
                onChanged: (p) {
                  password = p;
                },
              ),
            ),
            Container(
              margin: const EdgeInsets.only(top: 5, bottom: 5),
              child: MaterialButton(
                color: Colors.blue,
                textColor: Colors.white,
                child: const Text('Login'),
                onPressed: () async {
                  final state = context.read<AppState>();

                  try {
                    final res = await http.post(
                        Uri.parse('${state.baseUrl}/login'),
                        headers: {'content-type': 'application/json'},
                        body: jsonEncode(
                            {'username': username, 'password': password}));

                    if (res.statusCode != HttpStatus.ok) throw res.body;

                    final token = res.body;
                    state.setApiKey(token);
                    (await SharedPreferences.getInstance())
                        .setString('apiKey', token);

                    initLogin(state, token);
                  } catch (e) {
                    print('[ERROR LOGIN] $e');

                    showDialog(
                        context: context,
                        builder: (_) => AlertDialog(
                              title: const Text('Error!'),
                              content: Text('$e'),
                            ));
                  }

                  // showDialog(
                  //     context: context,
                  //     builder: (_) => AlertDialog(
                  //           title: const Text('Values'),
                  //           content: Text(jsonEncode(
                  //               {'username': username, 'password': password})),
                  //         ));
                },
              ),
            )
          ],
        ),
      ),
    );
  }
}
