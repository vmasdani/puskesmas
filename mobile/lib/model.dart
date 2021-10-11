import 'package:json_annotation/json_annotation.dart';

part 'model.g.dart';

@JsonSerializable()
class UserBody {
  UserBody(
      {this.id,
      this.name,
      this.username,
      this.changePassword,
      this.newPassword});

  int? id;
  String? name;
  String? username;
  bool? changePassword;
  String? newPassword;

  factory UserBody.fromJson(Map<String, dynamic> json) =>
      _$UserBodyFromJson(json);
  Map<String, dynamic> toJson() => _$UserBodyToJson(this);
}
