// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserBody _$UserBodyFromJson(Map<String, dynamic> json) => UserBody(
      id: json['id'] as int?,
      name: json['name'] as String?,
      username: json['username'] as String?,
      changePassword: json['changePassword'] as bool?,
      newPassword: json['newPassword'] as String?,
    );

Map<String, dynamic> _$UserBodyToJson(UserBody instance) => <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'username': instance.username,
      'changePassword': instance.changePassword,
      'newPassword': instance.newPassword,
    };
