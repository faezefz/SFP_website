import 'package:flutter/material.dart';
import 'package:flutter_project/pages/colors.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class EditProfilePage extends StatefulWidget {
  final String userId;
  const EditProfilePage({super.key, required this.userId});

  @override
  _EditProfilePageState createState() => _EditProfilePageState();
}

class _EditProfilePageState extends State<EditProfilePage> {
  final _formKey = GlobalKey<FormState>();
  String _email = '';
  String _password = '';
  String _fullName = '';

  Future<void> _updateUser() async {
    final url = Uri.parse('http://localhost:8081/updateUser');
    final response = await http.put(
      url,
      headers: <String, String>{
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'id': widget.userId,
        'email': _email,
        'password_hash': _password,
        'full_name': _fullName,
      }),
    );

    if (response.statusCode == 200) {
      print('User updated successfully');
      Navigator.pop(context);
    } else {
      print('Failed to update user: ${response.body}');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawerScrimColor: AppColors.lightOr,
      appBar: AppBar(
        backgroundColor: AppColors.darkGr,
        foregroundColor: AppColors.lightGr,
        surfaceTintColor: AppColors.lightGr,
        title: const Text('Edit Profile'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              TextFormField(
                initialValue: _email,
                decoration: InputDecoration(
                  labelText: 'Email',
                  hintText: 'Enter your email',
                  labelStyle: TextStyle(color: AppColors.darkGr),
                  hintStyle: TextStyle(color: AppColors.lightGry),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: AppColors.lightOr),
                  ),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter your email';
                  }
                  return null;
                },
                onSaved: (value) {
                  _email = value ?? '';
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                initialValue: _password,
                decoration: InputDecoration(
                  labelText: 'Password',
                  hintText: 'Enter your password',
                  labelStyle: TextStyle(color: AppColors.darkGr),
                  hintStyle: TextStyle(color: AppColors.lightGry),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: AppColors.lightOr),
                  ),
                ),
                obscureText: true,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter your password';
                  }
                  return null;
                },
                onSaved: (value) {
                  _password = value ?? '';
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                initialValue: _fullName,
                decoration: InputDecoration(
                  labelText: 'Full Name',
                  hintText: 'Enter your full name',
                  labelStyle: TextStyle(color: AppColors.darkGr),
                  hintStyle: TextStyle(color: AppColors.lightGry),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: AppColors.lightOr),
                  ),
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter your full name';
                  }
                  return null;
                },
                onSaved: (value) {
                  _fullName = value ?? '';
                },
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () {
                  if (_formKey.currentState!.validate()) {
                    _formKey.currentState!.save();
                    _updateUser();
                  }
                },
                style: ElevatedButton.styleFrom(
                  primary: AppColors.lightOr,
                  onPrimary: AppColors.darkGr,
                ),
                child: const Text('Save Changes'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
