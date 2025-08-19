import 'package:flutter/material.dart';
import 'api_service.dart';
import 'dashboard_page.dart';

class SignupPage extends StatefulWidget {
  const SignupPage({super.key});

  @override
  _SignupPageState createState() => _SignupPageState();
}

class _SignupPageState extends State<SignupPage> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _fullNameController = TextEditingController();

  void _signup() async {
    try {
      Map<String, dynamic> response = await ApiService().signup(
        _emailController.text,
        _passwordController.text,
        _fullNameController.text,
      );

      if (response.containsKey('user_id')) {
        String userId = response['user_id'].toString();

        await ApiService().saveUserId(userId);

        ScaffoldMessenger.of(context)
            .showSnackBar(const SnackBar(content: Text('Signup Successful')));

        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => DashboardPage(userId: userId),
          ),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context)
          .showSnackBar(const SnackBar(content: Text('Signup Failed')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 7, 82, 82),
      appBar: AppBar(title: const Text('Signup')),
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
                  'Signup',
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
                  controller: _emailController,
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
                  controller: _passwordController,
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
              const SizedBox(height: 12),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: TextFormField(
                  controller: _fullNameController,
                  decoration: InputDecoration(
                    labelText: 'Full name',
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
                  onPressed: _signup,
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
                      'Signup'),
                ),
              ),
              const Spacer(),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Text(
                    'Already have an account? ',
                    style: TextStyle(color: Color.fromARGB(255, 7, 82, 82)),
                  ),
                  TextButton(
                    onPressed: () {
                      Navigator.pushReplacementNamed(context, '/login');
                    },
                    child: const Text(
                      'Login',
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
