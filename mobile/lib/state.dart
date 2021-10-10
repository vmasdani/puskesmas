import 'package:flutter/cupertino.dart';

class AppState with ChangeNotifier {
  String? _apiKey;
  String? get apiKey => _apiKey;

  String? _fcmToken;
  String? get fcmToken => _fcmToken;

  String? _baseUrl;
  String? get baseUrl => _baseUrl;

  void setApiKey(String? apiKey) {
    _apiKey = apiKey;
    notifyListeners();
  }

  void setFcmToken(String? fcmToken) {
    _fcmToken = fcmToken;
    notifyListeners();
  }

  void setBaseUrl(String? baseUrl) {
    _baseUrl = baseUrl;
    notifyListeners();
  }
}
