import 'package:flutter/material.dart';
import 'package:flutter_project/pages/colors.dart';
import 'package:flutter_project/pages/project_detail_page.dart';
import 'api_service.dart';

class MyProjectsPage extends StatefulWidget {
  final String userId;

  const MyProjectsPage({super.key, required this.userId});

  @override
  _MyProjectsPageState createState() => _MyProjectsPageState();
}

class _MyProjectsPageState extends State<MyProjectsPage> {
  List<dynamic> projects = [];
  bool isLoading = true;

  @override
  void initState() {
    super.initState();
    _fetchProjects();
  }

  Future<void> _fetchProjects() async {
    try {
      List<dynamic> fetchedProjects =
          await ApiService().fetchProjects(widget.userId);
      setState(() {
        projects = fetchedProjects;
        isLoading = false;
      });
    } catch (e) {
      setState(() {
        isLoading = false;
      });
    }
  }

  Future<void> _createProject(String projectName, String description) async {
    bool success = await ApiService()
        .createProject(widget.userId, projectName, description);

    if (success) {
      _fetchProjects();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Project created successfully')),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Failed to create project')),
      );
    }
  }

  Future<void> _updateProject(
      String projectId, String projectName, String description) async {
    bool success =
        await ApiService().updateProject(projectId, projectName, description);

    if (success) {
      _fetchProjects();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Project updated successfully')),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Failed to update project')),
      );
    }
  }

  Future<void> _deleteProject(String projectId) async {
    bool success = await ApiService().deleteProject(projectId);

    if (success) {
      _fetchProjects();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Project deleted successfully')),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Failed to delete project')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawerScrimColor: AppColors.lightOr,
      appBar: AppBar(
        backgroundColor: AppColors.lightGr,
        surfaceTintColor: AppColors.lightOr,
        title: const Text(
          'My Projects',
          style: TextStyle(color: AppColors.darkGr),
        ),
      ),
      body: isLoading
          ? const Center(child: CircularProgressIndicator())
          : Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                children: [
                  ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      primary: AppColors.darkGr,
                      padding: const EdgeInsets.symmetric(
                          horizontal: 100, vertical: 12),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                    onPressed: () async {
                      Map<String, String> result =
                          await _showProjectDialog(context);
                      String projectName = result['projectName'] ?? '';
                      String description = result['description'] ?? '';
                      if (projectName.isNotEmpty) {
                        _createProject(projectName, description);
                      }
                    },
                    child: const Text('Create New Project',
                        style: TextStyle(color: AppColors.lightGr)),
                  ),
                  const SizedBox(height: 16),
                  Expanded(
                    child: ListView.builder(
                      itemCount: projects.length,
                      itemBuilder: (context, index) {
                        return Card(
                          elevation: 4.0,
                          margin: const EdgeInsets.symmetric(vertical: 8.0),
                          color: AppColors.lightGry,
                          child: ListTile(
                            iconColor: AppColors.darkGr,
                            title: Text(
                              projects[index]['name'],
                              style: TextStyle(color: AppColors.darkGr),
                            ),
                            subtitle: Text(
                              projects[index]['description'] ??
                                  'No description',
                              style: TextStyle(color: AppColors.darkGr),
                            ),
                            onTap: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (context) => ProjectDetailsPage(
                                    projectId: projects[index]['id'].toString(),
                                    userId: int.parse(widget.userId),
                                  ),
                                ),
                              );
                            },
                            trailing: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                IconButton(
                                  color: AppColors.lightOr,
                                  icon: const Icon(Icons.edit),
                                  onPressed: () {
                                    String projectId =
                                        projects[index]['id'].toString();
                                    String projectName =
                                        projects[index]['name'];
                                    String description =
                                        projects[index]['description'] ?? '';
                                    _updateProject(
                                        projectId, projectName, description);
                                  },
                                ),
                                IconButton(
                                  color: AppColors.lightOr,
                                  icon: const Icon(Icons.delete),
                                  onPressed: () {
                                    String projectId =
                                        projects[index]['id'].toString();
                                    _deleteProject(projectId);
                                  },
                                ),
                              ],
                            ),
                          ),
                        );
                      },
                    ),
                  ),
                ],
              ),
            ),
    );
  }

  Future<Map<String, String>> _showProjectDialog(BuildContext context) async {
    String projectName = '';
    String description = '';

    await showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          backgroundColor: AppColors.darktGry,
          title: const Text('Enter Project Details',
              style: TextStyle(color: AppColors.lightGr)),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: <Widget>[
              TextField(
                onChanged: (value) {
                  projectName = value;
                },
                decoration: InputDecoration(
                  hintText: 'Project Name',
                  hintStyle: TextStyle(color: AppColors.lightGry),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: AppColors.lightOr),
                  ),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                onChanged: (value) {
                  description = value;
                },
                decoration: InputDecoration(
                  hintText: 'Project Description',
                  hintStyle: TextStyle(color: AppColors.lightGry),
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: AppColors.lightOr),
                  ),
                ),
              ),
            ],
          ),
          actions: <Widget>[
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
              },
              child: const Text('Cancel',
                  style: TextStyle(color: AppColors.lightGr)),
            ),
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
              },
              child: const Text('Submit',
                  style: TextStyle(color: AppColors.lightGr)),
            ),
          ],
        );
      },
    );

    return {
      'projectName': projectName,
      'description': description,
    };
  }
}
