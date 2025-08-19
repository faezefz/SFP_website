import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

// True if user is logged in
  Future<bool> _isLoggedIn() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    return prefs.getBool('isLoggedIn') ?? false;
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<bool>(
      future: _isLoggedIn(),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(child: CircularProgressIndicator());
        }

        if (snapshot.hasData && snapshot.data!) {
          // if logged in go to dashboard
          Future.delayed(Duration.zero, () {
            Navigator.pushReplacementNamed(context, '/dashboard');
          });
          return const Center(child: CircularProgressIndicator());
        }

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
                children: [
                  const Text('Welcome to Software Fault Prediction Website!'),
                  const Text('Please Login to Continue!'),
                  const SizedBox(height: 20),
                  ElevatedButton(
                    onPressed: () {
                      Navigator.pushNamed(context, '/login');
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color.fromARGB(
                          255, 7, 82, 82), // تغییر رنگ پس‌زمینه دکمه
                      padding: const EdgeInsets.symmetric(
                          horizontal: 100, vertical: 12),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                    child: const Text(
                      'Login',
                      style: TextStyle(
                        color: Color.fromARGB(255, 207, 240, 177), // رنگ متن
                      ),
                    ),
                  ),
                  const SizedBox(height: 20),
                  ElevatedButton(
                    onPressed: () {
                      Navigator.pushNamed(context, '/signup');
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color.fromARGB(
                          255, 7, 82, 82), // تغییر رنگ پس‌زمینه دکمه
                      padding: const EdgeInsets.symmetric(
                          horizontal: 100, vertical: 12),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                    child: const Text(
                      'Signup',
                      style: TextStyle(
                        color: Color.fromARGB(255, 207, 240, 177), // رنگ متن
                      ),
                    ),
                  )
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
