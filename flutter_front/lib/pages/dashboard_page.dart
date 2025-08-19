import 'package:flutter/material.dart';
import 'package:flutter_project/pages/api_service.dart';
import 'package:flutter_project/pages/my_projects_page.dart';
import 'colors.dart';
import 'edit_profile_page.dart';

class DashboardPage extends StatefulWidget {
  final String userId;

  const DashboardPage({super.key, required this.userId});

  @override
  _DashboardPageState createState() => _DashboardPageState();
}

class _DashboardPageState extends State<DashboardPage> {
  int _selectedIndex = 0;
  String? userId;

  @override
  void initState() {
    super.initState();
    _loadUserId();
  }

  Future<void> _loadUserId() async {
    userId = await ApiService().getUserId();
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    if (userId == null) {
      return const Center(child: CircularProgressIndicator());
    }

    final List<Widget> pages = [
      MyProjectsPage(userId: userId!),
      EditProfilePage(userId: userId!),
    ];

    return Scaffold(
      appBar: AppBar(
        title:
            const Text('Dashboard', style: TextStyle(color: AppColors.lightGr)),
        backgroundColor: AppColors.darkGr,
        foregroundColor: AppColors.lightGr,
        surfaceTintColor: AppColors.lightGr,
      ),
      drawerScrimColor: AppColors.lightOr,
      body: pages[_selectedIndex],
      drawer: Drawer(
        child: Container(
          color: AppColors.lightGry,
          child: ListView(
            padding: EdgeInsets.zero,
            children: <Widget>[
              ListTile(
                title: const Text('Edit Profile',
                    style: TextStyle(color: AppColors.lightGr)),
                iconColor: AppColors.darkGr,
                selectedColor: AppColors.lightOr,
                selectedTileColor: AppColors.darkGr,
                onTap: () {
                  _onItemTapped(1);
                  Navigator.pop(context);
                },
              ),
              ListTile(
                iconColor: AppColors.darkGr,
                selectedColor: AppColors.lightOr,
                selectedTileColor: AppColors.darkGr,
                title: const Text('My Projects',
                    style: TextStyle(color: AppColors.lightGr)),
                onTap: () {
                  _onItemTapped(0);
                  Navigator.pop(context);
                },
              ),
              ListTile(
                iconColor: AppColors.darkGr,
                selectedColor: AppColors.lightOr,
                selectedTileColor: AppColors.darkGr,
                title: const Text('Logout',
                    style: TextStyle(color: AppColors.lightGr)),
                onTap: () async {
                  Navigator.pop(context);
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }
}
