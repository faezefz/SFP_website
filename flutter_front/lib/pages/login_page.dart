// ignore_for_file: use_build_context_synchronously

import 'package:flutter/material.dart';
import 'package:flutter_project/pages/api_service.dart';
import 'package:flutter_project/pages/dashboard_page.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  // ignore: library_private_types_in_public_api
  _LoginPageState createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

  // تابع برای لاگین
  void _login() async {
    String email = _emailController.text;
    String password = _passwordController.text;

    // بررسی اگر ایمیل یا پسورد خالی باشد
    if (email.isEmpty || password.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text("Email and password cannot be empty")),
      );
      return;
    }

    try {
      Map<String, dynamic> response = await ApiService().login(
        email,
        password,
      );

      if (response.containsKey('user_id')) {
        String userId = response['user_id'].toString();

        // ذخیره user_id در shared_preferences
        await ApiService().saveUserId(userId);

        ScaffoldMessenger.of(context)
            .showSnackBar(const SnackBar(content: Text('Login Successful')));

        // هدایت کاربر به صفحه داشبورد با استفاده از آیدی کاربر
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) =>
                DashboardPage(userId: userId), // هدایت به داشبورد با آیدی کاربر
          ),
        );
      } else {
        // اگر user_id در پاسخ نیست
        ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Invalid response from server')));
      }
    } catch (e) {
      ScaffoldMessenger.of(context)
          .showSnackBar(const SnackBar(content: Text('Login Failed')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 7, 82, 82),
      appBar: AppBar(title: const Text('Login')),
      body: Center(
        child: Container(
          width: 350,
          height: 400,
          decoration: BoxDecoration(
            color: const Color.fromARGB(255, 207, 240, 177),
            borderRadius: BorderRadius.circular(10),
            boxShadow: const [
              BoxShadow(
                color: Color.fromARGB(255, 207, 240, 177),
                offset: Offset(0, 4),
                blurRadius: 6,
              ),
            ],
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Padding(
                padding: EdgeInsets.all(16.0),
                child: Text(
                  'LOGIN',
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                    color: Color.fromARGB(255, 7, 82, 82),
                  ),
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: TextFormField(
                  controller: _emailController, // اتصال کنترلر به فیلد
                  decoration: InputDecoration(
                    labelText: 'Email',
                    labelStyle: const TextStyle(
                        color: Color.fromARGB(255, 207, 240, 177)),
                    fillColor: const Color.fromARGB(255, 7, 82, 82),
                    filled: true,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: TextFormField(
                  controller: _passwordController, // اتصال کنترلر به فیلد
                  obscureText: true,
                  decoration: InputDecoration(
                    labelText: 'Password',
                    labelStyle: const TextStyle(
                        color: Color.fromARGB(255, 207, 240, 177)),
                    fillColor: const Color.fromARGB(255, 7, 82, 82),
                    filled: true,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 20),
              Center(
                child: ElevatedButton(
                  onPressed: _login,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color.fromARGB(255, 7, 82, 82),
                    padding: const EdgeInsets.symmetric(
                        horizontal: 100, vertical: 12),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(8),
                    ),
                  ),
                  child: const Text(
                      style:
                          TextStyle(color: Color.fromARGB(255, 207, 240, 177)),
                      'LOGIN'),
                ),
              ),
              const SizedBox(height: 10),
              TextButton(
                onPressed: () {},
                child: const Text(
                  'Forgot Password?',
                  style: TextStyle(color: Color.fromARGB(255, 7, 82, 82)),
                ),
              ),
              const Spacer(),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Text(
                    'Don\'t have an account? ',
                    style: TextStyle(color: Color.fromARGB(255, 7, 82, 82)),
                  ),
                  TextButton(
                    onPressed: () {
                      Navigator.pushReplacementNamed(context, '/signup');
                    },
                    child: const Text(
                      'Register',
                      style: TextStyle(
                          backgroundColor: Color.fromARGB(255, 7, 82, 82),
                          color: Color.fromARGB(255, 207, 240, 177)),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
