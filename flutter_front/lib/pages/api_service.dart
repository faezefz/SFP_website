import 'dart:convert';
import 'dart:typed_data';
import 'package:file_picker/file_picker.dart';
import 'package:flutter_project/models/dataset.dart';
import 'package:flutter_project/models/project.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'dart:io';
import 'package:dio/dio.dart';

class ApiService {
  final Dio _dio = Dio();
  final String baseUrl = 'http://localhost:8081'; // backend url
  final String pyUrl = 'http://localhost:5001';

  // POST for signup
  Future<Map<String, dynamic>> signup(
      String email, String password, String fullName) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/signup'),
        headers: <String, String>{
          'Content-Type': 'application/json',
        },
        body: jsonEncode(<String, String>{
          'email': email,
          'password': password,
          'full_name': fullName,
        }),
      );

      if (response.statusCode == 201) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Failed to sign up: ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to sign up');
    }
  }

  // POST for login
  Future<Map<String, dynamic>> login(String email, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/login'),
        headers: <String, String>{
          'Content-Type': 'application/json',
        },
        body: jsonEncode(<String, String>{
          'email': email,
          'password': password,
        }),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        throw Exception('Failed to log in: ${response.body}');
      }
    } catch (e) {
      throw Exception('Failed to log in');
    }
  }

  Future<void> saveUserId(String userId) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('userId', userId);
  }

  Future<void> handleLogin(String email, String password) async {
    try {
      final Map<String, dynamic> response =
          await ApiService().login(email, password);

      String userId = response['userId'];
      saveUserId(userId);
    } catch (e) {}
  }

  Future<String?> getUserId() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('userId');
  }

  Future<List<dynamic>> fetchProjects(String userId) async {
    final url = Uri.parse('$baseUrl/projects/$userId');

    try {
      final response = await http.get(url);

      if (response.statusCode == 200) {
        return json.decode(response.body);
      } else {
        throw Exception('Failed to load projects');
      }
    } catch (e) {
      throw Exception('Error: $e');
    }
  }

  Future<bool> createProject(
      String userId, String projectName, String description) async {
    final url = Uri.parse('$baseUrl/projects');
    int userIdInt = int.parse(userId);
    final response = await http.post(
      url,
      headers: <String, String>{
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'owner_user_id': userIdInt,
        'name': projectName,
        'description': description
      }),
    );

    if (response.statusCode == 201) {
      return true;
    } else {
      return false;
    }
  }

  Future<bool> updateProject(
      String projectId, String projectName, String description) async {
    final url = Uri.parse('$baseUrl/projects/$projectId');

    final response = await http.put(
      url,
      headers: <String, String>{
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'name': projectName,
        'description': description,
      }),
    );

    if (response.statusCode == 200) {
      return true;
    } else {
      return false;
    }
  }

  Future<bool> deleteProject(String projectId) async {
    final url = Uri.parse('$baseUrl/projects/$projectId');

    final response = await http.delete(url);

    if (response.statusCode == 200) {
      return true;
    } else {
      return false;
    }
  }

  Future<List<Project>> getProjectsByOwnerID(int ownerUserID) async {
    final response = await http.get(
      Uri.parse('$baseUrl/projects/owner/$ownerUserID'),
    );

    if (response.statusCode == 200) {
      List<dynamic> data = json.decode(response.body);
      return data.map((item) => Project.fromJson(item)).toList();
    } else {
      throw Exception('Failed to load projects');
    }
  }

  Future<List<Dataset>> getDatasetsByProjectID(int projectID) async {
    final response = await http.get(
      Uri.parse('$baseUrl/projects/$projectID/datasets'),
    );

    if (response.statusCode == 200) {
      List<dynamic> data = json.decode(response.body);
      return data.map((item) => Dataset.fromJson(item)).toList();
    } else {
      throw Exception('Failed to load datasets');
    }
  }

  Future<Map<String, dynamic>> getProjectById(int projectId) async {
    final response =
        await http.get(Uri.parse('$baseUrl/projects/project/$projectId'));

    if (response.statusCode == 200) {
      return json.decode(response.body);
    } else {
      throw Exception('Failed to load project');
    }
  }

  Future<String> uploadCustomModel() async {
    File? modelFile = await _pickModelFile();

    if (modelFile == null) {
      return 'No file selected';
    }

    try {
      String uploadUrl = '$baseUrl/upload-model';

      FormData formData = FormData.fromMap({
        'model': await MultipartFile.fromFile(modelFile.path,
            filename: 'custom_model.pkl'),
      });

      Response response = await _dio.post(uploadUrl, data: formData);

      if (response.statusCode == 200) {
        return 'Model uploaded successfully';
      } else {
        return 'Failed to upload model: ${response.statusCode}';
      }
    } catch (e) {
      return 'Error: $e';
    }
  }

  Future<File?> _pickModelFile() async {
    final result = await FilePicker.platform.pickFiles(
      allowedExtensions: ['pkl'],
      type: FileType.custom,
    );

    if (result == null) return null;

    return File(result.files.single.path!);
  }

  Future<List<String>> getDatasetsByUserId(int userId) async {
    try {
      final response = await _dio.get('$baseUrl/datasets/$userId');

      if (response.statusCode == 200) {
        // response is list of strings : {id - name}
        List<dynamic> datasets = response.data;
        return datasets
            .map((dataset) =>
                '${dataset['id']} - ${dataset['name']} : ${dataset['description']}')
            .toList();
      }
      if (response.statusCode == 404) {
        return [];
      } else {
        throw Exception('Failed to load datasets');
      }
    } catch (e) {
      return [];
    }
  }

// Upload a dataset
  Future<bool> importDataset(
      int userId, String name, String description, Uint8List content) async {
    try {
      String uploadUrl = '$baseUrl/import-dataset';

      FormData formData = FormData.fromMap({
        'user_id': userId.toString(),
        'name': name,
        'description': description,
        'content': MultipartFile.fromBytes(content, filename: 'dataset.csv')
      });

      Response response = await _dio.post(uploadUrl, data: formData);

      if (response.statusCode == 200) {
        return true;
      } else {
        return false;
      }
    } catch (e) {
      return false;
    }
  }

  Future<http.Response> startPrediction(
      int projectId,
      List<String> selectedTrainDatasets,
      List<String> selectedTestDatasets,
      String selectedModel,
      Uint8List? customModel, // model is nullable
      bool isCustomModel) async {
    String? base64Model;
    if (isCustomModel && customModel != null) {
      base64Model = base64Encode(customModel);
    }

    // Python URL
    var url = Uri.parse('$pyUrl/predict');

    var response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        'project_id': projectId,
        'train_datasets': selectedTrainDatasets,
        'test_datasets': selectedTestDatasets,
        'model': selectedModel,
        'custom_model': base64Model,
        'is_custom_model': isCustomModel,
      }),
    );
    return response;
  }
}
