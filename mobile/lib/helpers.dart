import 'dart:convert';
import 'dart:io';

import 'package:android/state.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:http/http.dart' as http;

Future<void> initLogin(AppState? state, String? apiKey) async {
  if (apiKey != null) {
    state?.setApiKey(apiKey);
    (await SharedPreferences.getInstance()).setString('apiKey', apiKey);
  }

  try {
    final res = await http.post(Uri.parse('${state?.baseUrl}/save-fcm-token'),
        headers: {
          'content-type': 'application/json',
          'authorization': apiKey ?? ''
        },
        body:
            jsonEncode({'token': await FirebaseMessaging.instance.getToken()}));

    if (res.statusCode != HttpStatus.ok) throw res.body;
  } catch (e) {
    print('[ERROR] $e');
  }
}
